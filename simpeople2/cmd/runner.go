package main

import (
	"log"

	"github.com/Flokey82/genideas/simpeople2"
	"github.com/hajimehoshi/ebiten"
)

func main() {
	// New world.
	w, err := simpeople2.NewWorld(25, 25)
	if err != nil {
		log.Fatal(err)
	}

	// Add some people.
	p1 := w.NewPerson("Person 1")
	w.People = append(w.People, p1)

	p2 := w.NewPerson("Person 2")
	w.People = append(w.People, p2)

	ebiten.SetWindowTitle("Simpeople2 (Ebitengine Demo)")
	ebiten.SetWindowSize(500, 600)
	ebiten.SetWindowResizable(true)

	if err := ebiten.RunGame(w); err != nil {
		log.Fatal(err)
	}
}
