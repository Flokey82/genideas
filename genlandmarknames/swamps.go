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

// NewSwampGenerator returns a new generator for swamp names.
func NewSwampGenerator(seed int64) *BasicGenerator {
	return NewBasicGenerator(seed,
		SwampAdjectives,
		SwampSubjects,
		WordPair{
			A: GenitivePhraseAdjNegWaterStill,
			B: GenitivePhraseSubject,
		})
}
