package genvegetation

import (
	"math"

	"github.com/Flokey82/go_gens/vectors"
)

type Sun struct {
	elevation float64
	azimuth   float64
}

func NewSun(elevation float64, azimuth float64) *Sun {
	return &Sun{
		elevation: elevation,
		azimuth:   azimuth,
	}
}

// ConvertToUvCoordinates transforms the polar coordinates of the sun to cartesian coordinates.
// Returns cartesian coordinates.
func (s *Sun) ConvertToUvCoordinates() vectors.Vec3 {
	u := math.Cos(toRadians(s.azimuth)) * math.Cos(toRadians(s.elevation))
	v := math.Sin(toRadians(s.azimuth)) * math.Cos(toRadians(s.elevation))
	w := math.Sin(toRadians(s.elevation))
	return vectors.Vec3{X: u, Y: v, Z: w}
}

func round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func toRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}
