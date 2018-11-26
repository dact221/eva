package dist

import "math"

type WeibullMax struct {
	Loc   float64
	Scale float64
	Shape float64
}

func (w WeibullMax) CDF(x float64) float64 {
	if x > w.Loc {
		return 1
	}
	return math.Exp(-math.Pow((w.Loc-x)/w.Scale, w.Shape))
}

func (w WeibullMax) TransformX(x float64) float64 {
	if x >= w.Loc {
		panic("Values must be less than location parameter for a WeibullMax x-scale.")
	}
	return -math.Log(w.Loc - x)
}

func (WeibullMax) TransformY(y float64) float64 {
	if y <= 0 {
		panic("Values must be greater than 0 for a WeibullMax y-scale.")
	}
	return -math.Log(-math.Log(y))
}

type WeibullMin struct {
	Loc   float64
	Scale float64
	Shape float64
}

func (w WeibullMin) CDF(x float64) float64 {
	if x < w.Loc {
		return 0
	}
	return 1 - math.Exp(-math.Pow((x-w.Loc)/w.Scale, w.Shape))
}

func (w WeibullMin) TransformX(x float64) float64 {
	if x <= w.Loc {
		panic("Values must be greater than location parameter for a WeibullMin y-scale.")
	}
	return -math.Log(x - w.Loc)
}

func (WeibullMin) TransformY(y float64) float64 {
	if y >= 1 {
		panic("Values must be less than 1 for a WeibullMin y-scale.")
	}
	return -math.Log(-math.Log(1 - y))
}
