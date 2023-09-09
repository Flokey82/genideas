package simpeople2

import (
	"fmt"
	"log"
	"math"
)

// MotiveTypeSleep is the motive for sleep.
var MotiveTypeSleep = &MotiveType{
	Name:  "Sleep",
	Curve: CurveTypeExponential,
	Decay: 5.0,
	OnMax: func() {},
	OnMin: func() {
		log.Println("Person is very tired!")
	},
}

// MotiveTypeFood is the motive for food.
var MotiveTypeFood = &MotiveType{
	Name:  "Food",
	Curve: CurveTypeSigmoid,
	Decay: 15.0,
	OnMax: func() {},
	OnMin: func() {
		log.Println("Person is starving to death!")
	},
}

// MotiveTypeFun is the motive for fun stuff.
var MotiveTypeFun = &MotiveType{
	Name:  "Fun",
	Curve: CurveTypeParabola,
	Decay: 5.0,
	OnMax: func() {},
	OnMin: func() {
		log.Println("Person is very bored!")
	},
}

// MotiveTypeBladder is the motive for bladder.
// TODO: If bladder is too high, the person should pee themselves,
// which should decay hygiene but restore bladder.
var MotiveTypeBladder = &MotiveType{
	Name:  "Bladder",
	Curve: CurveTypeExponential,
	Decay: 15.0,
	OnMax: func() {},
	OnMin: func() {
		log.Println("Person is about to pee themselves!")
	},
}

// MotiveHygiene is the motive for hygiene.
var MotiveHygiene = &MotiveType{
	Name:  "Hygiene",
	Curve: CurveTypeExponential,
	Decay: 5.0,
	OnMax: func() {},
	OnMin: func() {
		log.Println("Person is very stinky!")
	},
}

// MotiveType is a motive for a person.
type MotiveType struct {
	Name  string
	Curve CurveType // How the muliplier changes based on the current value
	Decay float64   // How much the value decays per tick
	OnMax func()    // Called when the motive reaches the maximum value
	OnMin func()    // Called when the motive reaches the minimum value
}

// New creates a new motive.
func (m *MotiveType) New() *Motive {
	return &Motive{
		Type: m,
		Val:  100,
	}
}

const (
	minMotiveValue = -100.0
	maxMotiveValue = 100.0
)

// Motive is a motive for a person.
type Motive struct {
	Type *MotiveType
	Val  float64 // Current value of the motive
}

// String returns a string representation of the motive.
func (m *Motive) String() string {
	return m.Type.Name + ": " + fmt.Sprintf("%.2f", m.Val)
}

// Tick decays the motive value.
func (m *Motive) Tick(elapsed float64) {
	m.Change(-m.Type.Decay * elapsed)
	if m.Val >= maxMotiveValue {
		m.Type.OnMax()
	} else if m.Val <= minMotiveValue {
		m.Type.OnMin()
	}
}

// Change changes the motive value by the given amount.
func (m *Motive) Change(amount float64) {
	m.Val += amount
	if m.Val > maxMotiveValue {
		m.Val = maxMotiveValue
	} else if m.Val < minMotiveValue {
		m.Val = minMotiveValue
	}
}

// MissingToMax returns how much the motive is missing to reach the maximum value.
func (m *Motive) MissingToMax() float64 {
	return maxMotiveValue - m.Val
}

// Multiplier returns the current multiplier for the motive.
func (m *Motive) Multiplier() float64 {
	return m.Type.Curve.Multiplier(m.Val)
}

// Log logs the current value (and multiplier) of the motive.
func (m *Motive) Log() {
	log.Printf("%s: %.2f (%.2f)", m.Type.Name, m.Val, m.Multiplier())
}

// CurveType is a type of curve.
type CurveType int

const (
	CurveTypeLinear CurveType = iota
	CurveTypeExponential
	CurveTypeSigmoid
	CurveTypeParabola
)

// Multiplier returns the multiplier for the given value.
// NOTE: The input value is between -100 and 100 and the output is between 0 and 10.
func (c CurveType) Multiplier(val float64) float64 {
	// Scale the value to be between 1 and 0 (if the original value is -100, the scaled value is 1)
	val = 1 - (val+100)/200
	var mul float64
	switch c {
	case CurveTypeLinear:
		// 1
		// |  /
		// 0 /  1
		mul = val
	case CurveTypeExponential:
		// 1
		// |   /
		// 0 -  1
		mul = val * val
	case CurveTypeSigmoid:
		// Proper sigmoid function: 1 / (1 + e^-x)
		// 1     -
		// |   /
		// 0 -    1
		mul = 1 / (1 + math.Pow(math.E, -8*val+4))
	case CurveTypeParabola:
		// 1
		// | \   /
		// 0   -  1
		mul = 4*val*val - 4*val + 1
	}
	return mul * 10
}
