// Package genclutter provides means of randomly generating clutter and furnishings for a room
// based on a seed value, a profession that is associated with the room, or the room's purpose.
package genclutter

func GenRoom(name string) *Item {
	return nameToItemBase[name].Generate(nil)
}

var nameToItemBase = map[string]*ItemBase{
	"bedroom": &ItemBase{
		Name:     "bedroom",
		Capacity: 5,
		Rarity:   RarityCommon,
		Required: []string{
			"bed",
		},
		Optional: []string{
			"nightstand",
			"bookcase",
			"chair",
			"carpet",
			"desk",
		},
	},
	"desk": &ItemBase{
		Name:     "desk",
		Capacity: 5,
		Rarity:   RarityCommon,
		OptionalMulti: []string{
			"drawer",
		},
	},
	"drawer": &ItemBase{
		Name:     "drawer",
		Capacity: 5,
		Rarity:   RarityCommon,
		Variants: []*Variant{
			NewVariant("regular", RarityCommon),
			NewVariant("secret", RarityRare),
		},
		OptionalMulti: []string{
			"clothing",
			"book",
			"letter",
			"map",
		},
	},
	"bed": &ItemBase{
		Name:   "bed",
		Rarity: RarityCommon,
		Variants: []*Variant{
			NewVariant("single", RarityCommon),
			NewVariant("double", RarityCommon),
			NewVariant("queen", RarityCommon),
			NewVariant("king", RarityCommon),
			NewVariant("bedroll", RarityCommon),
			NewVariant("cot", RarityCommon),
			NewVariant("rag pile", RarityCommon),
		},
	},
	"bench": &ItemBase{
		Name:   "bench",
		Rarity: RarityCommon,
		Variants: []*Variant{
			NewVariant("wood", RarityCommon),
			NewVariant("cushioned", RarityCommon),
			NewVariant("stone", RarityCommon),
		},
	},
	"bookcase": &ItemBase{
		Name:          "bookcase",
		Rarity:        RarityCommon,
		Capacity:      15,
		Required:      []string{"book"},
		OptionalMulti: []string{"book"},
		Variants: []*Variant{
			NewVariant("wood", RarityCommon),
			NewVariant("metal", RarityCommon),
			NewVariant("bone", RarityRare),
		},
	},
	"book": &ItemBase{
		Name:     "book",
		Rarity:   RarityCommon,
		Capacity: 1,
		Variants: []*Variant{
			NewVariant("book", RarityCommon),
			NewVariant("spellbook", RarityRare),
			NewVariant("journal", RarityCommon),
			NewVariant("ledger", RarityCommon),
			NewVariant("novel", RarityCommon),
			NewVariant("poetry", RarityCommon),
			NewVariant("manual", RarityCommon),
			NewVariant("historical account", RarityCommon),
		},
		Optional: []string{
			"secret note",
			"map",
			"letter",
		},
	},
	"cabinet": &ItemBase{
		Name:     "cabinet",
		Required: []string{"container"},
	},
	"carpet": &ItemBase{
		Name: "carpet",
	},
	"chair": &ItemBase{
		Name: "chair",
		Variants: []*Variant{
			NewVariant("chair", RarityCommon),
			NewVariant("armchair", RarityCommon),
			NewVariant("wood", RarityCommon),
			NewVariant("cushioned", RarityCommon),
			NewVariant("stone", RarityCommon),
			NewVariant("stool", RarityCommon),
		},
	},
	"secret note": &ItemBase{
		Name:   "secret note",
		Rarity: RarityRare,
	},
	"map": &ItemBase{
		Name:   "map",
		Rarity: RarityRare,
	},
	"letter": &ItemBase{
		Name: "letter",
		Variants: []*Variant{
			NewVariant("letter", RarityCommon),
			NewVariant("recommendation", RarityCommon),
			NewVariant("invitation", RarityCommon),
			NewVariant("notice", RarityCommon),
			NewVariant("bill", RarityCommon),
			NewVariant("contract", RarityCommon),
			NewVariant("deed", RarityRare),
			NewVariant("will", RarityUncommon),
		},
		Rarity: RarityRare,
	},
	"container": &ItemBase{
		Name: "container",
		Variants: []*Variant{
			NewVariant("chest", RarityCommon), // TODO: Add capacity.
			NewVariant("barrel", RarityCommon),
			NewVariant("crate", RarityCommon),
			NewVariant("jar", RarityCommon),
			NewVariant("urn", RarityUncommon),
			NewVariant("box", RarityCommon),
			NewVariant("basket", RarityCommon),
			NewVariant("bag", RarityCommon),
			NewVariant("trunk", RarityCommon),
			NewVariant("coffer", RarityCommon),
			NewVariant("locker", RarityUncommon),
			NewVariant("safe", RarityUncommon),
			NewVariant("vault", RarityUncommon),
		},
	},
}
