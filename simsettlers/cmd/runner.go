package main

import (
	"fmt"

	"github.com/Flokey82/genideas/simsettlers"
)

func main() {
	m := simsettlers.NewMap(200, 200)
	m.Settle()
	for i := 0; i < 100*365; i++ {
		m.Tick()
	}
	m.ExportPNG("test.png")

	// Log all houses and their occupants.
	for _, b := range m.Buildings {
		if b.Type == simsettlers.BuildingTypeHouse {
			fmt.Println(b.String())
			for _, p := range b.Occupants {
				fmt.Printf("\t%v\n", p)
				if p.Home != b {
					panic("Person not living in own home")
				}
			}
			// Log all owners.
			for _, p := range b.Owners {
				fmt.Printf("\tOwner: %v\n", p)
			}
		}
	}

	// Log all people and their home.
	for _, p := range m.RealPop {
		str := p.String()
		if p.Home == nil {
			str += " (homeless)"
		}
		if p.Job == simsettlers.JobTypeUnemployed {
			str += " (unemployed)"
		}
		fmt.Println(str)
	}

	// Log all the people in the cemetery.
	for _, p := range m.Cemetery.Occupants {
		fmt.Printf("Cemetery: %v - %s\n", p, p.Goals.String())
	}
}
