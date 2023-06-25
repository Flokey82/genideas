package gamestrategy

import (
	"log"

	"github.com/ojrac/opensimplex-go"
)

type Grid struct {
	Width   int
	Height  int
	Cells   []Cell
	Players []*Player
	AIs     []*AI
	*webpExport
}

func NewGrid(width, height int) *Grid {
	g := &Grid{
		Width:      width,
		Height:     height,
		webpExport: newWebPExport(width, height),
	}
	g.Cells = make([]Cell, width*height)
	// Use noise bands to set the type of each cell
	noise := opensimplex.NewNormalized(0)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := g.Cell(x, y)
			c.X = x
			c.Y = y
			val := noise.Eval2(float64(x)*2/float64(width), float64(y)*2/float64(height)) * 0.5
			val += noise.Eval2(float64(x)*4/float64(width), float64(y)*4/float64(height)) * 0.25
			val += noise.Eval2(float64(x)*8/float64(width), float64(y)*8/float64(height)) * 0.25
			c.Value = val
			if val < 0.2 {
				c.Type = &TypeWater
			} else if val < 0.4 {
				c.Type = &TypeMeadow
			} else if val < 0.6 {
				c.Type = &TypeForest
			} else {
				c.Type = &TypeMountain
			}
		}
	}
	log.Printf("Grid created with %d cells", len(g.Cells))
	return g
}

func (g *Grid) Tick() {
	// Loop over all cells and deduct cost from player
	for _, c := range g.Cells {
		if c.ControlledBy != nil {
			c.ControlledBy.Gold += c.Yield() - c.Cost()
		}
	}

	log.Printf("Tick! Players: %d, AIs: %d", len(g.Players), len(g.AIs))

	// AI actions
	for _, ai := range g.AIs {
		log.Println("AI tick for player", ai.Player.Name, "with gold", ai.Player.Gold)
		ai.Act()
	}

	// Check who is bankrupt
	for _, p := range g.Players {
		if p.Gold < 0 {
			log.Printf("Player %s is bankrupt!", p.Name)
		}
	}
	g.storeWebPFrame()
}

func (g *Grid) Cell(x, y int) *Cell {
	if x < 0 || x >= g.Width || y < 0 || y >= g.Height {
		return nil
	}
	return &g.Cells[y*g.Width+x]
}

func (g *Grid) CellNeighbors(x, y int) []*Cell {
	c := g.Cell(x, y)
	if c == nil {
		return nil
	}
	neighbors := make([]*Cell, 0, 8)
	for _, d := range Directions {
		n := g.Cell(c.X+d.X, c.Y+d.Y)
		if n != nil {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

func (g *Grid) Occupy(x, y int, p *Player) bool {
	c := g.Cell(x, y)
	if c == nil {
		return false
	}
	return c.Occupy(p)
}

var Directions = []struct {
	X int
	Y int
}{
	{X: -1, Y: -1},
	{X: 0, Y: -1},
	{X: 1, Y: -1},
	{X: -1, Y: 0},
	{X: 1, Y: 0},
	{X: -1, Y: 1},
	{X: 0, Y: 1},
	{X: 1, Y: 1},
}
