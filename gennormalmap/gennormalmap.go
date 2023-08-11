// Package gennormalmap provides a function to generate a normal map from a height map.
// NOTE: This is a port of https://github.com/8bittree/normal_heights to Go.
package gennormalmap

import (
	"image"
	"image/color"
	"math"
)

// Why 6.0? Because that was the whole number that gave the closest results to
// the topographic map and normal map I was using as reference material.
// Considering my primary intent for creating this library is to create
// alternatives to those two files to use in the program they came with, it
// seemed like a good idea to match them, at least approximately.
const DefaultStrength = 6.0

// adjPixels represents the 8 pixels surrounding the current pixel.
type adjPixels struct {
	NW, N, NE float64
	W, E      float64
	SW, S, SE float64
}

func newAdjPixels(x, y int, img *image.Gray) *adjPixels {
	// Edge pixels are duplicated when necessary.
	// TODO: Cache max values.
	maxX := img.Bounds().Max.X
	maxY := img.Bounds().Max.Y

	// North coordinate
	n := y - 1
	if n < 0 {
		n = 0
	}

	// South coordinate
	s := y + 1
	if s >= maxY {
		s = maxY - 1
	}

	// West coordinate
	w := x - 1
	if w < 0 {
		w = 0
	}

	// East coordinate
	e := x + 1
	if e >= maxX {
		e = maxX - 1
	}

	return &adjPixels{
		NW: fetchPixel(w, n, img),
		N:  fetchPixel(x, n, img),
		NE: fetchPixel(e, n, img),
		W:  fetchPixel(w, y, img),
		E:  fetchPixel(e, y, img),
		SW: fetchPixel(w, s, img),
		S:  fetchPixel(x, s, img),
		SE: fetchPixel(e, s, img),
	}
}

// xNormals calculates the normals along the x-axis. Usually used for the red
// channel after normalization.
func (a *adjPixels) xNormals() float64 {
	return -(a.SE - a.SW + 2.0*(a.E-a.W) + a.NE - a.NW)
}

// yNormals calculates the normals along the y-axis. Usually used for the green
// channel after normalization.
func (a *adjPixels) yNormals() float64 {
	return -(a.NW - a.SW + 2.0*(a.N-a.S) + a.NE - a.SE)
}

// fetchPixel fetches the pixel at (x,y) and returns its value as an float64 scaled to
// between 0.0 and 1.0. Coordinate parameters are reversed from usual to better
// match compass directions.
func fetchPixel(x, y int, img *image.Gray) float64 {
	return float64(img.GrayAt(x, y).Y) / 255.0
}

// MapNormals creates the normal mapping from the given image with
// DEFAULT_STRENGTH.
func MapNormals(img *image.Gray) *image.RGBA {
	return MapNormalsWithStrength(img, DefaultStrength)
}

// MapNormalsWithStrength creates the normal mapping from the given image with the
// given strength.
func MapNormalsWithStrength(img *image.Gray, strength float64) *image.RGBA {
	bounds := img.Bounds()
	maxX := bounds.Max.X
	maxY := bounds.Max.Y

	normalMap := image.NewRGBA(bounds)
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			s := newAdjPixels(x, y, img)

			newP := [3]float64{s.xNormals(), s.yNormals(), 1.0 / strength}
			newP = scaleNormalizedTo0To1(normalize(newP))

			normalMap.SetRGBA(x, y, color.RGBA{
				R: uint8(newP[0] * 255.0),
				G: uint8(newP[1] * 255.0),
				B: uint8(newP[2] * 255.0),
				A: 255,
			})
		}
	}
	return normalMap
}

// normalize normalizes the given vector.
func normalize(v [3]float64) [3]float64 {
	vMag := math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
	return [3]float64{v[0] / vMag, v[1] / vMag, v[2] / vMag}
}

// scaleNormalizedTo0To1 scales the given vector from -1 to 1 to 0 to 1.
func scaleNormalizedTo0To1(v [3]float64) [3]float64 {
	return [3]float64{
		v[0]*0.5 + 0.5,
		v[1]*0.5 + 0.5,
		v[2]*0.5 + 0.5,
	}
}
