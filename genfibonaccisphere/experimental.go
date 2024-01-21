package genfibonaccisphere

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"

	"github.com/Flokey82/genworldvoronoi/various"
	"github.com/mazznoer/colorgrad"
	"github.com/ojrac/opensimplex-go"
)

type SphereWithContinents struct {
	SeedPointLatLons [][2]float64
	noise            opensimplex.Noise
}

func NewSphereWithContinents(seedPointLatLons [][2]float64, seed int64) *SphereWithContinents {
	return &SphereWithContinents{
		SeedPointLatLons: seedPointLatLons,
		noise:            opensimplex.NewNormalized(seed),
	}
}

// FindIndexToPoint uses the ratio between the two closest seed points and their noise values to identify
// the index of the seed point that is closest to the given coordinates.
// This can be used to generate continents on the sphere.
func (sphere *SphereWithContinents) FindIndexToPoint(lat, lon float64) int {

	// Find the two closest seed points.
	index1, index2 := sphere.FindTwoClosestPoints(lat, lon)

	// TODO: Use bearing?

	// Calculate the angle from the first point and the given lat, lon.
	angle := sphere.GetAngleFromPoint(lat, lon, sphere.SeedPointLatLons[index1][0], sphere.SeedPointLatLons[index1][1]) - 90.0
	if angle < 0 {
		angle += 360.0
	}

	// Calculate the angle from the second point and the given lat, lon.
	angle2 := sphere.GetAngleFromPoint(lat, lon, sphere.SeedPointLatLons[index2][0], sphere.SeedPointLatLons[index2][1]) - 90.0
	if angle2 < 0 {
		angle2 += 360.0
	}

	// Get both noise values.
	noise1 := sphere.GetNoiseValueAngleFromIndex(index1, angle)
	noise2 := sphere.GetNoiseValueAngleFromIndex(index2, angle2)

	// Now given the distance between seed points and the lat, lon, we can can decide which index
	// should be assigned to the given lat, lon.
	// We use the ratio between the two distances and the ratio between the two noise values to
	// decide which index to use.
	dist1 := sphere.GreatArcDistanceLatLon(lat, lon, sphere.SeedPointLatLons[index1][0], sphere.SeedPointLatLons[index1][1])
	dist2 := sphere.GreatArcDistanceLatLon(lat, lon, sphere.SeedPointLatLons[index2][0], sphere.SeedPointLatLons[index2][1])
	//dist1 := sphere.GetCartesianDistance(lat, lon, sphere.SeedPointLatLons[index1][0], sphere.SeedPointLatLons[index1][1])
	//dist2 := sphere.GetCartesianDistance(lat, lon, sphere.SeedPointLatLons[index2][0], sphere.SeedPointLatLons[index2][1])

	// Calculate the ratio between the distances.
	distRatio1 := dist1 / (dist1 + dist2)
	distRatio2 := dist2 / (dist1 + dist2)

	// Calculate the ratio between the noise values.
	noiseRatio1 := (noise1 + 1) / 2
	noiseRatio2 := (noise2 + 1) / 2

	//log.Printf("dist1: %f, dist2: %f, distRatio: %f, noise1: %f, noise2: %f, noiseRatio: %f", dist1, dist2, distRatio, noise1, noise2, noiseRatio)

	if distRatio1*noiseRatio1 < distRatio2*noiseRatio2 {
		return index1
	}

	return index2
}

// GetAngleFromPoint returns the angle in degrees from the given point to the given point.
// The points are given in degrees.
func (sphere *SphereWithContinents) GetAngleFromPoint(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert to radians.
	lat1 = lat1 * math.Pi / 180.0
	lon1 = lon1 * math.Pi / 180.0
	lat2 = lat2 * math.Pi / 180.0
	lon2 = lon2 * math.Pi / 180.0

	// Calculate angle.
	return math.Atan2(math.Sin(lon2-lon1)*math.Cos(lat2), math.Cos(lat1)*math.Sin(lat2)-math.Sin(lat1)*math.Cos(lat2)*math.Cos(lon2-lon1)) * 180.0 / math.Pi
}

// GetNoiseValueAngleFromIndex returns the noise value at the given index and angle.
// The angle is given in degrees.
func (sphere *SphereWithContinents) GetNoiseValueAngleFromIndex(index int, angle float64) float64 {
	// The noise value will work across 360 degrees seamlessly.
	// We just use angle functions to calculate the noise value, so it alternates between
	// -1 and 1 as we go around the circle.
	// The noise value will be between -1 and 1.
	angleVal := math.Cos(angle * math.Pi / 180.0)
	noise1 := sphere.noise.Eval2(angleVal, float64(index))

	// Add octaves.
	noise2 := sphere.noise.Eval2(angleVal*2, float64(index)) / 2
	noise3 := sphere.noise.Eval2(angleVal*4, float64(index)) / 4
	noise4 := sphere.noise.Eval2(angleVal*8, float64(index)) / 8

	// Add octaves.
	noise := noise1 + noise2 + noise3 + noise4

	// Normalize the noise value to be between 0 and 1.
	noise = (noise + 1) / 2
	return noise
}

// FindTwoClosestPoints finds the two closest points on the sphere to the given point.
// It returns the indices of the two closest points.
func (sphere *SphereWithContinents) FindTwoClosestPoints(lat, lon float64) (index1, index2 int) {
	var minDistance1, minDistance2 float64
	for i := range sphere.SeedPointLatLons {
		lat2, lon2 := sphere.SeedPointLatLons[i][0], sphere.SeedPointLatLons[i][1]
		distance := sphere.GreatArcDistanceLatLon(lat, lon, lat2, lon2)
		if distance < minDistance1 || i == 0 {
			minDistance2 = minDistance1
			minDistance1 = distance
			index2 = index1
			index1 = i
		} else if distance < minDistance2 || i == 1 {
			minDistance2 = distance
			index2 = i
		}
	}
	return index1, index2
}

// FindClosestPoint finds the closest point on the sphere to the given point.
// It returns the index of the closest point.
func (sphere *SphereWithContinents) FindClosestPoint(lat, lon float64) int {
	var minDistance float64
	var index int
	for i := range sphere.SeedPointLatLons {
		lat2, lon2 := sphere.SeedPointLatLons[i][0], sphere.SeedPointLatLons[i][1]
		distance := sphere.GreatArcDistanceLatLon(lat, lon, lat2, lon2)
		if distance < minDistance || i == 0 {
			minDistance = distance
			index = i
		}
	}
	return index
}

// GreatArcDistanceLatLon calculates the great arc distance between two points on the sphere.
// The points are given in degrees.
func (sphere *SphereWithContinents) GreatArcDistanceLatLon(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert to radians.
	lat1 = lat1 * math.Pi / 180.0
	lon1 = lon1 * math.Pi / 180.0
	lat2 = lat2 * math.Pi / 180.0
	lon2 = lon2 * math.Pi / 180.0

	// Calculate great arc distance.
	return math.Acos(math.Sin(lat1)*math.Sin(lat2)+math.Cos(lat1)*math.Cos(lat2)*math.Cos(lon1-lon2)) * 180.0 / math.Pi
}

func (sphere *SphereWithContinents) GetCartesianDistance(lat1, lon1, lat2, lon2 float64) float64 {
	c1 := various.LatLonToCartesian(lat1, lon1)
	x1, y1, z1 := c1[0], c1[1], c1[2]
	c2 := various.LatLonToCartesian(lat2, lon2)
	x2, y2, z2 := c2[0], c2[1], c2[2]
	return math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2) + (z1-z2)*(z1-z2))
}

func Gen2dContinents() {

	// Generate 6 seed points:
	// North pole, south pole, and 4 points around the equator.
	seedPoints := [][2]float64{
		{90, 0},
		{-90, 0},
		{0, 0},
		{0, 90},
		{0, 180},
		{0, 270},
	}

	s := NewSphereWithContinents(seedPoints, 1000)
	w, h := 1000, 1000

	var points [][2]float64

	numPoints := 20

	// generate numPoints points around the center.
	/*
		centerX, centerY := float64(w)/2, float64(h)/2
		for i := 0; i < numPoints; i++ {
			angle := float64(i) * 360 / float64(numPoints)
			x := centerX + 250*math.Cos(angle*math.Pi/180)
			y := centerY + 250*math.Sin(angle*math.Pi/180)
			points = append(points, [2]float64{x, y})
		}
	*/

	// generate random points within the image.
	for i := 0; i < numPoints; i++ {
		x := rand.Intn(w)
		y := rand.Intn(h)
		points = append(points, [2]float64{float64(x), float64(y)})
	}

	var colors []color.Color

	// Generate a color gradient.
	colorGrad := colorgrad.Rainbow()
	for i := 0; i < len(points); i++ {
		colors = append(colors, colorGrad.At(float64(i)/float64(len(points))))
	}

	oldstuff(s, w, h, points, colors)
	newstuff(s, w, h, points, colors)

}

func oldstuff(s *SphereWithContinents, w, h int, points [][2]float64, colors []color.Color) {
	// Create a new image that we will draw our circle on.
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// Fill the image with white.
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, color.White)
		}
	}

	drawAtAngle := func(idx int, angle, radius, cx, cy float64, color color.Color) {
		// calculate the x, y position of the pixel
		// based on the angle and radius
		r := radius * s.GetNoiseValueAngleFromIndex(idx, angle)
		x := cx + r*math.Cos(angle*math.Pi/180)
		y := cy + r*math.Sin(angle*math.Pi/180)
		// set the pixel at that position to black
		img.Set(int(x), int(y), color)
	}

	findTwoClosest := func(x, y float64) (int, int) {
		// Find the two nearest points.
		var p1, p2 int
		var dist1, dist2 float64
		for i, p := range points {
			dist := math.Sqrt(math.Pow(float64(x)-p[0], 2) + math.Pow(float64(y)-p[1], 2))
			if dist < dist1 || dist1 == 0 {
				p2 = p1
				dist2 = dist1
				p1 = i
				dist1 = dist
			} else if dist < dist2 || dist2 == 0 {
				p2 = i
				dist2 = dist
			}
		}
		return p1, p2
	}

	// For all pixels, find the associated point.
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			// Find the two nearest points.
			p1Idx, p2Idx := findTwoClosest(float64(x), float64(y))
			p1 := points[p1Idx]
			p2 := points[p2Idx]

			// Get the angle between the point and p1 and p2.
			angle1 := angleBetweenXAndY(float64(x), float64(y), p1[0], p1[1])
			angle2 := angleBetweenXAndY(float64(x), float64(y), p2[0], p2[1])

			// Get the noise value for the angle.
			noise1 := s.GetNoiseValueAngleFromIndex(p1Idx, angle1)
			noise2 := s.GetNoiseValueAngleFromIndex(p2Idx, angle2)

			// Get the color for the noise value.
			color1 := colors[p1Idx]
			color2 := colors[p2Idx]

			// get the distance between the point and p1 and p2.
			dist1 := math.Sqrt(math.Pow(float64(x)-p1[0], 2) + math.Pow(float64(y)-p1[1], 2))
			dist2 := math.Sqrt(math.Pow(float64(x)-p2[0], 2) + math.Pow(float64(y)-p2[1], 2))

			distRatio1 := dist1 / (dist1 + dist2)
			distRatio2 := dist2 / (dist1 + dist2)

			// Get the noise ratio for the distance.
			noiseRatio1 := noise1
			noiseRatio2 := noise2

			// Get the final color.
			if noiseRatio1*distRatio1 < noiseRatio2*distRatio2 {
				img.Set(x, y, color1)
			} else {
				img.Set(x, y, color2)
			}
		}
	}

	// Log one full rotation noise values.
	for i := 0; i < 360; i++ {
		for j, p := range points {
			drawAtAngle(j, float64(i), 100, p[0], p[1], color.Black)
		}
	}

	// Save to file.
	f, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}

}

func angleBetweenXAndY(x1, y1, x2, y2 float64) float64 {
	// Calculate the angle between two points.
	// https://stackoverflow.com/questions/9614109/how-to-calculate-an-angle-from-points
	return math.Atan2(y2-y1, x2-x1) * 180 / math.Pi
}

func newstuff(s *SphereWithContinents, w, h int, points [][2]float64, colors []color.Color) {
	// Create a new image that we will draw our circle on.
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// Fill the image with white.
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, color.White)
		}
	}

	drawAtAngle := func(idx int, angle, radius, cx, cy float64, color color.Color) {
		// calculate the x, y position of the pixel
		// based on the angle and radius
		r := radius * (s.GetNoiseValueAngleFromIndex(idx, angle) + 1) / 2
		x := cx + r*math.Cos(angle*math.Pi/180)
		y := cy + r*math.Sin(angle*math.Pi/180)
		// set the pixel at that position to black
		img.Set(int(x), int(y), color)
	}

	findThreeClosest := func(x, y float64) (int, int, int) {
		// Find the three nearest points.
		var p1, p2, p3 int
		var dist1, dist2, dist3 float64
		for i, p := range points {
			dist := math.Sqrt(math.Pow(float64(x)-p[0], 2) + math.Pow(float64(y)-p[1], 2))
			if dist < dist1 || dist1 == 0 {
				p3 = p2
				dist3 = dist2
				p2 = p1
				dist2 = dist1
				p1 = i
				dist1 = dist
			} else if dist < dist2 || dist2 == 0 {
				p3 = p2
				dist3 = dist2
				p2 = i
				dist2 = dist
			} else if dist < dist3 || dist3 == 0 {
				p3 = i
				dist3 = dist
			}
		}
		return p1, p2, p3
	}
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			// Find the three nearest points.
			p1Idx, p2Idx, p3Idx := findThreeClosest(float64(x), float64(y))
			p1 := points[p1Idx]
			p2 := points[p2Idx]
			p3 := points[p3Idx]

			// Get the angles between the point and p1, p2, and p3.
			angle1 := angleBetweenXAndY(float64(x), float64(y), p1[0], p1[1])
			angle2 := angleBetweenXAndY(float64(x), float64(y), p2[0], p2[1])
			angle3 := angleBetweenXAndY(float64(x), float64(y), p3[0], p3[1])

			// Get the noise values for the angles.
			noise1 := s.GetNoiseValueAngleFromIndex(p1Idx, angle1)
			noise2 := s.GetNoiseValueAngleFromIndex(p2Idx, angle2)
			noise3 := s.GetNoiseValueAngleFromIndex(p3Idx, angle3)

			// Get the color for the noise values.
			color1 := colors[p1Idx]
			color2 := colors[p2Idx]
			color3 := colors[p3Idx]

			// Get the distances between the point and p1, p2, and p3.
			dist1 := math.Sqrt(math.Pow(float64(x)-p1[0], 2) + math.Pow(float64(y)-p1[1], 2))
			dist2 := math.Sqrt(math.Pow(float64(x)-p2[0], 2) + math.Pow(float64(y)-p2[1], 2))
			dist3 := math.Sqrt(math.Pow(float64(x)-p3[0], 2) + math.Pow(float64(y)-p3[1], 2))

			// Calculate the ratios between the distances.
			distRatio1 := dist1 / (dist1 + dist2 + dist3)
			distRatio2 := dist2 / (dist1 + dist2 + dist3)
			distRatio3 := dist3 / (dist1 + dist2 + dist3)

			// Calculate the ratios between the noise values.
			noiseRatio1 := noise1
			noiseRatio2 := noise2
			noiseRatio3 := noise3

			// Get the final color.
			var finalColor color.Color
			if noiseRatio1*distRatio1 < noiseRatio2*distRatio2 && noiseRatio1*distRatio1 < noiseRatio3*distRatio3 {
				finalColor = color1
			} else if noiseRatio2*distRatio2 < noiseRatio1*distRatio1 && noiseRatio2*distRatio2 < noiseRatio3*distRatio3 {
				finalColor = color2
			} else {
				finalColor = color3
			}

			img.Set(x, y, finalColor)
		}
	}

	// Log one full rotation noise values.
	for i := 0; i < 360; i++ {
		for j, p := range points {
			drawAtAngle(j, float64(i), 100, p[0], p[1], color.Black)
		}
	}

	// Save to file.
	f, err := os.Create("out2.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
