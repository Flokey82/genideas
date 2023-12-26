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
		p := &Person{
			LastName: m.lastGen.String(),
			Birthday: m.Day,
			Age:      rand.Intn(20) + 10,
			Gender:   rand.Intn(2),
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
		if p.Birthday == m.Day {
			p.Age++
		}
		// Check if anyone dies.
		if gameconstants.DiesAtAgeWithinNDays(p.Age, 1) {
			m.Population--
			// TODO: Remove from spouse?
			log.Printf("Died: %v", p)
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
		if p.isMarried() || p.Age < 18 {
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

// Person represents a person in the village.
type Person struct {
	Home      *Building
	Mother    *Person
	Father    *Person
	FirstName string
	LastName  string
	Birthday  int
	Age       int
	Gender    int
	Pregnant  int     // The number of days the person will still be pregnant.
	Spouse    *Person // The spouse of this person.
	Children  []*Person
}

// String returns the string representation of the person.
func (p *Person) String() string {
	return fmt.Sprintf("%s %s (%d %s)", p.FirstName, p.LastName, p.Age, p.genderString())
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
