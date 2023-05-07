package main

import (
	"fmt"

	"github.com/Flokey82/genideas/simciv"
)

func main() {
	m := simciv.NewMap(100, 100, 0)
	for i := 0; i < 30000; i++ {
		m.Tick()
	}
	for _, s := range m.Settlements {
		fmt.Println(s)
	}
}
