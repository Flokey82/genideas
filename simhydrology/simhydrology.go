package simhydrology

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"

	"github.com/Flokey82/go_gens/vectors"
	"github.com/fogleman/delaunay"
)

type Heightmap interface {
	Elevation(idx int) float64          // Returns the elevation of the given index.
	IdxToXY(idx int) (float64, float64) // Returns the x and y coordinates of the given index.
	Neighbors(idx int) []int            // Returns the indices of the neighbors of the given index.
	NumRegions() int                    // Returns the total number of points.
	Width() int                         // Returns the width of the heightmap.
	Height() int                        // Returns the height of the heightmap.
}

// WHAT IF: We, instead of using the downhill neighbors to determine the local pressure vector, we
// sort the regions by height and go from top to bottom.
// We sum up the river vector of r to downhill neighbor times the flux. This will give us the overall
// water flow vector for the downhill region. We could use this vector to determine the pressure vector
// for the downhill region by subtracting the downhill direction vector (the direction of the river) from the overall flow vector (the sum of all vectors that flow into the region),
// resulting in a vector that pushes away from the flow direction, which is the direction in which the river would erode the river bank
//
// This way we can kinda impose a sort of "inertia" as we have it in lagrangian fluid simulation on the
// Eulerian fluid simulation grid and instead of only taking in account the immediate neighbors and the
// changes in water flow. If we erode all neighbors with a fraction depending on "how much they lie in the direction of the eroding vector"
// we can simulate the behavior of water in the same way as with a particle (or raindrop)

// Map represents the heightmap of the terrain.
type Map struct {
	Soil            []float64 // Deposited or eroded soil (positive or negative).
	Suspension      []float64 // Suspended soil.
	Flux            []float64 // Flux of water.
	Precipitation   []float64 // Precipitation.
	Downhill        []int     // Downhill neighbor.
	VerticalScaling float64   // Vertical scaling factor.
	Heightmap
	ErosionAmount        float64        // Amount of soil to erode from the river bed.
	BankErosionAmount    float64        // Amount of soil to erode from the river bank.
	BankDepositionAmount float64        // Amount of soil to deposit on the river bank.
	UseSlopeModifier     bool           // Slope influences erosion and deposition. (Flat and high flux = less erosion, more meander)
	FlowVector           []vectors.Vec3 // Flow vector.
	FlowVector2d         []vectors.Vec2 // Flow vector.
}

// NewMap creates a new map.
func NewMap(hm Heightmap) *Map {
	m := &Map{
		Heightmap:            hm,
		VerticalScaling:      100,
		ErosionAmount:        0.05,
		BankErosionAmount:    0.05,
		BankDepositionAmount: 0.05,
		UseSlopeModifier:     false,
	}
	m.Reset()
	// Fill sinks.
	m.Soil = m.fillSinks(true)

	for i := 0; i < 400; i++ {
		// TODO: Copy to soil (minus the elevation).
		m.generateDownhill()
		m.calculateFlux()
		m.ExportErosionRatePNG(fmt.Sprintf("test_erosion_rate_%d.png", i))
		m.ExportPNG(fmt.Sprintf("test_%d.png", i))
		m.ExportFluxPNG(fmt.Sprintf("test_flux_%d.png", i))

		// Erode.
		m.Soil = m.erode()
		m.generateDownhill()
		m.calculateFlux()

		diff := m.calculateSuspension()
		// Write the difference to a file.
		exportFloatSliceToPNG(fmt.Sprintf("test_diff_%d.png", i), m.Width(), m.Height(), diff)
		// diff := m.calculateSuspension()
		// Write the difference to a file.
		//exportFloatSliceToPNG(fmt.Sprintf("test_diff_%d.png", i), m.Width(), m.Height(), diff)

		m.generateDownhill()
		m.calculateFlux()

		// Fill sinks.
		m.Soil = m.fillSinks(true)
		//exportFloatSliceToPNG(fmt.Sprintf("test_suspension_%d.png", i), m.Width(), m.Height(), m.Suspension)
	}
	return m
}

// Reset resets the map.
func (m *Map) Reset() {
	m.Soil = make([]float64, m.NumRegions())
	m.Suspension = make([]float64, m.NumRegions())
	m.Flux = make([]float64, m.NumRegions())
	m.Precipitation = make([]float64, m.NumRegions())
	m.Downhill = make([]int, m.NumRegions())
	for i := range m.Soil {
		m.Soil[i] = 0
		m.Flux[i] = 0
		m.Precipitation[i] = 0
		m.Downhill[i] = -1
	}
}

// Elevation returns the elevation of the given index (including soil).
func (m *Map) Elevation(idx int) float64 {
	return m.Heightmap.Elevation(idx) + m.Soil[idx]
}

func (m *Map) generateDownhill() {
	doLog := false
	newDownhill := make([]int, m.NumRegions())
	for i := range m.Downhill {
		newDownhill[i] = m.downhill(i)

		// Check if the downhill neighbor has changed.
		if newDownhill[i] != m.Downhill[i] && doLog {
			log.Printf("downhill[%d] = %d -> %d", i, m.Downhill[i], newDownhill[i])
		}
	}
	m.Downhill = newDownhill
}

// downhill returns the index of the downhill neighbor.
// If there is no downhill neighbor, -1 is returned, which indicates that the current point is a sink.
func (m *Map) downhill(idx int) int {
	neighbors := m.Neighbors(idx)
	minIdx := -1
	minElevation := m.Elevation(idx)
	for _, n := range neighbors {
		if m.Elevation(n) < minElevation {
			minIdx = n
			minElevation = m.Elevation(n)
		}
	}
	return minIdx
}

// FillSinks is an implementation of the algorithm described in
// https://www.researchgate.net/publication/240407597_A_fast_simple_and_versatile_algorithm_to_fill_the_depressions_of_digital_elevation_models
// and a partial port of the implementation in:
// https://github.com/Rob-Voss/Learninator/blob/master/js/lib/Terrain.js
//
// Returns a slice containing the height difference of each region compared to the
// original heightmap.
//
// If randEpsilon is true, a randomized epsilon value is added to the elevation
// during each iteration. This is to prevent the algorithm from being too
// uniform.
func (m *Map) fillSinks(randEpsilon bool) []float64 {
	inf := math.Inf(0)
	baseEpsilon := 1.0 / float64(m.NumRegions())
	newHeight := make([]float64, m.NumRegions())
	// Usually you'd set some of the elevation to be below sea level.
	log.Println("Warning: Hacky fix for fillSinks.")
	minElev, maxElev := minMaxElevation(m)
	elevThreshold := minElev + (maxElev-minElev)*0.001
	for i := range newHeight {
		// NOTE: Originally this was <= 0, but it seems like this algorithm
		// fails if we have no regions below or at sea level.
		if m.Elevation(i) <= elevThreshold {
			// Set the elevation at or below sea level to the current
			// elevation.
			newHeight[i] = m.Elevation(i)
		} else {
			// Set the elevation above sea level to infinity.
			newHeight[i] = inf
		}
	}

	// Loop until no more changes are made.
	var epsilon float64
	for {
		if randEpsilon {
			// Variation.
			//
			// In theory we could use noise or random values to slightly
			// alter epsilon here. It should still work, albeit a bit slower.
			// The idea is to make the algorithm less destructive and more
			// natural looking.
			//
			// NOTE: I've decided to use m.rand.Float64() instead of noise.
			epsilon = baseEpsilon * rand.Float64()
		}
		changed := false

		// By shuffling the order in which we parse regions,
		// we ensure a more natural look.
		for _, r := range rand.Perm(m.NumRegions()) {
			// Skip all regions that have the same elevation as in
			// the current heightmap.
			if newHeight[r] == m.Elevation(r) {
				continue
			}

			// Iterate over all neighbors.
			// NOTE: This used to be in a random order, but that
			// had a high cost, so I dropped it for now.
			for _, nb := range m.Neighbors(r) {
				// Since we have set all inland regions to infinity,
				// we will only succeed here if the newHeight of the neighbor
				// is either below sea level or if the newHeight has already
				// been set AND if the elevation is higher than the neighbors.
				//
				// This means that we're working our way inland, starting from
				// the coast, comparing each region with the processed / set
				// neighbors (that aren't set to infinity) in the new heightmap
				// until we run out of regions that need change.
				if m.Elevation(r) >= newHeight[nb]+epsilon {
					newHeight[r] = m.Elevation(r)
					changed = true
					break
				}

				// If we reach this point, the neighbor in the new heightmap
				// is higher than the current elevation of 'r'.
				// This can mean two things. Either the neighbor is set to infinity
				// or the current elevation might indicate a sink.

				// So we check if the newHeight of r is larger than the
				// newHeight of the neighbor (plus epsilon), which will ensure that
				// the newHeight of neighbor is not set to infinity.
				//
				// Additionally we check if the newHeight of the neighbor
				// is higher than the current height of r, which ensures that if the
				// current elevation indicates a sink, we will fill up the sink to the
				// new neighbor height plus epsilon.
				//
				// TODO: Simplify this comment word salad.
				oh := newHeight[nb] + epsilon
				if newHeight[r] > oh && oh > m.Elevation(r) {
					newHeight[r] = oh
					changed = true
				}
			}
		}
		if !changed {
			break
		}
	}

	// Calculate the remaining soil after erosion.
	newSoil := make([]float64, m.NumRegions())
	for i := range newSoil {
		originalElevation := m.Heightmap.Elevation(i)
		if newHeight[i] < 0 {
			newSoil[i] = -originalElevation
		} else {
			newSoil[i] = newHeight[i] - originalElevation
		}
	}
	return newSoil
}

// calculateFlux calculates the flux of water for each region.
// We create a slice of all regions, sort them by elevation from high to low
// and then iterate over them. For each region we add the precipitation and
// add the flux value to the downhill neighbor.
func (m *Map) calculateFlux() {
	m.FlowVector = make([]vectors.Vec3, m.NumRegions())
	m.FlowVector2d = make([]vectors.Vec2, m.NumRegions())
	// Create a slice of all regions.
	regions := make([]int, m.NumRegions())
	for r := range regions {
		regions[r] = r
		m.Flux[r] = 1 / float64(m.NumRegions())

		// If the current region is a sink, we're done.
		if m.Downhill[r] == -1 {
			continue
		}
		m.FlowVector[r] = m.regionVector(r, m.Downhill[r]).Normalize().Mul(m.Flux[r])
		m.FlowVector2d[r] = m.regionVector2d(r, m.Downhill[r]).Normalize().Mul(m.Flux[r])
	}

	// Sort the regions by elevation from high to low.
	sortRegionsByElevation(regions, m)

	// Iterate over all regions.
	for _, r := range regions {
		// If the current region is a sink, we're done.
		if m.Downhill[r] == -1 {
			continue
		}

		// Add the flux to the downhill neighbor.
		m.Flux[m.Downhill[r]] += m.Flux[r]

		// Add the vector to the downhill neighbor.
		m.FlowVector[m.Downhill[r]] = m.FlowVector[m.Downhill[r]].Add(m.FlowVector[r])
		m.FlowVector2d[m.Downhill[r]] = m.FlowVector2d[m.Downhill[r]].Add(m.FlowVector2d[r])
	}
}

// calculateSuspension attempts to simulate sediment erosion and suspension to generate
// river banks and meandering. This variant inspects the change in direction of the river
// between the current region (uphill neighbor and r) and r and the downhill neighbor.
// This variant uses 2d vectors to determine the force of the water pressing against the
// river banks and therefore in what direction the river would erode the river banks (with
// depositing soil on the opposite side).
func (m *Map) calculateSuspension() []float64 {
	// WARNING: This code doesn't work yet.
	// First calculate erosion rate.
	// erosionRate := m.erosionRate()

	// Get slope.
	// slope := m.calcSlopes()

	// Create a slice of all regions.
	regions := make([]int, m.NumRegions())
	for i := range regions {
		regions[i] = i
	}

	// Copy the soil before erosion.
	soilBefore := make([]float64, m.NumRegions())
	copy(soilBefore, m.Soil)

	soilNew := make([]float64, m.NumRegions())
	copy(soilNew, m.Soil)
	//soilNew = m.Soil

	erosionAmounts := make([]float64, m.NumRegions())

	// Sort the regions by erosion rate from high to low.
	sortRegionsByElevation(regions, m)

	// Iterate over all regions, calculate the amount of soil that can be suspended
	// and if we exceed the amount of soil that can be suspended, we distribute
	// the soil to the neighbors. Then we add the suspended soil to the suspension
	// of the downhill neighbor. We do that from high to low erosion rate.
	for _, r := range regions {
		// m.Soil[r] += toDeposit
		// Get the vector of region to downhill neighbor.
		if m.Downhill[r] != -1 && m.Downhill[m.Downhill[r]] != -1 {
			// Get vector to downhill neighbor.
			vecDownhill := m.regionVector2d(r, m.Downhill[r]).Normalize()

			// TODO: Find a better way to determine the deposition region.
			// Usually the river will eat into the terrain when changing direction,
			// and deposit soil on the oppsite side.
			//          ,- erosion
			//       ,---.
			// / /  / ,-. \
			// \ \_/ / | \ \
			//  ` - ´   deposition
			//    `- erosion
			//
			// So in theory we need to find the downhill neighbor region that has the
			// highest dot product (which would be the region that the river would
			// naturally eat into) and the region with the smallest dot product (which
			// would be the region that the river would naturally deposit soil on).
			//
			// We calculate the dot product using a vector that we calculate using
			// V1 = r -> dh[r] and V2 = flowVector[r]
			//
			//  \ /    ___
			//   r    /   \
			//    \  /     \    V1 = r -> dh[r]
			//     dh[r]
			//
			//  \ /    ___
			//   |    /   \
			//   |\  /     \    V2 = flowVector[r]
			//   V V2
			//
			// V3 = V2 - V1
			//
			// The resulting vector will point to the direction of the outward force
			// generated by the river.
			// Then we need to find the downhill neighbor's downhill neighbor that
			// has the highest and lowest dot product with V3.
			//
			//    /    ___
			//   r    /   \
			//    \n1/     \    V4 = dh[r] -> n1 (n2, n3, n4 ...)
			//  n4 ** n2        minDot = V(dh[r]->n1)
			//	   n3           maxDot = V(dh[r]->n3)

			// Calculate V1, V2, and the pressure / force vector (V3).
			// V1 = r -> dh[r]
			v1 := vecDownhill.Normalize()
			// V2 = flowVector[r]
			v2 := m.FlowVector2d[r].Normalize()
			// V3 = V2 - V1
			v3 := v2.Sub(v1).Normalize()

			dotV1V2 := 1 - math.Abs(v1.Dot(v2))

			// calculate angle between v1 and v2 and log it along with the dot product.
			dotLimited := math.Max(-1, math.Min(1, v1.Dot(v2)))
			angle := math.Acos(dotLimited) * 180 / math.Pi
			log.Printf("angle %f, dot %f", angle, dotV1V2)

			// Get the neighbors.
			nbs := m.Neighbors(r)
			nbDots := make([]float64, len(nbs))
			for i, nb := range nbs {
				// Skip the downhill region.
				if nb == m.Downhill[r] {
					continue
				}
				// Get the vector from the downhill neighbor to the current neighbor.
				vec := m.regionVector2d(r, nb).Normalize()

				// Calculate the dot product between the the pressure vectur v3
				// and the vector from the downhill neighbor to the neighbor 'nb'.
				//
				// If they point in the same direction, the force of the river
				// will erode the soil of the neighbor. If the dot product is
				// <= 0, no erosion will occur.
				dot := vec.Dot(v3)

				// Store the dot product.
				nbDots[i] = dot
			}

			// TODO: Maybe I need to erode / deposit soil for all neighbors, not just the
			// maximum and minimum dot product neighbors?

			// The secondary mode applies erosion and deposition to all neighbors, depending
			// on the dot product between the force vector and the vector to the neighbor.
			// Enable excess soil dumping?
			// Iterate over all neighbors.
			for i, nb := range nbs {
				// If the dot product is 0, we skip the neighbor.
				if nbDots[i] == 0 {
					continue
				}

				// Since we want to erode river banks, we only look at neighbors that are
				// higher than the current region.
				heightDiff := m.Elevation(nb) - m.Elevation(r)
				if heightDiff <= 0 {
					continue
				}

				// If the dot product is positive, we erode soil from the neighbor 'nb'
				// and deposit it on the current region 'm.Downhill[r]'.
				// This simulates the river eating into the terrain at the bend where the
				// river changes direction.
				if nbDots[i] > 0 {
					// If the dot product is positive, we erode soil from the current region.
					erosionFactor := nbDots[i] * m.Flux[r] * dotV1V2 // * erosionRate[r]

					// TODO: Take in account slope. The less steep the slope, the more soil we erode while meandering.
					erodeAmount := heightDiff * erosionFactor

					totalErodeAmount := erodeAmount // * m.BankErosionAmount
					if totalErodeAmount > heightDiff {
						totalErodeAmount = heightDiff
					}
					totalDepositAmount := erodeAmount // * m.BankDepositionAmount
					if totalDepositAmount > heightDiff {
						totalDepositAmount = heightDiff
					}

					// Log percentage of erosion and deposition.
					log.Printf("%.2f%% erosion, %.2f%% deposition", totalErodeAmount/heightDiff*100, totalDepositAmount/heightDiff*100)

					log.Printf("dot %f, erosionFactor %f, erodeAmount %f, totalErodeAmount %f, totalDepositAmount %f heightdiff %f", nbDots[i], erosionFactor, erodeAmount, totalErodeAmount, totalDepositAmount, heightDiff)

					soilNew[nb] -= totalErodeAmount
					soilNew[r] += totalDepositAmount

					erosionAmounts[nb] += dotV1V2
					//erosionAmounts[r] -= totalDepositAmount
				}
			}
			// TODO: Deposit excess suspended soil!
		}
	}

	// Calculate the difference in soil before and after erosion.
	soilDiff := make([]float64, m.NumRegions())
	for i := range soilDiff {
		soilDiff[i] = soilNew[i] - soilBefore[i]
	}

	m.Soil = soilNew

	// Return the difference in soil.
	return erosionAmounts
}

// calculateSuspension2 attempts to simulate sediment erosion and suspension to generate
// river banks and meandering. This variant inspects the change in direction of the river
// between the current region and the downhill neighbor and the downhill neighbor's downhill
// neighbor. If the river changes direction, we assume that the river will eat into the
// terrain and deposit soil on the opposite side.
func (m *Map) calculateSuspension2() []float64 {
	// WARNING: This code doesn't work yet.

	secondary := true
	// First calculate erosion rate.
	erosionRate := m.erosionRate()

	// Get slope.
	// slope := m.calcSlopes()

	// Create a slice of all regions.
	suspension := make([]float64, m.NumRegions())
	regions := make([]int, m.NumRegions())
	for i := range regions {
		regions[i] = i

		// Start with the eroded amount in suspension for each region.
		suspension[i] = erosionRate[i] * m.BankErosionAmount
	}

	// Copy the soil before erosion.
	soilBefore := make([]float64, m.NumRegions())
	copy(soilBefore, m.Soil)

	soilNew := make([]float64, m.NumRegions())
	copy(soilNew, m.Soil)

	slopeModValues := make([]float64, m.NumRegions())

	erosionAmounts := make([]float64, m.NumRegions())

	// Sort the regions by erosion rate from high to low.
	sortRegionsByElevation(regions, m)

	// Iterate over all regions, calculate the amount of soil that can be suspended
	// and if we exceed the amount of soil that can be suspended, we distribute
	// the soil to the neighbors. Then we add the suspended soil to the suspension
	// of the downhill neighbor. We do that from high to low erosion rate.
	for _, r := range regions {
		// m.Soil[r] += toDeposit
		// Get the vector of region to downhill neighbor.
		if m.Downhill[r] != -1 && m.Downhill[m.Downhill[r]] != -1 {
			// Get vector to downhill neighbor.
			vecDownhill := m.regionVector(r, m.Downhill[r]).Normalize()

			// Get vector to downhill neighbor's downhill neighbor.
			vec2Downhill := m.regionVector(m.Downhill[r], m.Downhill[m.Downhill[r]]).Normalize()

			// Calculate the dot product between the two vectors.
			// This will give us an indication as to how much the river changes direction.
			// If both vectors point in the same direction, the dot product will be 1.
			dot := vecDownhill.Dot(vec2Downhill)

			// Now calculate how strong the effect is based on how much the river changes direction.
			// No change would mean a dot product of 1, with the maximum change being -1 (opposite direction).
			// In order to get the strength of the effect, we subtract the dot product from 1 and divide by 2.
			// This gives us the strength as a value between 0 (weakest) and 1 (strongest).
			dotMeanderStrength := (1.0 - dot) / 2.0
			if dotMeanderStrength < 0 {
				dotMeanderStrength = 0
			}
			log.Printf("dot = %f, dotStrength = %f", dot, dotMeanderStrength)

			// TODO: Find a better way to determine the deposition region.
			// Usually the river will eat into the terrain when changing direction,
			// and deposit soil on the oppsite side.
			//          ,- erosion
			//       ,---.
			// / /  / ,-. \
			// \ \_/ / | \ \
			//  ` - ´   deposition
			//    `- erosion
			//
			// So in theory we need to find the downhill neighbor region that has the
			// highest dot product (which would be the region that the river would
			// naturally eat into) and the region with the smallest dot product (which
			// would be the region that the river would naturally deposit soil on).
			//
			// We calculate the dot product using a vector that we calculate using
			// V1 = r -> dh[r] and V2 =  dh[dh[r]] -> dh[r].
			//
			//    /    ___
			//   r    /   \
			//    \  /     \    V1 = r -> dh[r]
			//     dh[r]
			//
			//    /    ___
			//   |    dh[dh[r]]
			//    \  /     \    V2 = dh[dh[r]] -> dh[r]
			//     dh[r]
			//
			// V3 = V1 + V2
			//
			// The resulting vector will point to the direction of the outward force
			// generated by the river.
			// Then we need to find the downhill neighbor's downhill neighbor that
			// has the highest and lowest dot product with V3.
			//
			//    /    ___
			//   r    /   \
			//    \n1/     \    V4 = dh[r] -> n1 (n2, n3, n4 ...)
			//  n4 ** n2        minDot = V(dh[r]->n1)
			//	   n3           maxDot = V(dh[r]->n3)

			minDotRegion := -1
			minDot := 1.0
			maxDotRegion := -1
			maxDot := -1.0

			// Calculate V1, V2, and the pressure / force vector (V3).
			// V1 = r -> dh[r]
			v1 := vecDownhill
			// V2 = dh[dh[r]] -> dh[r]
			v2 := m.regionVector(m.Downhill[m.Downhill[r]], m.Downhill[r]).Normalize()
			// V3 = V1 + V2
			v3 := v1.Add(v2).Normalize()

			// Get the downhill neighbor's neighbors.
			nbs := m.Neighbors(m.Downhill[r])
			nbDots := make([]float64, len(nbs))
			for i, nb := range nbs {
				// Skip the current region.
				if nb == r || nb == m.Downhill[m.Downhill[r]] {
					continue
				}
				// Get the vector from the downhill neighbor to the current neighbor.
				vec := m.regionVector(m.Downhill[r], nb).Normalize()

				// Calculate the dot product between the the pressure vectur v3
				// and the vector from the downhill neighbor to the neighbor 'nb'.
				//
				// If they point in the same direction, the force of the river
				// will erode the soil of the neighbor. If the dot product is
				// <= 0, no erosion will occur.
				dot := vec.Dot(v3)

				// Store the dot product.
				nbDots[i] = dot

				// Check if the dot product is smaller than the current minimum.
				if dot < minDot {
					minDot = dot
					minDotRegion = nb
				}

				// Check if the dot product is larger than the current maximum.
				if dot > maxDot {
					maxDot = dot
					maxDotRegion = nb
				}
			}

			log.Printf("minDot = %f, maxDot = %f", minDot, maxDot)

			// TODO: Maybe I need to erode / deposit soil for all neighbors, not just the
			// maximum and minimum dot product neighbors?

			// The secondary mode applies erosion and deposition to all neighbors, depending
			// on the dot product between the force vector and the vector to the neighbor.
			if secondary {
				// Enable excess soil dumping?
				doDeposit := false

				// We calculate the amount of soil that can be suspended, which is the maximum erosion rate.
				suspendable := erosionRate[r]
				toDeposit := suspension[r] - suspendable

				// Iterate over all neighbors.
				for i, nb := range nbs {
					// If the dot product is 0, we skip the neighbor.
					if nbDots[i] == 0 {
						continue
					}

					// Since we want to erode river banks, we only look at neighbors that are
					// higher than the current region.
					heightDiff := m.Elevation(nb) - m.Elevation(m.Downhill[r])
					if heightDiff < 0 {
						continue
					}

					// If the dot product is positive, we erode soil from the neighbor 'nb'
					// and deposit it on the current region 'm.Downhill[r]'.
					// This simulates the river eating into the terrain at the bend where the
					// river changes direction.
					if nbDots[i] > 0 {
						slopeMod := 1.0

						// If the flux is high and the slope is low, we start to erode the river banks.
						if m.UseSlopeModifier {
							//* (1 - math.Pow(m.slope(m.Downhill[r]), 0.25))
							slopeMod = 1 - math.Pow(1-((1-m.slope(m.Downhill[r]))+m.Flux[m.Downhill[r]])/2, 0.5)

							slopeModValues[nb] += dotMeanderStrength
							//slopeModValues[nb] += nbDots[i]
							//log.Printf("slopeMod = %f slope = %f flux = %f", slopeMod, m.slope(m.Downhill[r]), m.Flux[m.Downhill[r]])
						}

						// If the dot product is positive, we erode soil from the current region.
						erosionFactor := math.Sqrt(nbDots[i]) * dotMeanderStrength
						// erosionFactor := math.Sqrt(nbDots[i]) * math.Sqrt(dotStrength)
						// erosionFactor := nbDots[i] * dotStrength

						// TODO: Take in account slope. The less steep the slope, the more soil we erode while meandering.
						erodeAmount := erosionRate[r] * erosionFactor * slopeMod

						// Get the height difference between the current region and the neighbor.
						heightDiff := m.Elevation(nb) - m.Elevation(m.Downhill[r])

						totalErodeAmount := erodeAmount * m.BankErosionAmount
						if totalErodeAmount > heightDiff/2 {
							totalErodeAmount = heightDiff / 2
						}
						totalDepositAmount := erodeAmount * m.BankDepositionAmount
						if totalDepositAmount > heightDiff/2 {
							totalDepositAmount = heightDiff / 2
						}

						// Calculate the total difference we remove through both erosion and deposition.
						// totalDiff := totalErodeAmount + totalDepositAmount

						// Log what percentage of the height difference we're eroding.
						// log.Printf("Eroding %f%% of height difference between %d and %d (difference before %f, after %f)", totalDiff/heightDiff*100, m.Downhill[r], nb, heightDiff, heightDiff-totalDiff)

						soilNew[nb] -= totalErodeAmount
						soilNew[m.Downhill[r]] += totalDepositAmount

						erosionAmounts[nb] += totalErodeAmount
						erosionAmounts[m.Downhill[r]] -= totalDepositAmount
						// log.Println("Erode", nb, "by", erodeAmount, "factor", erosionFactor, "erosionRate", erosionRate[m.Downhill[r]]*m.ErosionAmount)
						// log.Println("Deposit", m.Downhill[r], "by", erodeAmount)
					} else if toDeposit > 0 && doDeposit {
						// If the dot product is negative, we deposit soil on the current region
						// and erode it from the neighbor.
						depositionFactor := math.Sqrt(-nbDots[i]) * math.Sqrt(dotMeanderStrength)
						// depositionFactor := math.Sqrt(-nbDots[i]) * math.Sqrt(dotStrength)
						depositionAmount := toDeposit * depositionFactor
						soilNew[nb] += depositionAmount * m.BankDepositionAmount
						log.Println("Deposit", nb, "by", depositionAmount*dotMeanderStrength*0.1)
					}
				}
				// TODO: Deposit excess suspended soil!
			} else {
				// The maximum dot product is the region that the river would naturally
				// eat into.
				if maxDotRegion != -1 && maxDot > 0 {
					// We take the erosion rate at the downhill neighbor and multiply it by the
					// dot product representing how closely the vector points in the direction
					// of the outward force generated by the river.
					// Then we multiply it by the erosion amount and a factor that represents
					// how extreme the change in direction is.
					erodeAmount := erosionRate[m.Downhill[r]] * maxDot * m.BankErosionAmount * dotMeanderStrength // * 0.1
					soilNew[maxDotRegion] -= erodeAmount                                                          // Remove soil from the region that the river would naturally eat into.
					soilNew[m.Downhill[r]] += erodeAmount                                                         // Add the soil to the downhill neighbor.
					//toDeposit += erodeAmount
					log.Println("Erode", maxDotRegion, "by", erodeAmount*dotMeanderStrength*0.1)
				}

				// We calculate the amount of soil that can be suspended, which is the maximum erosion rate.
				suspendable := erosionRate[r]

				// If the amount of soil that can be suspended is less than the amount of
				// soil that is already suspended, we simply dump the excess soil.
				if suspendable < suspension[r] {
					toDeposit := suspension[r] - suspendable
					log.Printf("Dumping %f from %d", toDeposit, r)

					// The minimum dot product is the region that the river would naturally
					// deposit soil on.
					if minDotRegion != -1 && minDot < 0 {
						depositionAmount := toDeposit * dotMeanderStrength * 0.1 * m.BankDepositionAmount * (1 - minDot) / 2
						soilNew[minDotRegion] += depositionAmount
						log.Println("Deposit", minDotRegion, "by", depositionAmount*dotMeanderStrength*0.1)
					}

					// The more the river changes direction, the more soil we deposit.
					// The dot product is between -1 and 1, so we add 1 and divide by 2
					// to get a value between 0 and 1, which will determine how much soil
					// we deposit.
					//soilNew[depositionRegion] += toDeposit * (1 - dot) / 2
					// suspendable += toDeposit * dot
					//log.Printf("dot = %f, deposit %f", dot, dot*toDeposit)
					suspension[r] = suspendable
				}
			}
		}

		// If the amount of soil that can be suspended is less than the amount of
		// soil that is already suspended, we distribute the soil to the neighbors with
		// lower flux
		/*
			if suspendable < suspension[r] {
				log.Printf("Distributing %f to neighbors of %d", suspension[r]-suspendable, r)
				rem := (suspension[r] - suspendable) * m.DepositionAmount // Distribute 1% of the suspended soil.
				var nbs []int
				for _, nb := range m.Neighbors(r) {
					// Check if the neighbor has a lower flux and is lower than the current region.
					if m.Flux[nb] < m.Flux[r] && m.Elevation(nb) < m.Elevation(r) {
						nbs = append(nbs, nb)
					}
				}
				for _, nb := range nbs {
					soilNew[nb] += rem
				}
				suspension[r] = suspendable
			}
		*/

		// Add the suspended soil to the downhill neighbor.
		if m.Downhill[r] != -1 {
			suspension[m.Downhill[r]] += suspension[r]
		}
	}
	m.Suspension = suspension

	// Calculate the difference in soil before and after erosion.
	soilDiff := make([]float64, m.NumRegions())
	for i := range soilDiff {
		soilDiff[i] = soilNew[i] - soilBefore[i]
	}

	m.Soil = soilNew

	// Return the difference in soil.
	return erosionAmounts
}

// calculateSuspension3 attempts to simulate sediment erosion and suspension to generate
// river banks and meandering. This variant is very similar to calculateSuspension, but
// uses 3d vectors to determine the force of the water pressing against the river banks
// and therefore in what direction the river would erode the river banks (with depositing
// soil on the opposite side).
func (m *Map) calculateSuspension3() []float64 {
	// WARNING: This code doesn't work yet.
	// First calculate erosion rate.
	erosionRate := m.erosionRate()

	// Get slope.
	// slope := m.calcSlopes()

	// Create a slice of all regions.
	regions := make([]int, m.NumRegions())
	for i := range regions {
		regions[i] = i
	}

	// Copy the soil before erosion.
	soilBefore := make([]float64, m.NumRegions())
	copy(soilBefore, m.Soil)

	soilNew := make([]float64, m.NumRegions())
	copy(soilNew, m.Soil)
	//soilNew = m.Soil

	erosionAmounts := make([]float64, m.NumRegions())

	// Sort the regions by erosion rate from high to low.
	sortRegionsByElevation(regions, m)

	// Iterate over all regions, calculate the amount of soil that can be suspended
	// and if we exceed the amount of soil that can be suspended, we distribute
	// the soil to the neighbors. Then we add the suspended soil to the suspension
	// of the downhill neighbor. We do that from high to low erosion rate.
	for _, r := range regions {
		// m.Soil[r] += toDeposit
		// Get the vector of region to downhill neighbor.
		if m.Downhill[r] != -1 && m.Downhill[m.Downhill[r]] != -1 {
			// Get vector to downhill neighbor.
			vecDownhill := m.regionVector(r, m.Downhill[r]).Normalize()

			// TODO: Find a better way to determine the deposition region.
			// Usually the river will eat into the terrain when changing direction,
			// and deposit soil on the oppsite side.
			//          ,- erosion
			//       ,---.
			// / /  / ,-. \
			// \ \_/ / | \ \
			//  ` - ´   deposition
			//    `- erosion
			//
			// So in theory we need to find the downhill neighbor region that has the
			// highest dot product (which would be the region that the river would
			// naturally eat into) and the region with the smallest dot product (which
			// would be the region that the river would naturally deposit soil on).
			//
			// We calculate the dot product using a vector that we calculate using
			// V1 = r -> dh[r] and V2 = flowVector[r]
			//
			//  \ /    ___
			//   r    /   \
			//    \  /     \    V1 = r -> dh[r]
			//     dh[r]
			//
			//  \ /    ___
			//   |    /   \
			//   |\  /     \    V2 = flowVector[r]
			//   V V2
			//
			// V3 = V2 - V1
			//
			// The resulting vector will point to the direction of the outward force
			// generated by the river.
			// Then we need to find the downhill neighbor's downhill neighbor that
			// has the highest and lowest dot product with V3.
			//
			//    /    ___
			//   r    /   \
			//    \n1/     \    V4 = dh[r] -> n1 (n2, n3, n4 ...)
			//  n4 ** n2        minDot = V(dh[r]->n1)
			//	   n3           maxDot = V(dh[r]->n3)

			// Calculate V1, V2, and the pressure / force vector (V3).
			// V1 = r -> dh[r]
			v1 := vecDownhill
			v1.Z = 0
			v1 = v1.Normalize()
			// V2 = flowVector[r]
			v2 := m.FlowVector[r]
			v2.Z = 0
			v2 = v2.Normalize()
			// V3 = V2 - V1
			v3 := v2.Sub(v1)
			// Remove the z component so we can compare directions on a 2D plane.
			v3.Z = 0
			if v3.Len() != 0 {
				v3 = v3.Normalize()
			}

			// Get the neighbors.
			nbs := m.Neighbors(r)
			nbDots := make([]float64, len(nbs))
			for i, nb := range nbs {
				// Skip the downhill region.
				if nb == m.Downhill[r] {
					continue
				}
				// Get the vector from the downhill neighbor to the current neighbor.
				vec := m.regionVector(r, nb)
				vec.Z = 0
				vec = vec.Normalize()

				// Calculate the dot product between the the pressure vectur v3
				// and the vector from the downhill neighbor to the neighbor 'nb'.
				//
				// If they point in the same direction, the force of the river
				// will erode the soil of the neighbor. If the dot product is
				// <= 0, no erosion will occur.
				dot := vec.Dot(v3)

				// Store the dot product.
				nbDots[i] = dot
			}

			// TODO: Maybe I need to erode / deposit soil for all neighbors, not just the
			// maximum and minimum dot product neighbors?

			// The secondary mode applies erosion and deposition to all neighbors, depending
			// on the dot product between the force vector and the vector to the neighbor.
			// Enable excess soil dumping?
			// Iterate over all neighbors.
			for i, nb := range nbs {
				// If the dot product is 0, we skip the neighbor.
				if nbDots[i] == 0 {
					continue
				}

				// Since we want to erode river banks, we only look at neighbors that are
				// higher than the current region.
				heightDiff := m.Elevation(nb) - m.Elevation(r)
				if heightDiff < 0 {
					continue
				}

				// If the dot product is positive, we erode soil from the neighbor 'nb'
				// and deposit it on the current region 'm.Downhill[r]'.
				// This simulates the river eating into the terrain at the bend where the
				// river changes direction.
				if nbDots[i] > 0 {
					// If the dot product is positive, we erode soil from the current region.
					erosionFactor := nbDots[i] //math.Sqrt(nbDots[i])

					// TODO: Take in account slope. The less steep the slope, the more soil we erode while meandering.
					erodeAmount := erosionRate[r] * erosionFactor

					totalErodeAmount := erodeAmount * m.BankErosionAmount
					if totalErodeAmount > heightDiff/2 {
						totalErodeAmount = heightDiff / 2
					}
					totalDepositAmount := erodeAmount * m.BankDepositionAmount
					if totalDepositAmount > heightDiff/2 {
						totalDepositAmount = heightDiff / 2
					}

					log.Printf("dot %f, erosionFactor %f, erodeAmount %f, totalErodeAmount %f, totalDepositAmount %f heightdiff %f", nbDots[i], erosionFactor, erodeAmount, totalErodeAmount, totalDepositAmount, heightDiff)

					soilNew[nb] -= totalErodeAmount
					soilNew[r] += totalDepositAmount

					erosionAmounts[nb] += totalErodeAmount
					erosionAmounts[r] -= totalDepositAmount
				}
			}
			// TODO: Deposit excess suspended soil!
		}
	}

	// Calculate the difference in soil before and after erosion.
	soilDiff := make([]float64, m.NumRegions())
	for i := range soilDiff {
		soilDiff[i] = soilNew[i] - soilBefore[i]
	}

	m.Soil = soilNew

	// Return the difference in soil.
	return erosionAmounts
}

func sortRegionsByElevation(regions []int, m *Map) {
	sort.Slice(regions, func(i, j int) bool {
		return m.Elevation(regions[i]) > m.Elevation(regions[j])
	})
}

func (m *Map) calcSlopes() []float64 {
	slopes := make([]float64, m.NumRegions())
	for i := range slopes {
		slopes[i] = m.slope(i)
	}
	return slopes
}

func (m *Map) slope(idx int) float64 {
	// Get the surface normal of the current region and calculate the slope.
	normal := m.surfaceNormal(idx)
	slope := vectors.Vec2{X: normal.X, Y: normal.Y}.Len()
	return slope
}

var zVec = vectors.Vec3{X: 0, Y: 0, Z: 1}

func (m *Map) surfaceNormal(idx int) vectors.Vec3 {
	// Get the neighbors of the current region.
	nbs := m.Neighbors(idx)

	// Calculate the surface normal taking into account the neighbor elevations
	// and the distance between the current region and the neighbors.
	var sumNormals vectors.Vec3
	for _, nb := range nbs {
		// Get the region vector between the current region and the neighbor.
		vecIdxNb := m.regionVector(idx, nb)
		dist := vecIdxNb.Len()
		vecIdxNb = vecIdxNb.Normalize()

		// Get the perpendicular vector of the region vector.
		localNormal := vecIdxNb.Cross(zVec).Normalize()

		// Now calculate the normal by taking the cross product of the perpendicular
		// vector and the region vector.
		normal := localNormal.Cross(vecIdxNb).Normalize()

		// Add the normal to the sum of normals and multiply its length by the distance
		// between the current region and the neighbor.
		sumNormals = sumNormals.Add(normal.Mul(dist))
	}

	// Normalize the normal vector.
	sumNormals.Normalize()

	// If the normal vector is zero, we return the z-vector.
	if sumNormals.X == 0 && sumNormals.Y == 0 && sumNormals.Z == 0 || sumNormals.Len() == 0 {
		return zVec
	}
	return sumNormals
}

// regionVector returns the vector between two regions pointing from A to B.
func (m *Map) regionVector(idxA, idxB int) vectors.Vec3 {
	x1, y1 := m.IdxToXY(idxA)
	x2, y2 := m.IdxToXY(idxB)
	return vectors.Vec3{X: x2 - x1, Y: y2 - y1, Z: m.Elevation(idxB) - m.Elevation(idxA)}
}

// regionVector2d returns the vector between two regions pointing from A to B.
func (m *Map) regionVector2d(idxA, idxB int) vectors.Vec2 {
	x1, y1 := m.IdxToXY(idxA)
	x2, y2 := m.IdxToXY(idxB)
	return vectors.Vec2{X: x2 - x1, Y: y2 - y1}
}

// erode erodes the terrain.
// Returns a slice containing the difference in elevation for each region compared
// to the original heightmap.
// Negative values indicate erosion below the original height, positive values
// indicate deposition above the original height.
func (m *Map) erode() []float64 {
	onlyErodeDelta := false
	erosionRate := m.erosionRate()
	newHeight := make([]float64, m.NumRegions())
	flux := m.Flux
	for i := range newHeight {
		if math.IsNaN(erosionRate[i]) {
			erosionRate[i] = 0
		}
		slopeMod := 1.0
		if m.UseSlopeModifier {
			// With decreasing slope and increasing flux, the erosion rate decreases.
			// Starting at a certain water depth, the erosion rate drops off, especially
			// if the slope is low, leading to shallow, wide rivers.
			slopeMod = math.Pow(1-((1-m.slope(i))+flux[i])/2, 0.5)
			//log.Printf("slopeMod (meander) %f slope %f flux %f", slopeMod, m.slope(i), flux[i])
		}
		delta := 1.0
		if m.Downhill[i] != -1 && onlyErodeDelta {
			// Only erode up to the elevation of the downhill neighbor.
			// This will prevent sinks from forming.
			// Get the delta between current and downhill neighbor.
			delta = m.Elevation(i) - m.Elevation(m.Downhill[i])
		}
		newHeight[i] = m.Elevation(i) - erosionRate[i]*m.ErosionAmount*delta*slopeMod
		if newHeight[i] < 0 {
			newHeight[i] = 0
		}
	}

	// Calculate the remaining soil after erosion.
	newSoil := make([]float64, m.NumRegions())
	for i := range newSoil {
		originalElevation := m.Heightmap.Elevation(i)
		if newHeight[i] < 0 {
			newSoil[i] = -originalElevation
		} else {
			newSoil[i] = newHeight[i] - originalElevation
		}
	}
	return newSoil
}

func (m *Map) erosionRate() []float64 {
	erodeNeighbors := true
	flux := m.Flux
	slope := m.calcSlopes()
	newHeight := make([]float64, m.NumRegions())
	for i := range newHeight {
		// NOTE: This was directly taken from mewo2's code.
		river := math.Sqrt(flux[i]) * slope[i]
		creep := slope[i] * slope[i]
		total := river + creep/1000
		// log.Printf("river = %f, creep = %f, total = %f", river, creep, total)
		if total > 0.200 {
			total = 0.200
		}

		// Apply the a fraction of the erosion partially to the neighbours.
		if erodeNeighbors {
			nbs := m.Neighbors(i)
			for _, nb := range nbs {
				// TODO: Not sure if we should add to the existing erosion rate
				// or rather set it to whatever is the higher value?
				// ... alternatively we could average it out.
				newHeight[nb] += total * 0.25
			}
		}

		newHeight[i] += total
	}

	// TODO: Normalize the erosion rate.

	// TODO:
	// We create a slice of all regions, sort them by elevation from high to low.
	// Depending on the flux in the region, we can calculate how much soil can be
	// suspended. We then distribute the suspended soil to the neighbors, depending
	// on the flux and the difference in elevation.
	// TODO: Steepness, water speed, meandering (uphill->region region->downhill vector dot product).
	return newHeight
}

func normalizeFloatSlice(values []float64) []float64 {
	min, max := minMaxFloat64(values)
	newSlice := make([]float64, len(values))
	for i := range values {
		newSlice[i] = (values[i] - min) / (max - min)
	}
	return newSlice
}

func (m *Map) Triangulate() (*delaunay.Triangulation, error) {
	var pts []delaunay.Point
	for i := 0; i < m.NumRegions(); i++ {
		x, y := m.IdxToXY(i)
		pts = append(pts, delaunay.Point{X: x, Y: y})
	}
	return delaunay.Triangulate(pts)
}

// ExportOBJ returns a Wavefront OBJ file representing the heightmap.
func (m *Map) ExportOBJ(path string) error {
	tr, err := m.Triangulate()
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for i, p := range tr.Points {
		w.WriteString(fmt.Sprintf("v %f %f %f \n", p.X, m.Elevation(i)*m.VerticalScaling, p.Y)) //
	}
	for i := 0; i < len(tr.Triangles); i += 3 {
		w.WriteString(fmt.Sprintf("f %d %d %d \n", tr.Triangles[i]+1, tr.Triangles[i+1]+1, tr.Triangles[i+2]+1))
	}
	return nil
}

// ExportPNG exports the heightmap as a PNG image.
func (m *Map) ExportPNG(path string) error {
	elevations := make([]float64, m.NumRegions())
	for i := range elevations {
		elevations[i] = m.Elevation(i)
	}
	return exportFloatSliceToPNG(path, m.Width(), m.Height(), elevations)
}

func (m *Map) ExportSoilPNG(path string) error {
	return exportFloatSliceToPNG(path, m.Width(), m.Height(), m.Soil)
}

func (m *Map) ExportFluxPNG(path string) error {
	return exportFloatSliceToPNG(path, m.Width(), m.Height(), m.Flux)
}

func (m *Map) ExportSinksPNG(path string) error {
	sinks := make([]float64, m.NumRegions())
	for i := range sinks {
		if m.Downhill[i] == -1 {
			sinks[i] = 1
		}
	}
	return exportFloatSliceToPNG(path, m.Width(), m.Height(), sinks)
}

func (m *Map) ExportErosionRatePNG(path string) error {
	return exportFloatSliceToPNG(path, m.Width(), m.Height(), m.erosionRate())
}

func exportFloatSliceToPNG(path string, width, height int, values []float64) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	vMin, vMax := minMaxFloat64(values)
	img := image.NewGray(image.Rect(0, 0, width, height))
	for i := 0; i < len(values); i++ {
		x := i % width
		y := i / width
		img.Set(x, y, color.Gray{Y: uint8(255 * (values[i] - vMin) / (vMax - vMin))})
	}
	if err := png.Encode(w, img); err != nil {
		return err
	}
	w.Flush()
	return nil
}

func minMaxFloat64(values []float64) (float64, float64) {
	min, max := math.Inf(1), math.Inf(-1)
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}

type elevationFetcher interface {
	Elevation(idx int) float64
	NumRegions() int
}

func minMaxElevation(h elevationFetcher) (float64, float64) {
	min, max := math.Inf(1), math.Inf(-1)
	for i := 0; i < h.NumRegions(); i++ {
		e := h.Elevation(i)
		if e < min {
			min = e
		}
		if e > max {
			max = e
		}
	}
	return min, max
}
