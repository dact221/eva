package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/dact221/eva"
	"github.com/dact221/eva/dist"
	"gonum.org/v1/plot/vg"
)

var (
	fName   = flag.String("in", "", "Input data file.")
	pltName = flag.String("out", "", "Output plot file.")
	isMax   = flag.Bool("max", true, "Model the maximum?")
	loc     = flag.Float64("loc", 0, "Location parameter.")
	xLabel  = flag.String("xlabel", "x", "x-axis label.")
	yLabel  = flag.String("ylabel", "P", "y-axis label.")
	pltPos  = flag.String("pltpos", "", "Plotting position (weibull, blom, hazen, gringorten).")
	width   = flag.Float64("width", 10, "Plot width (cm).")
	height  = flag.Float64("height", 10, "Plot height (cm).")
)

func main() {
	flag.Parse()

	fmt.Println("Read file:", *fName)
	is, xs, err := readData(*fName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	n := len(is)
	fmt.Println("Sample size:", n)

	var w dist.Distribution
	if *isMax {
		w = &dist.WeibullMax{Loc: *loc}
	} else {
		w = &dist.WeibullMin{Loc: *loc}
	}
	fmt.Println("Distribution:", w.GetName())
	fmt.Println("Location parameter:", *loc)

	f := eva.GetPlottingPosition(*pltPos)
	fmt.Println("Plotting position:", *pltPos)
	imax := is[n-1]
	var (
		pr  = make([]float64, n)
		xi  = make([]float64, n)
		eta = make([]float64, n)
	)
	for i := 0; i < n; i++ {
		pr[i] = f(is[i], imax)
		xi[i] = w.TransformX(xs[i])
		eta[i] = w.TransformY(pr[i])
	}

	slope, intercept, r := eva.FitDist(&xi, &eta)
	fmt.Println("Slope:", slope)
	fmt.Println("Intercept:", intercept)
	w.SetParams(slope, intercept)
	wParams := w.GetParams()

	fmt.Printf(report, *loc, wParams["scale"], wParams["shape"], r)

	p, err := eva.NewProbPlot(&xs, &pr, w)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p.X.Label.Text = *xLabel
	p.Y.Label.Text = *yLabel

	fmt.Println("Save plot:", *pltName)
	err = p.Save(vg.Length(*width)*vg.Centimeter, vg.Length(*height)*vg.Centimeter, *pltName)
	if err != nil {
		panic(err)
	}
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

var report = `
Weibull distribution parameters
===============================
Location = %.2f
Scale    = %.2f
Shape    = %.2f
R value  = %.2f

`
