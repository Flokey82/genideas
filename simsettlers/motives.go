package simsettlers

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
)

// pickMotive picks the most urgent motive to satisfy.
func (m *Map) pickMotive(p *Person, elapsed float64) {
	// TODO: Pick the top motive to satisfy and keep it as the current motive.
	// - If we keep picking the same motive, we should continue with the same plan.
	// - If we pick a different motive, we should re-evaluate the plan, and set up a new one.
	// TODO: Once we have picked a motive to satisfy, we should
	// set up a plan to satisfy it. Maybe we set a "current motive"
	// and if we continue picking the same top motive, we should
	// continue with the same plan. If we pick a different motive,
	// we should re-evaluate the plan, and set up a new one.
	if len(p.Motives) == 0 {
		return
	}

	// Check if we have unsatisfied motives.
	var remMotives []*Motive
	for _, motive := range p.Motives {
		// If the motive is no longer satisfied,
		// satisfaction decays over time.
		if !motive.Type.IsSatisfied(p, m) {
			motive.Tick(1.0)
			remMotives = append(remMotives, motive)
		}
	}

	// No unsatisfied motives, we are done.
	if len(remMotives) == 0 {
		// If the current motive is satisfied, we can remove it.
		p.CurrentMotive = nil

		// Set position to home (if any).
		// This is a hack for now. We can be anywhere, since we can have multi-day
		// tasks etc.
		if p.Home != nil {
			p.X = float64(p.Home.X)
			p.Y = float64(p.Home.Y)
		} else {
			// Set position to the root building.
			p.X = float64(m.Root.X)
			p.Y = float64(m.Root.Y)
		}
		return
	}

	// Sort the motives by increasing value.
	sort.Slice(remMotives, func(i, j int) bool {
		return remMotives[i].Val*remMotives[i].Multiplier() < remMotives[j].Val*remMotives[j].Multiplier()
	})

	// Log all values
	for _, m := range remMotives {
		m.Log()
	}

	// If the top motive is the same as the current motive,
	// we continue with the same plan.
	var chosen *Motive
	if p.CurrentMotive == remMotives[0] {
		// TODO: Continue with the same plan.
		chosen = p.CurrentMotive
		log.Printf("Continue with the same plan: %v (%v)", chosen, p)
		if p.CurrentTree == nil {
			log.Printf("ERROR: Current tree is nil: %v", p)
			p.CurrentTree = chosen.Type.GetTree(p, m)
			log.Printf("ERROR: Current tree is STILL nil: %v", p)
		}
	} else {
		// Pick a new motive from the top 3 motives.
		if rand.Intn(100) < 10 {
			limitRand := min(3, len(p.Motives))
			// Limit to the first 3 (or less) motives.
			chosen = p.Motives[rand.Intn(limitRand)]
		} else {
			chosen = remMotives[0]
		}
		p.CurrentMotive = chosen
		log.Printf("Pick a new motive: %v (%v)", chosen, p)
		p.CurrentTree = chosen.Type.GetTree(p, m)

		// Set position to home (if any).
		// This is a hack for now. We can be anywhere, since we can have multi-day
		// tasks etc.
		if p.Home != nil {
			p.X = float64(p.Home.X)
			p.Y = float64(p.Home.Y)
		} else {
			// Set position to the root building.
			p.X = float64(m.Root.X)
			p.Y = float64(m.Root.Y)
		}
	}

	// Satisfy this motive through the current tree.
	if p.CurrentTree != nil {
		// If the tree is done, we are done... satisfy the motive.
		if p.CurrentTree.Step(elapsed) {
			chosen.Change(200)
		}
	} else {
		// HACK: If we have nothing to do, we just satisfy the motive.
		chosen.Type.Satisfy(p, m)
		chosen.Change(200)
	}
}

// TODO: We need to define the base motives.
// - Health / Survival
// - Hunger / Thirst
// - Sleep
// - Aggression
// - Social
const (
	BaseMotiveHealth = iota // Life or death
	// BaseMotiveThirst
	BaseMotiveHunger
	BaseMotiveSleep // Maybe replace with "Energy"?
	// BaseMotiveFun
	// BaseMotiveHygiene
	// BaseMotiveBladder
	// BaseMotiveEnergy
	// BaseMotiveAggression
	BaseMotiveMax
)

const (
	CategoryBasic     = "Basic"     // Basic motives like hunger, sleep, etc.
	CategoryLife      = "Life"      // Life goals like getting a job, getting married, etc.
	CategoryEphemeral = "Ephemeral" // Ephemeral motives like "I want to eat this apple"
)

// Some motives should be based on the environmental factors, such as a
// need to provide shelter for the family.

// MotiveType is a motive for a person.
// TODO: The motive type should also define if it can be interrupted or not,
// and how long plans are valid before they need to be re-evaluated.
type MotiveType struct {
	Name        string
	Goal        Goal                          // The goal of the motive
	Curve       CurveType                     // How the muliplier changes based on the current value
	Decay       float64                       // How much the value decays per second
	OnMax       func()                        // Called when the motive reaches the maximum value
	OnMin       func()                        // Called when the motive reaches the minimum value
	IsSatisfied func(p *Person, m *Map) bool  // Returns true if the motive is satisfied
	Satisfy     func(p *Person, m *Map) bool  // Called when the motive is satisfied
	GetTree     func(p *Person, m *Map) *Tree // Returns the tree to use for this motive
}

// New creates a new motive.
func (m *MotiveType) New() *Motive {
	return &Motive{
		Type: m,
		Val:  motiveValueStart,
	}
}

const (
	motiveValueStart = 100.0
	motiveValueMin   = -100.0
	motiveValueMax   = motiveValueStart
)

// Motive is a motive instance for a person.
// TODO: We should store the current tree (plan) for the motive,
// so we can continue where we left off if a motive is interrupted.
type Motive struct {
	Type *MotiveType
	Val  float64 // Current value of the motive
}

// String returns a string representation of the motive.
func (m *Motive) String() string {
	return m.Type.Name + ": " + fmt.Sprintf("%.2f", m.Val)
}

// Tick decays the motive value.
func (m *Motive) Tick(elapsed float64) {
	m.Change(-m.Type.Decay * elapsed)
	if m.Val >= motiveValueMax {
		m.Type.OnMax()
	} else if m.Val <= motiveValueMin {
		m.Type.OnMin()
	}
}

// Change changes the motive value by the given amount.
func (m *Motive) Change(amount float64) {
	m.Val += amount
	if m.Val > motiveValueMax {
		m.Val = motiveValueMax
	} else if m.Val < motiveValueMin {
		m.Val = motiveValueMin
	}
}

// MissingToMax returns how much the motive is missing to reach the maximum value.
func (m *Motive) MissingToMax() float64 {
	return motiveValueMax - m.Val
}

// Multiplier returns the current multiplier for the motive.
func (m *Motive) Multiplier() float64 {
	return m.Type.Curve.Multiplier(m.Val)
}

// Log logs the current value (and multiplier) of the motive.
func (m *Motive) Log() {
	log.Printf("%s: %.2f (%.2f)", m.Type.Name, m.Val, m.Multiplier())
}
