package main

import (
	"fmt"

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
}
