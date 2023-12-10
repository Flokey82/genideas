package genlandmarknames

// MountainAdjective is a list of prefixes suitable for mountains.
var MountainAdjective = []string{
	"rocky",
	"mountainous",
	"spiked",
	"steep",
	"rough",
	"craggy",
	"craggly",
	"toothy",
	"jagged",
	"broken",
	"crumbling",
}

// MountainRangeSubjects is a list of subjects for mountain ranges.
var MountainRangeSubjects = []string{
	"mountains",
	"teeth",
	"spikes",
	"peaks",
	"rocks",
	"thorns",
	"jags",
	"spurs",
	"spires",
	"pinnacles",
	"stones",
	"backs",
	"needles",
	"nails",
	"knives",
	"scythes",
	"cliffs",
	"ridges",
	"sickles",
	"spikes",
	"spurs",
	"spires",
	"swords",
	"edges",
	"pikes",
}

// DangerGenitivePhraseAdjMountain is a list of prefixes suitable for dangerous mountains.
var DangerGenitivePhraseAdjMountain = []string{
	"shattered",
	"lost",
	"petrified",
	"forgotten",
	"abandoned",
	"broken",
	"buried",
	"fallen",
	"crushed",
	"sundered",
	"sharpened",
	"sharp",
	"cutting",
}

// NewMountainRangeGenerator returns a new generator for mountain range names.
func NewMountainRangeGenerator(seed int64) *BasicGenerator {
	return NewBasicGenerator(seed,
		MountainAdjective,
		MountainRangeSubjects,
		WordPair{
			A: DangerGenitivePhraseAdjMountain,
			B: DangerGenitivePhraseSubject,
		})
}

// MountainSubjects is a list of subjects for mountains.
var MountainSubjects = []string{
	"mountain",
	"tooth",
	"spike",
	"peak",
	"rock",
	"thorn",
	"jag",
	"spur",
	"spire",
	"pinnacle",
	"stone",
	"back",
	"ridge",
	"crest",
	"summit",
	"needle",
	"nail",
	"knife",
	"scythe",
	"cliff",
	"sickle",
	"spike",
	"spur",
	"spire",
	"sword",
	"edge",
	"pike",
}

// NewMountainGenerator returns a new generator for mountain names.
func NewMountainGenerator(seed int64) *BasicGenerator {
	return NewBasicGenerator(seed,
		MountainAdjective,
		MountainSubjects,
		WordPair{
			A: DangerGenitivePhraseAdjMountain,
			B: DangerGenitivePhraseSubject,
		})
}
