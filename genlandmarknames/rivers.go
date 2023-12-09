package genlandmarknames

// NewRiverGenerator returns a new generator for river names.
type RiverGenerator struct {
	*namer
	Prefix          []string
	SuffixSmall     []string
	SuffixLarge     []string
	DangerousSuffix WordPair
}

// NewRiverGenerator returns a new generator for river names.
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

// Generate generates a river name.
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
