package main

import (
	"github.com/Flokey82/genideas/simerosion"
)

func main() {
	m, _ := simerosion.NewHeightMapFromPNG("heightmap.png")
	vals := m.Elevations
	for i := 0; i < 1000; i++ {
		vals = m.ThermalErosion(vals)
	}
	m.ExportOBJ("test.obj", m.Elevations)
	m.ExportOBJ("test_eroded.obj", vals)
}
