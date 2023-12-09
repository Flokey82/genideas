package genlandmarknames

import (
	"math"
	"math/rand"
	"strings"
)

type NameGenerators struct {
	Desert        *BasicGenerator
	Mountain      *BasicGenerator
	MountainRange *BasicGenerator
	Forest        *BasicGenerator
	Swamp         *BasicGenerator
	River         *RiverGenerator
}

func NewNameGenerators(seed int64) *NameGenerators {
	return &NameGenerators{
		Desert:        NewDesertGenerator(seed),
		Mountain:      NewMountainGenerator(seed),
		MountainRange: NewMountainRangeGenerator(seed),
		Forest:        NewForestGenerator(seed),
		Swamp:         NewSwampGenerator(seed),
		River:         NewRiverGenerator(seed),
	}
}

var LargeAreaSuffix = []string{
	"land",
	"plains",
	"expanse",
	"region",
}

var FertileLandPrefix = []string{
	"green",
	"lush",
	"bountiful",
	"fruitful",
	"rich",
	"abundant",
	"fertile",
}

func NewDesertGenerator(seed int64) *BasicGenerator {
	return NewBasicGenerator(seed, []string{
		"charred",
		"desolate",
		"empty",
		"arid",
		"bleak",
		"scorched",
		"burnt",
		"forsaken",
	}, []string{
		"desert",
		"wasteland",
		"sands",
		"barrens",
		"expanse",
		"region",
	}, WordPair{
		A: []string{
			"burned",
			"scorched",
			"charred",
			"lost",
			"ashen",
		},
		B: DangerousSuffixB,
	})
}

var MountainPrefix = []string{
	"rocky",
	"mountainous",
	"spiked",
	"steep",
	"rough",
	"craggy",
	"toothy",
	"jagged",
	"broken",
}

var MountainDangerPrefix = []string{
	"shattered",
	"lost",
	"petrified",
	"forgotten",
	"abandoned",
	"broken",
	"buried",
	"fallen",
}

func NewMountainRangeGenerator(seed int64) *BasicGenerator {
	return NewBasicGenerator(seed, MountainPrefix, []string{
		"mountains",
		"teeth",
		"spikes",
		"peaks",
		"rocks",
		"thorns",
		"jags",
		"spurs",
		"spires",
		"pinnacles",
	}, WordPair{
		A: MountainDangerPrefix,
		B: DangerousSuffixB,
	})
}

func NewMountainGenerator(seed int64) *BasicGenerator {
	return NewBasicGenerator(seed, MountainPrefix, []string{
		"mountain",
		"tooth",
		"spike",
		"peak",
		"rock",
		"thorn",
		"jag",
		"spur",
		"spire",
		"pinnacle",
		"stone",
	}, WordPair{
		A: MountainDangerPrefix,
		B: DangerousSuffixB,
	})
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

var DangerousPrefixes = []string{
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

var DangerousSuffixB = []string{
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

func NewForestGenerator(seed int64) *BasicGenerator {
	return NewBasicGenerator(seed, []string{
		"dark",
		"shadowy",
		"twisted",
		"broad",
		"thick",
		"dense",
		"overgrown",
		"lush",
		"green",
		"mossy",
		"moldy",
		"mold-covered",
		"leafy",
		"leaf-covered",
		"leaf-strewn",
		"leaf-littered",
		"woody",
		"wooded",
		"wood-strewn",
		"wood-littered",
		"wood-covered",
		"wooden",
		"bark-covered",
		"bark-strewn",
		"bark-littered",
		"barky",
	}, []string{
		"forest",
		"woods",
		"wood",
		"grove",
		"groves",
		"thicket",
		"thickets",
	}, WordPair{
		A: []string{
			"rotten",
			"lost",
			"petrified",
			"forgotten",
			"abandoned",
			"consumed",
			"buried",
		},
		B: DangerousSuffixB,
	})
}

func NewSwampGenerator(seed int64) *BasicGenerator {
	return NewBasicGenerator(seed, []string{
		"muddy",
		"mud-covered",
		"swampy",
		"marshy",
		"sticky",
		"humid",
		"muggy",
		"wet",
		"waterlogged",
		"water-covered",
		"moist",
		"mossy",
		"moldy",
		"decaying",
		"rotten",
		"rotting",
		"putrid",
		"stinking",
		"stenchy",
	}, []string{
		"swamp",
		"marsh",
		"mire",
		"bog",
		"quagmire",
		"quag",
		"porridge",
		"muck",
		"mud",
		"phlegm",
		"sewer",
		"sludge",
		"puddles",
	}, WordPair{
		A: []string{
			"drowned",
			"lost",
			"sunken",
			"forgotten",
			"abandoned",
			"flushed",
			"buried",
		},
		B: DangerousSuffixB,
	})
}

type RiverGenerator struct {
	*namer
	Prefix          []string
	SuffixSmall     []string
	SuffixLarge     []string
	DangerousSuffix WordPair
}

func NewRiverGenerator(seed int64) *RiverGenerator {
	return &RiverGenerator{
		namer: newNamer(seed),
		Prefix: []string{
			"clear",
			"clean",
			"pure",
			"fresh",
			"cold",
			"cool",
			"bracing",
			"refreshing",
			"crisp",
			"fast",
			"swift",
			"running",
			"rushing",
			"flowing",
			"rippling",
			"lively",
			"snaking",
			"meandering",
			"bubbling",
			"sparkling",
			"glittering",
			"shimmering",
			"shining",
		},
		SuffixSmall: []string{
			"brook",
			"run",
			"brooklet",
			"waters",
			"flow",
			"wee",
		},
		SuffixLarge: []string{
			"river",
			"stream",
			"current",
			"torrent",
			"flow",
		},
		DangerousSuffix: WordPair{
			A: []string{
				"drowned",
				"lost",
				"sunken",
				"forgotten",
				"abandoned",
				"flushed",
			},
			B: DangerousSuffixB,
		},
	}
}

func (g *RiverGenerator) Generate(seed int64, small, dangerous bool) string {
	// TODO: Add chance to generate a custom name using fantasyname.
	g.resetToSeed(seed)
	prefix := g.randomString(g.Prefix)
	var suffix string
	if small {
		suffix = g.randomString(g.SuffixSmall)
	} else {
		suffix = g.randomString(g.SuffixLarge)
	}
	// There is a chance that we simply merge the two words.
	var name string
	if g.randomChance(0.5) && !dangerous {
		name = prefix + suffix
	} else {
		name = prefix + " " + suffix
	}
	if !dangerous {
		return name
	}
	return name + " of " + g.randomPair(g.DangerousSuffix.A, g.DangerousSuffix.B)
}

type BasicGenerator struct {
	*namer
	Prefix          []string
	Suffix          []string
	DangerousSuffix WordPair
}

func NewBasicGenerator(seed int64, prefix, suffix []string, danger WordPair) *BasicGenerator {
	return &BasicGenerator{
		namer:           newNamer(seed),
		Prefix:          prefix,
		Suffix:          suffix,
		DangerousSuffix: danger,
	}
}

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
