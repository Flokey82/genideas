package simsettlers

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/Flokey82/go_gens/gameconstants"
)

// Person represents a person in the village.
// TODO:
// - Add a personality.
// - Add a job.
// - Allow multiple construction projects.
type Person struct {
	FirstName    string
	LastName     string
	Birthday     int
	Age          int
	Gender       int
	Pregnant     int // The number of days the person will still be pregnant.
	Resources    int // The amount of resources the person has.
	Dead         bool
	Home         *Building   // The home of the person.
	Constructing []*Building // The building the person is currently constructing. TODO: Allow multiple construction projects.
	Owns         []*Building // The buildings the person owns.
	Mother       *Person
	Father       *Person
	Spouse       *Person // The spouse of this person.
	Children     []*Person
}

// OwnsOwnHome returns true if the person owns their own home.
func (p *Person) OwnsOwnHome() bool {
	if p.Home == nil {
		return false
	}
	for _, o := range p.Home.Owners {
		if o == p {
			return true
		}
	}
	return false
}

// LivesWithParents returns true if the person lives with their parents.
func (p *Person) LivesWithParents() bool {
	if p.Home == nil {
		return false
	}
	if p.Mother != nil && p.Mother.Home == p.Home {
		return true
	}
	if p.Father != nil && p.Father.Home == p.Home {
		return true
	}
	return false
}

func (p *Person) heir() *Person {
	// TODO: If there are no children, then siblings should inherit.
	// TODO: If there are no siblings, then parents should inherit.
	// TODO: If there are no children, then siblings should inherit.
	if p.Spouse != nil && !p.Spouse.Dead {
		return p.Spouse
	}
	for _, c := range p.Children {
		if !c.Dead {
			return c
		}
	}

	if p.Father != nil {
		return p.Father
	}
	if p.Mother != nil {
		return p.Mother
	}
	return nil
}

func (p *Person) SetHome(b *Building) {
	// Remove person from the occupants list of the previous home.
	if p.Home != nil {
		p.Home.RemoveOccupant(p)
	}
	p.Home = b
	b.AddOccupant(p)
}

// String returns the string representation of the person.
func (p *Person) String() string {
	str := fmt.Sprintf("%s %s (%d %s - %d)", p.FirstName, p.LastName, p.Age, p.genderString(), p.Resources)
	if p.Dead {
		str += " (dead)"
	}
	return str
}

func (p *Person) genderString() string {
	if p.Gender == GenderMale {
		return "M"
	}
	return "F"
}

func (p *Person) isMarried() bool {
	return p.Spouse != nil
}

const numDaysPregnant = 270

const (
	GenderFemale = iota
	GenderMale
)

func (m *Map) addNRandomPeople(n int) {
	for i := 0; i < n; i++ {
		p := &Person{
			LastName:  m.lastGen.String(),
			Birthday:  m.Day,
			Age:       rand.Intn(20) + 18,
			Gender:    rand.Intn(2),
			Resources: rand.Intn(10),
		}
		if p.Gender == GenderMale {
			p.FirstName = m.firstGen[1].String()
		} else {
			p.FirstName = m.firstGen[0].String()
		}
		m.RealPop = append(m.RealPop, p)
	}
}

func (m *Map) agePop() {
	// Let's see who has a birthday today and check if anyone dies.
	var remPop []*Person
	for _, p := range m.RealPop {
		if p.Dead {
			continue
		}
		if p.Birthday == m.Day {
			p.Age++
		}
		// Check if anyone dies.
		if gameconstants.DiesAtAgeWithinNDays(p.Age, 1) {
			m.Population--
			// TODO: Remove from spouse?
			log.Printf("Died: %v", p)
			p.Dead = true

			// TODO:
			// - Identify who will inherit all buildings, resources, etc.
			// - Find all buildings that we own and transfer ownership to the spouse or children.
			// - What if we are constructing a building? Transfer ownership to the spouse or children.
			// - Allow multiple heirs.

			heir := p.heir()
			if heir != nil {
				// Move resources to the heir.
				heir.Resources += p.Resources

				// Move buildings to the heir.
				for _, b := range p.Owns {
					// Remove from the occupants list of the building.
					b.RemoveOccupant(p)

					// Remove from the owners list of the building.
					b.RemoveOwner(p)

					// Add to the occupants list of the building.
					b.AddOwner(heir)
				}

				// Move constructing buildings to the heir.
				for _, b := range p.Constructing {
					// Remove the deceased from the constructing list of the building.
					b.RemoveOwner(p)

					// Add heir to the owners list of the building.
					b.AddOwner(heir)

					// Add the building to the constructing list of the heir.
					heir.Constructing = append(heir.Constructing, b)
				}

				// Move a homeless heir to the home of the deceased.
				if p.Home != nil && heir.Home == nil {
					heir.SetHome(p.Home)
				}
			}

			if p.Home != nil {
				// Remove from the occupants list of the home.
				p.Home.RemoveOccupant(p)

				// Remove from the owners list of the home.
				p.Home.RemoveOwner(p)
			}

			// Move to the cemetery.
			p.SetHome(m.Cemetery)
		} else {
			remPop = append(remPop, p)
		}
	}
	// TODO: Find a more efficient way to do this besides re-allocation.
	m.RealPop = remPop
}

func (m *Map) matchSingles() {
	var men, women []*Person
	for _, p := range m.RealPop {
		if p.Dead || p.isMarried() || p.Age < 18 {
			continue
		}
		if p.Gender == GenderMale {
			men = append(men, p)
		} else if p.Gender == GenderFemale {
			women = append(women, p)
		}
	}

	// Sort by age.
	sort.Slice(men, func(i, j int) bool {
		return men[i].Age < men[j].Age
	})
	sort.Slice(women, func(i, j int) bool {
		return women[i].Age < women[j].Age
	})

	// Match the singles with a max age difference of 25%.
	// TODO:
	// - Personalities matter.
	// - Matching should not just happen, there should be a probability and some randomness.
	// - There should be a chance that a person stays single (by choice, or because they just can't find the right person).
	const ageDiffFactor = 0.25
	const menTakeFemaleName = true
	for _, wp := range women {
		// Ladies' choice.
		for _, mp := range men {
			if mp.isMarried() {
				continue
			}
			if ageDiff := math.Abs(float64(wp.Age - mp.Age)); ageDiff < ageDiffFactor*float64(wp.Age) {
				wp.Spouse = mp
				mp.Spouse = wp
				if menTakeFemaleName {
					mp.LastName = wp.LastName
				} else {
					wp.LastName = mp.LastName
				}
				log.Printf("Matched %v and %v", wp, mp)
				break
			}
		}
	}
}

func (m *Map) advancePregnancies() {
	// Advance all pregnancies.
	for _, p := range m.RealPop {
		if p.Dead {
			continue
		}
		if p.Pregnant > 0 {
			p.Pregnant--
			if p.Pregnant == 0 {
				child := &Person{
					Home:     p.Home,
					Mother:   p,
					Father:   p.Spouse,
					LastName: p.LastName, // Last name of the mother.
					Birthday: m.Day,
					Age:      0,
					Gender:   rand.Intn(2),
				}
				if child.Gender == GenderMale {
					child.FirstName = m.firstGen[1].String()
				} else {
					child.FirstName = m.firstGen[0].String()
				}
				log.Printf("Born: %v", child)
				p.Children = append(p.Children, child)
				if p.Home != nil {
					child.SetHome(p.Home)
				}
				if p.Spouse != nil {
					p.Spouse.Children = append(p.Spouse.Children, child)
				}
				m.RealPop = append(m.RealPop, child)
			}
		}
	}

	// All couples have a chance of getting pregnant, which decreases
	// by the number of children they already have.
	for _, p := range m.RealPop {
		if p.isMarried() && p.Pregnant == 0 && p.Gender == GenderFemale {
			if rand.Intn(100) < 3-len(p.Children) {
				p.Pregnant = numDaysPregnant
			}
		}
	}
}

func (m *Map) tickPeople() {
	// If true, owners that live in a house will repair their homes.
	repairHome := false

Loop:
	for _, p := range m.RealPop {
		if p.Dead {
			continue
		}
		// TODO: Once we are old enough, we can save up for a house.
		// Maybe start with 16 years old?
		if p.Age > 16 && rand.Intn(100) < 10 {
			// TODO: Maybe use fractional resources?
			p.Resources += 1 // 1 is quite a lot for a single day?
		}
		if p.Age < 18 {
			continue
		}

		// Check if we are already building a house or if our spouse is.
		if p.Constructing != nil || p.Spouse != nil && p.Spouse.Constructing != nil {
			continue
		}

		// Check if we live with our parents and if we have a spouse.
		// TODO: Check if we still live with our parents and move out if we have a spouse that is pregnant.
		if p.LivesWithParents() && p.Spouse == nil {
			// Wait until we have a spouse.
			continue
		}

		// Check if we own our own home.
		if p.OwnsOwnHome() {
			// Check if we need to repair our home.
			// TODO:
			// - Make the decision to repair or not based on personality
			// and the condition of the building.
			// - Any occupant should be able to repair the home if they are
			// old enough, have enough resources, and are crafty enough.
			if repairHome && p.Home.Condition < 100 && p.Resources > 0 {
				p.Home.Condition++
				p.Resources--
			}
			continue
		}

		// Check if our spouse owns their own home, and if so, move in.
		if p.Spouse != nil && p.Spouse.OwnsOwnHome() {
			p.SetHome(p.Spouse.Home)
			p.Spouse.Home.AddOwner(p)
			continue
		}

		// We want to move to a new place if:
		// - we don't have a home.
		// - we live with our parents, have a spouse, and/or can afford to build a house.
		budget := p.Resources
		if p.Spouse != nil {
			budget += p.Spouse.Resources
		}

		// TODO: If there is a house in bad condition available,
		// we should be able to buy it for a lower price!

		// Check if there is an existing house that we can move into.
		// Either, because it is abandoned, or because there are no occupants.
		// (Inherited from parents, or because the owners died.)
		for _, b := range m.Buildings {
			if b.Type != BuildingTypeHouse || b.IsOccupied() {
				continue
			}
			if b.IsOwned() {
				// Buy it from the owner.
				// Calculate the purchase price based on the condition of the building.
				purchasePrice := b.PurchasePrice()
				if budget < purchasePrice {
					//log.Printf("Not enough resources to buy a house")
					continue
				}

				// Split the cost between the owners.
				cost := purchasePrice / len(b.Owners)
				for _, o := range b.Owners {
					o.Resources += cost
				}
				b.Occupants = nil
				b.Owners = nil

				// Deduct the cost from the resources of the buyer.
				p.Resources -= cost

				// Move in.
				p.SetHome(b)
				b.AddOwner(p)

				// If there is a spouse, add them as an owner and move them in.
				// TODO: If we are pooling resources with our spouse, we should split the cost.
				if p.Spouse != nil {
					b.AddOwner(p.Spouse)
					p.Spouse.SetHome(b)
				}
			} else {
				// No owners, so move in.
				// TODO: Repair costs etc.
				p.SetHome(b)
				b.AddOwner(p)
			}
			continue Loop
		}

		// Can we afford to build a new house?
		if budget < buildingCosts[BuildingTypeHouse] {
			//log.Printf("Not enough resources to build a house")
			continue
		}

		// TODO: Find a suitable location for the house.
		best, score := m.getHighestHouseFitness()
		if !math.IsInf(score, -1) {
			log.Printf("Building a house for %v", p)

			// Construct the house.
			home := m.AddBuilding(best%m.Width, best/m.Width, BuildingTypeHouse)
			home.AddOwner(p)
			p.Constructing = append(p.Constructing, home)

			// If there is a spouse, add them as an owner and set their constructing building.
			if p.Spouse != nil {
				home.AddOwner(p.Spouse)
				if p.Spouse.Constructing != nil {
					panic("Spouse is already constructing a building")
				}
				p.Spouse.Constructing = append(p.Spouse.Constructing, home)
			}
			// TODO: Deduct from resources from both partners.
			p.Resources -= buildingCosts[BuildingTypeHouse]
		} else {
			log.Printf("No suitable location for a house found")
		}
	}
}
