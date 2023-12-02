package simciv

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/ojrac/opensimplex-go"
)

const (
	cellSquareKm = 125.0                  // 125 square km per cell
	cellAcres    = cellSquareKm * 247.105 // 125 square km in acres
)

type Map struct {
	score                   []float64         // Suitability
	levelWidth, levelHeight int               // Dimensions in tiles
	Settlements             []*Settlement     // Settlements
	seed                    int64             // Seed for the noise
	opensimplex             opensimplex.Noise // Noise generator
	day                     int               // Day of the year
	year                    int               // Year
}

// NewMap returns a new map.
func NewMap(dimX, dimY int, seed int64) *Map {
	m := &Map{
		score:       make([]float64, dimX*dimY),
		levelWidth:  dimX,
		levelHeight: dimY,
		seed:        seed,
		opensimplex: opensimplex.NewNormalized(seed),
	}
	m.initMap()
	m.placeNSettlements(10)
	return m
}

// initMap initializes the map tiles with a random score based on noise.
func (m *Map) initMap() {
	for y := 0; y < m.levelHeight; y++ {
		for x := 0; x < m.levelWidth; x++ {
			m.score[y*m.levelWidth+x] = m.opensimplex.Eval2(4*float64(x)/float64(m.levelWidth), 4*float64(y)/float64(m.levelHeight))
		}
	}
}

// Tick advances the simulation by one day.
func (m *Map) Tick(days int) {
	m.day += days
	if m.day > 365 {
		m.year += m.day / 365
		m.day = m.day % 365
	}

	// Tick all settlements (produce, consume, grow).
	for _, s := range m.Settlements {
		s.tick(days)
	}

	// Now we have to check if there are any settlements that can trade.
	// If there are, we have to check if there are any settlements that can
	// trade with them.
	var activeSettlements []*Settlement
	for _, s := range m.Settlements {
		// Reset trading partners and protected settlements.
		s.tradingSettlements = s.tradingSettlements[:0]
		s.protectedSettlements = s.protectedSettlements[:0]

		// Reset protection if expired. (Maybe do this after updating the settlements?)
		if s.protectedBy != nil && !s.protectedBy.canProtect(s) {
			s.protectedBy = nil
		}
		if s.pop == 0 {
			continue
		}
		activeSettlements = append(activeSettlements, s)

		// Check if we can trade with any other settlement.
		for _, s2 := range m.Settlements {
			if s == s2 || s.pop == 0 {
				continue
			}
			// Check if we can trade with this settlement.
			if s.canTrade(s2) {
				// We can trade with this settlement.
				s.trade(s2)
			}

			// Check if we can protect this settlement.
			if s.canProtect(s2) {
				// We can protect this settlement.
				s.protect(s2)
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
		} else if s.pop > 500 && rand.Float64() < 0.01 {
			// TODO: Calculate prosperity of the settlement... if it is low, people might leave.
			// People just might want to move to a new place.
			refugees = rand.Intn(s.pop / 20)
			log.Println("Migration from", s, "refugees:", refugees)
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
		// TODO:
		// - Distribute refugees to multiple settlements.
		// - The longer the distance, the more likely it is that someone will die.
		var found bool

		// Check if we can reach any settlement within 16 tiles.
		for _, s2 := range candidates {
			if s.dist(s2.x, s2.y) > 16 {
				break
			}

			// Check if we can sustain the refugees or if the settlement is abandoned but more prosperous than the one that we are
			// fleeing from.
			if s2.getResourceBalance() > float64(refugees)*consumptionPerPop || (s2.pop == 0 && s2.score >= s.score) {
				s2.pop += refugees
				log.Printf("Famine! %d Refugees from %s to %s\n", refugees, s, s2)
				found = true
				break
			}
		}

		// If we haven't found a suitable settlement, we have to found a new one or
		// people might die.
		if !found {
			// Found a new settlement.
			if rand.Float32() < 0.7 {
				// TODO: Find a suitable place nearby.
				m.placeSettlement(refugees, s)
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
		m.placeSettlement(100+rand.Intn(100), nil)
	}
}

func (m *Map) placeSettlement(pop int, origin *Settlement) {
	// Find the best place to place a settlement.
	var bestScore float64
	var bestX, bestY int
	for y := 0; y < m.levelHeight; y++ {
		for x := 0; x < m.levelWidth; x++ {
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
			// Calculate distance to origin (if any).
			// Penalize large distances for refugees.
			distOrigin := 0.0
			if origin != nil {
				distOrigin = origin.dist(x, y)
			}

			// Calculate the score.
			score := m.score[y*m.levelWidth+x] * (1 + dist - distOrigin)
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
		score:       m.score[bestY*m.levelWidth+bestX],
		pop:         pop,
		foundedDay:  m.day,
		foundedYear: m.year,
	})
}
