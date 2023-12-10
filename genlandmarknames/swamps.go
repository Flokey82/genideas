package genlandmarknames

// SwampAdjectives is a list of adjectives for swamps.
var SwampAdjectives = []string{
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
}

// SwampSubjects is a list of subjects for swamps.
var SwampSubjects = []string{
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
}

// DangerGenitivePhraseAdjSwamp returns a danger genitive phrase adjective for swamps.
var DangerGenitivePhraseAdjSwamp = []string{
	"drowned",
	"lost",
	"sunken",
	"forgotten",
	"abandoned",
	"flushed",
}

// NewSwampGenerator returns a new generator for swamp names.
func NewSwampGenerator(seed int64) *BasicGenerator {
	return NewBasicGenerator(seed,
		SwampAdjectives,
		SwampSubjects,
		WordPair{
			A: DangerGenitivePhraseAdjSwamp,
			B: DangerGenitivePhraseSubject,
		})
}
