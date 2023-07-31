package genstatblock5e

import (
	"math"
	"math/rand"
	"time"
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

	// E.g. Damage Resistances bludgeoning, piercing, and slashing from nonmagical attacks
	Resistances     []DamageType
	Immunities      []DamageType
	Vulnerabilities []DamageType

	// E.g. Darkvision 60 ft., passive Perception 13
	Senses map[Sense]int

	// TODO: Change these to be more structured (like attack, damage, range, etc.)
	Traits           []string
	Actions          []string
	Reactions        []string
	LegendaryActions []string
	Languages        []Language
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
