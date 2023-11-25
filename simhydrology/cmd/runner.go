package main

import (
	"image/png"
	"log"
	"os"

	"github.com/Flokey82/genideas/simhydrology"
)

func main() {
	hm := NewHeightFromPNG("heightmap4.png")
	m := simhydrology.NewMap(hm)
	m.ExportOBJ("test.obj")
	m.ExportPNG("test.png")
	m.ExportFluxPNG("test_flux.png")
	m.ExportSoilPNG("test_soil.png")
	m.ExportSinksPNG("test_sinks.png")
}

type FakeHeightmap struct {
	width, height int
	Elevations    [][]float64
}

func NewHeightFromPNG(filename string) *FakeHeightmap {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	elevations := make([][]float64, height)
	for y := 0; y < height; y++ {
		elevations[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			elevations[y][x] = (float64(r)/65535.0 + float64(g)/65535.0 + float64(b)/65535.0) / 3.0
		}
	}
	return &FakeHeightmap{
		width:      width,
		height:     height,
		Elevations: elevations,
	}

}

func (m *FakeHeightmap) Width() int {
	return m.width
}

func (m *FakeHeightmap) Height() int {
	return m.height
}

func (m *FakeHeightmap) Elevation(idx int) float64 {
	x := idx % m.width
	y := idx / m.width
	return m.Elevations[y][x]
}

func (m *FakeHeightmap) IdxToXY(idx int) (float64, float64) {
	x := idx % m.width
	y := idx / m.width
	return float64(x), float64(y)
}

func (m *FakeHeightmap) Neighbors(idx int) []int {
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

func (m *FakeHeightmap) NumRegions() int {
	return m.height * m.width
}
