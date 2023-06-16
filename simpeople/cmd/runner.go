package main

import (
	"github.com/Flokey82/genideas/simpeople"
)

func main() {
	factory := simpeople.NewPeopleFactory()
	var people []*simpeople.Person
	for i := 0; i < 40; i++ {
		p := factory.NewPerson()
		p.Log()
		st := p.Personality.Stats()
		st.Log()
		people = append(people, p)
	}

	// Now check if people get along.
	for _, a := range people {
		for _, b := range people {
			if a == b {
				continue
			}
			// Print the plan of A.
			plan := a.PickAction(b)
			if plan == nil {
				println(a.Name, "has no plan for", b.Name)
			} else {
				println("!", a.Name, "wants to", plan.Name, b.Name)
			}
		}
	}
}
