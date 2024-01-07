package main

import (
	"fmt"
	"time"

	"github.com/Flokey82/gameloop"
	"github.com/Flokey82/genideas/simsettlers"
)

func main() {
	m := simsettlers.NewMap(200, 200)
	m.Settle()

	loop := gameloop.New(time.Second/60, m.Tick)
	loop.Start()
	// Wait for a keypress.
	fmt.Scanln()
	loop.Stop()

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
		// Log opinions.
		for opp, o := range p.Opinions {
			fmt.Printf("\t%v: %v\n", opp, o)
		}
	}

	// Log all the people in the cemetery.
	for _, p := range m.Cemetery.Occupants {
		fmt.Printf("Cemetery: %v - %s\n", p, p.Goals.String())
	}

	m.Export.ExportWebp("test.webp")
}
