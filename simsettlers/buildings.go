package simsettlers

import (
	"log"
	"math"
)

func (m *Map) advanceConstruction() {
	var stillBuilding []*Building
	// Build new buildings.
	for _, b := range m.Construction {
		b.Remaining--
		if b.Remaining == 0 {
			log.Printf("Finished building %v", b)
			m.Buildings = append(m.Buildings, b)
		} else {
			log.Printf("Still building %v", b)
			stillBuilding = append(stillBuilding, b)
		}
	}
	m.Construction = stillBuilding
}

func (m *Map) constructMoreHouses() {
	// Check if we have enough housing capacity.
	// TODO: Instead, let individual settlers decide if they want to build a house or live with their parents.
	// The parents might also decide to expand their house, increasing the housing capacity.
	gotCap := m.getHousingCapacity()
	if gotCap < m.Population {
		// We need to build more houses (if we can afford it).
		theoCap := m.getTheoreticalHousingCapacity()
		if theoCap < m.Population {
			if m.Resources >= buildingCosts[BuildingTypeHouse] {
				// TODO: Find a suitable location for the house.
				best, score := m.getHighestHouseFitness()
				if !math.IsInf(score, -1) {
					log.Printf("Building a house")
					// Which is close to the market, but not too close.
					m.AddBuilding(best%m.Width, best/m.Width, BuildingTypeHouse)
					m.Resources -= buildingCosts[BuildingTypeHouse]
				} else {
					log.Printf("No suitable location for a house found")
				}
			} else {
				log.Printf("Not enough resources to build a house")
			}
		}
	}
}

// Building represents a building on the map.
type Building struct {
	X, Y      int
	Remaining int
	Type      string
}

// NewBuilding creates a new building of the given type at the given position.
func NewBuilding(x, y int, t string) *Building {
	return &Building{
		X:         x,
		Y:         y,
		Remaining: buildingCosts[t], // The number of ticks remaining until the building is finished.
		Type:      t,
	}
}

// AddBuilding adds a new building to the map.
func (m *Map) AddBuilding(x, y int, t string) *Building {
	b := NewBuilding(x, y, t)
	m.Construction = append(m.Construction, b)
	return b
}

// Yield returns the resource yield of the given building in a tick.
func (b *Building) Yield() int {
	return buildingYield[b.Type]
}

const (
	BuildingTypeMarket = "market"
	BuildingTypeHouse  = "house"
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

func (m *Map) calcFitnessScoreHouse() []float64 {
	// Calculate the fitness score for each point.
	// We want the lowest flux in the cell, proximity to the market, and proximity to other houses,
	// but not too close. If there is already a house in the cell, the cell is not suitable.
	fitness := make([]float64, len(m.Flux))
Loop:
	for i := range fitness {
		// Can't build on water.
		if m.Flux[i] > fluxRiverThreshold {
			// This cell is not suitable.
			fitness[i] = math.Inf(-1)
			continue
		}

		// Lower flux is better.
		// fitness[i] = -m.Flux[i]

		// Similar elevation to the market is better.
		fitness[i] -= 1 - 1.0/(1.0+math.Abs(m.Elevation[i]-m.Elevation[m.Root.X+m.Root.Y*m.Width]))

		// Add distance to the market. (closer is better)
		x := i % m.Width
		y := i / m.Width
		fitness[i] += 1 - 1.0/(1.0+distanceToBuilding(x, y, m.Root))

		// TODO: No direct neighbors.

		// Add distance to other houses. (further away is better)
		for _, b := range m.Buildings {
			// Check if the building is in the current cell.
			if b.X == x && b.Y == y {
				// This cell is not suitable.
				fitness[i] = math.Inf(-1)
				continue Loop
			}
			if b.Type == BuildingTypeHouse {
				distToBuild := distanceToBuilding(x, y, b)

				// Can't build too close to another house.
				if distToBuild < 2.5 {
					// Too close to another house.
					fitness[i] = math.Inf(-1)
					continue Loop
				}
				// fitness[i] += 1.0 / (1.0 + distToBuild)
			}

			// TODO: We should prefer staying at the same river bank.
		}

		// Check for construction sites.
		// TODO: Merge with the loop above.
		for _, b := range m.Construction {
			// Check if the building is in the current cell.
			if b.X == x && b.Y == y {
				// This cell is not suitable.
				fitness[i] = math.Inf(-1)
				continue Loop
			}
			if b.Type == BuildingTypeHouse {
				distToBuild := distanceToBuilding(x, y, b)

				// Can't build too close to another house.
				if distToBuild < 2.5 {
					// Too close to another house.
					fitness[i] = math.Inf(-1)
					continue Loop
				}
				// fitness[i] += 1.0 / (1.0 + distToBuild)
			}
		}

		// Add distance to the border of the map.
		// NOTE: Yuk, this is a hack. We should use a proper distance function.
		x -= m.Width / 2
		y -= m.Height / 2

		// Now invert the distance, so that the center of the map has the highest score.
		distVal := math.Pow((1.0 - 1.0/(1.0+math.Sqrt(float64(x*x+y*y)))), 2)
		fitness[i] -= distVal
	}
	return fitness
}

func (m *Map) getHighestHouseFitness() (int, float64) {
	fitness := m.calcFitnessScoreHouse()
	best := 0
	for i := range fitness {
		if fitness[i] > fitness[best] {
			best = i
		}
	}
	return best, fitness[best]
}
