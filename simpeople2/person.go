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
}

// NewPerson creates a new person.
func (w *World) NewPerson(name string) *Person {
	return &Person{
		Name: name,
		Motives: []*Motive{
			MotiveTypeSleep.New(),
			MotiveTypeFood.New(),
			MotiveTypeFun.New(),
		},
		w: w,
		Position: vectors.Vec2{
			X: rand.Float64() * 50,
			Y: rand.Float64() * 50,
		},
	}
}

// Tick ticks the person.
func (p *Person) Tick() {
	for _, m := range p.Motives {
		m.Tick()
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
	actions := make([]*actionRank, 0)
	for _, o := range p.w.Objects {
		for _, a := range o.Actions {
			// MaxEffect is the maximum effect of the action taking into accoount the
			// difference between the current value of the motive and the max value which
			// limits the effect of the action.
			maxEffect := math.Min(a.Effect, missingToMax[a.Motive])

			r := &actionRank{
				Action:   a,
				Object:   o,
				Priority: maxEffect * multipliers[a.Motive],
			}
			r.Log()
			actions = append(actions, r)
		}
	}

	// Sort the actions by priority.
	sort.Slice(actions, func(i, j int) bool {
		return actions[i].Priority > actions[j].Priority
	})

	allowInterruption := false

	// Perform the action with the highest priority.
	ac := actions[0]
	if p.Action == nil || (allowInterruption && p.Action != ac.Action) {
		p.Action = ac.Action
		p.Destination = ac.Object

		log.Printf("%s: %s %s", p.Name, p.Action.Name, p.Destination.Name)
	}

	p.PerformAction()
}

const walkSpeed = 10.0 // How far the person can walk per tick

// PerformAction performs the current action.
func (p *Person) PerformAction() {
	if p.Action == nil {
		return
	}

	// If we haven't reached the destination yet, move towards it.
	if !p.Position.Equalish(p.Destination.Position) {
		log.Printf("%s: moving towards %s (distance %.2f)", p.Name, p.Destination.Name, p.Position.DistanceTo(p.Destination.Position))
		// Set the speed
		p.Speed = vectors.Normalize(p.Destination.Position.Sub(p.Position)).Mul(walkSpeed)

		// If we are faster than the distance to the destination, we set the speed to the distance.
		if p.Speed.Len() > p.Position.DistanceTo(p.Destination.Position) {
			p.Speed = vectors.Normalize(p.Destination.Position.Sub(p.Position)).Mul(p.Position.DistanceTo(p.Destination.Position))
		}

		// Move towards the destination
		p.Position = p.Position.Add(p.Speed)
		return
	}

	// We have reached the destination, perform the action.
	log.Printf("%s: performing %s", p.Name, p.Action.Name)
	findMotive := p.Action.Motive
	for _, m := range p.Motives {
		if m.Type == findMotive {
			m.Change(p.Action.Effect)
			// We have performed the action, reset the action.
			p.Action = nil
			p.Destination = nil
			break
		}
	}
}

// Log logs the current state of the person.
func (p *Person) Log() {
	log.Printf("%s: %.2f %.2f", p.Name, p.Position.X, p.Position.Y)
	for _, m := range p.Motives {
		m.Log()
	}
}

type actionRank struct {
	Action   *Action
	Object   *Object
	Priority float64
}

func (a *actionRank) Log() {
	log.Printf("%s %s: %.2f", a.Action.Name, a.Object.Name, a.Priority)
}
