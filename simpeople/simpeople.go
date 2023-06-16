package simpeople

import (
	"math"
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
	Name        string
	Personality *aifiver.Personality
}

func NewPerson(name string, traiter *aifiver.Traiter) *Person {
	return &Person{
		Name:        name,
		Personality: aifiver.NewPersonalityRandomized(traiter),
	}
}

// Log logs the persons name and personality.
func (p *Person) Log() {
	println(p.Name)
	p.Personality.Log()
}

// Compatibility returns the compatibility between two people.
func (p *Person) Compatibility(b *Person) int {
	// TODO: This should be different from a to b and b to a.
	// ... because some personalities consider certain traits more attractive or repulsive
	// than others. So disagreeable people might consider agreeable people less attractive
	// than the other way around, just because they are disagreeable.
	return aifiver.Compatibility(p.Personality, b.Personality)
}

// PickAction returns an action that the person wants to do.
func (p *Person) PickAction(b *Person) *Action {
	// Determine how much we like the other person.
	compat := p.Compatibility(b)

	// The sum of compatibility and the relevant facets determine if we meet or exceed
	// the threshold for the action. We pick the action with the highest score.
	// Sort the actions by cost and pick the first one that meets the threshold.
	sort.Slice(Actions, func(i, j int) bool {
		return math.Abs(float64(Actions[i].Cost)) < math.Abs(float64(Actions[j].Cost))
	})
	for _, action := range Actions {
		// Determine the cost of the action.
		cost := action.Cost

		// Sum the cost and the compatibility.
		eval := action.Eval(p, b) + compat

		// If the cost is negative, the eval must be lower than (or eqal to) the cost.
		if cost < 0 {
			if eval <= cost {
				return action
			}
		} else if eval >= cost {
			return action
		}
	}
	return nil
}

// Action represents an action that a person can take.
type Action struct {
	Name string                     // Name of the action
	Cost int                        // Positive or negative (negative harms the other person)
	Eval func(*Person, *Person) int // Evaluation function of the action
}

var (
	ActionHug = &Action{
		Name: "Hug",
		Cost: 4,
		Eval: func(a, b *Person) int {
			// TODO: Implement a good evaluation function.
			// For now, let's say that the value is the average of the agreeableness of the two people
			// and the average of the openness of the two people.
			avgAgree := (a.Personality.Get(aifiver.FactorAgreeableness) + b.Personality.Get(aifiver.FactorAgreeableness)) / 2
			avgOpen := (a.Personality.Get(aifiver.FactorOpenness) + b.Personality.Get(aifiver.FactorOpenness)) / 2
			return int(avgAgree+avgOpen) / 2
		},
	}
	ActionKiss = &Action{
		Name: "Kiss",
		Cost: 6,
		Eval: func(a, b *Person) int {
			// Same as above but with impulsivity.
			avgAgree := (a.Personality.Get(aifiver.FactorAgreeableness) + b.Personality.Get(aifiver.FactorAgreeableness)) / 2
			avgOpen := (a.Personality.Get(aifiver.FactorOpenness) + b.Personality.Get(aifiver.FactorOpenness)) / 2
			avgExc := a.Personality.GetFacet(aifiver.FacetExtrExcitementSeeking)
			avgImp := a.Personality.GetFacet(aifiver.FacetNeurImpulsiveness)
			return int(avgAgree+avgOpen+avgExc+avgImp) / 4
		},
	}
	ActionSteal = &Action{
		Name: "Steal",
		Cost: -6,
		Eval: func(a, b *Person) int {
			// Low altruism, lack of self discipline and high impulsivity make people more likely to steal.
			alt := a.Personality.GetFacet(aifiver.FacetAgreAltruism)
			disc := a.Personality.GetFacet(aifiver.FacetConsSelfDicipline)
			imp := a.Personality.GetFacet(aifiver.FacetNeurImpulsiveness) * -1 // We invert this facet because it is negative.
			return int(alt+disc+imp) / 3
		},
	}
	ActionInsult = &Action{
		Name: "Insult",
		Cost: -4,
		Eval: func(a, b *Person) int {
			// Low positive emotions, lack of self discipline and high hostility make people more likely to insult others.
			pos := a.Personality.GetFacet(aifiver.FacetExtrPositiveEmotions)
			disc := a.Personality.GetFacet(aifiver.FacetConsSelfDicipline)
			imp := a.Personality.GetFacet(aifiver.FacetNeurAngryHostility) * -1 // We invert this facet because it is negative.
			return int(pos+disc+imp) / 3
		},
	}
)

var Actions = []*Action{
	ActionHug,
	ActionKiss,
	ActionSteal,
	ActionInsult,
}
