package genfibonaccisphere

import (
	"math"
)

type FibonacciSphere struct {
	numPoints   int
	stepSize    float64
	goldenRatio float64
}

func NewFibonacciSphere(numPoints int) *FibonacciSphere {
	sphere := &FibonacciSphere{numPoints: numPoints}
	sphere.goldenRatio = (1 + math.Sqrt(5)) / 2
	sphere.stepSize = 2 * math.Pi / (sphere.goldenRatio * float64(sphere.numPoints))
	return sphere
}

func (sphere *FibonacciSphere) IndexToCoordinates(index int) (x, y, z float64) {
	phi := sphere.goldenRatio * float64(index)
	theta := 2 * math.Pi * phi
	y = 1 - 2.0/float64(sphere.numPoints)*float64(index)
	r := math.Sqrt(1 - y*y)
	x = math.Cos(theta) * r
	z = math.Sin(theta) * r
	return x, y, z
}

func (sphere *FibonacciSphere) IndexToLatLon(index int) (lat, lon float64) {
	x, y, z := sphere.IndexToCoordinates(index)
	lon = math.Atan2(x, y)
	lat = math.Asin(z)
	return lat, lon
}

func (sphere *FibonacciSphere) FindNearestNeighbors(index int) (above, below, left, right int) {
	// Calculate circumference at index
	_, idxY, idxZ := sphere.IndexToCoordinates(index)
	circumference := 2 * math.Pi * math.Sqrt(idxY*idxY+idxZ*idxZ)

	// Estimate neighbors based on circumference and step size
	above = index + int(circumference/sphere.stepSize)
	below = index - int(circumference/sphere.stepSize)
	left = index - 1
	right = index + 1

	// Refine nearest neighbors using Euclidean distance for above and below.
	// This function will search starting from the candidate index in the
	// positive or negative direction until it finds a point that is closest
	// to the index point than the step size. It will then return the index
	// of that point.
	findClosest := func(index, candidate int) int {
		// Calculate distance between index and candidate (just use the square, which is faster)
		distance := sphere.EuclideanDistanceSquare(index, candidate)

		for _, dir := range []int{1, -1} {
			// Now search until distance increases.
			for idx := candidate + dir; idx >= 0 && idx < sphere.numPoints; idx += dir {
				newDistance := sphere.EuclideanDistanceSquare(index, idx)
				if newDistance > distance {
					break
				}
				distance = newDistance
				candidate = idx
			}
		}
		return candidate
	}

	above = findClosest(index, above)
	below = findClosest(index, below)

	// TODO: Allow up to two results above and below each.
	// var above2 int
	// if sphere.EuclideanDistanceSquare(index, above+1) < sphere.EuclideanDistanceSquare(index, above-1) {
	// 	above2 = above+1
	// } else {
	// 	above2 = above-1
	// }
	//
	// var below2 int
	// if sphere.EuclideanDistanceSquare(index, below+1) < sphere.EuclideanDistanceSquare(index, below-1) {
	// 	below2 = below+1
	// } else {
	// 	below2 = below-1
	// }

	return above, below, left, right
}

func (sphere *FibonacciSphere) EuclideanDistance(index1, index2 int) float64 {
	x1, y1, z1 := sphere.IndexToCoordinates(index1)
	x2, y2, z2 := sphere.IndexToCoordinates(index2)
	return math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2) + (z1-z2)*(z1-z2))
}

func (sphere *FibonacciSphere) EuclideanDistanceSquare(index1, index2 int) float64 {
	x1, y1, z1 := sphere.IndexToCoordinates(index1)
	x2, y2, z2 := sphere.IndexToCoordinates(index2)
	return (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2) + (z1-z2)*(z1-z2)
}

func (sphere *FibonacciSphere) GreatArcDistance(index1, index2 int) float64 {
	lat1, lon1 := sphere.IndexToLatLon(index1)
	lat2, lon2 := sphere.IndexToLatLon(index2)
	return sphere.GreatArcDistanceLatLon(lat1, lon1, lat2, lon2)
}

func (sphere *FibonacciSphere) GreatArcDistanceLatLon(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert to radians.
	lat1 = lat1 * math.Pi / 180.0
	lon1 = lon1 * math.Pi / 180.0
	lat2 = lat2 * math.Pi / 180.0
	lon2 = lon2 * math.Pi / 180.0

	// Calculate great arc distance.
	return math.Acos(math.Sin(lat1)*math.Sin(lat2)+math.Cos(lat1)*math.Cos(lat2)*math.Cos(lon1-lon2)) * 180.0 / math.Pi
}
