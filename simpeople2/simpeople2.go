package simpeople2

import (
	"github.com/hajimehoshi/ebiten"
)

// World represents the world of the simulation.
type World struct {
	*Level
	*Render
	People []*Person
}

// NewWorld creates a new world.
func NewWorld(width, height int) (*World, error) {
	r, err := newRender(width, height)
	if err != nil {
		return nil, err
	}

	return &World{
		Level:  NewLevel(width, height),
		Render: r,
	}, nil
}

// Tick ticks the world.
func (w *World) Tick(elapsed float64) {
	w.count++ // Move this to render?
	for _, p := range w.People {
		p.Tick(elapsed)
		p.Log()
	}
}

func (g *World) Update() error {
	g.Tick(1.0 / float64(ebiten.TPS()))

	// Handle input.
	g.handleInput()
	return nil
}
