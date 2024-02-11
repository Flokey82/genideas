package main

import (
	"log"
	"os"

	"github.com/Flokey82/genideas/maptiles"
	"github.com/hajimehoshi/ebiten"
)

func main() {
	g, err := maptiles.NewGame()
	if err != nil {
		os.Exit(1)
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
