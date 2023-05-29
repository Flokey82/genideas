package genclutter

import "math/rand"

type Rarity struct {
	Name   string
	Rarity int
}

func (r *Rarity) Roll() bool {
	return rand.Intn(101) >= r.Rarity
}

var RarityAbundant = &Rarity{
	Name:   "abundant",
	Rarity: 25,
}

var RarityCommon = &Rarity{
	Name:   "common",
	Rarity: 45,
}

var RarityAverage = &Rarity{
	Name:   "average",
	Rarity: 65,
}

var RarityUncommon = &Rarity{
	Name:   "uncommon",
	Rarity: 80,
}

var RarityRare = &Rarity{
	Name:   "rare",
	Rarity: 93,
}

var RarityExotic = &Rarity{
	Name:   "exotic",
	Rarity: 99,
}

var RarityLegendary = &Rarity{
	Name:   "legendary",
	Rarity: 100,
}

// Condition represents the condition of an item.
type Condition struct {
	Name   string
	Rarity int
}

func (c *Condition) Roll() bool {
	return rand.Intn(101) >= c.Rarity
}

var ConditionDecaying = &Condition{
	Name:   "decaying",
	Rarity: 95,
}

var ConditionBusted = &Condition{
	Name:   "busted",
	Rarity: 85,
}

var ConditionPoor = &Condition{
	Name:   "poor",
	Rarity: 75,
}

var ConditionAverage = &Condition{
	Name:   "average",
	Rarity: 50,
}

var ConditionGood = &Condition{
	Name:   "good",
	Rarity: 60,
}

var ConditionExquisite = &Condition{
	Name:   "exquisite",
	Rarity: 100,
}

func RollCondition() *Condition {
	val := rand.Intn(101)
	if val >= 100 {
		return ConditionExquisite
	}
	if val >= 95 {
		return ConditionDecaying
	}
	if val >= 85 {
		return ConditionBusted
	}
	if val >= 75 {
		return ConditionPoor
	}
	if val >= 60 {
		return ConditionGood
	}
	return ConditionAverage
}
