package genlandmarknames

// FertileLandSuffix is a list of suffixes suitable for fertile land.
var FertileLandPrefix = []string{
	"green",
	"lush",
	"bountiful",
	"fruitful",
	"rich",
	"abundant",
	"fertile",
	"verdant",
	"flourishing",
}

// LargeAreaSuffix is a list of suffixes suitable for large areas.
var LargeAreaSuffix = []string{
	"land",
	"plains",
	"expanse",
	"region",
	"territory",
	"wilderness",
	"wilds",
}

// DangerousTerrainDescriptor returns a random descriptor for a dangerous terrain.
// This suffix is supposed to be used with "... of ...".
var DangerousTerrainDescriptor = []string{
	"death",
	"doom",
	"despair",
	"darkness",
	"evil",
	"chaos",
	"madness",
	"loss",
	"pain",
	"anguish",
	"terror",
	"horror",
	"lost souls",
	"the dead",
	"the damned",
	"the cursed",
	"the forsaken",
	"the lost",
	"the forgotten",
	"the abandoned",
	"the unknown",
	"the doomed",
}

func NewPlainsGenerator(seed int64) *BasicGenerator {
	return NewBasicGenerator(seed,
		FertileLandPrefix,
		LargeAreaSuffix,
		WordPair{
			A: GenitivePhraseAdjNeg,
			B: GenitivePhraseSubject,
		})
}
