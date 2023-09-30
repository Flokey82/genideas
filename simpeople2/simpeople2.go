package simpeople2

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
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
	for _, p := range w.People {
		if len(p.path) > 0 {
			p.count++ // For the walking animation.
		}
		p.Tick(elapsed)
		p.Log()
	}
}

func (g *World) Update() error {
	g.Tick(1.0 / float64(ebiten.TPS()))

	// Handle input.
	g.handleInput()

	// If we click, print the tile we clicked on.
	// TODO: Implement demonstration of pathfinding.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := g.getTileXY()
		fmt.Printf("Clicked on tile %d,%d\n", x, y)
	}
	return nil
}
