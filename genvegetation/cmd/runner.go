package main

import (
	"fmt"

	"github.com/Flokey82/genideas/genvegetation"
)

func main() {
	fmt.Println("Starting genvegetation...")
	c := genvegetation.NewController()
	c.LoadHeightAndSoilMap("default")
	c.CalculateAll("default", 12, 0, 0, 80, 0.1)
	c.PrepareProbabilitesCalculation("Beech", "default")
}
