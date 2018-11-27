package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"github.com/dact221/eva"
	"github.com/dact221/eva/dist"
)

var (
	fname = flag.String("f", "data.csv", "Data file.")
	loc = flag.Float64("l", 0, "Location parameter.")
)

func readData(fname string) ([]float64, []float64, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	var (
		i, x float64
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

func main() {

	flag.Parse()

	is, hs, err := readData(*fname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	w := &dist.WeibullMax{Loc: *loc}

	n := len(is)
	imax := is[n-1]

	var (
		pr = make([]float64, n)
		xi = make([]float64, n)
		eta = make([]float64, n)
	)

	for i := 0; i < n; i++ {
		pr[i] = eva.Hazen(is[i], imax)
		xi[i] = w.TransformX(hs[i])
		eta[i] = w.TransformY(pr[i])
	}

	slope, intercept, r := eva.FitDist(xi, eta)
	w.Scale = math.Exp(intercept / slope)
	w.Shape = slope

	fmt.Printf(report, w.Loc, w.Scale, w.Shape, r)
}

var report = `Weibull distribution parameters
===============================
Location = %.2f
Scale    = %.2f
Shape    = %.2f
R value  = %.2f
`