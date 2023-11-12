package gamestrategy

import (
	"log"
	"math/rand"
)

func (g *Grid) AddPlayer(p *Player) {
	p.ID = len(g.Players)
	log.Printf("Adding player %s with ID %d", p.Name, p.ID)
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
	ai := NewAI(p, g)
	g.AIs = append(g.AIs, ai)
	g.Messenger.Register(p.ID, ai)

	// TODO: Register with messenger
	log.Printf("Player %s added to grid", p.Name)
}

type Player struct {
	ID   int
	Name string
	Gold float64
}

func NewPlayer(name string) *Player {
	return &Player{
		Name: name,
		Gold: 100.0,
	}
}
