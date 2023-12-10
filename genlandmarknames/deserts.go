package genlandmarknames

// DesertAdjectives is a list of adjectives suitable for deserts.
var DesertAdjectives = []string{
	"charred",
	"desolate",
	"empty",
	"arid",
	"bleak",
	"scorched",
	"burnt",
	"forsaken",
}

// DesertSubjects is a list of subjects suitable for deserts.
var DesertSubjects = []string{
	"desert",
	"wasteland",
	"sands",
	"barrens",
	"expanse",
	"region",
}

// DangerGenitivePhraseAdjDesert returns a danger genitive phrase adjective for deserts.
var DangerGenitivePhraseAdjDesert = []string{
	"burned",
	"scorched",
	"charred",
	"lost",
	"ashen",
}

// NewDesertGenerator returns a new generator for desert names.
func NewDesertGenerator(seed int64) *BasicGenerator {
	return NewBasicGenerator(seed,
		DesertAdjectives,
		DesertSubjects,
		WordPair{
			A: DangerGenitivePhraseAdjDesert,
			B: GenitivePhraseSubject,
		})
}
