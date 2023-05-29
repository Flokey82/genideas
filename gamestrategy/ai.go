package gamestrategy

import (
	"log"
	"math/rand"
	"sort"
)

type AI struct {
	*Player
	*Grid
}

func NewAI(p *Player, g *Grid) *AI {
	return &AI{
		Player: p,
		Grid:   g,
	}
}

const (
	ActionBuild  = "build"
	ActionExpand = "expand"
	ActionAttack = "attack"
)

func (a *AI) Act() {
	log.Printf("AI %s acts", a.Name)
	var possibleActions []Task
	// Look at all cells that we own and decide what to do
	for i := range a.Cells {
		c := &a.Cells[i]
		if c.ControlledBy == a.Player {
			// We own this cell
			// We might expand, build, or do nothing

			// Check if we can build anything.
			if c.AllowedFeatures&^c.Features != 0 {
				for _, f := range splitFeatures(c.AllowedFeatures &^ c.Features) {
					possibleActions = append(possibleActions, Task{
						Action:  ActionBuild,
						Cell:    c,
						Feature: f,
					})
				}
			}

			// Check if we can expand anywhere
			for _, n := range a.CellNeighbors(c.X, c.Y) {
				if n.ControlledBy == nil {
					// We can expand here
					possibleActions = append(possibleActions, Task{
						Action: ActionExpand,
						Cell:   n,
					})
				} else if n.ControlledBy != a.Player {
					// We can attack here
					possibleActions = append(possibleActions, Task{
						Action: ActionAttack,
						Cell:   n,
					})
				}
			}
		}
	}

	// Randomize the order of possible actions to avoid artifacts.
	rand.Shuffle(len(possibleActions), func(i, j int) {
		possibleActions[i], possibleActions[j] = possibleActions[j], possibleActions[i]
	})

	// Evaluate all possible actions and pick the best one
	// Sort by cost
	sort.Slice(possibleActions, func(i, j int) bool {
		return possibleActions[i].Cost() < possibleActions[j].Cost()
	})

	// Pick the first one
	if len(possibleActions) > 0 {
		for _, ac := range possibleActions {
			log.Printf("AI %s possible action %s on %d,%d for %f", a.Name, ac.Action, ac.Cell.X, ac.Cell.Y, ac.Cost())
		}
		if possibleActions[0].Cost() > a.Gold {
			// We can't afford this action
			return
		}
		a.Do(possibleActions[0])
	}
}

func (a *AI) Do(t Task) {
	switch t.Action {
	case ActionBuild:
		a.Build(t.Cell)
	case ActionExpand:
		a.Expand(t.Cell)
	case ActionAttack:
		a.Attack(t.Cell)
	}
}

func (a *AI) Build(c *Cell) {
	log.Printf("AI %s builds on %d,%d", a.Name, c.X, c.Y)
	// Pick a random feature to build
	for f := int64(1); f <= c.AllowedFeatures; f <<= 1 {
		if f&c.AllowedFeatures != 0 {
			c.Features |= f
			cost := c.CostToBuild(c.AllowedFeatures &^ c.Features)
			log.Printf("AI %s builds %d on %d,%d for %f", a.Name, f, c.X, c.Y, cost)
			a.Gold -= cost
			break
		}
	}
}

func (a *AI) Expand(c *Cell) {
	log.Printf("AI %s expands to %d,%d", a.Name, c.X, c.Y)
	c.ControlledBy = a.Player
	cost := c.Cost()
	log.Printf("AI %s expands to %d,%d for %f", a.Name, c.X, c.Y, cost)
	a.Gold -= cost
}

func (a *AI) Attack(c *Cell) {
	log.Printf("AI %s attacks %d,%d", a.Name, c.X, c.Y)
	c.ControlledBy = a.Player
	cost := c.Cost() * 2.0
	log.Printf("AI %s attacks %d,%d for %f", a.Name, c.X, c.Y, cost)
	a.Gold -= cost
}

type Task struct {
	Action  string
	Cell    *Cell
	Feature int64
}

func (t *Task) Cost() float64 {
	switch t.Action {
	case ActionBuild:
		return t.Cell.Type.CostToBuild(t.Feature)
	case ActionExpand:
		return t.Cell.Type.Cost
	case ActionAttack:
		return t.Cell.Type.Cost * 2.0
	}
	return 0.0
}
