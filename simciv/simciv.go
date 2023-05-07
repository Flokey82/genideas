// Package simciv is a playground for simulating the spread of civilization.
package simciv

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/ojrac/opensimplex-go"
)

type Map struct {
	score       []float64 // Suitability
	dimX        int
	dimY        int
	Settlements []*Settlement
	seed        int64
	opensimplex opensimplex.Noise
	day         int
	year        int
}

func NewMap(dimX, dimY int, seed int64) *Map {
	m := &Map{
		score:       make([]float64, dimX*dimY),
		dimX:        dimX,
		dimY:        dimY,
		seed:        seed,
		opensimplex: opensimplex.NewNormalized(seed),
	}
	m.initMap()
	m.placeNSettlements(10)
	return m
}

// initMap initializes the map with a random score based on noise.
func (m *Map) initMap() {
	for y := 0; y < m.dimY; y++ {
		for x := 0; x < m.dimX; x++ {
			m.score[y*m.dimX+x] = m.opensimplex.Eval2(float64(x), float64(y))
		}
	}
}

func (m *Map) generateCivilization() {
	m.initMap()
}

func (m *Map) Tick() {
	m.day++
	if m.day > 365 {
		m.day = 0
		m.year++
	}
	// Tick all settlements (produce, consume, grow).
	for _, s := range m.Settlements {
		s.tick()
	}

	// Now we have to check if there are any settlements that can trade.
	// If there are, we have to check if there are any settlements that can
	// trade with them.
	for _, s := range m.Settlements {
		// Check if we can trade with any other settlement.
		for _, s2 := range m.Settlements {
			if s == s2 {
				continue
			}
			// Check if we can trade with this settlement.
			if s.canTrade(s2) {
				// Check if this settlement can trade with us.
				if s2.canTrade(s) {
					// We can trade with this settlement.
					// Trade with this settlement.
					s.trade(s2)
				}
			}
		}
	}

	// If there are any settlements that have a deficit, there will be a famine.
	for _, s := range m.Settlements {
		if s.deficit > 0 {
			// There is a famine.

			// TODO: If there is a famine, people might also move to other
			// settlements or found new settlements.
			// We'd have to check how big the deficit is and how many people
			// would not be able to eat. If we exceed a certain percentage of
			// what is required to found a new settlement, we should found a
			// new settlement, given that there is a suitable place nearby.
			refugees := int(math.Min(math.Ceil(s.deficit*2), 10))
			s.pop -= refugees

			// TODO: Find a prosperous settlement nearby that can take refugees.
			// Found a new settlement.
			if rand.Float32() < 0.7 {
				var found bool
				for _, i := range rand.Perm(len(m.Settlements)) {
					s2 := m.Settlements[i]
					if s2.getResourceBalance() > float64(refugees)*consumptionPerPop {
						log.Println("Famine! Refugees from %s to %s", s, s2)
						s2.pop += refugees
						found = true
						break
					}
				}
				if !found {
					// If relocation fails some people might unfortunately die.
					// The population will decrease by 1%.
					s.pop -= int(float64(s.pop) * 0.01)
					log.Println("Famine in", s, "population decreased by 1%")
				}
			} else {
				m.placeSettlement(refugees)
				log.Println("Famine! A new settlement was founded by refugees from", s)
			}
		}
	}

	// If a settlement becomes large (and wealthy) enough, someone might call themselves
	// a lord (as long as there is no other lord in the area).
	// Other settlements within the area of influence might request protection from the
	// lord, which will increase the lord's power. The lord will then be able to tax
	// the settlements under his protection.
	// If the lord is a despot, he might also demand tribute from the settlements that
	// are not under his protection.
}

func (m *Map) placeNSettlements(n int) {
	for i := 0; i < n; i++ {
		m.placeSettlement(10 + rand.Intn(10))
	}
}

func (m *Map) placeSettlement(pop int) {
	// Find the best place to place a settlement.
	bestScore := 0.0
	bestX := 0
	bestY := 0
	for y := 0; y < m.dimY; y++ {
		for x := 0; x < m.dimX; x++ {
			// Check the distance to the nearest settlement.
			var dist float64
			for _, s := range m.Settlements {
				if d := s.dist(x, y); d < dist || dist == 0 {
					dist = d
				}
			}
			// Calculate the score.
			score := m.score[y*m.dimX+x] * (1 + dist)
			if score > bestScore {
				bestScore = score
				bestX = x
				bestY = y
			}
		}
	}
	// Place the settlement.
	m.Settlements = append(m.Settlements, &Settlement{
		name:        fmt.Sprintf("Settlement %d", len(m.Settlements)),
		x:           bestX,
		y:           bestY,
		score:       m.score[bestY*m.dimX+bestX],
		pop:         pop,
		foundedDay:  m.day,
		foundedYear: m.year,
	})
}

type Settlement struct {
	name           string
	x, y           int     // Position
	pop            int     // Population
	score          float64 // Suitability
	foundedDay     int
	foundedYear    int
	resourceStores float64
	deficit        float64
}

func (s *Settlement) String() string {
	return fmt.Sprintf("%s (%d)", s.name, s.pop)
}

func (s *Settlement) dist(x, y int) float64 {
	return math.Sqrt(math.Pow(float64(x-s.x), 2) + math.Pow(float64(y-s.y), 2))
}

func (s *Settlement) tick() {
	// There is a 0.19% growth rate per year.
	// So the probability of growth is 0.19% * 1 day / 365 days.
	// TODO: The growth rate should also be based on the score and the
	// current population. Fix this!
	factor := float64(s.pop) * 0.19 / 365
	if factor >= 1 {
		s.pop += int(math.Ceil(factor * rand.Float64()))
	} else if rand.Float64() < factor {
		s.pop++
	}

	// If we have a resource balance, we can store it.
	balance := s.getResourceBalance()
	s.deficit = 0
	if balance > 0 {
		s.resourceStores += balance
	} else if s.resourceStores > 0 {
		// If we have a resource store, we can consume it.
		s.resourceStores += balance
		if s.resourceStores < 0 {
			s.deficit = -s.resourceStores
			s.resourceStores = 0
		}
	} else {
		// If we have no resource store, we have a deficit.
		s.deficit = -balance
	}
}

func (s *Settlement) getTradeRadius() float64 {
	return math.Sqrt(float64(s.pop)) * 10
}

func (s *Settlement) getProtectionRadius() float64 {
	return math.Sqrt(float64(s.pop)) * 5
}

func (s *Settlement) getResourceProduction() float64 {
	// The resource production is based on the population and the score.
	// TODO: After a certain population, the resource production should
	// decrease since there is not enough space to grow.
	// We converted 125km^2 to acres, check how much is used up by housing,
	// and then calculate the resource production potential based on the
	// population and remaining area. Or something like that.
	return math.Min(
		(30888.2-float64(s.pop)/2), // Total theoretical resource production.
		float64(s.pop)*4.0,         // Resource production based on existing population.
	) * productionPerPop * s.score // One acre can produce 1 unit of resources, enough to feed 2 people
}

func (s *Settlement) getResourceConsumption() float64 {
	// The resource consumption is based on the population.
	// The consumption is about 1/2 of the production potential.
	return float64(s.pop) * consumptionPerPop
}

func (s *Settlement) getResourceBalance() float64 {
	return s.getResourceProduction() - s.getResourceConsumption()
}

func (s *Settlement) canTrade(s2 *Settlement) bool {
	// Check if we are within the trade radius of the other settlement.
	return s.dist(s2.x, s2.y) < s.getTradeRadius()
}

func (s *Settlement) trade(s2 *Settlement) {
	// If no one has a deficit, we don't have to trade.
	// If both have a deficit, we can't trade.
	if (s.deficit == 0 && s2.deficit == 0) || (s.deficit > 0 && s2.deficit > 0) {
		return
	}
	var giver, receiver *Settlement
	if s.deficit == 0 {
		giver = s
		receiver = s2
	} else if s2.deficit == 0 {
		giver = s2
		receiver = s
	}

	if giver.resourceStores == 0 {
		return
	}

	// Calculate the amount to trade.
	amount := math.Min(giver.resourceStores, receiver.deficit)
	// Trade the amount.
	giver.resourceStores -= amount
	receiver.deficit -= amount
	log.Println(giver, "traded", amount, "resources with", receiver)
}

const (
	productionPerPop  = 1.0
	consumptionPerPop = 0.5
)
