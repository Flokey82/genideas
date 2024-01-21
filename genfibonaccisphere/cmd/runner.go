package main

import (
	"fmt"
	"log"

	"github.com/Flokey82/genideas/genfibonaccisphere"
)

func main() {
	sphere := genfibonaccisphere.NewFibonacciSphere(1000)
	index := 500
	above, below, left, right := sphere.FindNearestNeighbors(index)
	fmt.Println("Nearest neighbors of index", index)
	fmt.Println("Above:", above)
	fmt.Println("Below:", below)
	fmt.Println("Left:", left)
	fmt.Println("Right:", right)

	// Now let's check if we can convert the index back to coordinates and back to index.
	var mismatches int
	for i := 0; i < sphere.NumPoints; i++ {
		lat, lon := sphere.IndexToLatLonDeg(i)
		log.Println("lat:", lat, "lon:", lon)
		index := sphere.CoordinatesToIndex(lat, lon)
		fmt.Println(i, index)
		if i != index {
			mismatches++
		}
	}
	fmt.Println("Mismatches:", mismatches)
	genfibonaccisphere.Gen2dContinents()
}
