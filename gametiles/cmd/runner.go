package main

import (
	"log"

	"github.com/Flokey82/genideas/gametiles"
	"github.com/hajimehoshi/ebiten"
)

func main() {
	ebiten.SetWindowSize(gametiles.ScreenWidth, gametiles.ScreenHeight)
	ebiten.SetWindowTitle("Tiles (Ebitengine Demo)")
	g := gametiles.NewGame(25, 25)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
