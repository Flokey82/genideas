// Package genlandmarknames generates names for landmarks like mountains, rivers, forests, etc.
package genlandmarknames

import (
	"math"
	"math/rand"
	"strings"
)

// NameGenerators contains all the generators for the different types of landmarks.
type NameGenerators struct {
	Desert        *BasicGenerator
	Mountain      *BasicGenerator
	MountainRange *BasicGenerator
	Forest        *BasicGenerator
	Swamp         *BasicGenerator
	River         *RiverGenerator
	Lake          *BasicGenerator
}

// NewNameGenerators returns a new NameGenerators.
func NewNameGenerators(seed int64) *NameGenerators {
	return &NameGenerators{
		Desert:        NewDesertGenerator(seed),
		Mountain:      NewMountainGenerator(seed),
		MountainRange: NewMountainRangeGenerator(seed),
		Forest:        NewForestGenerator(seed),
		Swamp:         NewSwampGenerator(seed),
		River:         NewRiverGenerator(seed),
		Lake:          NewLakeGenerator(seed),
	}
}

// LargeAreaSuffix is a list of suffixes suitable for large areas.
var LargeAreaSuffix = []string{
	"land",
	"plains",
	"expanse",
	"region",
}

// FertileLandSuffix is a list of suffixes suitable for fertile land.
var FertileLandPrefix = []string{
	"green",
	"lush",
	"bountiful",
	"fruitful",
	"rich",
	"abundant",
	"fertile",
}

// DangerousTerrainDescriptor returns a random descriptor for a dangerous terrain.
// This suffix is supposed to be used with "... of ...".
var DangerousTerrainDescriptor = []string{
	"death",
	"doom",
	"despair",
	"darkness",
	"evil",
	"chaos",
	"madness",
	"loss",
	"pain",
	"anguish",
	"terror",
	"horror",
	"lost souls",
	"the dead",
	"the damned",
	"the cursed",
	"the forsaken",
	"the lost",
	"the forgotten",
	"the abandoned",
	"the unknown",
	"the doomed",
}

// DangerGenitivePhraseAdj is a list of prefixes suitable for dangerous areas.
var DangerGenitivePhraseAdj = []string{
	"cursed",
	"forsaken",
	"forbidden",
	"lost",
	"abandoned",
	"unknown",
	"doomed",
	"haunted",
	"dark",
	"evil",
	"chaotic",
	"mad",
	"lost",
	"lost souls",
	"the dead",
	"the damned",
	"the cursed",
	"the forsaken",
	"the lost",
	"the forgotten",
	"the abandoned",
	"the unknown",
	"the doomed",
}

// DangerGenitivePhraseSubject is a list of suffixes suitable for dangerous areas.
var DangerGenitivePhraseSubject = []string{
	"hope",
	"souls",
	"fate",
	"dreams",
	"love",
	"life",
	"joy",
	"peace",
	"serenity",
	"calm",
	"tranquility",
	"corpses",
	"heroes",
	"villains",
	"princes",
	"princesses",
	"queens",
	"kings",
	"emperors",
	"empresses",
	"lords",
	"ladies",
	"knights",
	"lovers",
}

// NewBasicGenerator returns a new generator for basic names.
type BasicGenerator struct {
	*namer
	Prefix          []string
	Suffix          []string
	DangerousSuffix WordPair
}

// NewBasicGenerator returns a new generator for basic names.
func NewBasicGenerator(seed int64, prefix, suffix []string, danger WordPair) *BasicGenerator {
	return &BasicGenerator{
		namer:           newNamer(seed),
		Prefix:          prefix,
		Suffix:          suffix,
		DangerousSuffix: danger,
	}
}

// Generate generates a basic name.
func (g *BasicGenerator) Generate(seed int64, dangerous bool) string {
	// TODO: Add chance to generate a custom name using fantasyname.
	g.resetToSeed(seed)
	name := "The " + g.randomPair(g.Prefix, g.Suffix)
	if !dangerous {
		return name
	}
	return name + " of " + g.randomPair(g.DangerousSuffix.A, g.DangerousSuffix.B)
}

type namer struct {
	rand *rand.Rand
}

func newNamer(seed int64) *namer {
	return &namer{
		rand: rand.New(rand.NewSource(seed)),
	}
}

func (n *namer) resetToSeed(seed int64) {
	n.rand.Seed(seed)
}

func (n *namer) randomPair(a, b []string) string {
	// Make sure that the string chosen from a is not contained in b and vice versa.
	// This is to avoid names like "The muddy mud" or "The swampy swamp".
	for i := 0; i < 100; i++ {
		s1 := n.randomString(a)
		s2 := n.randomString(b)
		if !strings.Contains(s2, s1) && !strings.Contains(s1, s2) {
			return s1 + " " + s2
		}
	}
	// If we can't find a pair that doesn't contain the other, just return a random pair.
	return n.randomString(a) + " " + n.randomString(b)
}

func (n *namer) randomString(list []string) string {
	return list[n.rand.Intn(len(list))]
}

func (n *namer) randomChance(chance float64) bool {
	return math.Abs(n.rand.Float64()) < chance
}

type WordPair struct {
	A []string
	B []string
}
