package genlandmarknames

// NewSwampGenerator returns a new generator for swamp names.
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
