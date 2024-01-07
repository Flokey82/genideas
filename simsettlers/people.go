package simsettlers

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/Flokey82/go_gens/gameconstants"
)

func (m *Map) addNRandomPeople(n int) {
	for i := 0; i < n; i++ {
		p := m.newPerson("", m.lastGen.String(), byte(rand.Intn(2)), uint16(rand.Intn(20)+18))
		p.Resources = rand.Intn(10)

		// TODO: Fix assignment of goals.
		p.assignChildhoodGoals()
		p.assignAdultGoals()
		// p.assignElderlyGoals()

		m.RealPop = append(m.RealPop, p)
	}
}

func (m *Map) newPerson(firstName, lastName string, gender byte, age uint16) *Person {
	// If no first name is given, pick a random one.
	if firstName == "" {
		if gender == GenderMale {
			firstName = m.firstGen[1].String()
		} else {
			firstName = m.firstGen[0].String()
		}
	}
	return &Person{
		Health:    healthMax,
		FirstName: firstName,
		LastName:  lastName,
		Birthday:  m.Day,
		Age:       age,
		Gender:    gender,
		Opinions:  make(Opinions),
	}
}

const healthMax = 100.0

// Person represents a person in the village.
// TODO:
// - Add a personality.
type Person struct {
	FirstName      string
	LastName       string
	Age            uint16  // age of the person in years (pretty generous with the uint16)
	Gender         byte    //
	Birthday       uint16  // day of the year the person was born
	Pregnant       uint16  // days the person will still be pregnant
	Health         float64 // current health of the person
	Dead           bool    // true if the person is dead
	LocationPerson         // location and speed of the person

	// Actions, jobs, tasks
	// TODO: Move currentTree into the motive, so a plan can be resumed if we switch motives
	// temporarily.
	Motives       []*Motive // List of motives that the person has.
	CurrentMotive *Motive   // The current motive that we are trying to satisfy.
	CurrentTree   *Tree     // The current task tree that we are executing.
	Goals         Goal      // personal goals
	Job           JobType   // current job

	// Real estate and wealth
	Home         *Building   // home of the person
	Constructing []*Building // buildings under construction
	Owns         []*Building // buildings owned
	Resources    int         // wealth in resources

	// Family
	Mother       *Person
	Father       *Person
	Spouse       *Person
	FormerSpouse []*Person
	Children     []*Person

	// Opinions
	Opinions Opinions // opinions of other people
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
		return false // No home associated with the person.
	}

	// Lives with either parent?
	return p.Mother != nil && p.Mother.Home == p.Home || p.Father != nil && p.Father.Home == p.Home
}

func (p *Person) heir() *Person {
	// If we have a living spouse, they inherit.
	if p.Spouse != nil && !p.Spouse.Dead {
		return p.Spouse
	}

	// If we have a living child, the first one inherits.
	for _, c := range p.Children {
		if !c.Dead {
			return c
		}
	}

	// If we have a living parent, they inherit.
	if p.Father != nil {
		return p.Father
	}
	if p.Mother != nil {
		return p.Mother
	}

	// TODO: Grandchildren, siblings, etc?
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

func (m *Map) handleDeath(p *Person) {
	m.Population--
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

	// Remove from spouse so they can remarry?
	// TODO:
	// - Adopt children?
	// - Add mourning period?
	if p.Spouse != nil {
		p.Spouse.FormerSpouse = append(p.Spouse.FormerSpouse, p)
		p.Spouse.Spouse = nil
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

			// Check if we have to decide life goals.
			if p.Age == 1 {
				p.assignChildhoodGoals()
			} else if p.Age == 18 {
				p.assignAdultGoals()
			} else if p.Age == 65 {
				p.assignElderlyGoals()
			}
		}

		// Check if anyone dies.
		if gameconstants.DiesAtAgeWithinNDays(int(p.Age), 1) {
			m.handleDeath(p)
		} else {
			remPop = append(remPop, p)
		}
	}
	// TODO: Find a more efficient way to do this besides re-allocation.
	m.RealPop = remPop
}

func (m *Map) matchSingles() {
	const ageDiffFactor = 0.25
	const menTakeFemaleName = true
	var men, women []*Person
	for _, p := range m.RealPop {
		if p.Dead || p.isMarried() || p.Age < 18 {
			continue
		}
		// Check if we even want a partner.
		if !p.Goals.IsSet(GoalAdultPartner) {
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
	for _, wp := range women {
		// Ladies' choice.
		for _, mp := range men {
			if mp.isMarried() {
				continue
			}
			// Check if we're on the same page regarding children.
			// TODO: Check other goals as well.
			if wp.Goals.IsSet(GoalAdultChildren) != mp.Goals.IsSet(GoalAdultChildren) {
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

	// TODO: Instead, we should go through both men and women and find their respective best match
	// based on their goals, opinions, etc.
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
				// Create a new child.
				child := m.newPerson("", p.LastName, byte(rand.Intn(2)), 0)
				child.X = p.X
				child.Y = p.Y
				child.Mother = p
				child.Father = p.Spouse

				log.Printf("Born: %v", child)

				// Assign the child to the home.
				if p.Home != nil {
					child.SetHome(p.Home)
				}

				// Assign the child to the parents.
				p.Children = append(p.Children, child)
				if p.Spouse != nil {
					p.Spouse.Children = append(p.Spouse.Children, child)
				}

				// Add the child to the population.
				m.RealPop = append(m.RealPop, child)
			}
		}
	}

	// All couples have a chance of getting pregnant, which decreases
	// by the number of children they already have.
	for _, p := range m.RealPop {
		if p.isMarried() && p.Pregnant == 0 && p.Gender == GenderFemale {
			// TODO: Check if both want children.
			if !p.Goals.IsSet(GoalAdultChildren) && !p.Spouse.Goals.IsSet(GoalAdultChildren) && rand.Intn(1000) > 1 {
				// If neither wants children, there is a very small chance of getting pregnant.
				continue
			} else if (!p.Goals.IsSet(GoalAdultChildren) || !p.Spouse.Goals.IsSet(GoalAdultChildren)) && rand.Intn(100) > 1 {
				// If only one or wants children, there is a slight chance of getting pregnant.
				continue
			}
			if rand.Intn(100) < 3-len(p.Children) {
				p.Pregnant = uint16(numDaysPregnant)
			}
		}
	}
}

func (m *Map) tickPeople(elapsed float64) {
	// TODO: Tick AI less often.
	for _, p := range m.RealPop {
		if p.Dead {
			continue
		}

		// Pick the motive to satisfy.
		// NOTE: This is the new, continuous way of doing things.
		m.pickMotive(p, elapsed)

		// Tick all the personal goals.
		// NOTE: This is the old, turn-based way of doing things.
		// m.tickPersonGoals(p, elapsed)
	}
}
