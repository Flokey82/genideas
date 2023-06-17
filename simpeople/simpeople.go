package simpeople

import (
	"math"
	"math/rand"
	"sort"

	"github.com/Flokey82/aifiver"
	"github.com/Flokey82/go_gens/genlanguage"
)

type PeopleFactory struct {
	lang    *genlanguage.Language
	traiter *aifiver.Traiter
}

func NewPeopleFactory() *PeopleFactory {
	lang := genlanguage.GenLanguage(12345)
	traiter := aifiver.NewTraiter()
	aifiver.DefaultTraits(traiter)
	return &PeopleFactory{
		lang:    lang,
		traiter: traiter,
	}
}

func (pf *PeopleFactory) NewPerson() *Person {
	// Generate a new name.
	name := pf.lang.MakeName()
	return NewPerson(name, pf.traiter)
}

type Person struct {
	Name           string
	Personality    *aifiver.Personality
	Insight        float64             // How well we can read the other person.
	FamiliarityVal map[*Person]float64 // How well we know the other person.
	ExperienceVal  map[*Person]float64 // Our personal experience with the other person.
}

func NewPerson(name string, traiter *aifiver.Traiter) *Person {
	return &Person{
		Name:           name,
		Personality:    aifiver.NewPersonalityRandomized(traiter),
		Insight:        rand.Float64(), // TODO: Make this dependent on the personality.
		FamiliarityVal: make(map[*Person]float64),
		ExperienceVal:  make(map[*Person]float64),
	}
}

// Log logs the persons name and personality.
func (p *Person) Log() {
	println(p.Name)
	p.Personality.Log()
}

// Compatibility returns the compatibility between two people.
func (p *Person) Compatibility(b *Person) int {
	// Sum up the compatibility of the two personalities.
	return p.Interaction(b) + p.Opinion(b) + p.Experience(b)
}

// Interaction returns a value for how well the personalities interact face to face.
func (p *Person) Interaction(b *Person) int {
	return int(float64(p.Personality.Interaction(b.Personality)) * p.FamiliarityVal[b])
}

// Opinion returns the opinion of the other person, depending on how well we know them,
// just based on their reputation.
func (p *Person) Opinion(b *Person) int {
	return int(float64(p.Personality.Opinion(b.Personality)) * p.FamiliarityVal[b])
}

// Experience returns the experience of the other person, depending on how well we know them.
func (p *Person) Experience(b *Person) int {
	return int(b.ExperienceVal[p] * 10)
}

// PickAction returns an action that the person wants to do.
func (p *Person) PickAction(b *Person) *Action {
	// Determine how well we can interact with the other person.
	compat := p.Compatibility(b)

	// The sum of compatibility and the relevant facets determine if we meet or exceed
	// the threshold for the action. We pick the action with the highest score.
	// Sort the actions by cost and pick the first one that meets the threshold.
	sort.Slice(Actions, func(i, j int) bool {
		return math.Abs(float64(Actions[i].Cost)) < math.Abs(float64(Actions[j].Cost))
	})
	var bestAction *Action
	for _, action := range Actions {
		// Determine the cost of the action.
		cost := action.Cost

		// Sum the cost and the compatibility.
		eval := action.Eval(p, b) + compat

		// If the cost is negative, the eval must be lower than (or eqal to) the cost.
		if cost < 0 {
			if eval <= cost {
				bestAction = action
				break
			}
		} else if eval >= cost {
			bestAction = action
			break
		}
	}

	// TODO: Only do this if we actually complete the action.
	// Maybe also have different values for success and failure.
	if bestAction != nil {
		// TODO: Dedupe and improve this code.

		// We have found an action that we want to do.
		// Increase the familiarity with the other person.
		p.FamiliarityVal[b] += bestAction.FamiliarityChange * p.Insight
		if p.FamiliarityVal[b] > 1 {
			p.FamiliarityVal[b] = 1
		}

		// Change the opinion.
		// Change the opinion of the actor.
		p.ExperienceVal[p] += (bestAction.OpinionChangeActor + bestAction.OpinionChange) * p.Insight
		if p.ExperienceVal[p] > 1 {
			p.ExperienceVal[p] = 1
		} else if p.ExperienceVal[p] < -1 {
			p.ExperienceVal[p] = -1
		}

		// Change the opinion of the target.
		b.ExperienceVal[p] += (bestAction.OpinionChangeTarget + bestAction.OpinionChange) * b.Insight
		if b.ExperienceVal[p] > 1 {
			b.ExperienceVal[p] = 1
		} else if b.ExperienceVal[p] < -1 {
			b.ExperienceVal[p] = -1
		}
		return bestAction
	}
	return nil
}

// Action represents an action that a person can take.
type Action struct {
	Name              string                     // Name of the action
	Cost              int                        // Positive or negative (negative harms the other person)
	Eval              func(*Person, *Person) int // Evaluation function of the action
	FamiliarityChange float64                    // How much the familiarity changes after the action.
	// TODO: Distinct opinion change on success vs failure and actor vs target?
	OpinionChange       float64 // How much the opinion changes after the action.
	OpinionChangeActor  float64 // How much the opinion changes after the action.
	OpinionChangeTarget float64 // How much the opinion changes after the action.
}

var (
	ActionTalk = &Action{
		Name: "talk to",
		Cost: 2,
		Eval: func(a, b *Person) int {
			return a.Personality.GetFacet(aifiver.FacetExtrGregariousness)
		},
		FamiliarityChange: 0.1,
		OpinionChange:     0.1,
	}
	ActionHug = &Action{
		Name: "hug",
		Cost: 4,
		Eval: func(a, b *Person) int {
			// TODO: Implement a good evaluation function.
			// For now, let's say that the value is the average of the agreeableness of the two people
			// and the average of the openness of the two people.
			avgAgree := (a.Personality.Get(aifiver.FactorAgreeableness) + b.Personality.Get(aifiver.FactorAgreeableness)) / 2
			avgOpen := (a.Personality.Get(aifiver.FactorOpenness) + b.Personality.Get(aifiver.FactorOpenness)) / 2
			return int(avgAgree+avgOpen) / 2
		},
		FamiliarityChange: 0.2,
		OpinionChange:     0.2,
	}
	ActionKiss = &Action{
		Name: "smooch with",
		Cost: 6,
		Eval: func(a, b *Person) int {
			// Same as above but with impulsivity.
			avgAgree := (a.Personality.Get(aifiver.FactorAgreeableness) + b.Personality.Get(aifiver.FactorAgreeableness)) / 2
			avgOpen := (a.Personality.Get(aifiver.FactorOpenness) + b.Personality.Get(aifiver.FactorOpenness)) / 2
			avgExc := a.Personality.GetFacet(aifiver.FacetExtrExcitementSeeking)
			avgImp := a.Personality.GetFacet(aifiver.FacetNeurImpulsiveness)
			return int(avgAgree+avgOpen+avgExc+avgImp) / 4
		},
		FamiliarityChange: 0.3,
		OpinionChange:     0.3,
	}
	ActionPunch = &Action{
		Name: "punch",
		Cost: -4,
		Eval: func(a, b *Person) int {
			ang := a.Personality.GetFacet(aifiver.FacetNeurAngryHostility)
			return int(ang) * -1 // We invert this facet because it is negative.
		},
		FamiliarityChange:   0.2,
		OpinionChangeTarget: -0.4,
	}
	ActionSteal = &Action{
		Name: "steal from",
		Cost: -6,
		Eval: func(a, b *Person) int {
			// Low altruism, lack of self discipline and high impulsivity make people more likely to steal.
			alt := a.Personality.GetFacet(aifiver.FacetAgreAltruism)
			disc := a.Personality.GetFacet(aifiver.FacetConsSelfDicipline)
			imp := a.Personality.GetFacet(aifiver.FacetNeurImpulsiveness) * -1 // We invert this facet because it is negative.
			return int(alt+disc+imp) / 3
		},
		FamiliarityChange:   0.1,
		OpinionChangeTarget: -0.3,
	}
	ActionInsult = &Action{
		Name: "insult",
		Cost: -4,
		Eval: func(a, b *Person) int {
			// Low positive emotions, lack of self discipline and high hostility make people more likely to insult others.
			pos := a.Personality.GetFacet(aifiver.FacetExtrPositiveEmotions)
			disc := a.Personality.GetFacet(aifiver.FacetConsSelfDicipline)
			imp := a.Personality.GetFacet(aifiver.FacetNeurAngryHostility) * -1 // We invert this facet because it is negative.
			return int(pos+disc+imp) / 3
		},
		FamiliarityChange:   0.2,
		OpinionChangeTarget: -0.2,
	}
)

var Actions = []*Action{
	ActionHug,
	ActionKiss,
	ActionPunch,
	ActionSteal,
	ActionInsult,
	ActionTalk,
}

// MovingAvg is a moving average.
type MovingAvg struct {
	Value float64
	Count int
}

// Add adds a value to the moving average.
func (m *MovingAvg) Add(v float64) {
	m.Value = (m.Value*float64(m.Count) + v) / float64(m.Count+1)
	m.Count++
}
