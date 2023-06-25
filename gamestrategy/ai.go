package gamestrategy

import (
	"log"
	"math"
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
						At:      c,
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
						From:   c,
						At:     n,
					})
				} else if n.ControlledBy != a.Player {
					// We can attack here
					possibleActions = append(possibleActions, Task{
						Action: ActionAttack,
						From:   c,
						At:     n,
					})
				}
			}
		}
	}

	// Randomize the order of possible actions to avoid artifacts.
	rand.Shuffle(len(possibleActions), func(i, j int) {
		possibleActions[i], possibleActions[j] = possibleActions[j], possibleActions[i]
	})

	// Evaluate all possible actions and pick the best one.
	// Sort by cost
	// TODO: Calculate utility instead by calculating the resources yield before and after the action,
	// the before and after per-turn cell cost, and the cost of the action.
	sort.Slice(possibleActions, func(i, j int) bool {
		return possibleActions[i].Cost() < possibleActions[j].Cost()
	})

	// Pick the first one
	if len(possibleActions) > 0 {
		for _, ac := range possibleActions {
			log.Printf("AI %s possible action %s on %d,%d for %f", a.Name, ac.Action, ac.At.X, ac.At.Y, ac.Cost())
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
		a.Build(t.At)
	case ActionExpand:
		a.Expand(t.At)
	case ActionAttack:
		a.Attack(t.At)
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
	Action  string // build, expand, attack
	From    *Cell  // From which cell to initiate the action
	At      *Cell  // In which cell to do the action
	Feature int64  // Which feature to build
}

// HeightDiff returns the height difference between the from and at cells.
func (t *Task) HeightDiff() float64 {
	if t.From != nil && t.At != nil {
		return t.At.Value - t.From.Value
	}
	return 0.0
}

// Cost returns the cost of the task.
func (t *Task) Cost() float64 {
	switch t.Action {
	case ActionBuild:
		return t.At.Type.CostToBuild(t.Feature)
	case ActionExpand:
		// Account for elevation difference as extra cost.
		// Expanding is preferred on the same elevation
		extra := math.Abs(t.HeightDiff())
		return t.At.Type.Cost + extra
	case ActionAttack:
		// Account for elevation difference as extra cost (positive or negative)
		// Attacking lower ground is cheaper, higher ground is more expensive.
		// NOTE: This is part of military strategy though, so should this knowledge be available to the AI
		// by default?
		extra := t.HeightDiff()
		return t.At.Type.Cost*2.0 + extra
	}
	return 0.0
}
