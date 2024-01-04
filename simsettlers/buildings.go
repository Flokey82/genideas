package simsettlers

import (
	"fmt"
	"log"
	"math"
	"math/rand"
)

func (m *Map) tickBuildings() {
	// Any unoccupied buildings have a chance of decaying.
	for _, b := range m.Buildings {
		if !b.IsOccupied() && b.Condition > 0 {
			if rand.Intn(100) < 10 {
				b.Condition--
			}
		}
	}
}

func (m *Map) advanceConstruction() {
	var stillBuilding []*Building

	// Build new buildings.
	for _, b := range m.Construction {
		if b.AdvanceConstruction() {
			log.Printf("Finished building %v", b)
			m.Buildings = append(m.Buildings, b)

			// Go through the new owners and assign them to the building.
			// TODO: This doesn't make sense for anything but houses.
			if b.Type != BuildingTypeHouse {
				continue
			}

			// Assign the new home to the owners.
			for _, p := range b.Owners {
				// Assign the new home to the person.
				prevHome := p.Home
				p.SetHome(b)

				// Remove the building from the constructing list.
				// TODO: Factor this out into a function.
				for i, c := range p.Constructing {
					if c == b {
						p.Constructing = append(p.Constructing[:i], p.Constructing[i+1:]...)
						break
					}
				}

				for _, c := range p.Children {
					// Check if the child still lives with the parents.
					// If so, move the child to the new home.
					// TODO:
					// - What if the child doesn't want to move?
					// - Maybe they want to stay with the grandparents?
					// - What if the child has a spouse living with them?
					// - What if the child has children living with them?
					if c.Home == prevHome {
						c.SetHome(b)
					}
				}
			}
		} else {
			log.Printf("Still building %v", b)
			stillBuilding = append(stillBuilding, b)
		}
	}
	m.Construction = stillBuilding
}

// Building represents a building on the map.
type Building struct {
	X, Y      int       // position on the map
	BuiltDay  uint16    // day the building was built
	BuiltYear int       // year the building was built
	Remaining uint16    // ticks until construction is complete
	Condition byte      // condition of the building
	Type      string    // building type
	Owners    []*Person // people who own the building
	Occupants []*Person // people who live in the building
}

// NewBuilding creates a new building of the given type at the given position.
func NewBuilding(x, y int, t string) *Building {
	return &Building{
		X:         x,
		Y:         y,
		Remaining: uint16(buildingCosts[t]), // The number of ticks remaining until the building is finished.
		Condition: 100,                      // The condition of the building.
		Type:      t,
	}
}

// AddBuilding adds a new building to the map.
func (m *Map) AddBuilding(x, y int, t string) *Building {
	b := NewBuilding(x, y, t)
	b.BuiltDay = m.Day
	b.BuiltYear = m.Year
	if b.Remaining == 0 {
		m.Buildings = append(m.Buildings, b)
	} else {
		m.Construction = append(m.Construction, b)
	}
	return b
}

// PurchasePrice returns the purchase price of the building.
func (b *Building) PurchasePrice() int {
	return int(float64(buildingCosts[BuildingTypeHouse]) * float64(b.Condition) / 100.0)
}

// String returns a string representation of the building.
func (b *Building) String() string {
	return fmt.Sprintf("%v at (%v, %v cond %d)", b.Type, b.X, b.Y, b.Condition)
}

// AdvanceConstruction advances the construction of the given building by one tick
// and returns true if the building is finished.
func (b *Building) AdvanceConstruction() bool {
	b.Remaining--
	return b.Remaining == 0
}

// IsOccupied returns true if the building is occupied.
func (b *Building) IsOccupied() bool {
	return len(b.Occupants) > 0
}

// IsOwned returns true if the building is owned.
func (b *Building) IsOwned() bool {
	return len(b.Owners) > 0
}

// Yield returns the resource yield of the given building in a tick.
func (b *Building) Yield() int {
	return buildingYield[b.Type]
}

// AddOccupant adds the given person to the occupants list of the building (if not already present).
func (b *Building) AddOccupant(p *Person) {
	for _, o := range b.Occupants {
		if o == p {
			return
		}
	}
	b.Occupants = append(b.Occupants, p)
}

// RemoveOccupant removes the given person from the occupants list of the building.
func (b *Building) RemoveOccupant(p *Person) {
	for i, o := range b.Occupants {
		if o == p {
			b.Occupants = append(b.Occupants[:i], b.Occupants[i+1:]...)
			break
		}
	}
}

// AddOwner adds the given person to the owners list of the building (if not already present).
func (b *Building) AddOwner(p *Person) {
	for _, o := range b.Owners {
		if o == p {
			return
		}
	}
	b.Owners = append(b.Owners, p)
	p.Owns = append(p.Owns, b)
}

// RemoveOwner removes the given person from the owners list of the building.
func (b *Building) RemoveOwner(p *Person) {
	for i, o := range b.Owners {
		if o == p {
			b.Owners = append(b.Owners[:i], b.Owners[i+1:]...)
			break
		}
	}
	for i, o := range p.Owns {
		if o == b {
			p.Owns = append(p.Owns[:i], p.Owns[i+1:]...)
			break
		}
	}
}

const (
	BuildingTypeMarket   = "market"
	BuildingTypeHouse    = "house"
	BuildingTypeCemetery = "cemetery"
	BuildingTypeDungeon  = "dungeon"
)

var buildingCosts = map[string]int{
	BuildingTypeMarket: 10,
	BuildingTypeHouse:  5,
}

var buildingCapacity = map[string]int{
	BuildingTypeMarket: 10,
	BuildingTypeHouse:  5,
}

var buildingYield = map[string]int{
	BuildingTypeMarket: 10,
	BuildingTypeHouse:  5,
}

func (m *Map) getHousingCapacity() int {
	capacity := 0
	for _, b := range m.Buildings {
		if b.Type == BuildingTypeHouse {
			capacity += buildingCapacity[b.Type]
		}
	}
	return capacity
}

func (m *Map) getTheoreticalHousingCapacity() int {
	baseCap := m.getHousingCapacity()
	for _, b := range m.Construction {
		if b.Type == BuildingTypeHouse {
			baseCap += buildingCapacity[b.Type]
		}
	}
	return baseCap
}

func distanceToBuilding(x, y int, b *Building) float64 {
	return math.Sqrt(float64((x-b.X)*(x-b.X) + (y-b.Y)*(y-b.Y)))
}

// fitnessScoreMarketProximity returns the fitness score for any point on the map as a value between 0 and 1.
func (m *Map) fitnessScoreMarketProximity(i int) float64 {
	// Add distance to the market. (closer is better)
	x := i % m.Width
	y := i / m.Width
	return 1 - 1.0/(1.0+distanceToBuilding(x, y, m.Root))
}

func (m *Map) fitnessScoreBuildingProximity(i int, buildings []*Building) float64 {
	if len(buildings) == 0 {
		return -1
	}

	// Get distance to the closest house.
	x := i % m.Width
	y := i / m.Width
	dist := math.Inf(1)
	for _, b := range buildings {
		dist = min(dist, distanceToBuilding(x, y, b))
	}
	return dist
}

func (m *Map) fitnessScoreFlux(i int) float64 {
	if m.Flux[i] > fluxRiverThreshold {
		// This cell is not suitable, can't build on water.
		return -1
	}

	// Lower flux is better.
	return 1 - m.Flux[i]
}

func (m *Map) calcSteepness() []float64 {
	// Calculate the steepness of each cell.
	steepness := make([]float64, len(m.Flux))
	for i := range steepness {
		// Calculate the steepness of the cell by comparing the elevation of the cell
		// to the elevation of the surrounding cells.
		elev := m.Elevation[i]
		var maxDiff float64
		for _, n := range m.Neighbors(i%m.Width, i/m.Width) {
			diff := math.Abs(elev - m.Elevation[n])
			if diff > maxDiff {
				maxDiff = diff
			}
		}
		steepness[i] = maxDiff
	}

	return steepness
}

func bellCurve(x, mu, sigma float64) float64 {
	return 1 / (sigma * math.Sqrt(2*math.Pi)) * math.Exp(-(x-mu)*(x-mu)/(2*sigma*sigma))
}

func (m *Map) calcDistanceToBorder() []float64 {
	// Calculate the distance to the border of the map.
	dist := make([]float64, len(m.Flux))
	for i := range dist {
		x := float64(i % m.Width)
		y := float64(i / m.Width)
		x -= float64(m.Width / 2)
		y -= float64(m.Height / 2)

		x /= float64(m.Width / 2)
		y /= float64(m.Height / 2)

		// Now invert the distance, so that the center of the map has the highest score.
		distVal := (1.0 - 1.0/(1.0+math.Sqrt(x*x+y*y)))
		dist[i] = distVal
	}
	return dist
}

func (m *Map) calcDistanceToCenter() []float64 {
	// Calculate the distance to the center of the map.
	dist := make([]float64, len(m.Flux))
	for i := range dist {
		x := float64(i % m.Width)
		y := float64(i / m.Width)
		x -= float64(m.Width / 2)
		y -= float64(m.Height / 2)

		x /= float64(m.Width / 2)
		y /= float64(m.Height / 2)

		distVal := math.Sqrt(x*x + y*y)
		dist[i] = distVal
	}
	return dist
}

func (m *Map) calcFitnessScoreDungeon() []float64 {
	// Calculate the fitness score for each point.
	// We want to put the dungeon at a higher elevation or the foot at a mounain.
	// We might calculate the steepness of a tile and use that as a fitness score.
	fitness := make([]float64, len(m.Flux))
	for i := range fitness {
		fitness[i] = -1
	}

	steepness := m.calcSteepness()

	// distToBorder := m.calcDistanceToBorder()
	// distToCenter := m.calcDistanceToCenter()

	for i := range fitness {
		// Can't build on water.
		if m.Flux[i] > fluxRiverThreshold {
			// This cell is not suitable.
			continue
		}
		var fit float64
		// Steeper is better.
		fit += steepness[i]

		// We use a bell curve to calculate the fitness score related
		// to the elevation of the cell.
		// The peak of the curve is at 0.5, so that the fitness score is
		// 1.0 at the average elevation of the map.

		var curveVal float64
		x := m.Elevation[i]
		mu := 0.5
		sigma := 0.1
		curveVal = bellCurve(x, mu, sigma)
		log.Printf("Elevation: %v, curveVal: %v", x, curveVal)

		fit += curveVal

		// TODO:
		// - Maximize the distance to other dungeons.
		for _, b := range m.Dungeons {
			fit *= distanceToBuilding(i%m.Width, i/m.Width, b)
		}
		// - Avoid corners of the map.
		// fit += (1 - distToBorder[i]) // * distToCenter[i]
		// - Avoid the center of the map.

		fitness[i] = fit
	}

	normalize(fitness)

	return fitness
}

func (m *Map) calcFitnessScoreHouse(closerIsBetter bool) []float64 {
	// If closerIsBetter is true, the fitness score will be higher for cells that are closer to other houses.
	// This way, if a person isn't social, they can choose to live further away from other people.

	// Calculate the fitness score for each point.
	// We want the lowest flux in the cell, proximity to the market, and proximity to other houses,
	// but not too close. If there is already a house in the cell, the cell is not suitable.
	fitness := make([]float64, len(m.Flux))
	for i := range fitness {
		fitness[i] = -1
	}
	for i := range fitness {
		// Can't build on water.
		if m.fitnessScoreFlux(i) == -1 {
			// This cell is not suitable.
			continue
		}
		var fit float64

		// Lower flux is better.
		// fitness[i] = -m.Flux[i]

		// Similar elevation to the market is better.
		fit += 1.0 / (1.0 + math.Abs(m.Elevation[i]-m.Elevation[m.Root.X+m.Root.Y*m.Width]))

		// Add distance to the market. (closer is better)
		x := i % m.Width
		y := i / m.Width
		if distToBuild := distanceToBuilding(x, y, m.Root); distToBuild < 2.5 {
			continue // Too close to the market.
		} else {
			fit += 1.0 / (1.0 + distToBuild)
		}
		// Add distance to other houses. (further away is better)
		if distToBuild := m.fitnessScoreBuildingProximity(i, m.Buildings); distToBuild >= 0 {
			// No direct neighbors.
			if distToBuild < 2.5 {
				continue // Too close to another house.
			} else {
				if closerIsBetter {
					fit += 1.0 / (1.0 + distToBuild)
				} else {
					fit += 1.0 - 1.0/(1.0+distToBuild)
				}
			}
		}

		// TODO: We should prefer staying at the same river bank.

		// Check for construction sites.
		// TODO: Merge with the loop above.
		if distToBuild := m.fitnessScoreBuildingProximity(i, m.Construction); distToBuild >= 0 {
			// No direct neighbors.
			if distToBuild < 2.5 {
				continue // Too close to another house under construction.
			} else {
				if closerIsBetter {
					fit += 1.0 / (1.0 + distToBuild)
				} else {
					fit += 1.0 - 1.0/(1.0+distToBuild)
				}
			}
		}
		// Add distance to the border of the map.
		// NOTE: Yuk, this is a hack. We should use a proper distance function.
		x -= m.Width / 2
		y -= m.Height / 2

		// Now invert the distance, so that the center of the map has the highest score.
		distVal := math.Pow((1.0 - 1.0/(1.0+math.Sqrt(float64(x*x+y*y)))), 2)
		fit += 1 - distVal

		fitness[i] = fit
	}
	return fitness
}

func (m *Map) getHighestHouseFitness(closerIsBetter bool) (int, float64) {
	fitness := m.calcFitnessScoreHouse(closerIsBetter)
	best := 0
	for i := range fitness {
		if fitness[i] > fitness[best] {
			best = i
		}
	}
	return best, fitness[best]
}
