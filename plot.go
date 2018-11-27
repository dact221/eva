package eva

import (
	"fmt"
	"math"

	"github.com/dact221/eva/dist"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

// ScaleX can be used as the value of an Axis.Scale function to set the x-axis to a custom probability scale.
type ScaleX func(x float64) float64

// Normalize returns the fractional transformed distance of x between min and max.
func (sx ScaleX) Normalize(min, max, x float64) float64 {
	tXMin := sx(min)
	return (sx(x) - tXMin) / (sx(max) - tXMin)
}

// ScaleY can be used as the value of an Axis.Scale function to set the y-axis to a custom probability scale.
type ScaleY func(y float64) float64

// Normalize returns the fractional transformed distance of y between min and max.
func (sy ScaleY) Normalize(min, max, x float64) float64 {
	tYMin := sy(min)
	return (sy(x) - tYMin) / (sy(max) - tYMin)
}

// NewConstantTicks returns custom Ticks suitable for the Tick.Marker field of an Axis.
func NewConstantTicks(n, prec int, min, max float64) plot.ConstantTicks {
	ticks := make(plot.ConstantTicks, n)
	label := fmt.Sprintf("%%.%df", prec)
	delta := (max - min) / float64(n-1)

	ticks[0].Value = min
	ticks[0].Label = fmt.Sprintf(label, min)

	for i := 1; i < n-1; i++ {
		x := min + float64(i)*delta
		ticks[i].Value = x
		ticks[i].Label = fmt.Sprintf(label, x)
	}

	ticks[n-1].Value = max
	ticks[n-1].Label = fmt.Sprintf(label, max)

	return ticks
}

// NewProbPlot returns a probability plot for the given data and Distribution.
func NewProbPlot(xs, pr *[]float64, w dist.Distribution) (*plot.Plot, error) {
	p, err := plot.New()
	if err != nil {
		return nil, err
	}

	p.Title.Text = fmt.Sprintf("%s Plot", w.GetName())
	p.X.Scale = ScaleX(w.TransformX)
	p.Y.Scale = ScaleY(w.TransformY)

	n := len(*xs)

	pts := make(plotter.XYs, n)
	for i := range *xs {
		pts[i].X = (*xs)[i]
		pts[i].Y = (*pr)[i]
	}

	s, err := plotter.NewScatter(pts)
	if err != nil {
		return nil, err
	}

	f := plotter.NewFunction(w.CDF)
	f.XMin = (*xs)[0]
	f.XMax = (*xs)[n-1]

	p.Add(s, f, plotter.NewGrid())

	p.X.Min = math.Min(pts[0].X, f.XMin)
	p.X.Max = math.Max(pts[n-1].X, f.XMax)
	p.Y.Min = math.Min(pts[0].Y, w.CDF(f.XMin))
	p.Y.Max = math.Max(pts[n-1].Y, w.CDF(f.XMax))

	p.X.Tick.Marker = NewConstantTicks(7, 1, p.X.Min, p.X.Max)
	p.Y.Tick.Marker = NewConstantTicks(7, 3, p.Y.Min, p.Y.Max)

	return p, nil
}
