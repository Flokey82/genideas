package simsettlers

import "math"

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
