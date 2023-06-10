package genstatblock5e

import (
	"math"
	"math/rand"
	"time"
)

type Size int

const (
	Tiny Size = iota
	Small
	Medium
	Large
	Huge
	Gargantuan
	SizeMax
)

type Alignment [2]int

const (
	AlignChaotic = iota
	AlignNeutral
	AlignLawful
	AlignGood
	AlignEvil
)

type MonsterType string

const (
	TypeAberration  MonsterType = "aberration"
	TypeBeast       MonsterType = "beast"
	TypeCelestial   MonsterType = "celestial"
	TypeConstruct   MonsterType = "construct"
	TypeDragon      MonsterType = "dragon"
	TypeUndead      MonsterType = "undead"
	TypeElemental   MonsterType = "elemental"
	TypeFiend       MonsterType = "fiend"
	TypeFey         MonsterType = "fey"
	TypeGiant       MonsterType = "giant"
	TypeHumanoid    MonsterType = "humanoid"
	TypeMonstrosity MonsterType = "monstrosity"
	TypeOoze        MonsterType = "ooze"
	TypePlant       MonsterType = "plant"
)

var MonsterTypes = []MonsterType{
	TypeAberration,
	TypeBeast,
	TypeCelestial,
	TypeConstruct,
	TypeDragon,
	TypeUndead,
	TypeElemental,
	TypeFiend,
	TypeFey,
	TypeGiant,
	TypeHumanoid,
	TypeMonstrosity,
	TypeOoze,
	TypePlant,
}

type AbilityScore int

const (
	AbilityStrength AbilityScore = iota
	AbilityDexterity
	AbilityConstitution
	AbilityIntelligence
	AbilityWisdom
	AbilityCharisma
	AbilityScoreCount
)

type Sense string

const (
	SenseBlindsight        Sense = "blindsight"
	SenseDarkvision        Sense = "darkvision"
	SenseTremorsense       Sense = "tremorsense"
	SenseTruesight         Sense = "truesight"
	SensePassivePerception Sense = "passive perception"
)

type Skill string

const (
	// Strength
	SkillAthletics Skill = "athletics"
	// Dexterity
	SkillAcrobatics Skill = "acrobatics"
	SkillSleight    Skill = "sleight of hand"
	SkillStealth    Skill = "stealth"
	// Intelligence
	SkillArcana   Skill = "arcana"
	SkillHistory  Skill = "history"
	SkillInvest   Skill = "investigation"
	SkillNature   Skill = "nature"
	SkillReligion Skill = "religion"
	// Wisdom
	SkillAnimal     Skill = "animal handling"
	SkillInsight    Skill = "insight"
	SkillMedicine   Skill = "medicine"
	SkillPerception Skill = "perception"
	SkillSurvival   Skill = "survival"
	// Charisma
	SkillDeception   Skill = "deception"
	SkillIntimidate  Skill = "intimidation"
	SkillPerformance Skill = "performance"
	SkillPersuasion  Skill = "persuasion"
)

type Monster struct {
	Name          string
	Description   string
	Type          string
	Size          Size
	Alignment     Alignment
	ChallengeRate float64
	// Combat stats
	ArmorClass    int
	HitPoints     int
	Speed         int
	AbilityScores [AbilityScoreCount]int
	SavingThrows  [AbilityScoreCount]int

	// E.g. Deception +5, Perception +3
	Skills map[Skill]int

	// E.g. Darkvision 60 ft., passive Perception 13
	Senses map[Sense]int

	// TODO: Change these to be more structured (like attack, damage, range, etc.)
	Traits           []string
	Actions          []string
	Reactions        []string
	LegendaryActions []string
	Languages        []string
}

// GenerateAttributes randomly generates a set of attributes for a character.
//
// 'numAttrs' is the number of attributes to generate.
// 'maxAttrPoints' is the maximum number of points that can be assigned to an attribute.
// 'numPoints' is the total number of points to assign.
// 'randomness' is a value between 0 and 1 that determines how random the attributes are.
// - 0 means all attributes will be equal.
// - 1 means all attributes will be random up to the maximum.
//
// TODO: Allow a custom random number generator to be passed in.
func GenerateAttributes(numAttrs, maxAttrPoints, numPoints int, randomness float64) []int {
	rand.Seed(time.Now().UnixNano())

	attributes := make([]int, numAttrs)
	totalPoints := numPoints
	avgPoints := int(math.Ceil(float64(numPoints) / float64(numAttrs)))

	for _, i := range rand.Perm(numAttrs) {
		// Calculate the maximum and minimum value for the attribute.
		// Depending on the randomness, the value will be closer to the average.
		maxValue := min(maxAttrPoints, avgPoints+int(float64(avgPoints)*randomness))
		minValue := max(1, avgPoints-int(float64(avgPoints)*randomness))
		// Generate a random value between the min and max and no more than the remaining points
		value := rand.Intn(maxValue-minValue+1) + minValue
		value = min(value, totalPoints)
		// Update the attribute value
		attributes[i] = value
		// Decrease the total points
		totalPoints -= value
	}

	// If there are still points left, add them evenly to the attributes.
	for totalPoints > 0 {
		for _, i := range rand.Perm(numAttrs) {
			if totalPoints > 0 {
				attributes[i]++
				totalPoints--
			}
		}
	}

	return attributes
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
