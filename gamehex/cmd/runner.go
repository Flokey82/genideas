package main

import (
	"log"

	"github.com/Flokey82/genideas/gamehex"
	"github.com/hajimehoshi/ebiten"
)

func main() {
	ebiten.SetWindowTitle("Hexagonal (Ebitengine Demo)")
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowResizable(true)

	g, err := gamehex.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	if err = ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
