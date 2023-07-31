package genstatblock5e

type Size int

const (
	Tiny Size = iota
	Small
	Medium
	Large
	Huge
	Gargantuan
	SizeMax
)

// String returns the string representation of the size.
func (s Size) String() string {
	switch s {
	case Tiny:
		return "tiny"
	case Small:
		return "small"
	case Medium:
		return "medium"
	case Large:
		return "large"
	case Huge:
		return "huge"
	case Gargantuan:
		return "gargantuan"
	default:
		return "unknown"
	}
}

// BaseArmorClass returns the base armor class for the given size.
func (s Size) BaseArmorClass() int {
	switch s {
	case Tiny:
		return 12
	case Small:
		return 13
	case Medium:
		return 13
	case Large:
		return 12
	case Huge:
		return 10
	case Gargantuan:
		return 10
	default:
		return 10
	}
}

type Alignment [2]int

const (
	AlignChaotic = iota
	AlignNeutral
	AlignLawful
	AlignGood
	AlignEvil
)

// String returns the string representation of the alignment.
func (a Alignment) String() string {
	switch {
	case a[0] == AlignNeutral && a[1] == AlignNeutral:
		return "true neutral"
	case a[0] == AlignChaotic && a[1] == AlignGood:
		return "chaotic good"
	case a[0] == AlignChaotic && a[1] == AlignNeutral:
		return "chaotic neutral"
	case a[0] == AlignChaotic && a[1] == AlignEvil:
		return "chaotic evil"
	case a[0] == AlignNeutral && a[1] == AlignGood:
		return "neutral good"
	case a[0] == AlignNeutral && a[1] == AlignEvil:
		return "neutral evil"
	case a[0] == AlignLawful && a[1] == AlignGood:
		return "lawful good"
	case a[0] == AlignLawful && a[1] == AlignNeutral:
		return "lawful neutral"
	case a[0] == AlignLawful && a[1] == AlignEvil:
		return "lawful evil"
	default:
		return "unknown"
	}
}

type MonsterType string

const (
	TypeAberration  MonsterType = "aberration"
	TypeBeast       MonsterType = "beast"
	TypeCelestial   MonsterType = "celestial"
	TypeConstruct   MonsterType = "construct"
	TypeDragon      MonsterType = "dragon"
	TypeUndead      MonsterType = "undead"
	TypeElemental   MonsterType = "elemental"
	TypeFiend       MonsterType = "fiend"
	TypeFey         MonsterType = "fey"
	TypeGiant       MonsterType = "giant"
	TypeHumanoid    MonsterType = "humanoid"
	TypeMonstrosity MonsterType = "monstrosity"
	TypeOoze        MonsterType = "ooze"
	TypePlant       MonsterType = "plant"
)

var MonsterTypes = []MonsterType{
	TypeAberration,
	TypeBeast,
	TypeCelestial,
	TypeConstruct,
	TypeDragon,
	TypeUndead,
	TypeElemental,
	TypeFiend,
	TypeFey,
	TypeGiant,
	TypeHumanoid,
	TypeMonstrosity,
	TypeOoze,
	TypePlant,
}

type Sense string

const (
	SenseBlindsight        Sense = "blindsight"
	SenseDarkvision        Sense = "darkvision"
	SenseTremorsense       Sense = "tremorsense"
	SenseTruesight         Sense = "truesight"
	SensePassivePerception Sense = "passive perception"
)

type AbilityScore int

const (
	AbilityStrength AbilityScore = iota
	AbilityDexterity
	AbilityConstitution
	AbilityIntelligence
	AbilityWisdom
	AbilityCharisma
	AbilityScoreCount
)

// CalcAbilityModifier calculates the ability modifier for the given ability score.
// The modifier is calculated as (score - 10) / 2.
func CalcAbilityModifier(score int) int {
	return (score - 10) / 2
}

type Skill struct {
	Label        string       // The name of the skill
	AbilityScore AbilityScore // The ability score associated with the skill
}

var (
	SkillAthletics   = Skill{Label: "Athletics", AbilityScore: AbilityStrength}
	SkillAcrobatics  = Skill{Label: "Acrobatics", AbilityScore: AbilityDexterity}
	SkillSleight     = Skill{Label: "Sleight of Hand", AbilityScore: AbilityDexterity}
	SkillStealth     = Skill{Label: "Stealth", AbilityScore: AbilityDexterity}
	SkillArcana      = Skill{Label: "Arcana", AbilityScore: AbilityIntelligence}
	SkillHistory     = Skill{Label: "History", AbilityScore: AbilityIntelligence}
	SkillInvest      = Skill{Label: "Investigation", AbilityScore: AbilityIntelligence}
	SkillNature      = Skill{Label: "Nature", AbilityScore: AbilityIntelligence}
	SkillReligion    = Skill{Label: "Religion", AbilityScore: AbilityIntelligence}
	SkillAnimal      = Skill{Label: "Animal Handling", AbilityScore: AbilityWisdom}
	SkillInsight     = Skill{Label: "Insight", AbilityScore: AbilityWisdom}
	SkillMedicine    = Skill{Label: "Medicine", AbilityScore: AbilityWisdom}
	SkillPerception  = Skill{Label: "Perception", AbilityScore: AbilityWisdom}
	SkillSurvival    = Skill{Label: "Survival", AbilityScore: AbilityWisdom}
	SkillDeception   = Skill{Label: "Deception", AbilityScore: AbilityCharisma}
	SkillIntimidate  = Skill{Label: "Intimidation", AbilityScore: AbilityCharisma}
	SkillPerformance = Skill{Label: "Performance", AbilityScore: AbilityCharisma}
	SkillPersuasion  = Skill{Label: "Persuasion", AbilityScore: AbilityCharisma}
)

var Skills = []Skill{
	SkillAthletics,
	SkillAcrobatics,
	SkillSleight,
	SkillStealth,
	SkillArcana,
	SkillHistory,
	SkillInvest,
	SkillNature,
	SkillReligion,
	SkillAnimal,
	SkillInsight,
	SkillMedicine,
	SkillPerception,
	SkillSurvival,
	SkillDeception,
	SkillIntimidate,
	SkillPerformance,
	SkillPersuasion,
}

type DamageType string

const (
	DamageTypeAcid     DamageType = "acid"
	DamageTypeBludgeon DamageType = "bludgeoning"
	DamageTypeCold     DamageType = "cold"
	DamageTypeFire     DamageType = "fire"
	DamageTypeForce    DamageType = "force"
	DamageTypeLight    DamageType = "lightning"
	DamageTypeNecrotic DamageType = "necrotic"
	DamageTypePiercing DamageType = "piercing"
	DamageTypePoison   DamageType = "poison"
	DamageTypePsychic  DamageType = "psychic"
	DamageTypeRadiant  DamageType = "radiant"
	DamageTypeSlashing DamageType = "slashing"
	DamageTypeThunder  DamageType = "thunder"
	DamageTypeMag      DamageType = "magical"
	DamageTypeNonMag   DamageType = "nonmagical"
)

var DamageTypes = []DamageType{
	DamageTypeAcid,
	DamageTypeBludgeon,
	DamageTypeCold,
	DamageTypeFire,
	DamageTypeForce,
	DamageTypeLight,
	DamageTypeNecrotic,
	DamageTypePiercing,
	DamageTypePoison,
	DamageTypePsychic,
	DamageTypeRadiant,
	DamageTypeSlashing,
	DamageTypeThunder,
	DamageTypeMag,
	DamageTypeNonMag,
}

type Language string

const (
	LanguageCommon Language = "common"
	LanguageDwarv  Language = "dwarvish"
	LanguageElvish Language = "elvish"
	LanguageGiant  Language = "giant"
	LanguageGnom   Language = "gnomish"
	LanguageGob    Language = "goblin"
	LanguageHalfl  Language = "halfling"
	LanguageOrc    Language = "orc"
	LanguageAby    Language = "abyssal"
	LanguageCele   Language = "celestial"
	LanguageDra    Language = "draconic"
	LanguageDeep   Language = "deep speech"
	LanguageInf    Language = "infernal"
	LanguagePrim   Language = "primordial"
	LanguageSylv   Language = "sylvan"
	LanguageUnder  Language = "undercommon"
)

var Languages = []Language{
	LanguageCommon,
	LanguageDwarv,
	LanguageElvish,
	LanguageGiant,
	LanguageGnom,
	LanguageGob,
	LanguageHalfl,
	LanguageOrc,
	LanguageAby,
	LanguageCele,
	LanguageDra,
	LanguageDeep,
	LanguageInf,
	LanguagePrim,
	LanguageSylv,
	LanguageUnder,
}

// CalcArmorClass calculates the armor class for the given properties.
// We use size, dexterity and armor to calculate the armor class.
//
// NOTE: This is not official, I just pulled it out of my rear end.
//
// Example:
// - Size: medium (base armor class 13)
// - Dexterity: 14 (+2)
// - Armor: 2 (chain mail)
// - Armor class: 17
func CalcArmorClass(size Size, dex, armor int) int {
	// First we use the size to calculate the base armor class.
	ac := size.BaseArmorClass()

	// Then we add the dexterity modifier.
	ac += CalcAbilityModifier(dex)

	// Finally we add the armor.
	ac += armor

	return ac
}
