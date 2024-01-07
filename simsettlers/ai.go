package simsettlers

import (
	"fmt"
)

// ATTENTION: Everything in this file is just a draft and unused.

// The AI:
// - Perception
// - Conditions
// - Status
// - Planning
// - Execution
//
// Perception:
// - What is the current state of the world that we can perceive?
//
// Status:
// - What is our current state?
// Note to self: Do we have some fixed states that take priority over
// other states? For example, if we are starving, we should eat, even
// if we are tired. Or will we generalize all states?
//
// Planning:
// - What do we want to do? (What is the most important thing to do?)
//
// Execution:
// - How do we do it? (What is the best way to do it?)
//
// For the first iteration, we just fake perception. Sta

type AI struct {
	*Person

	// Motives contains all base motives of the person
	// and potential life goals.
	Motives []*Motive2

	// CurrentMotive is the motive that is currently
	// being satisfied.
	CurrentMotive *Motive2

	// All basic motives, wich can be directly addressed
	// to increase hunger, sleep, etc.
	BaseMotives [BaseMotiveMax]*Motive2

	Plan *ActionPlan
}

func NewAI(p *Person) *AI {
	ai := &AI{
		Person: p,
	}

	ai.BaseMotives[BaseMotiveHealth] = MotiveType2Health.New()
	ai.BaseMotives[BaseMotiveHunger] = MotiveType2Hunger.New()
	ai.BaseMotives[BaseMotiveSleep] = MotiveType2Sleep.New()

	ai.Motives = append(ai.Motives, ai.BaseMotives[:]...)

	return ai
}

// Tick ticks the AI.
func (ai *AI) Tick(elapsed float64) {
	// TODO: Update perception and change any related motive values.
	// For example, if we are threatened, we could increase the
	// aggression motive or the flight motive.
	var highestPriority float64
	var highestPriorityMotive *Motive2
	for _, m := range ai.Motives {
		m.Tick(elapsed)
		if m.Val > highestPriority {
			highestPriority = m.Val
			highestPriorityMotive = m
		}
	}

	// TODO: Instead, we should randomly select a motive
	// from the top 3 motives.
	if ai.CurrentMotive != highestPriorityMotive {
		ai.CurrentMotive = highestPriorityMotive

		// TODO: Determine the best action to satisfy the motive.

		// The Sims has a similar system, where the environment (furniture)
		// is broadcasting the available actions. For example, a fridge
		// broadcasts the "Eat" action, which is then picked up by the person,
		// if the person is in range of the fridge and hungry.

		// But what do we do if we are hungry and don't have food, or can
		// reach a fridge? We could rob someone. Would that mean that people
		// would broadcast the "Rob" action? Or would items that don't
		// belong to us broadcast a "steal" action? Stealing would not
		// directly satiate the hunger, but if we steal food, we can then
		// eat it after.
		// How do we get to that point? We could have items broadcast a
		// eat action, but items that do not belong to us would have
		// a precondition that we have to steal them first.
		// Food -> Eat -> Precond: owned by us
		// If not owned by us: Buy or steal

		// For this, we loop through all available actions and determine
		// which is the cheapest and the most effective.

		// Here are the options wrt. action selection:
		// - Actions are broadcasted by the environment and picked up
		//   by the person.
		// - Actions are

		// For example, if we are hungry, we could eat an apple from
		// our inventory, or we could go to the market and buy some
		// food. The first option is quicker, but depends on us having
		// an apple in our inventory. The second option is slower, but
		// we can always buy food at the market (if we have money).

		// In theory, we could also reward opportunistic behaviour, if it
		// would satisfy another high priority motive. For example, if we
		// we are quite sleepy, but also hungry, we might just grab a snack
		// if we are close to the fridge before going to bed.

		// For now, an action will represent all required steps
		// to satisfy the motive.
		// For example, sleep could be:
		// - Determine the best place to sleep
		// - Go to the determined place
		// - Sleep
		// - Wake up when the motive is satisfied (or when the
		//   person is woken up by something else)

		// Invalidate the current plan since we have a new motive.
		ai.Plan = nil
	} else {
		// Check if our current action is still valid.
		// TODO: determine preconditions for plans.
		ai.Plan = nil
	}
	var objects []*Object

	if ai.Plan == nil && ai.CurrentMotive != nil {
		// Loop through all the objects and find the object that helps
		// us satisfy the motive the quickest.
		// Loop through all actions of the object and determine
		// the best action to satisfy the motive.

		// TODO: For scoring, also take distance into account.
		var plan *ActionPlan
		for _, o := range objects {
			for _, a := range o.Actions {
				if a.Motive == ai.CurrentMotive.MotiveType2 {
					if plan == nil {
						plan = &ActionPlan{
							Action: a,
							Object: o,
						}
					} else {
						if a.Amount > plan.Action.Amount {
							plan.Action = a
							plan.Object = o
						}
					}
				}
			}
		}
		ai.Plan = plan
	}

	if ai.Plan != nil {
		// Execute the plan.
		pos := ai.Position()
		if pos.DistanceTo(ai.Plan.Object.Position()) > 1 {
			// Set direction towards the object.
			objectPos := ai.Plan.Object.Position()
			ai.SetDirection(objectPos.X, objectPos.Y)

			// Move towards the object.
			ai.Move(1.0) // TODO: Move should take elapsed time into account.
			return
		}

		// Perform the action.
		// Update the current motive.
		ai.ChangeMotive(ai.Plan.Action.Motive, ai.Plan.Action.Amount)
	}
}

// ChangeBaseMotive changes the given base motive by the given amount.
func (ai *AI) ChangeBaseMotive(motive int, amount float64) {
	ai.BaseMotives[motive].Change(amount)
}

// ChangeMotive changes the given motive by the given amount.
func (ai *AI) ChangeMotive(motive *MotiveType2, amount float64) {
	for _, m := range ai.Motives {
		if m.MotiveType2 == motive {
			m.Change(amount)
			return
		}
	}
}

type ActionPlan struct {
	Action *Action // The action to perform
	Object *Object // The object to perform the action on
}

type Action struct {
	Name   string
	Motive *MotiveType2
	Amount float64
}

var ActionEat = &Action{
	Name:   "Eat",
	Motive: MotiveType2Hunger,
	Amount: -10,
}

type MotiveType2 struct {
	Name     string
	Category string
	Decay    float64 // Base decay per second
	// Curve    CurveType
}

func (m *MotiveType2) New() *Motive2 {
	return &Motive2{
		MotiveType2: m,
		Val:         motiveValueStart,
	}
}

type Motive2 struct {
	*MotiveType2
	Val float64
}

func (m *Motive2) String() string {
	return m.Name + ": " + fmt.Sprintf("%.2f", m.Val)
}

func (m *Motive2) Tick(elapsed float64) {
	m.Change(-m.Decay * elapsed)
}

func (m *Motive2) Change(amount float64) {
	m.Val += amount

	// Clamp the value.
	if m.Val < motiveValueMin {
		m.Val = motiveValueMin
	} else if m.Val > motiveValueMax {
		m.Val = motiveValueMax
	}
}

var MotiveType2Health = &MotiveType2{
	Name:     "Health",
	Category: CategoryBasic,
	Decay:    0.1,
}

var MotiveType2Hunger = &MotiveType2{
	Name:     "Hunger",
	Category: CategoryBasic,
	Decay:    0.1,
}

var MotiveType2Sleep = &MotiveType2{
	Name:     "Sleep",
	Category: CategoryBasic,
	Decay:    0.1,
}
