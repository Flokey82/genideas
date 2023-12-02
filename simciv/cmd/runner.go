package main

import (
	"log"

	"github.com/Flokey82/genideas/simciv"
	"github.com/hajimehoshi/ebiten"
)

func main() {
	g := simciv.NewGame(100, 100)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
