package genlandmarknames

// ForestAdjectives are adjectives for forests.
var ForestAdjectives = []string{
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
}

// ForestSubjects are subjects for forests.
var ForestSubjects = []string{
	"forest",
	"woods",
	"wood",
	"grove",
	"groves",
	"thicket",
	"thickets",
}

// DangerGenitivePhraseAdjForest returns a danger genitive phrase adjective for forests.
var DangerGenitivePhraseAdjForest = []string{
	"rotten",
	"lost",
	"petrified",
	"forgotten",
	"abandoned",
	"consumed",
	"buried",
}

// NewForestGenerator returns a new generator for forest names.
func NewForestGenerator(seed int64) *BasicGenerator {
	return NewBasicGenerator(seed,
		ForestAdjectives,
		ForestSubjects,
		WordPair{
			A: DangerGenitivePhraseAdjForest,
			B: DangerGenitivePhraseSubject,
		})
}
