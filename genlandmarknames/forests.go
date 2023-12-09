package genlandmarknames

// NewForestGenerator returns a new generator for forest names.
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
