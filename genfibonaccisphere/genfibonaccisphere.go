package genfibonaccisphere

import (
	"math"

	"github.com/Flokey82/genworldvoronoi/various"
)

type FibonacciSphere struct {
	NumPoints          int
	s                  float64
	dlong              float64
	dz                 float64
	circumferenceSteps int
}

func NewFibonacciSphere(numPoints int) *FibonacciSphere {
	sphere := &FibonacciSphere{NumPoints: numPoints}

	// Second algorithm from http://web.archive.org/web/20120421191837/http://www.cgafaq.info/wiki/Evenly_distributed_points_on_sphere
	sphere.s = 3.6 / math.Sqrt(float64(numPoints))
	sphere.dz = 2.0 / float64(numPoints)

	// Calculate the longitude step size.
	sphere.dlong = math.Pi * (3 - math.Sqrt(5)) // ~2.39996323

	// Calculate the number of steps in the circumference.
	// This is the number of steps needed for nodes to complete a full circle.
	sphere.circumferenceSteps = int(2 * math.Pi / (math.Pi - sphere.dlong))
	return sphere
}

func (sphere *FibonacciSphere) IndexToCoordinates(index int) (x, y, z float64) {
	xyz := various.LatLonToCartesian(sphere.IndexToLatLonRad(index))
	return xyz[0], xyz[1], xyz[2]
}

func (sphere *FibonacciSphere) IndexToLatLonRad(index int) (lat, lon float64) {
	// Calculate latitude as z value from -1 to 1.
	z := 1 - (sphere.dz / 2) - float64(index)*sphere.dz

	// Calculate longitude in rad.
	long := float64(index) * sphere.dlong

	// Calculate the radius at the given z.
	// r := math.Sqrt(1 - z*z)

	// Calculate latitude and longitude in degrees.
	lat = math.Asin(z)
	lon = long

	return lat, math.Mod(lon, 2*math.Pi)
}

func (sphere *FibonacciSphere) IndexToLatLonDeg(index int) (lat, lon float64) {
	lat, lon = sphere.IndexToLatLonRad(index)
	lat = lat * 180.0 / math.Pi
	lon = lon * 180.0 / math.Pi
	return lat, lon
}

func (sphere *FibonacciSphere) CoordinatesToIndex(lat, lon float64) int {
	// Convert to radians.
	lat = lat * math.Pi / 180.0
	lon = lon * math.Pi / 180.0

	latSteps := int((1 - (sphere.dz / 2) - math.Sin(lat)) / sphere.dz)

	// Now search between latStep - circumferenceSteps and latStep + circumferenceSteps
	// for the longitude that is closest to the given longitude.
	var index int
	minDistance := math.Inf(1)

	// NOTE: This could be improved by searching in two directions and stopping when the distance
	// increases.
	useOldMethod := true

	if useOldMethod {
		for i := latSteps - sphere.circumferenceSteps; i < latSteps+sphere.circumferenceSteps; i++ {
			// Calculate the longitude of the current index.
			_, lon2 := sphere.IndexToLatLonRad(i)

			// Calculate the distance between the given longitude and the current longitude.
			distance := math.Abs(lon - lon2)

			// If the distance is smaller than the current minimum distance, update the minimum distance
			// and the index.
			if distance < minDistance {
				minDistance = distance
				index = i
			}
		}
	} else {
		// Calculate distance between index and candidate (just use the square, which is faster)
		distance := sphere.EuclideanDistanceSquare(index, latSteps)
		for _, dir := range []int{1, -1} {
			// Now search until distance increases.
			for idx := latSteps + dir; idx >= 0 && idx < sphere.NumPoints; idx += dir {
				newDistance := sphere.EuclideanDistanceSquare(index, idx)
				if newDistance > distance {
					break
				}
				distance = newDistance
				latSteps = idx
			}
		}
		index = latSteps
	}

	return index
}

func (sphere *FibonacciSphere) FindNearestNeighbors(index int) (above, below, left, right int) {
	// Estimate neighbors based on circumference and step size
	above = index + sphere.circumferenceSteps
	below = index - sphere.circumferenceSteps
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
			for idx := candidate + dir; idx >= 0 && idx < sphere.NumPoints; idx += dir {
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

func (sphere *FibonacciSphere) GreatArcDistanceLatLon(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert to radians.
	lat1 = lat1 * math.Pi / 180.0
	lon1 = lon1 * math.Pi / 180.0
	lat2 = lat2 * math.Pi / 180.0
	lon2 = lon2 * math.Pi / 180.0

	// Calculate great arc distance.
	return math.Acos(math.Sin(lat1)*math.Sin(lat2)+math.Cos(lat1)*math.Cos(lat2)*math.Cos(lon1-lon2)) * 180.0 / math.Pi
}

func (sphere *FibonacciSphere) EuclideanDistanceSquare(index1, index2 int) float64 {
	x1, y1, z1 := sphere.IndexToCoordinates(index1)
	x2, y2, z2 := sphere.IndexToCoordinates(index2)
	return (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2) + (z1-z2)*(z1-z2)
}
