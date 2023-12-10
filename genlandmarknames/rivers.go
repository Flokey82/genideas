package genlandmarknames

import "github.com/Flokey82/go_gens/genstory"

const (
	TokenNoun           = "[NOUN]"
	TokenAdjGenNegative = "[ADJ_GEN_NEG]"
	TokenAdjGenPositive = "[ADJ_GEN_POS]"
)

// RiverTemplates is a list of templates for river names.
var RiverTemplates = []string{
	"[NOUN][SUBJECT]", // Shimmerbrook, Silverbrook, Diamondriver, etc.
	"[ADJ][SUBJECT]",  // Clearriver, Oldstream, etc.
	"[ADJ] [SUBJECT]", // Clear River, Cold Brook, etc.
	"[ADJ] [SUBJECT] of [ADJ_GEN_NEG] [SUBJECT_GEN]",
	"[ADJ] [SUBJECT] of [ADJ_GEN_POS] [SUBJECT_GEN]",
	"[ADJ] [SUBJECT] of [SUBJECT_GEN]",
	"[SUBJECT] of [ADJ_GEN_NEG] [SUBJECT_GEN]",
	"[SUBJECT] of [ADJ_GEN_POS] [SUBJECT_GEN]",
	"[SUBJECT] of [SUBJECT_GEN]",
	"[PLACE] [SUBJECT]",
}

// RiverPrefixNouns is a list of nouns used as prefixes for river names.
// Examples: Shimmer-, Silver-, Diamond-, etc.
var RiverNouns = []string{
	"gold",
	"silver",
	"bronze",
	"iron",
	"copper",
	"diamond",
	"ruby",
	"emerald",
	"jade",
	"pearl",
	"amber",
	"crystal",
	"opal",
	"topaz",
	"shimmer",
	"sparkle",
	"glitter",
	"glint",
	"gleam",
	"glow",
	"mirror",
	"glass",
	"ice",
	"water",
}

// TODO:
// - Prevent picking of closely related consecutive (or close in proximity) adjectives and subjects
// (except if aliteration is on).
// - Add chance to use wild adjectives for small rivers.
// - Find a way to control the implicit properties communicated through the name
//   - type of river (calm, fast, etc.)
//   - the size (small, large, etc.)
//   - positive or negative connotations (genetive phrases with positive or negative adjectives)
var RiverNameConfig = &genstory.TextConfig{
	TokenPools: map[string][]string{
		TokenNoun:           RiverNouns,
		TokenAdj:            RiverAdjectivesMild, // TODO: Add chance to use wild adjectives.
		TokenSubject:        RiverSubjectsSmall,
		TokenAdjGenNegative: GenitivePhraseAdjNegWaterFlowing,
		TokenAdjGenPositive: GenitivePhraseAdjPos,
		TokenSubjectGen:     GenitivePhraseSubject,
	},
	TokenIsMandatory: map[string]bool{},
	Tokens:           []string{TokenAdj, TokenSubject, TokenSubjectGen, TokenPlace, TokenNoun, TokenAdjGenPositive, TokenAdjGenNegative},
	Templates:        RiverTemplates,
	UseAllProvided:   true,
	UseAlliteration:  false,
	Title:            true,
}

// NewRiverGenerator returns a new generator for river names.
type RiverGenerator struct {
	*namer
	AdjectiveMild   []string
	AdjectiveWild   []string
	SubjectSmall    []string
	SubjectLarge    []string
	DangerousSuffix WordPair
	RegularSuffix   WordPair
}

// RiverAdjectivesMild is a list of adjectives for rivers.
// TODO: Split into two lists, one for calm rivers and one for fast rivers.
var RiverAdjectivesMild = []string{
	"calm",
	"quiet",
	"still",
	"slow",
	"gentle",
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
	"twinkling",
	"glinting",
	"gleaming",
	"glowing",
	"luminous",
}

// RiverAdjectivesWild is a list of adjectives for wild rivers.
var RiverAdjectivesWild = []string{
	"rushing",
	"frothing",
	"foaming",
	"churning",
	"boiling",
	"raging",
	"roaring",
	"thundering",
	"turbulent",
	"wild",
	"raging",
	"angry",
	"violent",
	"storming",
	"tempestuous",
	"tumultuous",
	"seething",
	"wrathful",
	"twisted",
	"tangled",
}

// RiverSubjectsSmall is a list of subjects for small rivers.
var RiverSubjectsSmall = []string{
	"brook",
	"run",
	"brooklet",
	"waters",
	"flow",
	"wee",
	"stream",
	"rivulet",
	"runnel",
	"runlet",
	"rill",
	"burn",
	"beck",
	"bourn",
}

// RiverSubjectsLarge is a list of subjects for large rivers.
var RiverSubjectsLarge = []string{
	"river",
	"stream",
	"current",
	"torrent",
	"flow",
	"water",
	"waters",
	"bend",
	"carve",
	"channel",
	"course",
}

// NewRiverGenerator returns a new generator for river names.
func NewRiverGenerator(seed int64) *RiverGenerator {
	return &RiverGenerator{
		namer:         newNamer(seed),
		AdjectiveMild: RiverAdjectivesMild,
		AdjectiveWild: RiverAdjectivesWild,
		SubjectSmall:  RiverSubjectsSmall,
		SubjectLarge:  RiverSubjectsLarge,
		DangerousSuffix: WordPair{
			A: GenitivePhraseAdjNegWaterFlowing,
			B: GenitivePhraseSubject,
		},
		RegularSuffix: WordPair{
			A: GenitivePhraseAdjPos,
			B: GenitivePhraseSubject,
		},
	}
}

// Generate generates a river name.
func (g *RiverGenerator) Generate(seed int64, small, dangerous bool) string {
	// TODO: Add chance to generate a custom name using fantasyname.

	// Select the adjective.
	g.resetToSeed(seed)
	var prefix string
	if dangerous {
		prefix = g.randomString(g.AdjectiveWild)
	} else {
		prefix = g.randomString(g.AdjectiveMild)
	}

	// Select the subject.
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

	// If the river is dangerous, we add a genitive phrase with a negative adjective.
	if dangerous {
		return name + " of " + g.randomPair(g.DangerousSuffix.A, g.DangerousSuffix.B)
	}

	// There is a chance that we add a genitive phrase.
	if g.randomChance(0.5) {
		return name + " of " + g.randomPair(g.RegularSuffix.A, g.RegularSuffix.B)
	}

	return name
}
