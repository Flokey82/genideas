package gamestrategy

import "image/color"

type Cell struct {
	X     int
	Y     int
	Value float64
	*Type
	ControlledBy *Player
	Features     int64
}

// Cost returns the cost of occupying (or attacking) this cell.
func (c *Cell) Cost() float64 {
	return c.Type.Cost
}

// Yield returns the amount of resources this cell yields.
func (c *Cell) Yield() float64 {
	return resourceYield(c.Features) + c.BaseYield
}

// IsOccupied returns true if the cell is occupied by a player.
func (c *Cell) IsOccupied() bool {
	return c.ControlledBy != nil
}

func (c *Cell) Occupy(p *Player) bool {
	if c.IsOccupied() || p.Gold < c.Cost() {
		return false
	}
	c.ControlledBy = p
	return true
}

// Type represents the type of a cell.
// The type determines the cost of occupying the cell and the base yield, as well as the features
// that can be built on the cell.
// TODO: Add offensive and defensive modifiers and features.
type Type struct {
	Name            string  // Water, Meadow, Forest, Mountain, Desert...
	Cost            float64 // Occupation cost and multiplier for actions
	BaseYield       float64 // Base yield for the cell
	AllowedFeatures int64   // Bitmask of allowed features
	Color           color.Color
}

var (
	TypeCapital = Type{
		Name:      "Capital",
		Cost:      0.0,
		BaseYield: 10.0,
		Color: color.RGBA{
			R: 0xff,
			G: 0xff,
			B: 0xff,
			A: 0xff,
		},
	}
	TypeWater = Type{
		Name: "Water",
		Cost: 5.0,
		Color: color.RGBA{
			R: 0x00,
			G: 0x00,
			B: 0xff,
			A: 0xff,
		},
	}
	TypeMeadow = Type{
		Name:            "Meadow",
		Cost:            1.0,
		AllowedFeatures: FeatureFarm | FeatureSettlement,
		Color: color.RGBA{
			R: 0x00,
			G: 0xff,
			B: 0x00,
			A: 0xff,
		},
	}
	TypeForest = Type{
		Name:            "Forest",
		Cost:            2.0,
		AllowedFeatures: FeatureLumber,
		Color: color.RGBA{
			R: 0x00,
			G: 0x80,
			B: 0x00,
			A: 0xff,
		},
	}
	TypeMountain = Type{
		Name:            "Mountain",
		Cost:            3.0,
		AllowedFeatures: FeatureQuarry | FeatureMine | FeatureSettlement,
		Color: color.RGBA{
			R: 0x80,
			G: 0x80,
			B: 0x80,
			A: 0xff,
		},
	}
	TypeDesert = Type{
		Name: "Desert",
		Cost: 4.0,
		Color: color.RGBA{
			R: 0xff,
			G: 0xff,
			B: 0x00,
			A: 0xff,
		},
	}
)

const (
	FeatureNone = 0
	FeatureFarm = 1 << iota
	FeatureLumber
	FeatureQuarry
	FeatureMine
	FeatureSettlement
)

func resourceYield(f int64) float64 {
	var yield float64
	if f&FeatureFarm != 0 {
		yield += 2.0
	}
	if f&FeatureLumber != 0 {
		yield += 3.0
	}
	if f&FeatureQuarry != 0 {
		yield += 5.0
	}
	if f&FeatureMine != 0 {
		yield += 5.0
	}
	if f&FeatureSettlement != 0 {
		yield += 2.0
	}
	return yield
}

func costToBuild(f int64) float64 {
	var cost float64
	switch f {
	case FeatureFarm:
		cost = 1.0
	case FeatureLumber:
		cost = 1.0
	case FeatureQuarry:
		cost = 4.0
	case FeatureMine:
		cost = 4.0
	case FeatureSettlement:
		cost = 10.0
	}
	return cost
}

func (t *Type) CostToBuild(f int64) float64 {
	return t.Cost * costToBuild(f)
}

func splitFeatures(f int64) []int64 {
	var features []int64
	if f&FeatureFarm != 0 {
		features = append(features, FeatureFarm)
	}
	if f&FeatureLumber != 0 {
		features = append(features, FeatureLumber)
	}
	if f&FeatureQuarry != 0 {
		features = append(features, FeatureQuarry)
	}
	if f&FeatureMine != 0 {
		features = append(features, FeatureMine)
	}
	if f&FeatureSettlement != 0 {
		features = append(features, FeatureSettlement)
	}

	return features
}
