package genshops

import "github.com/Flokey82/go_gens/genstory"

/*
Examples:
The best biscuits in the city
Great armor makes great warriors
Choice weapons for the discerning adventurer
The worse it is, the greater the fun
Good food, good friends, good times
Terrible wares for terrible people
*/

var sloganTemplates = []string{
	"[ADJ] [PRODUCT] for [ADJ] [PEOPLE]",
	"[ADJ_POSITIVE] [PRODUCT] for [ADJ_NEGATIVE] [PEOPLE]",
	"[ADJ_NEGATIVE] [PRODUCT] for [ADJ_NEGATIVE] [PEOPLE]",
	"If you want [ADJ_POSITIVE] [PRODUCT], you need [ADJ_POSITIVE] [PEOPLE]",
}

const (
	TokenAdjPositive = "[ADJ_POSITIVE]"
	TokenAdjNegative = "[ADJ_NEGATIVE]"
	TokenPeople      = "[PEOPLE]"
)

var SloganConfig = &genstory.TextConfig{
	TokenPools: map[string][]string{
		TokenAdjPositive: adjectivesPositive,
		TokenAdjNegative: adjectivesNegative,
		TokenPeople:      people,
		TokenProduct:     products,
		TokenAdj:         adjectives,
	},
	TokenIsMandatory: map[string]bool{},
	Tokens: []string{
		TokenAdjPositive,
		TokenAdjNegative,
		TokenPeople,
		TokenProduct,
		TokenAdj,
	},
	Templates:       sloganTemplates,
	UseAllProvided:  true,
	UseAlliteration: true,
}

var adjectivesPositive = []string{
	"amazing",
	"awesome",
	"brilliant",
	"choice",
	"deluxe",
	"excellent",
	"exceptional",
	"fantastic",
	"fine",
	"first-class",
	"first-rate",
	"good",
	"great",
	"high-class",
	"high-end",
	"high-grade",
	"high-quality",
	"luxury",
	"magnificent",
	"marvellous",
	"outstanding",
	"premium",
	"prime",
	"quality",
	"select",
	"super",
	"superb",
	"superior",
	"supreme",
	"terrific",
	"top-grade",
	"top-notch",
	"top-quality",
}

var adjectivesNegative = []string{
	"awful",
	"bad",
	"cheap",
	"crappy",
	"crummy",
	"horrible",
	"junky",
	"lousy",
	"low-class",
	"low-end",
	"low-grade",
	"low-quality",
	"poor",
	"rubbish",
	"second-class",
	"second-rate",
	"shabby",
	"shitty",
	"shoddy",
	"substandard",
	"tacky",
	"tasteless",
	"terrible",
	"trashy",
	"unacceptable",
	"unpleasant",
	"unsatisfactory",
	"unsatisfying",
	"unsavoury",
}

var people = []string{
	"adventurers",
	"arseholes",
	"bandits",
	"bravehearts",
	"brigands",
	"champions",
	"criminals",
	"explorers",
	"fighters",
	"followers",
	"fools",
	"friends",
	"goons",
	"guards",
	"heroes",
	"knights",
	"leaders",
	"low-lives",
	"people",
	"pirates",
	"scoundrels",
	"scum",
	"seekers",
	"souls",
	"thieves",
	"travellers",
	"userpers",
	"undesirables",
	"vagabonds",
	"villains",
	"warriors",
}
