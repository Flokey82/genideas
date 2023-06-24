package main

import (
	"log"

	"github.com/Flokey82/genideas/gametiles"
	"github.com/hajimehoshi/ebiten"
)

func main() {
	g := gametiles.NewGame()

	ebiten.SetWindowSize(gametiles.ScreenWidth*2, gametiles.ScreenHeight*2)
	ebiten.SetWindowTitle("Tiles (Ebitengine Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
