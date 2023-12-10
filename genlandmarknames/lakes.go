package genlandmarknames

import "github.com/Flokey82/go_gens/genstory"

const (
	TokenAdj        = "[ADJ]"
	TokenSubject    = "[SUBJECT]"
	TokenAdjGen     = "[ADJ_GEN]"     // Used for genitive phrases
	TokenSubjectGen = "[SUBJECT_GEN]" // Used for genitive phrases
	TokenPlace      = "[PLACE]"
)

var LakeTemplates = []string{
	"The [ADJ] [SUBJECT]",
	"The [ADJ] [SUBJECT] of the [ADJ_GEN] [SUBJECT_GEN]",
	"The [ADJ] [SUBJECT] of [ADJ_GEN] [SUBJECT_GEN]",
	"The [ADJ] [SUBJECT] of [SUBJECT_GEN]",
	"[PLACE] [SUBJECT]",
	"[SUBJECT] of [PLACE]",
}

var LakeNameConfig = &genstory.TextConfig{
	TokenPools: map[string][]string{
		TokenAdj:        LakeAdjectives,
		TokenSubject:    LakeSubjects,
		TokenAdjGen:     GenitivePhraseAdjNegWaterStill,
		TokenSubjectGen: GenitivePhraseSubject,
	},
	TokenIsMandatory: map[string]bool{},
	Tokens:           []string{TokenAdj, TokenSubject, TokenAdjGen, TokenSubjectGen, TokenPlace},
	Templates:        LakeTemplates,
	UseAllProvided:   true,
	UseAlliteration:  false,
	Title:            true,
}

var LakeAdjectives = []string{
	"silver",
	"golden",
	"crystal",
	"deep",
	"clear",
	"blue",
	"green",
}

var LakeSubjects = []string{
	"lake",
	"pond",
	"pool",
	"lagoon",
	"loch",
	"mere",
	"tarn",
	"reservoir",
	"basin",
	"bowl",
	"bath",
	"mirror",
	"puddle",
	"waters",
	"water",
	"pit",
}

var LakeDangerousSuffix = WordPair{
	A: GenitivePhraseAdjNegWaterStill,
	B: GenitivePhraseSubject,
}

// NewLakeGenerator returns a new generator for lake names.
func NewLakeGenerator(seed int64) *BasicGenerator {
	return NewBasicGenerator(seed, LakeAdjectives, LakeSubjects, LakeDangerousSuffix)
}
