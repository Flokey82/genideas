package simpeople2

import "github.com/Flokey82/go_gens/vectors"

// ObjectTypeFridge is the object type of a fridge.
var ObjectTypeFridge = &ObjectType{
	Name:     "Fridge",
	SpriteID: 146,
	Actions: []*Action{
		ActionEat,
	},
}

// ObjectTypeBed is the object type of a bed.
var ObjectTypeBed = &ObjectType{
	Name:     "Bed",
	SpriteID: 186,
	Actions: []*Action{
		ActionSleep,
	},
}

// ObjectTypeCouch is the object type of a couch.
var ObjectTypeCouch = &ObjectType{
	Name:     "Couch",
	SpriteID: 190,
	Actions: []*Action{
		ActionWatchTV,
	},
}

// ObjectTypeToilet is the object type of a toilet.
var ObjectTypeToilet = &ObjectType{
	Name:     "Toilet",
	SpriteID: 27,
	Actions: []*Action{
		ActionPee,
	},
}

// ObjectTypeShower is the object type of a shower.
var ObjectTypeShower = &ObjectType{
	Name:     "Shower",
	SpriteID: 425,
	Actions: []*Action{
		ActionShower,
	},
}

// ObjectType is an object in the simulation that can advertise available actions and their effects.
type ObjectType struct {
	Name     string
	SpriteID int
	Actions  []*Action
}

// New creates a new object.
func (ot *ObjectType) New(pos vectors.Vec2) *Object {
	return &Object{
		ObjectType: ot,
		Position:   pos,
	}
}

// Object is an object in the simulation.
type Object struct {
	*ObjectType
	Position vectors.Vec2
}

// Action is an action that can be performed on an object.
type Action struct {
	Name       string  // The name of the action
	Effect     *Effect // Primary effect of the action
	SideEffect *Effect // Secondary effect of the action
}

// Effect is the effect of an action.
type Effect struct {
	Motive *MotiveType // The motive that is affected by the action
	Effect float64     // How much the motive changes when the action is performed
}

// ActionEat is the action of eating.
var ActionEat = &Action{
	Name: "Eat",
	Effect: &Effect{
		Motive: MotiveTypeFood,
		Effect: 100.0,
	},
	SideEffect: &Effect{
		Motive: MotiveTypeBladder,
		Effect: -20.0,
	},
}

// ActionSleep is the action of sleeping.
var ActionSleep = &Action{
	Name: "Sleep",
	Effect: &Effect{
		Motive: MotiveTypeSleep,
		Effect: 50.0,
	},
	SideEffect: &Effect{ // Sleeping makes you less hungry.
		Motive: MotiveTypeFood,
		Effect: -4.0,
	},
}

// ActionWatchTV is the action of watching TV.
var ActionWatchTV = &Action{
	Name: "Watch TV",
	Effect: &Effect{
		Motive: MotiveTypeFun,
		Effect: 40.0,
	},
}

// ActionPee is the action of peeing.
var ActionPee = &Action{
	Name: "Pee",
	Effect: &Effect{
		Motive: MotiveTypeBladder,
		Effect: 100.0,
	},
}

// ActionShower is the action of showering.
var ActionShower = &Action{
	Name: "Shower",
	Effect: &Effect{
		Motive: MotiveHygiene,
		Effect: 100.0,
	},
	SideEffect: &Effect{
		Motive: MotiveTypeBladder,
		Effect: -10.0,
	},
}
