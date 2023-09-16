package simerosion

import (
	"bufio"
	"fmt"
	"image/png"
	"os"

	"github.com/fogleman/delaunay"
)

type HeightMap struct {
	width           int
	height          int
	Elevations      []float64
	VerticalScaling float64
}

func NewHeightMap(width, height int) *HeightMap {
	return &HeightMap{
		width:           width,
		height:          height,
		Elevations:      make([]float64, width*height),
		VerticalScaling: 100.0,
	}
}

func NewHeightMapFromPNG(filename string) (*HeightMap, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	elevations := make([]float64, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			elevations[y*width+x] = (float64(r)/65535.0 + float64(g)/65535.0 + float64(b)/65535.0) * 100 / 3.0
		}
	}
	return &HeightMap{
		width:           width,
		height:          height,
		Elevations:      elevations,
		VerticalScaling: 100.0,
	}, nil
}

func (m *HeightMap) Elevation(idx int) float64 {
	return m.Elevations[idx]
}

func (m *HeightMap) IdxToXY(idx int) (float64, float64) {
	x := idx % m.width
	y := idx / m.width
	return float64(x), float64(y)
}

func (m *HeightMap) XYToIdx(x, y int) int {
	return int(y)*m.width + int(x)
}

func (m *HeightMap) Neighbors(idx int) []int {
	x := idx % m.width
	y := idx / m.width
	neighbors := []int{}
	if x > 0 {
		neighbors = append(neighbors, idx-1)
	}
	if x < m.width-1 {
		neighbors = append(neighbors, idx+1)
	}
	if y > 0 {
		neighbors = append(neighbors, idx-m.width)
	}
	if y < m.height-1 {
		neighbors = append(neighbors, idx+m.width)
	}

	enableDiagonals := true
	if enableDiagonals {
		if x > 0 && y > 0 {
			neighbors = append(neighbors, idx-m.width-1)
		}
		if x < m.width-1 && y > 0 {
			neighbors = append(neighbors, idx-m.width+1)
		}
		if x > 0 && y < m.height-1 {
			neighbors = append(neighbors, idx+m.width-1)
		}
		if x < m.width-1 && y < m.height-1 {
			neighbors = append(neighbors, idx+m.width+1)
		}
	}

	return neighbors
}

// ExportOBJ returns a Wavefront OBJ file representing the heightmap.
func (m *HeightMap) ExportOBJ(path string, values []float64) error {
	tr, err := m.Triangulate()
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for i, p := range tr.Points {
		w.WriteString(fmt.Sprintf("v %f %f %f \n", p.X, values[i], p.Y)) //
	}
	for i := 0; i < len(tr.Triangles); i += 3 {
		w.WriteString(fmt.Sprintf("f %d %d %d \n", tr.Triangles[i]+1, tr.Triangles[i+1]+1, tr.Triangles[i+2]+1))
	}
	return nil
}
func (m *HeightMap) Triangulate() (*delaunay.Triangulation, error) {
	var pts []delaunay.Point
	for i := 0; i < len(m.Elevations); i++ {
		x, y := m.IdxToXY(i)
		pts = append(pts, delaunay.Point{X: x, Y: y})
	}
	return delaunay.Triangulate(pts)
}

func (m *HeightMap) ThermalErosion(vals []float64) []float64 {
	// Based on: https://aparis69.github.io/public_html/posts/terrain_erosion.html
	// And: https://aparis69.github.io/public_html/posts/terrain_erosion_2.html
	// Thermal erosion
	outData := make([]float64, len(vals))
	cellSize := 1.0
	amplitude := 1e-1

	const tanThresholdAngle = 0.6 // ~33°

	for x := 0; x < m.width; x++ {
		for y := 0; y < m.height; y++ {
			// Sample a 3x3 grid around the pixel
			samples := make([]float64, 9)
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					tapU, tapV := (x+i-1+m.width)%m.width, (y+j-1+m.height)%m.height
					samples[i*3+j] = vals[tapU*m.width+tapV]
				}
			}

			// Check stability with all neighbours
			id := m.XYToIdx(x, y)
			z := vals[id]
			willReceiveMatter := false
			willDistributeMatter := false
			for i := 0; i < 9; i++ {
				zd := samples[i] - z
				if zd/cellSize > tanThresholdAngle {
					willReceiveMatter = true
				}

				zd = z - samples[i]
				if zd/cellSize > tanThresholdAngle {
					willDistributeMatter = true
				}
			}

			// Add/Remove matter if necessary
			zOut := z
			if willReceiveMatter {
				zOut += amplitude
			}
			if willDistributeMatter {
				zOut -= amplitude
			}
			outData[id] = zOut
		}
	}

	return outData
}
