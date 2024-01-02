package simsettlers

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/Flokey82/go_gens/gameconstants"
)

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

// Person represents a person in the village.
// TODO:
// - Add a personality.
type Person struct {
	FirstName    string
	LastName     string
	Dead         bool
	Birthday     uint16      // day of the year the person was born
	Age          uint16      // age of the person in years (pretty generous with the uint16)
	Gender       byte        //
	Pregnant     uint16      // days the person will still be pregnant
	Goals        Goal        // personal goals
	Resources    int         // wealth in resources
	Job          JobType     // current job
	Home         *Building   // home of the person
	Constructing []*Building // buildings under construction
	Owns         []*Building // buildings owned
	Mother       *Person
	Father       *Person
	Spouse       *Person
	Children     []*Person
	Opinions     Opinions // opinions of other people
}

type Opinions map[*Person][2]int

func (o Opinions) Value(p *Person) int {
	return o[p][0]
}

func (o Opinions) Counter(p *Person) int {
	return o[p][1]
}

func (o Opinions) IncrementBy(p *Person, value int) {
	o[p] = [2]int{
		min(max(o[p][0]+value, 0), 255),
		min(o[p][1]+1, 255),
	}
}

func (o Opinions) Change(p *Person, value int) {
	o[p] = [2]int{min(max((o[p][0]*o[p][1]+value)/(o[p][1]+1), 0), 255), min(o[p][1]+1, 255)}
}

// OwnsOwnHome returns true if the person owns their own home.
func (p *Person) OwnsOwnHome() bool {
	if p.Home == nil {
		return false
	}
	for _, o := range p.Home.Owners {
		if o == p {
			return true
		}
	}
	return false
}

// LivesWithParents returns true if the person lives with their parents.
func (p *Person) LivesWithParents() bool {
	if p.Home == nil {
		return false
	}
	if p.Mother != nil && p.Mother.Home == p.Home {
		return true
	}
	if p.Father != nil && p.Father.Home == p.Home {
		return true
	}
	return false
}

func (p *Person) heir() *Person {
	// TODO: If there are no children, then siblings should inherit.
	// TODO: If there are no siblings, then parents should inherit.
	// TODO: If there are no children, then siblings should inherit.
	if p.Spouse != nil && !p.Spouse.Dead {
		return p.Spouse
	}
	for _, c := range p.Children {
		if !c.Dead {
			return c
		}
	}

	if p.Father != nil {
		return p.Father
	}
	if p.Mother != nil {
		return p.Mother
	}
	return nil
}

func (p *Person) SetHome(b *Building) {
	// Remove person from the occupants list of the previous home.
	if p.Home != nil {
		p.Home.RemoveOccupant(p)
	}
	p.Home = b
	b.AddOccupant(p)
}

// String returns the string representation of the person.
func (p *Person) String() string {
	str := fmt.Sprintf("%s %s (%d %s - %d)", p.FirstName, p.LastName, p.Age, p.genderString(), p.Resources)
	if p.Dead {
		str += " (dead)"
	}
	return str
}

func (p *Person) genderString() string {
	if p.Gender == GenderMale {
		return "M"
	}
	return "F"
}

func (p *Person) isMarried() bool {
	return p.Spouse != nil
}

func (p *Person) assignChildhoodGoals() {
	if rand.Intn(100) < 90 {
		p.Goals |= GoalChildhoodSocialize
	}
	if rand.Intn(100) < 5 {
		p.Goals |= GoalChildhoodBully
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
	}
	if rand.Intn(100) < 90 {
		p.Goals |= GoalAdultChildren
	}
	if rand.Intn(100) < 90 {
		p.Goals |= GoalAdultJob
	}
	if rand.Intn(100) < 5 {
		p.Goals |= GoalAdultAdventurer
	}
}

func (p *Person) assignElderlyGoals() {
}

const numDaysPregnant = 270

const (
	GenderFemale = iota
	GenderMale
)

func (m *Map) addNRandomPeople(n int) {
	for i := 0; i < n; i++ {
		p := &Person{
			LastName:  m.lastGen.String(),
			Birthday:  m.Day,
			Age:       uint16(rand.Intn(20) + 18),
			Gender:    byte(rand.Intn(2)),
			Resources: rand.Intn(10),
			Opinions:  make(Opinions),
		}
		if p.Gender == GenderMale {
			p.FirstName = m.firstGen[1].String()
		} else {
			p.FirstName = m.firstGen[0].String()
		}

		// TODO: Fix assignment of goals.
		p.assignChildhoodGoals()
		p.assignAdultGoals()
		// p.assignElderlyGoals()

		m.RealPop = append(m.RealPop, p)
	}
}

func (m *Map) handleDeath(p *Person) {
	m.Population--
	// TODO: Remove from spouse?
	log.Printf("Died: %v", p)
	p.Dead = true

	// TODO:
	// - Identify who will inherit all buildings, resources, etc.
	// - Find all buildings that we own and transfer ownership to the spouse or children.
	// - What if we are constructing a building? Transfer ownership to the spouse or children.
	// - Allow multiple heirs.

	heir := p.heir()
	if heir != nil {
		// Move resources to the heir.
		heir.Resources += p.Resources

		// Move buildings to the heir.
		for _, b := range p.Owns {
			// Remove from the occupants list of the building.
			b.RemoveOccupant(p)

			// Remove from the owners list of the building.
			b.RemoveOwner(p)

			// Add to the occupants list of the building.
			b.AddOwner(heir)
		}

		// Move constructing buildings to the heir.
		for _, b := range p.Constructing {
			// Remove the deceased from the constructing list of the building.
			b.RemoveOwner(p)

			// Add heir to the owners list of the building.
			b.AddOwner(heir)

			// Add the building to the constructing list of the heir.
			heir.Constructing = append(heir.Constructing, b)
		}

		// Move a homeless heir to the home of the deceased.
		if p.Home != nil && heir.Home == nil {
			heir.SetHome(p.Home)
		}
	}

	if p.Home != nil {
		// Remove from the occupants list of the home.
		p.Home.RemoveOccupant(p)

		// Remove from the owners list of the home.
		p.Home.RemoveOwner(p)
	}

	// Move to the cemetery.
	p.SetHome(m.Cemetery)
}

func (m *Map) agePop() {
	// Let's see who has a birthday today and check if anyone dies.
	var remPop []*Person
	for _, p := range m.RealPop {
		if p.Dead {
			continue
		}
		if p.Birthday == m.Day {
			p.Age++

			// Check if we have to decide life goals.
			if p.Age == 1 {
				p.assignChildhoodGoals()
			} else if p.Age == 18 {
				p.assignAdultGoals()
			} else if p.Age == 65 {
				p.assignElderlyGoals()
			}
		}

		// Check if anyone dies.
		if gameconstants.DiesAtAgeWithinNDays(int(p.Age), 1) {
			m.handleDeath(p)
		} else {
			remPop = append(remPop, p)
		}
	}
	// TODO: Find a more efficient way to do this besides re-allocation.
	m.RealPop = remPop
}

func (m *Map) matchSingles() {
	const ageDiffFactor = 0.25
	const menTakeFemaleName = true
	var men, women []*Person
	for _, p := range m.RealPop {
		if p.Dead || p.isMarried() || p.Age < 18 {
			continue
		}
		// Check if we even want a partner.
		if !p.Goals.IsSet(GoalAdultPartner) {
			continue
		}
		if p.Gender == GenderMale {
			men = append(men, p)
		} else if p.Gender == GenderFemale {
			women = append(women, p)
		}
	}

	// Sort by age.
	sort.Slice(men, func(i, j int) bool {
		return men[i].Age < men[j].Age
	})
	sort.Slice(women, func(i, j int) bool {
		return women[i].Age < women[j].Age
	})

	// Match the singles with a max age difference of 25%.
	// TODO:
	// - Personalities matter.
	// - Matching should not just happen, there should be a probability and some randomness.
	// - There should be a chance that a person stays single (by choice, or because they just can't find the right person).
	for _, wp := range women {
		// Ladies' choice.
		for _, mp := range men {
			if mp.isMarried() {
				continue
			}
			// Check if we're on the same page regarding children.
			// TODO: Check other goals as well.
			if wp.Goals.IsSet(GoalAdultChildren) != mp.Goals.IsSet(GoalAdultChildren) {
				continue
			}

			if ageDiff := math.Abs(float64(wp.Age - mp.Age)); ageDiff < ageDiffFactor*float64(wp.Age) {
				wp.Spouse = mp
				mp.Spouse = wp
				if menTakeFemaleName {
					mp.LastName = wp.LastName
				} else {
					wp.LastName = mp.LastName
				}
				log.Printf("Matched %v and %v", wp, mp)
				break
			}
		}
	}

	// TODO: Instead, we should go through both men and women and find their respective best match
	// based on their goals, opinions, etc.
}

func (m *Map) advancePregnancies() {
	// Advance all pregnancies.
	for _, p := range m.RealPop {
		if p.Dead {
			continue
		}
		if p.Pregnant > 0 {
			p.Pregnant--
			if p.Pregnant == 0 {
				child := &Person{
					Home:     p.Home,
					Mother:   p,
					Father:   p.Spouse,
					LastName: p.LastName, // Last name of the mother.
					Birthday: m.Day,
					Age:      0,
					Gender:   byte(rand.Intn(2)),
					Opinions: make(Opinions),
				}
				if child.Gender == GenderMale {
					child.FirstName = m.firstGen[1].String()
				} else {
					child.FirstName = m.firstGen[0].String()
				}
				log.Printf("Born: %v", child)
				p.Children = append(p.Children, child)
				if p.Home != nil {
					child.SetHome(p.Home)
				}
				if p.Spouse != nil {
					p.Spouse.Children = append(p.Spouse.Children, child)
				}
				m.RealPop = append(m.RealPop, child)
			}
		}
	}

	// All couples have a chance of getting pregnant, which decreases
	// by the number of children they already have.
	for _, p := range m.RealPop {
		if p.isMarried() && p.Pregnant == 0 && p.Gender == GenderFemale {
			// TODO: Check if both want children.
			if !p.Goals.IsSet(GoalAdultChildren) && !p.Spouse.Goals.IsSet(GoalAdultChildren) && rand.Intn(1000) > 1 {
				// If neither wants children, there is a very small chance of getting pregnant.
				continue
			} else if (!p.Goals.IsSet(GoalAdultChildren) || !p.Spouse.Goals.IsSet(GoalAdultChildren)) && rand.Intn(100) > 1 {
				// If only one or wants children, there is a slight chance of getting pregnant.
				continue
			}
			if rand.Intn(100) < 3-len(p.Children) {
				p.Pregnant = uint16(numDaysPregnant)
			}
		}
	}
}

func (m *Map) tickPeople() {
	// If true, owners that live in a house will repair their homes.
	repairHome := false
	for _, p := range m.RealPop {
		if p.Dead {
			continue
		}

		// Check if we are old enough to work or build a house.
		// TODO: We should vary this based on personality.
		// Some people might refuse to retire, others might retire early.
		if p.Age < 18 {
			m.tickChildhood(p)
		} else if p.Age < 65 {
			m.tickAdulthood(p, repairHome)
		} else {
			m.tickElderly(p)
		}
	}
}

func (m *Map) tickChildhood(p *Person) {
	// TODO: Check if goals are satisfied or not.

	// TODO: Factor this out into a function... This can be generalized as adults also bully other adults.
	// If we are a bully, we might bully other children.
	// NOTE: This breaks my heart to write this, but this is a simulation of life,
	// and there are asshole children in the world. Bullying can change people's
	// lives forever, and cause lasting damage. But it sometimes cause children to
	// develop a strong sense of justice, and become a better person overall.
	if p.Goals.IsSet(GoalChildhoodBully) && rand.Intn(100) < 1 && p.Age > 5 {
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
			log.Printf("%v is bullying %v", p, victim)

			// Change the victim's opinion of the bully.
			victim.Opinions.IncrementBy(p, -20)

			// Change the bully's opinion of the victim.
			// TODO: The bullying might backfire, then we should change the bully's opinion of the victim positively
			// to avoid bullying the same person unsuccessfully over and over again.
			p.Opinions.IncrementBy(victim, -10)

			// Depending on the personality, this might result in violence, personality changes, etc.
			if victim.Goals.IsSet(GoalChildhoodSocialize) && rand.Intn(100) < 5 {
				// If the victim is social, they might loose the will to socialize.
				victim.Goals &= ^GoalChildhoodSocialize
			}
		}
	}

	// If we socialize, let's see if we can make friends.
	if p.Goals.IsSet(GoalChildhoodSocialize) && rand.Intn(100) < 1 && p.Age > 5 {
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

		// Sort by opinion.
		sort.Slice(friends, func(i, j int) bool {
			return p.Opinions.Value(friends[i]) > p.Opinions.Value(friends[j])
		})

		// Pick the friend with the highest opinion by the person.
		if len(friends) > 0 {
			friend := friends[0]
			// TODO: Occasionally, we might mix it up and socialize with someone else.
			if len(friends) > 1 && rand.Intn(100) < 10 {
				friend = friends[min(rand.Intn(3), len(friends)-1)]
			}
			log.Printf("%v is socializing with %v", p, friend)

			// TODO: Based on personality, we might become friends, or we might become enemies.

			// Change the friend's opinion of the person.
			friend.Opinions.IncrementBy(p, 10)

			// Change the person's opinion of the friend.
			p.Opinions.IncrementBy(friend, 10)
		}
	}
}

func (m *Map) tickAdulthood(p *Person, repairHome bool) {
	// TODO: Check if goals are satisfied or not.

	handleJob := func() {
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

	handleHome := func() {
		// Check if we are already building a house or if our spouse is.
		if p.Constructing != nil || p.Spouse != nil && p.Spouse.Constructing != nil {
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

		// Check if we own our own home.
		if p.OwnsOwnHome() {
			// Check if we need to repair our home.
			// TODO:
			// - Make the decision to repair or not based on personality
			// and the condition of the building.
			// - Any occupant should be able to repair the home if they are
			// old enough, have enough resources, and are crafty enough.
			if repairHome && p.Home.Condition < 100 && p.Resources > 0 {
				p.Home.Condition++
				p.Resources--
			}
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

	handleAdventure := func() {
		// Check if there is an adventure available to scratch that itch.
		// ... or if we are already on an adventure.

		// NOTE: Adventures should be memorable, rare occurrences with
		// lasting effects on the person. High risk, high reward.
		// This might get us killed, we migh bear scars, or maybe gain
		// a new fear, a new skill, maybe an artifact, etc.
		if rand.Intn(100) < 1 {
			// We're going on an adventure!
			// There is a chance that we die on our adventure.
			// TODO: What if we just go missing? Could someone save us?
			if rand.Intn(100) < 10 {
				// We died on our adventure.
				// We might be declared missing, someone might find our body,
				// or we might never be found.
				log.Printf("%s died on an adventure", p.String())
				// TODO: Inheritence, ownership, legally dead, etc.
				m.handleDeath(p)
				return
			}
			// We survived the adventure.
			loot := rand.Intn(1000)
			p.Resources += loot
			log.Printf("%s went on an adventure and found %d resources", p.String(), loot)
		}
	}

	// Handle our job goal.
	if p.Goals.IsSet(GoalAdultJob) {
		handleJob()
	}

	// Check if we even want a home.
	if p.Goals.IsSet(GoalAdultHome) {
		handleHome()
	}

	// Check if we want to go on an adventure.
	if p.Goals.IsSet(GoalAdultAdventurer) {
		handleAdventure()
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
}

func (m *Map) tickElderly(p *Person) {
}
