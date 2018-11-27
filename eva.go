package eva

import "gonum.org/v1/gonum/stat"

// PlottingPosition represents a plotting position where i is the ordered rank of a sample value and n is the sample size.
type PlottingPosition func(i, n float64) float64

// GetPlottingPosition returns a PlottingPosition with given name.
func GetPlottingPosition(name string) PlottingPosition {
	switch name {
	case "blom":
		return Blom
	case "gringorten":
		return Gringorten
	case "hazen":
		return Hazen
	default:
		return Weibull
	}
}

// Blom is Blom's plotting position.
func Blom(i, n float64) float64 {
	return (i - 0.375) / (n + 0.25)
}

// Gringorten is Gringorten's plotting position.
func Gringorten(i, n float64) float64 {
	return (i - 0.44) / (n + 0.12)
}

// Hazen is Hazen's plotting position.
func Hazen(i, n float64) float64 {
	return (i - 0.5) / n
}

// Weibull is Weibull's plotting position.
func Weibull(i, n float64) float64 {
	return i / (n + 1)
}

// FitDist calculate a linear least-squares regression for TransformX(x)
// and TransformY(y). You can use these linear parameters to compute
// distribution parameters.
func FitDist(tx, ty *[]float64) (slope, intercept, rvalue float64) {
	intercept, slope = stat.LinearRegression(*tx, *ty, nil, false)
	rvalue = stat.Correlation(*tx, *ty, nil)
	return
}
