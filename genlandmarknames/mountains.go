package genlandmarknames

// MountainPrefix is a list of prefixes suitable for mountains.
var MountainPrefix = []string{
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

// MountainDangerPrefix is a list of prefixes suitable for dangerous mountains.
var MountainDangerPrefix = []string{
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
	return NewBasicGenerator(seed, MountainPrefix, []string{
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
	}, WordPair{
		A: MountainDangerPrefix,
		B: DangerousSuffixB,
	})
}

// NewMountainGenerator returns a new generator for mountain names.
func NewMountainGenerator(seed int64) *BasicGenerator {
	return NewBasicGenerator(seed, MountainPrefix, []string{
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
	}, WordPair{
		A: MountainDangerPrefix,
		B: DangerousSuffixB,
	})
}
