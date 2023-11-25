// Package simciv is a playground for simulating the spread of civilization.
package simciv

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/ojrac/opensimplex-go"
)

type Map struct {
	score       []float64         // Suitability
	dimX, dimY  int               // Dimensions in tiles
	Settlements []*Settlement     // Settlements
	seed        int64             // Seed for the noise
	opensimplex opensimplex.Noise // Noise generator
	day         int               // Day of the year
	year        int               // Year
}

// NewMap returns a new map.
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

// initMap initializes the map tiles with a random score based on noise.
func (m *Map) initMap() {
	for y := 0; y < m.dimY; y++ {
		for x := 0; x < m.dimX; x++ {
			m.score[y*m.dimX+x] = m.opensimplex.Eval2(float64(x), float64(y))
		}
	}
}

// Tick advances the simulation by one day.
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
	var activeSettlements []*Settlement
	for _, s := range m.Settlements {
		if s.pop == 0 {
			continue
		}
		activeSettlements = append(activeSettlements, s)

		// Check if we can trade with any other settlement.
		for _, s2 := range m.Settlements {
			if s == s2 || s.pop == 0 {
				continue
			}
			// Check if we can trade with this settlement
			// and if this settlement can trade with us.
			if s.canTrade(s2) && s2.canTrade(s) {
				// We can trade with this settlement.
				// Trade with this settlement.
				s.trade(s2)
			}
		}
	}

	// If there are any settlements that have a deficit, there will be a famine.
	for _, s := range activeSettlements {
		// TODO: Also move people when there is a negative balance.
		var refugees int

		// If there is a deficit, we have starving people.
		if s.deficit > 0 {
			// There is a famine. Calculate how many people are endangered.
			refugees = int(math.Min(math.Ceil(s.deficit*(2+rand.Float64())), float64(s.pop)))
			log.Println("Famine in", s, "refugees:", refugees, "deficit:", s.deficit)
		}

		// TODO: Check overall prosperity of the settlement.
		// If prosperity is low, people might leave to seek a better life elsewhere.

		// If there are no refugees, we don't have to do anything.
		if refugees == 0 {
			continue
		}

		// Find a prosperous settlement nearby that can take refugees.
		candidates := make([]*Settlement, 0, len(m.Settlements)-1)
		for _, s2 := range m.Settlements {
			if s2 == s {
				continue
			}
			candidates = append(candidates, s2)
		}

		// Sort by distance to s.
		sort.Slice(candidates, func(i, j int) bool {
			return s.dist(candidates[i].x, candidates[i].y) < s.dist(candidates[j].x, candidates[j].y)
		})

		// Check if there is a suitable settlement nearby that can take refugees
		// or an abandoned settlement that can be repopulated.
		// TODO: Distribute refugees to multiple settlements.
		var found bool
		for i, s2 := range candidates {
			// TODO: Add max (actual) distance that we can migrate.
			if i > 6 {
				break
			}

			// The longer the distance, the more likely it is that someone will die.
			if s2.getResourceBalance() > float64(refugees)*consumptionPerPop || (s2.pop == 0 && s2.score >= s.score) {
				s2.pop += refugees
				log.Printf("Famine! %d Refugees from %s to %s\n", refugees, s, s2)
				found = true
				break
			}
		}

		if !found {
			// Found a new settlement.
			if rand.Float32() < 0.7 {
				// TODO: Find a suitable place nearby.
				m.placeSettlement(refugees)
				log.Println("Famine! A new settlement was founded by refugees from", s)
			} else {
				// If relocation fails some people might unfortunately die.
				// The population will by up to the number of refugees.
				refugees = rand.Intn(refugees)
				log.Println("Famine in", s, "population decreased by", refugees)
			}
		}
		s.pop -= refugees
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
	var bestScore float64
	var bestX, bestY int
	for y := 0; y < m.dimY; y++ {
		for x := 0; x < m.dimX; x++ {
			// Check if there is already a settlement here.
			// Check the distance to the nearest settlement.
			var occupied bool
			var dist float64
			for _, s := range m.Settlements {
				if s.x == x && s.y == y {
					occupied = true
					break
				}
				if d := s.dist(x, y); d < dist || dist == 0 {
					dist = d
				}
			}
			if occupied {
				continue
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
	return fmt.Sprintf("%s (%d, deficit %.2f, stores %.2f)", s.name, s.pop, s.deficit, s.resourceStores)
}

func (s *Settlement) dist(x, y int) float64 {
	return math.Sqrt(math.Pow(float64(x-s.x), 2) + math.Pow(float64(y-s.y), 2))
}

func (s *Settlement) tick() {
	// If there is no population, we don't have to do anything.
	if s.pop == 0 {
		s.deficit = 0
		return
	}

	useLogistics := true

	// There is a 0.19% growth rate per year.
	// So the rate of growth is 0.19% * 1 day / 365 days.
	//
	// TODO:
	// - The growth rate should also be based on the score. Fix this!
	// - The growth rate should also be based on the prosperity of the settlement.
	// - Also consider using the logistics function to limit the growth
	//  rate to the maximum sustainable population.
	var newPop float64
	growthRate := s.score * 0.19 / 365
	if useLogistics {
		maxPop := float64(s.getMaxPopulation())
		curPop := float64(s.pop)
		newPop = (curPop * maxPop) / (curPop + (maxPop-curPop)*math.Pow(math.E, -growthRate))
	} else {
		newPop = float64(s.pop) * math.Pow(math.E, growthRate)
	}
	if newPeople := newPop - float64(s.pop); newPeople >= 1 {
		log.Println(s, "grew by", int(math.Ceil(newPeople)), "people")
		s.pop += int(math.Ceil(newPeople))
	} else if rand.Float64() < newPeople {
		log.Println(s, "grew by 1 person")
		s.pop++
	}

	// There is a chance a fire will destroy part of the storage.
	if s.resourceStores > 0 && rand.Float64() < 0.01 {
		dest := rand.Float64() * s.resourceStores
		log.Println(s, "lost", dest, "resources in a fire")
		s.resourceStores -= dest
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

// getMaxPopulation returns the maximum population that can be sustained by the settlement.
func (s *Settlement) getMaxPopulation() int {
	return int(s.getResourceProduction() / consumptionPerPop)
}

func (s *Settlement) getTradeRadius() float64 {
	return math.Sqrt(float64(s.pop)) * 10
}

func (s *Settlement) getProtectionRadius() float64 {
	return math.Sqrt(float64(s.pop)) * 5
}

func (s *Settlement) getResourceProduction() float64 {
	// The resource production is based on the population and the score.
	//
	// Assumptions:
	// - 1 acre can produce 1 unit of resources, enough to feed 2 people.
	// - The settlement size is 125km^2 max, which is 30888.2 acres.
	// - 1 person can tend to 4 acres.
	// - Living space is 1 acre per 2 people.
	//
	// TODO: Living space can be reduced by higher density housing?
	//
	// After reaching a certain population, the resource production should
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
	log.Println(giver, "traded", amount, "resources with", receiver)

	// Trade the amount.
	giver.resourceStores -= amount
	receiver.deficit -= amount
}

const (
	productionPerPop  = 1.0
	consumptionPerPop = 0.5
)
