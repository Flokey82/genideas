package simpeople2

import "github.com/Flokey82/go_gens/vectors"

// ObjectTypeFridge is the object type of a fridge.
var ObjectTypeFridge = &ObjectType{
	Name: "Fridge",
	Actions: []*Action{
		ActionEat,
	},
}

// ObjectTypeBed is the object type of a bed.
var ObjectTypeBed = &ObjectType{
	Name: "Bed",
	Actions: []*Action{
		ActionSleep,
	},
}

// ObjectTypeTV is the object type of a TV.
var ObjectTypeTV = &ObjectType{
	Name: "TV",
	Actions: []*Action{
		ActionWatchTV,
	},
}

// ObjectType is an object in the simulation that can advertise available actions and their effects.
type ObjectType struct {
	Name    string
	Actions []*Action
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
	Name   string      // The name of the action
	Motive *MotiveType // The motive that is affected by the action
	Effect float64     // How much the motive changes when the action is performed
}

// ActionEat is the action of eating.
var ActionEat = &Action{
	Name:   "Eat",
	Motive: MotiveTypeFood,
	Effect: 30.0,
}

// ActionSleep is the action of sleeping.
var ActionSleep = &Action{
	Name:   "Sleep",
	Motive: MotiveTypeSleep,
	Effect: 30.0,
}

// ActionWatchTV is the action of watching TV.
var ActionWatchTV = &Action{
	Name:   "Watch TV",
	Motive: MotiveTypeFun,
	Effect: 30.0,
}
