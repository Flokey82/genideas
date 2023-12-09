package genlandmarknames

// NewDesertGenerator returns a new generator for desert names.
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
