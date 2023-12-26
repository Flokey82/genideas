package main

import (
	"github.com/Flokey82/genideas/simsettlers"
)

func main() {
	m := simsettlers.NewMap(200, 200)
	m.Settle()
	for i := 0; i < 60000; i++ {
		m.Tick()
	}
	m.ExportPNG("test.png")
}
