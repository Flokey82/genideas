package simsettlers

import (
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/Flokey82/go_gens/genlanguage"
)

func (m *Map) tickPersonGoals(p *Person, elapsed float64) {
	// Check if we are old enough to work or build a house.
	// TODO: We should vary this based on personality.
	// TODO: Check if goals are satisfied or not.
	// Some people might refuse to retire, others might retire early.
	if p.Age < 18 {
		// If we are a bully, we might bully other children.
		// TODO: Factor this out into a function... This can be generalized as adults also bully other adults.
		if p.Goals.IsSet(GoalChildhoodBully) && rand.Intn(100) < 1 && p.Age > 5 {
			handleBully(p, m)
		}

		// If we socialize, let's see if we can make friends.
		if p.Goals.IsSet(GoalChildhoodSocialize) && rand.Intn(100) < 1 && p.Age > 5 {
			handleSocialize(p, m)
		}
	} else if p.Age < 65 {
		// Handle our job goal.
		if p.Goals.IsSet(GoalAdultJob) {
			handleJob(p, m)
		}

		// Check if we even want a home.
		if p.Goals.IsSet(GoalAdultHome) {
			handleHome(p, m)
		}

		// Check if we want to go on an adventure.
		if p.Goals.IsSet(GoalAdultAdventurer) && rand.Intn(100) < 1 {
			// NOTE: Adventures should be memorable, rare occurrences with
			// lasting effects on the person. High risk, high reward.
			// This might get us killed, we migh bear scars, or maybe gain
			// a new fear, a new skill, maybe an artifact, etc.
			handleAdventure(p, m)
		}

		// Check if we want to find a partner.
		if p.Goals.IsSet(GoalAdultPartner) {
			// TODO: Go through the list of people that we know,
			// and pick the person we consider the best match.
			// Depending on the personality, we might have different
			// criteria for what makes a good match. It might be who
			// we like the most, who we hate the least, who is the
			// most attractive / high status, who is the most wealthy,
			// etc.
		}
	} else {
		// m.tickElderly(p)
	}
}

func (p *Person) assignChildhoodGoals() {
	if rand.Intn(100) < 90 {
		p.Goals |= GoalChildhoodSocialize
		p.Motives = append(p.Motives, MotiveTypeSocialize.New())
	}
	if rand.Intn(100) < 5 {
		p.Goals |= GoalChildhoodBully
		p.Motives = append(p.Motives, MotiveTypeBully.New())
	}
}

func (p *Person) assignAdultGoals() {
	// TODO:
	// - Personality matters.
	// - Also parents matter.
	// If both parents have a trait, the chance of the child having it is higher.
	if rand.Intn(100) < 90 {
		p.Goals |= GoalAdultPartner
	}
	if rand.Intn(100) < 90 {
		p.Goals |= GoalAdultHome
		p.Motives = append(p.Motives, MotiveTypeBuildHouse.New())
	}
	if rand.Intn(100) < 90 {
		p.Goals |= GoalAdultChildren
	}
	if rand.Intn(100) < 90 {
		p.Goals |= GoalAdultJob
		p.Motives = append(p.Motives, MotiveTypeJob.New())
	}
	if rand.Intn(100) < 5 {
		p.Goals |= GoalAdultAdventurer
		p.Motives = append(p.Motives, MotiveTypeAdventure.New())
	}
}

func (p *Person) assignElderlyGoals() {
}

type Goal uint32

func (g Goal) IsSet(sg Goal) bool {
	return g&sg == sg
}

func (g *Goal) String() string {
	// TODO: Make it a nice short string, like permissions in Linux.
	var str string
	for i := 0; i < 8; i++ {
		if g.IsSet(1 << i) {
			switch Goal(1 << i) {
			case GoalChildhoodSocialize:
				str += "S"
			case GoalChildhoodBully:
				str += "B"
			case GoalAdultPartner:
				str += "P"
			case GoalAdultHome:
				str += "H"
			case GoalAdultChildren:
				str += "C"
			case GoalAdultJob:
				str += "J"
			case GoalAdultAdventurer:
				str += "A"
			default:
				str += "?"
			}
		} else {
			str += "-"
		}
	}
	return str
}

const (
	// TODO: Childhood goals.
	GoalChildhoodSocialize Goal = 1 << 0
	GoalChildhoodBully     Goal = 1 << 1
	// TODO: Teenage goals.
	// TODO: Adult goals.
	// TODO: Elderly goals.
	GoalAdultPartner    Goal = 1 << 2
	GoalAdultHome       Goal = 1 << 3
	GoalAdultChildren   Goal = 1 << 4
	GoalAdultJob        Goal = 1 << 5
	GoalAdultAdventurer Goal = 1 << 6
)

// newSocializeTree creates a "behavior" tree for socializing with a target.
// NOTE: This is not a real behavior tree, but a simple sequence of tasks.
func newSocializeTree(p *Person, target *Person, onSuccess, onFailure func(), m *Map) Tree {
	// - Move to the target.
	// - Socialize with the target.
	// - Move back home.
	homeX, homeY := int(p.X), int(p.Y)

	// Depending on how much we like the person, we pick a duration to socialize.
	dur := float64(p.Opinions.Value(target))

	// TODO: Avoid socializing with someone who is on an adventure.
	// This means we first select a target, and every time we check the "find target"
	// node, we retrieve the selected target and check if they are still valid.
	// If not, we find a new target.

	// Wouldn't it be better if we get two individuals to agree to socialize with each other?
	// In some form of negotiation. Maybe we should first decide what each person wants to do,
	// e.g. socialize, bully, etc. and then find a suitable target. But who would 'agree' to
	// be bullied? Maybe some motives can interrupt others. But what if we interrupt an activity
	// that is shared by two people? So, maybe we head to a specific location with the intention
	// of socializing, which might be the location of the target, or a location where we can find
	// a target. Then, we check if we can find a target, and if so, we socialize with them.

	// For example, let's say we don't know anyone well enough to socialize with them, we might
	// go to a tavern, and then we might find someone to socialize with. Or we might go to a

	// TODO: For this to work, both people should agree to socialize with each other,
	// otherwise one might be running after the other and the other might be running
	// after someone else...
	/*
		// NOTE: A person can move around, so we need to update the target's position or
		// use a reference instead.
		var tx, ty int
		if target != nil {
			tx, ty = int(target.X), int(target.Y)
		} else {
			tx, ty = int(p.X), int(p.Y)
		}
		tRoot := NewTaskGeneric(p, "FindTarget", func(elapsed float64) TaskStatus {
			if target != nil && !target.Dead {
				return TaskStatusCompleted
			}

			var friends []*Person
			for _, c := range m.RealPop {
				// Not dead and similar age.
				if c.Dead || c == p {
					continue
				}
				if math.Abs(float64(c.Age-p.Age)) < 2 {
					friends = append(friends, c)
				}
			}

			// Pick the friend with the highest opinion by the person.
			if len(friends) == 0 {
				// TODO: Move to the market and look in the sourroundings for someone to socialize with.
				return TaskStatusFailed
			}

			// Sort by opinion.
			sort.Slice(friends, func(i, j int) bool {
				return p.Opinions.Value(friends[i]) > p.Opinions.Value(friends[j])
			})

			// Occasionally, we might mix it up and socialize with someone else from
			// the top 3 or so.
			target = friends[0]
			if len(friends) > 1 && rand.Intn(100) < 10 {
				target = friends[min(rand.Intn(3), len(friends)-1)]
			}
			tx, ty = int(target.X), int(target.Y)
			return TaskStatusCompleted
		})
		t := tRoot.Then(NewTaskMoveToLocation(p, func() vectors.Vec2 {
			if target == nil {
				// move to our own position.
				// TODO: Move to the market and look in the sourroundings for someone to socialize with.
				return p.Position()
			}
			return target.Position()
		}))
	*/
	tRoot := NewTaskMoveToXY(p, int(target.X), int(target.Y))
	t := tRoot.Then(NewTaskGeneric(p, "Socialize", func(elapsed float64) TaskStatus {
		log.Printf("%v is socializing with %v", p, target)

		dur -= elapsed

		if dur <= 0 {
			return TaskStatusCompleted
		}

		// Change the friend's opinion of the person.
		target.Opinions.IncrementBy(p, int(10*elapsed))

		// Change the person's opinion of the friend.
		p.Opinions.IncrementBy(target, int(10*elapsed))
		// TODO: Based on personality, we might become friends, or we might become enemies.
		return TaskStatusInProgress
	}))
	t.Then(NewTaskMoveToXY(p, homeX, homeY)) // Move back home.

	return NewTree(tRoot, onSuccess, onFailure)
}

// handleSocialize executes / fullfilles the socialization goal.
func handleSocialize(p *Person, m *Map) {
	// TODO: Factor this out into a function... This can be generalized as adults also socialize with other adults.
	// Find a random individual of similar age to socialize with.
	var friends []*Person
	for _, c := range m.RealPop {
		// Not dead and similar age.
		if c.Dead || c == p {
			continue
		}
		if math.Abs(float64(c.Age-p.Age)) < 2 {
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

		if !st.Do(1.0) {
			log.Println("Socializing failed")
		}
	}
}

// newBullyTree creates a "behavior" tree for bullying a target.
// NOTE: This is not a real behavior tree, but a simple sequence of tasks.
func newBullyTree(p *Person, target *Person, onSuccess, onFailure func()) Tree {
	// - Move to the target.
	// - Bully the target.
	// - Move back home.
	homeX, homeY := int(p.X), int(p.Y)

	// Depending on how much we like the person, we pick a duration to bully.
	dur := float64(p.Opinions.Value(target))

	// NOTE: A person can move around, so we need to update the target's position or
	// use a reference instead.
	tRoot := NewTaskMoveToXY(p, int(target.X), int(target.Y))
	t := tRoot.Then(NewTaskGeneric(p, "Bully", func(elapsed float64) TaskStatus {
		log.Printf("%v is bullying %v", p, target)

		dur -= elapsed

		if dur <= 0 {
			return TaskStatusCompleted
		}

		// Change the victim's opinion of the bully.
		target.Opinions.IncrementBy(p, int(-20*elapsed))

		// Change the bully's opinion of the victim.
		// TODO: The bullying might backfire, then we should change the bully's opinion of the victim positively
		// to avoid bullying the same person unsuccessfully over and over again.
		p.Opinions.IncrementBy(target, int(-10*elapsed))

		// Depending on the personality, this might result in violence, personality changes, etc.
		if target.Goals.IsSet(GoalChildhoodSocialize) && rand.Intn(100) < 5 {
			// If the victim is social, they might loose the will to socialize.
			target.Goals &= ^GoalChildhoodSocialize
		}

		// If the victim is bullied too much, they might become a bully themselves.
		if !target.Goals.IsSet(GoalChildhoodBully) && rand.Intn(100) < 1 {
			target.Goals |= GoalChildhoodBully
		}

		// There might be a chance that the victim becomes violent and kills or injures the bully.
		if rand.Intn(100) < 1 {
			return TaskStatusFailed
		}
		return TaskStatusInProgress
	}))
	t.Then(NewTaskMoveToXY(p, homeX, homeY)) // Move back home.

	return NewTree(tRoot, onSuccess, onFailure)
}

// handleBully executes / fullfilles the bullying goal.
func handleBully(p *Person, m *Map) {
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

		if !bt.Do(1.0) {
			log.Println("Bullying failed")
		}
	}
}

// handleHome executes / fullfilles the home goal.
func handleHome(p *Person, m *Map) {
	repairHome := false

	// Check if we own our own home.
	if p.OwnsOwnHome() {
		// Check if we need to repair our home.
		// TODO:
		// - Make the decision to repair or not based on personality
		// and the condition of the building.
		// - Any occupant should be able to repair the home if they are
		// old enough, have enough resources, and are crafty enough.
		if repairHome && p.Home.Condition < 100 && p.Resources > 0 {

			// TREE:
			// - Move to house location.
			// - Repair house at a certain rate.
			// - Move back home.
			p.Home.Condition++
			p.Resources--
		}
		return
	}

	// Check if we are already building a house or if our spouse is.
	if p.Constructing != nil || p.Spouse != nil && p.Spouse.Constructing != nil {
		// We are building a house.

		// TREE:
		// - Move to house location.
		// - Build house at a certain rate.
		// - Move back home.
		return
	}

	// Check if we live with our parents and if we want and have a spouse.
	// If we don't want a spouse, we can move out on our own. If we plan
	// to get married, we should wait until we have a spouse.
	// TODO: Also factor in the wish to have children.
	if p.LivesWithParents() && p.Spouse == nil && p.Goals.IsSet(GoalAdultPartner) {
		// Wait until we have a spouse (if we want one).
		return
	}

	// Check if our spouse owns their own home, and if so, move in.
	if p.Spouse != nil && p.Spouse.OwnsOwnHome() {
		p.SetHome(p.Spouse.Home)
		p.Spouse.Home.AddOwner(p)
		return
	}

	// We want to move to a new place if:
	// - we don't have a home.
	// - we live with our parents, have a spouse, and/or can afford to build a house.
	budget := p.Resources
	if p.Spouse != nil {
		budget += p.Spouse.Resources
	}

	// TODO: If there is a house in bad condition available,
	// we should be able to buy it for a lower price!

	// Check if there is an existing house that we can move into.
	// Either, because it is abandoned, or because there are no occupants.
	// (Inherited from parents, or because the owners died.)
	for _, b := range m.Buildings {
		if b.Type != BuildingTypeHouse || b.IsOccupied() {
			continue
		}
		if b.IsOwned() {
			// Buy it from the owner.
			// Calculate the purchase price based on the condition of the building.
			purchasePrice := b.PurchasePrice()
			if budget < purchasePrice {
				//log.Printf("Not enough resources to buy a house")
				continue
			}

			// Split the cost between the owners.
			cost := purchasePrice / len(b.Owners)
			for _, o := range b.Owners {
				o.Resources += cost
			}
			b.Occupants = nil
			b.Owners = nil

			// Deduct the cost from the resources of the buyer.
			p.Resources -= cost

			// Move in.
			p.SetHome(b)
			b.AddOwner(p)

			// If there is a spouse, add them as an owner and move them in.
			// TODO: If we are pooling resources with our spouse, we should split the cost.
			if p.Spouse != nil {
				b.AddOwner(p.Spouse)
				p.Spouse.SetHome(b)
			}
		} else {
			// No owners, so move in.
			// TODO: Repair costs etc.
			p.SetHome(b)
			b.AddOwner(p)
		}
		return
	}

	// Can we afford to build a new house?
	if budget < buildingCosts[BuildingTypeHouse] {
		//log.Printf("Not enough resources to build a house")
		return
	}

	// TODO: Find a suitable location for the house.
	// If we are social, we want to live closer to other people.
	best, score := m.getHighestHouseFitness(p.Goals.IsSet(GoalChildhoodSocialize))
	if score != -1 {
		log.Printf("Building a house for %v", p)

		// Construct the house.
		home := m.AddBuilding(best%m.Width, best/m.Width, BuildingTypeHouse)
		home.AddOwner(p)
		p.Constructing = append(p.Constructing, home)

		// If there is a spouse, add them as an owner and set their constructing building.
		if p.Spouse != nil {
			home.AddOwner(p.Spouse)
			if p.Spouse.Constructing != nil {
				panic("Spouse is already constructing a building")
			}
			p.Spouse.Constructing = append(p.Spouse.Constructing, home)
		}
		// TODO: Deduct from resources from both partners.
		p.Resources -= buildingCosts[BuildingTypeHouse]
	} else {
		log.Printf("No suitable location for a house found")
	}
}

// newDungeonCrawl creates a "behavior" tree for going on an adventure.
// NOTE: This is not a real behavior tree, but a simple sequence of tasks.
func newDungeonCrawl(p *Person, d *Building, onSuccess, onFailure func()) Tree {
	// - Move to the dungeon.
	// - Fight the monster.
	// - Move back home.
	homeX, homeY := int(p.X), int(p.Y)

	// Pick a random enemy.
	enemies := []string{"goblin", "troll", "dragon", "bear", "wolf", "orc", "giant", "spider", "snake", "bandit"}
	enemy := enemies[rand.Intn(len(enemies))]

	// Pick a random duration for the fight.
	dur := rand.Float64() * 10

	// TODO: What if we simply fail to reach the dungeon?
	// In this case, we should probably just go back home.
	log.Printf("%s is going to the location at %d,%d", p.String(), d.X, d.Y)
	tRoot := NewTaskMoveToXY(p, d.X, d.Y) // Move to the dungeon.
	t := tRoot.Then(NewTaskGeneric(p, "FightMonster", func(elapsed float64) TaskStatus {
		dur -= elapsed
		if dur <= 0 {
			if rand.Intn(100) < 10 {
				log.Printf("%s died on an adventure, killed by %s %s", p.String(), genlanguage.GetArticle(enemy), enemy)
				return TaskStatusFailed
			}
			log.Printf("%s killed %s %s", p.String(), genlanguage.GetArticle(enemy), enemy)
			return TaskStatusCompleted
		}
		log.Printf("%s is fighting %s %s (%.2f/%.2f)", p.String(), genlanguage.GetArticle(enemy), enemy, dur, elapsed)
		return TaskStatusInProgress
	})) // Fight the monster.
	t.Then(NewTaskMoveToXY(p, homeX, homeY)) // Move back home.

	return NewTree(tRoot, onSuccess, onFailure)
}

// handleAdventure executes / fullfilles the adventure goal.
func handleAdventure(p *Person, m *Map) {
	// Check if there is an adventure available to scratch that itch.
	// ... or if we are already on an adventure.
	// We're going on an adventure!
	// There is a chance that we die on our adventure.

	// Pick a random dungeon (TODO: Handle case if there are no dungeons
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
	if !dcTree.Do(1.0) {
		log.Println("Adventure failed")
	}
}

// handleJob executes / fullfilles the job goal.
func handleJob(p *Person, m *Map) {
	// We have a job, so we should work.

	// Check if we don't have a job yet and want one.
	// TODO: This would require a place to work...
	// Like, working on another farm, or upgrade the house with a farm plot, etc.
	if p.Job == JobTypeUnemployed {
		p.Job = JobTypeFarmer
	}

	// If we have a job, we can earn resources.
	// We might have a job, even though it is not our goal.
	// Maybe start with 16 years old?
	if p.Job != JobTypeUnemployed && rand.Intn(100) < 10 {
		// TODO: Maybe use fractional resources?
		p.Resources += 1 // 1 is quite a lot for a single day?
	}
}
