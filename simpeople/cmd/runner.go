package main

import (
	"log"

	"github.com/Flokey82/genideas/simpeople"
)

func main() {
	factory := simpeople.NewPeopleFactory()
	var people []*simpeople.Person
	for i := 0; i < 30; i++ {
		people = append(people, factory.NewPerson())
	}

	// Now check if people get along.
	for _, a := range people {
		a.Log()
		st := a.Personality.Stats()
		st.Log()
		for _, b := range people {
			if a == b {
				continue
			}
			// Print the plan of A.
			plan := a.PickAction(b)
			if plan == nil {
				log.Printf("%s has no plan for %s (%d)\n", a.Name, b.Name, a.Compatibility(b))
			} else {
				log.Printf("! %s wants to %s %s (%d)\n", a.Name, plan.Name, b.Name, a.Compatibility(b))
			}
		}
	}
}
