package eva

import "gonum.org/v1/gonum/stat"

// Distribution represents an extreme value distribution. Currently
// there are three implementations: Fréchet, Gumbel and Weibull.
type Distribution interface {
	CDF(x float64) float64
	TransformX(x float64) float64
	TransformY(y float64) float64
}

// ScaleX implements plot.Normalizer.
type ScaleX struct {
	Dist Distribution
}

func (s ScaleX) Normalize(min, max, x float64) float64 {
	tXMin := s.Dist.TransformX(min)
	return (s.Dist.TransformX(x) - tXMin) / (s.Dist.TransformX(max) - tXMin)
}

// ScaleY implements plot.Normalizer.
type ScaleY struct {
	Dist Distribution
}

func (s ScaleY) Normalize(min, max, x float64) float64 {
	tYMin := s.Dist.TransformY(min)
	return (s.Dist.TransformY(x) - tYMin) / (s.Dist.TransformY(max) - tYMin)
}

// Blom's plotting position.
func Blom(i, n float64) float64 {
	return (i - 0.375) / (n + 0.25)
}

// Gringorten's plotting position.
func Gringorten(i, n float64) float64 {
	return (i - 0.44) / (n + 0.12)
}

// Hazen's plotting position.
func Hazen(i, n float64) float64 {
	return (i - 0.5) / n
}

// Weibull's plotting position.
func Weibull(i, n float64) float64 {
	return i / (n + 1)
}

// FitDist calculate a linear least-squares regression for TransformX(x)
// and TransformY(y). You can use these linear parameters to compute
// distribution parameters.
func FitDist(tx, ty []float64) (slope, intercept, rvalue float64) {
	intercept, slope = stat.LinearRegression(tx, ty, nil, false)
	rvalue = stat.Correlation(tx, ty, nil)
	return
}
