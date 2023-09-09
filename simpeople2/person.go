package simpeople2

import (
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/Flokey82/go_gens/vectors"
)

// Person is a person in the simulation.
type Person struct {
	Name     string       // Name of the person
	Motives  []*Motive    // Motives of the person
	Position vectors.Vec2 // Position of the person
	Speed    vectors.Vec2 // Speed of the person
	w        *World       // The world the person is in

	// Current action the person is performing.
	// If nil, the person is not performing an action.
	Action      *Action
	Destination *Object
	path        []*Node
	pathIdx     int
}

// NewPerson creates a new person.
func (w *World) NewPerson(name string) *Person {
	p := &Person{
		Name: name,
		Motives: []*Motive{
			MotiveTypeSleep.New(),
			MotiveTypeFood.New(),
			MotiveTypeBladder.New(),
			MotiveHygiene.New(),
			MotiveTypeFun.New(),
		},
		w: w,
		Position: vectors.Vec2{
			X: float64(rand.Intn(w.Width)),
			Y: float64(rand.Intn(w.Height)),
		},
	}
	if w.IsSolid(int(p.Position.X), int(p.Position.Y)) {
		p.Position.X = float64(rand.Intn(w.Width))
		p.Position.Y = float64(rand.Intn(w.Height))
	}
	return p
}

// Happiness returns the happiness of the person.
func (p *Person) Happiness() float64 {
	var sum float64
	for _, m := range p.Motives {
		sum += m.Val
	}
	return sum / float64(len(p.Motives))
}

// Tick ticks the person.
func (p *Person) Tick(elapsed float64) {
	for _, m := range p.Motives {
		m.Tick(elapsed)
	}

	// Get all the current multipliers for the motives.
	multipliers := make(map[*MotiveType]float64)
	missingToMax := make(map[*MotiveType]float64)
	for _, m := range p.Motives {
		multipliers[m.Type] = m.Multiplier()
		missingToMax[m.Type] = m.MissingToMax()
	}

	// Calculate the priority of each action by multiplying the effect
	// of the action with the current multiplier of the motive.
	var actions []*actionRank
	var current *actionRank

	evaluateSideEffect := true
	for _, o := range p.w.Objects {
		for _, a := range o.Actions {
			// MaxEffect is the maximum effect of the action taking into accoount the
			// difference between the current value of the motive and the max value which
			// limits the effect of the action.
			maxEffect := math.Min(a.Effect.Effect, missingToMax[a.Effect.Motive])

			// Priority is the effect of the action multiplied by the current multiplier
			// of the motive.
			priority := maxEffect * multipliers[a.Effect.Motive]
			if evaluateSideEffect && a.SideEffect != nil {
				maxSideEffect := math.Min(a.SideEffect.Effect, missingToMax[a.SideEffect.Motive])
				priority += maxSideEffect * multipliers[a.SideEffect.Motive]
			}

			// TODO: Also take distance into account.
			priority -= p.Position.DistanceTo(o.Position)

			r := &actionRank{
				Action:    a,
				Object:    o,
				Priority:  priority,
				MaxEffect: maxEffect,
			}
			r.Log()
			actions = append(actions, r)
			if p.Action != nil && p.Action == a {
				current = r
			}
		}
	}

	// Check if we have any actions.
	if len(actions) == 0 {
		return
	}

	// Sort the actions by decreasing priority.
	sort.Slice(actions, func(i, j int) bool {
		return actions[i].Priority > actions[j].Priority
	})

	// for debug log all current priorities.
	log.Println("Priorities:")
	for _, a := range actions {
		a.Log()
	}

	// Log multiplier of each motive.
	log.Println("Multipliers:")
	for _, m := range p.Motives {
		log.Printf("%s: %.2f", m.Type.Name, m.Multiplier())
	}

	// Check if we allow interruption and the factor by which the priority must be higher.
	allowInterruption := true
	interuptionMultiplier := 10.0

	// Perform the action with the highest priority.
	// TODO: Add randomness and pick from the top 3 or so, depending on how much the priority differs.
	// If the top priority would be miles ahead of the second, we should pick the top one regardless.
	ac := actions[0]

	const (
		minInterruptPriority  = 50.0
		minInterruptThreshold = 10.0
	)

	if p.Action == nil || (allowInterruption &&
		p.Action != ac.Action &&
		ac.Priority > current.Priority*interuptionMultiplier &&
		ac.Priority > minInterruptPriority &&
		current.Priority < minInterruptThreshold) {
		p.Action = ac.Action
		p.Destination = ac.Object
		log.Printf("%s: %s %s", p.Name, p.Action.Name, p.Destination.Name)
	} else {
		log.Printf("%s: continuing %s %s", p.Name, p.Action.Name, p.Destination.Name)
	}

	p.PerformAction(elapsed)
}

const walkSpeed = 15.0 // How far the person can walk per tick

// PerformAction performs the current action.
func (p *Person) PerformAction(elapsed float64) {
	if p.Action == nil {
		return
	}

	if p.path == nil {
		p.path = findPath(p.w, p, p.Destination)
		p.pathIdx = 0
	}

	// If we haven't reached the destination yet, move towards it.
	if !p.Position.Equalish(p.Destination.Position) {
		log.Printf("%s: moving towards %s (distance %.2f)", p.Name, p.Destination.Name, p.Position.DistanceTo(p.Destination.Position))
		// Set the speed to the direction of the destination.

		// Get the tile we want to move to. If we are close to the current index, move to the next one.
		tileX := int(p.path[p.pathIdx].X)
		tileY := int(p.path[p.pathIdx].Y)

		// Check if the distance to the next tile is less than the distance we walk in the
		// elapsed time. If so, move to the next tile (if there is one).
		walkDist := walkSpeed * elapsed
		if p.Position.DistanceTo(vectors.Vec2{
			X: float64(tileX),
			Y: float64(tileY),
		}) < walkDist && p.pathIdx < len(p.path)-1 {
			p.pathIdx++
			if p.pathIdx >= len(p.path) {
				p.path = nil
				p.pathIdx = 0
				return
			}
			tileX = int(p.path[p.pathIdx].X)
			tileY = int(p.path[p.pathIdx].Y)
		}

		// Set the speed to the direction of the destination.
		distVec := vectors.Normalize(vectors.Vec2{
			X: float64(tileX),
			Y: float64(tileY),
		}.Sub(p.Position))

		// If we are faster than the distance to the destination, we set the speed to the distance.
		if distVec.Len() > walkDist {
			p.Speed = distVec.Mul(walkDist)
		} else {
			p.Speed = distVec
		}
		// p.Speed = vectors.Normalize(p.Destination.Position.Sub(p.Position)).Mul(walkSpeed * elapsed)

		// If we are faster than the distance to the destination, we set the speed to the distance.
		if p.Speed.Len() > p.Position.DistanceTo(p.Destination.Position) {
			p.Speed = vectors.Normalize(p.Destination.Position.Sub(p.Position)).Mul(p.Position.DistanceTo(p.Destination.Position))
		}

		// Move towards the destination
		p.Position = p.Position.Add(p.Speed)
		return
	} else {
		p.path = nil
		p.pathIdx = 0
	}

	// We have reached the destination, perform the action.
	log.Printf("%s: performing %s", p.Name, p.Action.Name)

	// Apply primary motive change.
	p.ApplyEffect(p.Action.Effect, elapsed)

	// Apply secondary motive change.
	p.ApplyEffect(p.Action.SideEffect, elapsed)

	// Reset the action.
	// TODO: Continue action if not yet considered "completed"
	// p.Action = nil
	// p.Destination = nil
}

// ApplyEffect applies the effect to the person.
func (p *Person) ApplyEffect(e *Effect, elapsed float64) {
	if e == nil {
		return
	}
	for _, m := range p.Motives {
		if m.Type == e.Motive {
			m.Change(e.Effect * elapsed)
			break
		}
	}
}

// Log logs the current state of the person.
func (p *Person) Log() {
	log.Printf("%s: %.2f %.2f (Happiness: %.2f)", p.Name, p.Position.X, p.Position.Y, p.Happiness())
	log.Println("Motive values:")
	for _, m := range p.Motives {
		m.Log()
	}
}

type actionRank struct {
	Action    *Action
	Object    *Object
	Priority  float64
	MaxEffect float64
}

func (a *actionRank) Log() {
	log.Printf("%s %s: %.2f (max %2f)", a.Action.Name, a.Object.Name, a.Priority, a.MaxEffect)
}
