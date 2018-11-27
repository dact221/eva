package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"

	"github.com/dact221/eva"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func newConstantTicks(n, prec int, min, max float64) plot.ConstantTicks {
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

func readData(fName string) ([]float64, []float64, error) {
	f, err := os.Open(fName)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 2

	var (
		i, x   float64
		is, xs []float64
	)

	for {
		row, err := r.Read()
		if err == io.EOF {
			return is, xs, nil
		}
		if err != nil {
			return nil, nil, err
		}
		if i, err = strconv.ParseFloat(row[0], 64); err != nil {
			return nil, nil, err
		}
		if x, err = strconv.ParseFloat(row[1], 64); err != nil {
			return nil, nil, err
		}
		is = append(is, i)
		xs = append(xs, x)
	}
}

func newWeibullPlot(xs, pr []float64, w eva.Distribution) (*plot.Plot, error) {
	p, err := plot.New()
	if err != nil {
		return nil, err
	}

	p.Title.Text = fmt.Sprintf("%s Plot", w.GetName())
	p.X.Scale = eva.ScaleX(w.TransformX)
	p.Y.Scale = eva.ScaleY(w.TransformY)
	p.X.Label.Text = *xLabel
	p.Y.Label.Text = "P"

	n := len(xs)

	pts := make(plotter.XYs, n)
	for i := range xs {
		pts[i].X = xs[i]
		pts[i].Y = pr[i]
	}

	s, err := plotter.NewScatter(pts)
	if err != nil {
		return nil, err
	}

	f := plotter.NewFunction(w.CDF)
	f.XMin = xs[0]
	f.XMax = xs[n-1]

	p.Add(s, f, plotter.NewGrid())

	p.X.Min = math.Min(pts[0].X, f.XMin)
	p.X.Max = math.Max(pts[n-1].X, f.XMax)
	p.Y.Min = math.Min(pts[0].Y, w.CDF(f.XMin))
	p.Y.Max = math.Max(pts[n-1].Y, w.CDF(f.XMax))

	p.X.Tick.Marker = newConstantTicks(7, 1, p.X.Min, p.X.Max)
	p.Y.Tick.Marker = newConstantTicks(7, 3, p.Y.Min, p.Y.Max)

	return p, nil
}
