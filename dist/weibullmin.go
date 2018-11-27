package dist

import "math"

// WeibullMin represents a Weibull Minimum distribution.
type WeibullMin struct {
	Loc   float64
	Scale float64
	Shape float64
}

// CDF computes the value of the cumulative density function at x.
func (w WeibullMin) CDF(x float64) float64 {
	if x < w.Loc {
		return 0
	}
	return 1 - math.Exp(-math.Pow((x-w.Loc)/w.Scale, w.Shape))
}

// GetParams returns parameters map.
func (w WeibullMin) GetParams() Params {
	return Params{
		"loc":   w.Loc,
		"scale": w.Scale,
		"shape": w.Shape,
	}
}

// GetName returns the name of the distribution.
func (w WeibullMin) GetName() string {
	return "Weibull Minimum"
}

// SetParams calculates and sets distribution parameters from linear regression parameters.
func (w *WeibullMin) SetParams(slope, intercept float64) {
	w.Scale = math.Exp(intercept / slope)
	w.Shape = slope
}

// TransformX computes an x-axis linearization.
func (w WeibullMin) TransformX(x float64) float64 {
	if x <= w.Loc {
		panic("Values must be greater than location parameter for a WeibullMin y-scale.")
	}
	return -math.Log(x - w.Loc)
}

// TransformY computes an y-axis linearization.
func (WeibullMin) TransformY(y float64) float64 {
	if y >= 1 {
		panic("Values must be less than 1 for a WeibullMin y-scale.")
	}
	return -math.Log(-math.Log(1 - y))
}
