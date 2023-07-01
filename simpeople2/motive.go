package simpeople2

import (
	"log"
	"math"
)

// MotiveTypeSleep is the motive for sleep.
var MotiveTypeSleep = &MotiveType{
	Name:  "Sleep",
	Curve: CurveTypeExponential,
	Decay: 10.0,
}

// MotiveTypeFood is the motive for food.
var MotiveTypeFood = &MotiveType{
	Name:  "Food",
	Curve: CurveTypeExponential,
	Decay: 15.0,
}

// MotiveTypeFun is the motive for fun stuff.
var MotiveTypeFun = &MotiveType{
	Name:  "Fun",
	Curve: CurveTypeParabola,
	Decay: 5.0,
}

// MotiveType is a motive for a person.
type MotiveType struct {
	Name  string
	Curve CurveType // How the muliplier changes based on the current value
	Decay float64   // How much the value decays per tick
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

// Tick decays the motive value.
func (m *Motive) Tick() {
	m.Change(-m.Type.Decay)
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
	CurveTypeLogarithmic
	CurveTypeSigmoid
	CurveTypeParabola
)

// Multiplier returns the multiplier for the given value.
// NOTE: The input value is between -100 and 100 and the output is between 0 and 10.
func (c CurveType) Multiplier(val float64) float64 {
	// Scale the value to be between 1 and 0 (if the original value is -100, the scaled value is 1)
	val = 1 - (val+100)/200
	switch c {
	case CurveTypeLinear:
		return val * 10
	case CurveTypeExponential:
		return val * val * 10
	case CurveTypeLogarithmic:
		return -2 * 10 * math.Log(val)
	case CurveTypeSigmoid:
		return 10 / (1 + 9*math.Exp(-val))
	case CurveTypeParabola:
		return 4*val*val - 4*val + 1
	}
	return 0
}
