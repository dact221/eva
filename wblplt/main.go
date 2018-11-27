package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dact221/eva"
	"github.com/dact221/eva/dist"
	"gonum.org/v1/plot/vg"
)

var (
	fName   = flag.String("f", "", "Input data file.")
	isMax   = flag.Bool("m", true, "Model the maximum?")
	loc     = flag.Float64("l", 0, "Location parameter.")
	pltName = flag.String("o", "", "Output plot file.")
	xLabel  = flag.String("x", "", "x-axis label.")
)

func main() {
	flag.Parse()

	is, xs, err := readData(*fName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var w eva.Distribution
	switch *isMax {
	case true:
		w = &dist.WeibullMax{Loc: *loc}
	case false:
		w = &dist.WeibullMin{Loc: *loc}
	}

	n := len(is)
	imax := is[n-1]
	var (
		pr  = make([]float64, n)
		xi  = make([]float64, n)
		eta = make([]float64, n)
	)
	for i := 0; i < n; i++ {
		pr[i] = eva.Hazen(is[i], imax)
		xi[i] = w.TransformX(xs[i])
		eta[i] = w.TransformY(pr[i])
	}

	slope, intercept, r := eva.FitDist(xi, eta)
	w.SetParams(slope, intercept)
	wParams := w.GetParams()

	fmt.Printf(report, *loc, wParams["scale"], wParams["shape"], r)

	p, err := newWeibullPlot(xs, pr, w)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = p.Save(15*vg.Centimeter, 10*vg.Centimeter, *pltName)
	if err != nil {
		panic(err)
	}
}

var report = `Weibull distribution parameters
===============================
Location = %.2f
Scale    = %.2f
Shape    = %.2f
R value  = %.2f
`
