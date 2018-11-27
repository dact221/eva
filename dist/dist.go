package dist

// Params maps the names of distribution parameters to their values.
type Params map[string]float64

// Distribution represents an extreme value distribution. There are six implementations: FréchetMax/Min, GumbelMax/Min and WeibullMax/Min.
type Distribution interface {
	CDF(x float64) float64
	GetParams() Params
	GetName() string
	SetParams(slope, intercept float64)
	TransformX(x float64) float64
	TransformY(y float64) float64
}
