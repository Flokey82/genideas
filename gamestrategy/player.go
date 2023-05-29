package gamestrategy

import (
	"log"
	"math/rand"
)

func (g *Grid) AddPlayer(p *Player) {
	g.Players = append(g.Players, p)
	// Find a random cell to start in that is not occupied
	for {
		c := g.Cell(rand.Intn(g.Width), rand.Intn(g.Height))
		if c == nil {
			continue
		}
		if c.ControlledBy == nil && c.Type != &TypeWater {
			c.Occupy(p)
			c.Type = &TypeCapital
			break
		}
	}

	// Add AI
	g.AIs = append(g.AIs, NewAI(p, g))
	log.Printf("Player %s added to grid", p.Name)
}

type Player struct {
	Name string
	Gold float64
}

func NewPlayer(name string) *Player {
	return &Player{
		Name: name,
		Gold: 100.0,
	}
}
