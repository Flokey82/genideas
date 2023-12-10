package genlandmarknames

// NewRiverGenerator returns a new generator for river names.
type RiverGenerator struct {
	*namer
	Adjective       []string
	SubjectSmall    []string
	SubjectLarge    []string
	DangerousSuffix WordPair
}

// RiverAdjectives is a list of adjectives for rivers.
var RiverAdjectives = []string{
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
}

// RiverSubjectsSmall is a list of subjects for small rivers.
var RiverSubjectsSmall = []string{
	"brook",
	"run",
	"brooklet",
	"waters",
	"flow",
	"wee",
}

// RiverSubjectsLarge is a list of subjects for large rivers.
var RiverSubjectsLarge = []string{
	"river",
	"stream",
	"current",
	"torrent",
	"flow",
}

// DangerGenitivePhraseAdjRiver returns a danger genitive phrase adjective for rivers.
var DangerGenitivePhraseAdjRiver = []string{
	"drowned",
	"lost",
	"sunken",
	"forgotten",
	"abandoned",
	"flushed",
}

// NewRiverGenerator returns a new generator for river names.
func NewRiverGenerator(seed int64) *RiverGenerator {
	return &RiverGenerator{
		namer:        newNamer(seed),
		Adjective:    RiverAdjectives,
		SubjectSmall: RiverSubjectsSmall,
		SubjectLarge: RiverSubjectsLarge,
		DangerousSuffix: WordPair{
			A: DangerGenitivePhraseAdjRiver,
			B: DangerGenitivePhraseSubject,
		},
	}
}

// Generate generates a river name.
func (g *RiverGenerator) Generate(seed int64, small, dangerous bool) string {
	// TODO: Add chance to generate a custom name using fantasyname.
	g.resetToSeed(seed)
	prefix := g.randomString(g.Adjective)
	var suffix string
	if small {
		suffix = g.randomString(g.SubjectSmall)
	} else {
		suffix = g.randomString(g.SubjectLarge)
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
