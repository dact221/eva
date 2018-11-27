package dist

import "math"

// WeibullMax represents a Weibull Maximum distribution.
type WeibullMax struct {
	Loc   float64
	Scale float64
	Shape float64
}

// CDF computes the value of the cumulative density function at x.
func (w WeibullMax) CDF(x float64) float64 {
	if x > w.Loc {
		return 1
	}
	return math.Exp(-math.Pow((w.Loc-x)/w.Scale, w.Shape))
}

// GetParams returns parameters map.
func (w WeibullMax) GetParams() Params {
	return Params{
		"loc":   w.Loc,
		"scale": w.Scale,
		"shape": w.Shape,
	}
}

// GetName returns the name of the distribution.
func (w WeibullMax) GetName() string {
	return "Weibull Maximum"
}

// SetParams calculates, sets and returns distribution parameters from linear regression parameters.
func (w *WeibullMax) SetParams(slope, intercept float64) {
	w.Scale = math.Exp(intercept / slope)
	w.Shape = slope
}

// TransformX computes an x-axis linearization.
func (w WeibullMax) TransformX(x float64) float64 {
	if x >= w.Loc {
		panic("Values must be less than location parameter for a WeibullMax x-scale.")
	}
	return -math.Log(w.Loc - x)
}

// TransformY computes an y-axis linearization.
func (WeibullMax) TransformY(y float64) float64 {
	if y <= 0 {
		panic("Values must be greater than 0 for a WeibullMax y-scale.")
	}
	return -math.Log(-math.Log(y))
}
