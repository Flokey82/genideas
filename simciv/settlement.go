package simciv

import (
	"fmt"
	"log"
	"math"
	"math/rand"
)

const (
	productionPerPop      = 1.0 // 1 acre can produce 1 unit of resources, enough to feed 2 people.
	consumptionPerPop     = 0.5 // 1 person consumes 0.5 units of resources.
	productionAcresPerPop = 3.8 // One person can tend to 3.8 acres.
	productionAgedPop     = 0.5 // 50% of the population is able to work.
)

type Settlement struct {
	name                 string
	x, y                 int     // Position
	pop                  int     // Population
	score                float64 // Suitability
	foundedDay           int
	foundedYear          int
	resourceStores       float64
	deficit              float64
	protectedBy          *Settlement
	tradingSettlements   []*Settlement
	protectedSettlements []*Settlement
}

func (s *Settlement) String() string {
	return fmt.Sprintf("%s (%d, deficit %.2f, stores %.2f)", s.name, s.pop, s.deficit, s.resourceStores)
}

func (s *Settlement) dist(x, y int) float64 {
	return math.Sqrt(math.Pow(float64(x-s.x), 2) + math.Pow(float64(y-s.y), 2))
}

func (s *Settlement) tick(days int) {
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
	growthRate := float64(days) * (s.score/100 + 0.19/100) / 365 // Up to 1.19% growth rate per year.
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
	balance := s.getResourceBalance() * float64(days)
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
		(cellAcres-float64(s.pop)/2),                           // Total theoretical resource production.
		float64(s.pop)*productionAcresPerPop*productionAgedPop, // Resource production based on existing population of working age.
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

func (s *Settlement) getTradeRadius() float64 {
	return math.Sqrt(float64(s.pop))
}

func (s *Settlement) canTrade(s2 *Settlement) bool {
	// Check if we are within the trade radius of the other settlement.
	return s.dist(s2.x, s2.y) < s.getTradeRadius()
}

func (s *Settlement) trade(s2 *Settlement) {
	// Add the settlement to the list of trading settlements.
	var alreadyTrading bool
	for _, ts := range s.tradingSettlements {
		if ts == s2 {
			alreadyTrading = true
			break
		}
	}
	if !alreadyTrading {
		s.tradingSettlements = append(s.tradingSettlements, s2)
	}

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

func (s *Settlement) getProtectionRadius() float64 {
	return math.Sqrt(float64(s.pop) / 2)
}

func (s *Settlement) canProtect(s2 *Settlement) bool {
	// Check if we are within the protection radius of the other settlement.
	return s.dist(s2.x, s2.y) < s.getProtectionRadius()
}

func (s *Settlement) protect(s2 *Settlement) {
	// TODO: Check if they want to be protected or if they are already protected.
	if s2.protectedBy != nil && s2.protectedBy != s {
		// The other settlement is already protected by someone else.
		return
	}
	s2.protectedBy = s

	// If we are already protecting this settlement, we don't have to do anything.
	for _, ps := range s.protectedSettlements {
		if ps == s2 {
			return
		}
	}
	s.protectedSettlements = append(s.protectedSettlements, s2)
	log.Println(s, "protects", s2)
}
