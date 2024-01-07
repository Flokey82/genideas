package simsettlers

import (
	"log"
	"math"
	"math/rand"
	"sort"
)

var MotiveTypeAdventure = &MotiveType{
	Name:  "Adventure",
	Goal:  GoalAdultAdventurer,
	Curve: CurveTypeExponential,
	Decay: 5.0,
	OnMax: func() {},
	OnMin: func() {
		log.Println("Person is very bored!")
	},
	IsSatisfied: func(p *Person, m *Map) bool {
		return false
	},
	Satisfy: func(p *Person, m *Map) bool {
		handleAdventure(p, m)
		return true
	},
	GetTree: func(p *Person, m *Map) *Tree { // Pick a random dungeon (TODO: Handle case if there are no dungeons
		dcTree := newDungeonCrawl(p, m.Dungeons[rand.Intn(len(m.Dungeons))], func() {
			// We survived the adventure.
			loot := rand.Intn(1000)
			p.Resources += loot
			log.Printf("%s went on an adventure and found %d resources", p.String(), loot)
		}, func() {
			// TODO: What if we just go missing? Could someone save us?
			// We died on our adventure.
			// We might be declared missing, someone might find our body,
			// or we might never be found.
			// TODO: Inheritence, ownership, legally dead, etc.
			m.handleDeath(p)
		})
		return &dcTree
	},
}

var MotiveTypeBully = &MotiveType{
	Name:  "Bully",
	Goal:  GoalChildhoodBully,
	Curve: CurveTypeExponential,
	Decay: 0.5,
	OnMax: func() {},
	OnMin: func() {
		log.Println("Person is very bored!")
	},
	IsSatisfied: func(p *Person, m *Map) bool {
		// Find a random child to bully.
		var victims []*Person
		for _, c := range m.RealPop {
			// Not dead and similar age.
			if c.Dead || c == p {
				continue
			}
			if math.Abs(float64(c.Age-p.Age)) < 2 {
				victims = append(victims, c)
			}
		}
		// HACK: We are satisfied if we have no victims.
		if len(victims) == 0 {
			return true
		}
		return false
	},
	Satisfy: func(p *Person, m *Map) bool {
		if rand.Intn(100) < 1 && p.Age > 5 {
			handleBully(p, m)
		}
		return true
	},
	GetTree: func(p *Person, m *Map) *Tree {
		// Find a random child to bully.
		var victims []*Person
		for _, c := range m.RealPop {
			// Not dead and similar age.
			if c.Dead || c == p {
				continue
			}
			if math.Abs(float64(c.Age-p.Age)) < 2 {
				victims = append(victims, c)
			}
		}

		// Sort by opinion.
		sort.Slice(victims, func(i, j int) bool {
			return p.Opinions.Value(victims[i]) < p.Opinions.Value(victims[j])
		})

		// Pick the victim with the lowest opinion by the bully.
		// TODO: Usually we'd pick the weakest and most vulnerable victim.
		// Also, we'd prefer picking on the same victim over and over again.
		// Wow, I hate this so much.
		if len(victims) > 0 {
			victim := victims[0]
			// TODO: Occasionally, we might mix it up and bully someone else.
			if len(victims) > 1 && rand.Intn(100) < 10 {
				victim = victims[min(rand.Intn(3), len(victims)-1)]
			}
			bt := newBullyTree(p, victim, func() {
				// We bullied the victim.
			}, func() {
				// Bullying failed.
				if rand.Intn(100) < 10 {
					// TODO: This is a crime committed by the victim,
					// we need to find a way to handle this.
					log.Printf("%v killed %v", victim, p)
					m.handleDeath(p)
				} else {
					// The bully will be forced to respect the victim a little more.
					log.Printf("%v injured %v", victim, p)
					p.Opinions.IncrementBy(victim, 10)
				}
			})
			return &bt
		}
		return nil
	},
}

var MotiveTypeBuildHouse = &MotiveType{
	Name:  "BuildHouse",
	Goal:  GoalAdultHome,
	Curve: CurveTypeExponential,
	Decay: 5.0,
	OnMax: func() {},
	OnMin: func() {
		log.Println("Person is very bored!")
	},
	IsSatisfied: func(p *Person, m *Map) bool {
		return p.OwnsOwnHome()
	},
	Satisfy: func(p *Person, m *Map) bool {
		handleHome(p, m)
		return true
	},
	GetTree: func(p *Person, m *Map) *Tree {
		// TODO: implement
		handleHome(p, m)
		return nil
	},
}

var MotiveTypeSocialize = &MotiveType{
	Name:  "Socialize",
	Goal:  GoalChildhoodSocialize,
	Curve: CurveTypeExponential,
	Decay: 0.5,
	OnMax: func() {},
	OnMin: func() {
		log.Println("Person is very bored!")
	},
	IsSatisfied: func(p *Person, m *Map) bool {
		var friends []*Person
		for _, c := range m.RealPop {
			// Not dead and similar age.
			if c.Dead || c == p {
				continue
			}
			if math.Abs(float64(c.Age-p.Age)) < 3 {
				friends = append(friends, c)
			}
		}
		// HACK: We are satisfied if we have no friends.
		if len(friends) == 0 {
			return true
		}
		return false
	},
	Satisfy: func(p *Person, m *Map) bool {
		if rand.Intn(100) < 1 && p.Age > 5 {
			handleSocialize(p, m)
		}
		return true
	},
	GetTree: func(p *Person, m *Map) *Tree {
		// TODO: Factor this out into a function... This can be generalized as adults also socialize with other adults.
		// Find a random individual of similar age to socialize with.
		var friends []*Person
		for _, c := range m.RealPop {
			// Not dead and similar age.
			if c.Dead || c == p {
				continue
			}
			if math.Abs(float64(c.Age-p.Age)) < 3 {
				friends = append(friends, c)
			}
		}

		// Pick the friend with the highest opinion by the person.
		if len(friends) > 0 {
			// Sort by opinion.
			sort.Slice(friends, func(i, j int) bool {
				return p.Opinions.Value(friends[i]) > p.Opinions.Value(friends[j])
			})

			// Occasionally, we might mix it up and socialize with someone else from
			// the top 3 or so.
			friend := friends[0]
			if len(friends) > 1 && rand.Intn(100) < 10 {
				friend = friends[min(rand.Intn(3), len(friends)-1)]
			}
			st := newSocializeTree(p, friend, func() {
				// We socialized with the friend.
			}, func() {
				// Socializing failed.
			}, m)

			return &st
		}
		return nil
	},
}

var MotiveTypeJob = &MotiveType{
	Name:  "Job",
	Goal:  GoalAdultJob,
	Curve: CurveTypeLinear,
	Decay: 5.0,
	OnMax: func() {},
	OnMin: func() {
		log.Println("Person is very bored!")
	},
	IsSatisfied: func(p *Person, m *Map) bool {
		return false // p.Job != JobTypeUnemployed
	},
	Satisfy: func(p *Person, m *Map) bool {
		handleJob(p, m)
		return true
	},
	GetTree: func(p *Person, m *Map) *Tree {
		handleJob(p, m)
		return nil
	},
}
