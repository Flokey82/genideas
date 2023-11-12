package gamestrategy

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
)

type AI struct {
	*Player                                   // The player that the AI controls
	*Grid                                     // The grid that the AI is playing on
	DesirabilityModifiers map[string]float64  // How much the AI prefers a certain action
	Opinion               map[*Player]float64 // How much the AI likes or dislikes a player
}

func NewAI(p *Player, g *Grid) *AI {
	return &AI{
		Player: p,
		Grid:   g,
		DesirabilityModifiers: map[string]float64{
			ActionBuild:   (3 + rand.Float64()) / 3.0, // 1.0 - 1.33
			ActionExpand:  (3 + rand.Float64()) / 3.0, // 1.0 - 1.33
			ActionAttack:  (3 + rand.Float64()) / 3.0, // 1.0 - 1.33
			ActionAbandon: (3 + rand.Float64()) / 3.0, // 1.0 - 1.33
		},
		Opinion: map[*Player]float64{},
	}
}

const (
	ActionBuild   = "build"
	ActionExpand  = "expand"
	ActionAttack  = "attack"
	ActionAbandon = "abandon"
)

func (a *AI) Act() {
	// TODO:
	// - Separate short term and long term desirability.
	// - Short term desirability is based on the current state of the game.
	// - Long term desirability is based on the current and future state of the game (maximize future yield).
	// - Weight these depending on how financially stable (or savvy) the AI is, and how impulsive it is.

	// Calculate current balance per tick.
	var currentBalance float64
	for _, c := range a.Cells {
		if c.ControlledBy == a.Player {
			currentBalance += c.Yield() - c.Cost()
		}
	}

	calcDistancesToPlayer := func(p *Player) []float64 {
		distToPlayer := make([]float64, len(a.Cells))
		for i := range distToPlayer {
			distToPlayer[i] = -1
		}

		// Use a queue to do a breadth-first search
		queue := make([]*Cell, 0, 1000)
		// Start with all cells that the player owns
		for i := range a.Cells {
			c := &a.Cells[i]
			if c.ControlledBy == p {
				queue = append(queue, c)
				distToPlayer[i] = 0
			}
		}

		queuePos := 0
		truncateQueueAt := 1000

		for queuePos < len(queue) {
			c := queue[queuePos]
			queuePos++

			if queuePos >= truncateQueueAt {
				// Copy the remaining cells to the beginning of the queue
				copy(queue, queue[queuePos:])
				queue = queue[:len(queue)-queuePos]
				queuePos = 0
			}

			for _, n := range a.CellNeighbors(c.X, c.Y) {
				cellID := n.Y*a.Width + n.X
				if distToPlayer[cellID] == -1 {
					queue = append(queue, n)
					distToPlayer[cellID] = distToPlayer[cellID] + 1
				}
			}
		}

		return distToPlayer
	}

	log.Printf("AI %s acts", a.Name)

	// Calculate the distance to all players.
	distToPlayers := make([][]float64, len(a.Players))
	for i, p := range a.Players {
		distToPlayers[i] = calcDistancesToPlayer(p)
	}

	// Look at all cells that we own and decide what to do
	var possibleActions []Task
	for i := range a.Cells {
		c := &a.Cells[i]
		// Calculate the proximity value to the nearest hostile player.
		// TODO: Also calculate the proximity value to the nearest friendly player.
		var proximityToHostile float64
		for j, p := range a.Players {
			if p == a.Player {
				continue
			}
			if a.Opinion[p] < 0.0 {
				dist := distToPlayers[j][i]
				if dist == -1 {
					// This usually means that the player has been wiped out.
					continue
				}
				// If we are more than 10 cells away, we don't care about the proximity.
				if dist > 10 {
					continue
				}

				// The proximity is a value between 0 and 1, where 1 is the closest
				// and 0 is the furthest away.
				// TODO: Use opinion to determine how much we care about this player
				// being close to us.
				proxVal := (1 / max(dist, 1))
				proximityToHostile = max(proxVal, proximityToHostile-a.Opinion[p])
			}
		}

		// Check if we own this cell.
		if c.ControlledBy == a.Player {
			// We own this cell, so we might expand, build, or do nothing.

			// We can always abandon a cell if we're tight on money.
			if c.Type != &TypeCapital && currentBalance <= 0.0 {
				t := Task{
					Action:  ActionAbandon,
					At:      c,
					Payload: TaskAbandon{},
				}
				// We lose the yield, but also the cost of maintaining the cell.
				futureSavings := c.Cost() - c.Yield()
				if futureSavings <= 0.0 {
					continue
				}
				t.Desirability = a.DesirabilityModifiers[ActionAbandon] * (futureSavings + proximityToHostile)
				possibleActions = append(possibleActions, t)
				continue
			}

			// Check if we can build anything.
			if c.AllowedFeatures&^c.Features != 0 {
				// Subtract the current yield so we can calculate the future yield after building.
				currentBalanceExcl := currentBalance - c.Yield()
				for _, f := range splitFeatures(c.AllowedFeatures &^ c.Features) {
					// Build: Attempt to build a feature here
					// TODO:
					// - Also determine the one-time cost of building here.
					// - If this cell is close to a hostile player, we should build more military features.
					// E.g.:
					// - If we are threatened by the player, building defenses should be more desirable.
					// - If we are hostile towards the player, building offensive features should be more desirable.
					// - If we are friendly towards the player, building economic features should be more de
					t := Task{
						Action: ActionBuild,
						At:     c,
						Payload: TaskBuild{
							Feature: f,
						},
					}

					// Add potential future yield after building.
					futureBalance := currentBalanceExcl + resourceYield(c.Features|f)
					if futureBalance <= 0.0 || a.Gold < t.Cost() {
						continue
					}
					balanceDiff := futureBalance - currentBalance

					t.Desirability = a.DesirabilityModifiers[ActionBuild] * (balanceDiff - proximityToHostile)

					// We can build here
					possibleActions = append(possibleActions, t)
				}
			}

			// Check if we can expand anywhere around us.
			//
			// TODO:
			// - Prefer cells that have more than one neighbor that we own, this
			//   will create a more natural looking expansion.
			// - Cache a grid that represents the number of neighbors that we own.
			// - Depending on the personality, avoid expanding close to other players.
			for _, n := range a.CellNeighbors(c.X, c.Y) {
				if n.ControlledBy == a.Player {
					continue // We already own this cell
				}

				// Get the neighbors of the cell that we want to expand to or attack.
				nbs := a.CellNeighbors(n.X, n.Y)

				// Count the number of neighbors that we own.
				var numNeighborsOwned int
				for _, nb := range nbs {
					if nb.ControlledBy == a.Player {
						numNeighborsOwned++
					}
				}

				// Calculate the ownership factor, which is the fraction of neighbors that we own.
				// This is a value between 0 and 1, where 1 indicates that we own all neighbors.
				ownershipFactor := float64(numNeighborsOwned) / float64(len(nbs))

				// TODO: Calculate the desirability of expanding here.
				// Existing infrastructure, potential maximum yield, distance to capital, etc.
				if n.ControlledBy == nil {
					// Expansion: Attempt to claim this cell
					// TODO:
					// - Also determine the one-time cost of expanding here.
					// - If we have a hostile personality, we should expand more aggressively.
					t := Task{
						Action: ActionExpand,
						At:     n,
						Payload: TaskExpand{
							From: c,
						},
					}

					// Add potential future yield from this cell, subtract the cost of maintaining it.
					futureBalance := currentBalance + n.Yield() - n.Cost()
					// futureBalance += resourceYield(n.AllowedFeatures) / 2.0
					if futureBalance <= 0.0 || a.Gold < t.Cost() {
						continue
					}
					balanceDiff := futureBalance - currentBalance

					t.Desirability = a.DesirabilityModifiers[ActionExpand] * (balanceDiff + ownershipFactor - proximityToHostile)

					// We can expand here
					possibleActions = append(possibleActions, t)
				} else {
					// Attack: Attempt to take over this cell
					// TODO:
					// - Also determine the one-time cost of attacking here. (possible gain vs possible loss)
					// - If we already own many neighboring cells, attacking should be more desirable.
					// - If we have a hostile personality, attacking other hostiles or players that we are
					// hostile towards should be more desirable.
					t := Task{
						Action: ActionAttack,
						At:     n,
						Payload: TaskAttack{
							From: c,
							On:   n.ControlledBy,
						},
					}

					// Add potential future yield from this cell, subtract the cost of maintaining it.
					futureBalance := currentBalance + n.Yield() - n.Cost()
					if futureBalance <= 0.0 || a.Gold < t.Cost() {
						continue
					}
					balanceDiff := futureBalance - currentBalance

					t.Desirability = a.DesirabilityModifiers[ActionAttack] * (balanceDiff + ownershipFactor + proximityToHostile)

					// We can attack here
					possibleActions = append(possibleActions, t)
				}
			}
		}
	}

	// Randomize the order of possible actions to avoid artifacts.
	rand.Shuffle(len(possibleActions), func(i, j int) {
		possibleActions[i], possibleActions[j] = possibleActions[j], possibleActions[i]
	})

	// Evaluate all possible actions and pick the best one.
	// Sort by desirability.
	// TODO: Calculate utility instead by calculating the resources yield before and after the action,
	// the before and after per-turn cell cost, and the cost of the action.
	sort.Slice(possibleActions, func(i, j int) bool {
		pi := possibleActions[i]
		pj := possibleActions[j]

		// TODO: Depending on our finances, we might need to take cost into account.
		return pi.Cost()/pi.Desirability < pj.Cost()/pj.Desirability
	})

	doLog := false
	if doLog {
		for _, ac := range possibleActions {
			log.Printf("AI %s possible action %s on %d,%d for %f (%f)", a.Name, ac.Action, ac.At.X, ac.At.Y, ac.Cost(), ac.Desirability)
		}
	}

	// Pick the first one we can afford.
	for _, ac := range possibleActions {
		// Check if we can afford this action
		if ac.Cost() <= a.Gold {
			a.Do(ac)
			return
		}
	}
}

func (a *AI) Do(t Task) {
	switch t.Action {
	case ActionBuild:
		a.Build(t.At, t.Payload.(TaskBuild))
	case ActionExpand:
		a.Expand(t.At, t.Payload.(TaskExpand))
	case ActionAttack:
		a.Attack(t.At, t.Payload.(TaskAttack))
	case ActionAbandon:
		a.Abandon(t.At, t.Payload.(TaskAbandon))
	}
}

func (a *AI) Build(c *Cell, payload TaskBuild) {
	f := payload.Feature
	log.Printf("AI %s builds %d on %d,%d", a.Name, f, c.X, c.Y)

	// Only build features that are allowed and not already built.
	// NOTE: This should already be ensured by the task generation.
	if f&c.AllowedFeatures != 0 && f&c.Features == 0 {
		c.Features |= f
		cost := c.CostToBuild(f)
		log.Printf("AI %s builds %d on %d,%d for %f", a.Name, f, c.X, c.Y, cost)
		a.Gold -= cost

		// Broadcast the build to all players.
		a.Messenger.Broadcast(a.Player.ID, MsgBuild{
			FromID:  a.Player.ID,
			AtCell:  c,
			Feature: f,
		})
	}
}

func (a *AI) Expand(c *Cell, payload TaskExpand) {
	log.Printf("AI %s expands to %d,%d", a.Name, c.X, c.Y)
	cost := c.Cost()
	log.Printf("AI %s expands to %d,%d for %f", a.Name, c.X, c.Y, cost)
	a.Gold -= cost

	// Update the controlling player of the cell.
	c.ControlledBy = a.Player

	// Broadcast the expand action to all players.
	a.Messenger.Broadcast(a.Player.ID, MsgExpand{
		FromID:   a.Player.ID,
		FromCell: payload.From,
		AtCell:   c,
	})
}

func (a *AI) Attack(c *Cell, payload TaskAttack) {
	log.Printf("AI %s attacks %d,%d", a.Name, c.X, c.Y)
	cost := c.Cost() * 2.0
	log.Printf("AI %s attacks %d,%d for %f", a.Name, c.X, c.Y, cost)
	a.Gold -= cost

	// Check if we won
	// TODO: This should depend on ...
	// - how many neighboring cells we own vs. how many the opponent owns
	// - the number of offensive features that we have, and the number of defensive
	// features that the opponent has

	// Get the neighbors of the cell that we want to attack.
	nbs := a.CellNeighbors(c.X, c.Y)

	// Count the number of neighbors that we own vs. the opponent.
	var numNeighborsOwned, numNeighborsOpponent int
	for _, nb := range nbs {
		if nb.ControlledBy == a.Player {
			numNeighborsOwned++
		} else if nb.ControlledBy == payload.On {
			numNeighborsOpponent++
		}
	}

	// Ownership factor is the fraction of neighbors that we own vs. the opponent.
	ownershipFactor := float64(numNeighborsOwned) / float64(numNeighborsOpponent+1)

	var success bool
	if c.ControlledBy.Gold <= 0 || rand.Intn(int(c.ControlledBy.Gold)) < int(a.Gold*ownershipFactor) {
		log.Printf("AI %s wins against %s", a.Name, c.ControlledBy.Name)
		// Update the controlling player of the cell.
		c.ControlledBy = a.Player

		// We won!
		success = true
	} else {
		log.Printf("AI %s loses against %s", a.Name, c.ControlledBy.Name)
	}

	// TODO: Should our opinion of the attacked player actually change?
	// a.Opinion[c.ControlledBy] -= 0.1

	// Notify the opposite player that we are attacking.
	// TODO: Should this be a broadcast so allies can help?
	a.Messenger.Send(a.Player.ID, []int{payload.On.ID}, MsgAttack{
		FromID:  a.Player.ID,
		ToID:    c.ControlledBy.ID,
		Success: success,
	})
}

func (a *AI) Abandon(c *Cell, payload TaskAbandon) {
	log.Printf("AI %s abandons %d,%d", a.Name, c.X, c.Y)

	// Update the controlling player of the cell.
	c.ControlledBy = nil

	// Broadcast the abandon action to all players.
	a.Messenger.Broadcast(a.Player.ID, MsgAbandon{
		AtCell: c,
	})
}

// Receive a a message.
// TODO: Queue this up and process it in the next tick?
func (a *AI) Receive(from int, message any) {
	const maxDist = 10.0

	findClosestCell := func(start *Cell) *Cell {
		// Use a queue to do a breadth-first search
		queue := make([]*Cell, 0, 1000)
		queue = append(queue, start)
		queuePos := 0
		truncateQueueAt := 1000
		visited := map[*Cell]bool{}
		for queuePos < len(queue) {
			c := queue[queuePos]
			queuePos++
			if c.ControlledBy == a.Player {
				return c
			}
			if queuePos >= truncateQueueAt {
				// Copy the remaining cells to the beginning of the queue
				copy(queue, queue[queuePos:])
				queue = queue[:len(queue)-queuePos]
				queuePos = 0
			}
			for _, n := range a.CellNeighbors(c.X, c.Y) {
				if !visited[n] {
					queue = append(queue, n)
					visited[n] = true
				}
			}
		}
		return nil
	}

	log.Printf("AI %s received message from %d: %v", a.Name, from, message)
	switch message.(type) {
	case MsgBuild:
		// TODO: If the player builds some offensive feature close to us, we should like them less.
		msg := message.(MsgBuild)
		if msg.AtCell.ControlledBy == a.Player {
			return
		}

		// Get the minimum distance to our territory by expanding from the cell.
		foundCell := findClosestCell(msg.AtCell)
		if foundCell != nil {
			// We found a cell that is connected to our territory.
			// Calculate the distance to the cell.
			dist := math.Sqrt(float64((msg.AtCell.X-foundCell.X)*(msg.AtCell.X-foundCell.X) + (msg.AtCell.Y-foundCell.Y)*(msg.AtCell.Y-foundCell.Y)))
			if dist > maxDist {
				return // The cell is too far away to be considered.
			}
			modifier := 1.0
			if dist > 0.0 {
				modifier = 1 - dist/maxDist
			}
			// We like the player less the closer they are to our territory.
			a.ChangeOpinion(a.Players[msg.FromID], -0.1*modifier, fmt.Sprintf("building close to us (dist %f)", dist))
		}
	case MsgExpand:
		// TODO: If the player expands close to us, we should like them less, if they are not our ally.
		msg := message.(MsgExpand)
		if msg.FromCell.ControlledBy == a.Player {
			return
		}

		// Get the minimum distance to our territory by expanding from the cell.
		foundCell := findClosestCell(msg.FromCell)
		if foundCell != nil {
			// We found a cell that is connected to our territory.
			// Calculate the distance to the cell.
			dist := math.Sqrt(float64((msg.FromCell.X-foundCell.X)*(msg.FromCell.X-foundCell.X) + (msg.FromCell.Y-foundCell.Y)*(msg.FromCell.Y-foundCell.Y)))
			if dist > maxDist {
				return // The cell is too far away to be considered.
			}
			modifier := 1.0
			if dist > 0.0 {
				modifier = 1 - dist/maxDist
			}
			// We like the player less the closer they are to our territory.
			a.ChangeOpinion(a.Players[msg.FromID], -0.1*modifier, fmt.Sprintf("expanding close to us (dist %f)", dist))
		}
	case MsgAttack:
		msg := message.(MsgAttack)
		// If the attack was successful, we like the attacker even less.
		if msg.Success {
			a.ChangeOpinion(a.Players[msg.FromID], -0.2, "successful attack")
		} else {
			a.ChangeOpinion(a.Players[msg.FromID], -0.1, "failed attack")
		}
	case MsgAbandon:
		msg := message.(MsgAbandon)
		if msg.AtCell.ControlledBy == a.Player {
			return
		}
		// Take note of nearby abandoned cells.
		/*
			// Get the minimum distance to our territory by expanding from the cell.
			foundCell := findClosestCell(msg.AtCell)
			if foundCell != nil {
				// Calculate the distance to the cell.
				dist := math.Sqrt(float64((msg.AtCell.X-foundCell.X)*(msg.AtCell.X-foundCell.X) + (msg.AtCell.Y-foundCell.Y)*(msg.AtCell.Y-foundCell.Y)))
				if dist > maxDist {
					return // The cell is too far away to be considered.
				}
				modifier := 1.0
				if dist > 0.0 {
					modifier = 1 - dist/maxDist
				}
				// We like the player less the closer they are to our territory.
				a.ChangeOpinion(a.Players[msg.FromID], -0.1*modifier, fmt.Sprintf("abandoning close to us (dist %f)", dist))
			}
		*/
	default:
		log.Printf("AI %s received unknown message from %s: %v", a.Name, from, message)
	}
}

func (a *AI) ChangeOpinion(p *Player, amount float64, action string) {
	a.Opinion[p] += amount

	// Clamp the opinion to -1.0 to 1.0
	if a.Opinion[p] < -1.0 {
		a.Opinion[p] = -1.0
	} else if a.Opinion[p] > 1.0 {
		a.Opinion[p] = 1.0
	}

	log.Printf("AI %s changed opinion of %s by %f (%q), new %f", a.Name, p.Name, amount, action, a.Opinion[p])
}

type Task struct {
	Action       string  // build, expand, attack
	At           *Cell   // In which cell to do the action
	Desirability float64 // How desirable this action is
	Payload      any
}

type TaskBuild struct {
	Feature int64 // Which feature to build
}

type TaskExpand struct {
	From *Cell // From which cell to expand
}

type TaskAttack struct {
	From *Cell   // From which cell to attack
	On   *Player // Which player to attack
}

type TaskAbandon struct {
}

// HeightDiff returns the height difference between the from and at cells.
func (t *Task) HeightDiff() float64 {
	switch t.Action {
	case ActionExpand:
		taskExpand := t.Payload.(TaskExpand)
		return t.At.Value - taskExpand.From.Value
	case ActionAttack:
		taskAttack := t.Payload.(TaskAttack)
		return t.At.Value - taskAttack.From.Value
	}
	return 0.0
}

// Cost returns the cost of the task.
func (t *Task) Cost() float64 {
	switch t.Action {
	case ActionBuild:
		taskBuild := t.Payload.(TaskBuild)
		return t.At.Type.CostToBuild(taskBuild.Feature)
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
	case ActionAbandon:
		return 0
	}
	return 0.0
}
