package genclutter

import (
	"log"
	"math/rand"
)

type Item struct {
	*ItemBase
	*Condition
	//Tags     []string // Associated tags ("astronomy", "metal working", etc.)
	Variant  *Variant
	Parent   *Item // Parent item (if any).
	Contains []*Item
}

func (it *Item) Log(prefix string) {
	if it == nil {
		return
	}
	log.Println(prefix + it.Name())
	for _, c := range it.Contains {
		c.Log(prefix + "  ")
	}
}

func (it *Item) Name() string {
	var suffix string
	if it.Variant != nil {
		suffix = it.Variant.Name
	}
	if it.Condition != nil && it.Condition != ConditionAverage {
		if suffix != "" {
			suffix += ", "
		}
		suffix += it.Condition.Name + " condition"
	}
	if suffix != "" {
		return it.ItemBase.Name + " (" + suffix + ")"
	}
	return it.ItemBase.Name
}

type ItemBase struct {
	Name          string
	Rarity        *Rarity
	Capacity      int
	Variants      []*Variant
	Required      []string
	Optional      []string
	OptionalMulti []string   // Optional items that can be added multiple times.
	Sets          []*ItemSet // Optional item sets.
}

func (ib *ItemBase) Roll() bool {
	if ib == nil {
		return false
	}
	if ib.Rarity == nil {
		return true
	}
	return ib.Rarity.Roll()
}

func (ib *ItemBase) Generate(parent *Item) *Item {
	if ib == nil {
		return nil
	}
	it := &Item{
		ItemBase:  ib,
		Condition: RollCondition(), // TODO: Allow inherited condition.
		Parent:    parent,
	}
	if len(ib.Variants) > 0 {
		for _, i := range rand.Perm(len(ib.Variants)) {
			if ib.Variants[i].Rarity.Roll() {
				it.Variant = ib.Variants[i]
				break
			}
		}
	}

	// Generate a number of items from the set up to the capacity.
	var contains []*Item
	if ib.Capacity > 0 {
		cap := ib.Capacity
		for _, i := range rand.Perm(len(ib.Optional)) {
			if cit := nameToItemBase[ib.Optional[i]]; cit.Roll() {
				contains = append(contains, cit.Generate(it))
				cap--
			}
			if cap <= 0 {
				break
			}
		}
		if cap > 0 {
			for _, i := range rand.Perm(len(ib.OptionalMulti)) {
				toGen := rand.Intn(cap)
				for j := 0; j < toGen; j++ {
					if cit := nameToItemBase[ib.OptionalMulti[i]]; cit.Roll() {
						contains = append(contains, cit.Generate(it))
						cap--
					}
				}
				if cap <= 0 {
					break
				}
			}
		}
		if cap > 0 {
			for _, set := range ib.Sets {
				toGen := rand.Intn(cap)
				contains = append(contains, set.GenerateN(cap, it)...)
				cap -= toGen
				if cap <= 0 {
					break
				}
			}
		}
	}

	// Add all required items.
	for _, req := range ib.Required {
		contains = append(contains, nameToItemBase[req].Generate(it))
	}
	it.Contains = contains
	return it
}

type ItemSet struct {
	Name  string
	Items []*ItemBase
}

func (is *ItemSet) GenerateN(n int, parent *Item) []*Item {
	if len(is.Items) == 0 {
		return nil
	}
	var res []*Item
	for len(res) < n {
		res = append(res, is.Generate(parent))
	}
	return res
}

func (is *ItemSet) Generate(parent *Item) *Item {
	if len(is.Items) == 0 {
		return nil
	}
	for _, i := range rand.Perm(len(is.Items)) {
		item := is.Items[i]
		if item.Roll() {
			return item.Generate(parent)
		}
	}
	return is.Items[rand.Intn(len(is.Items))].Generate(parent)
}

type Variant struct {
	Name   string
	Rarity *Rarity
	// TODO: If we have a context tag that matches this variant,
	// we can increase the odds of this variant being picked.
	// Tags []string // Preferred tags ("astronomy", "metal working", etc.) for this variant.
}

func NewVariant(name string, rarity *Rarity) *Variant {
	return &Variant{
		Name:   name,
		Rarity: rarity,
	}
}
