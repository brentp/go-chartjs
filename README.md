chartjs
-------

[![GoDoc] (https://godoc.org/github.com/brentp/go-chartjs?status.png)](https://godoc.org/github.com/brentp/go-chartjs)
[![Build Status](https://travis-ci.org/brentp/go-chartjs.svg)](https://travis-ci.org/brentp/go-chartjs)

go wrapper for [chartjs](http://chartjs.org)

Since chartjs charts are defined purely in JSON, this library has a layer with minimal code;
structs and struct-tags indicate how to marshal to JSON. None of the currently
implemented parts are stringly-typed in this library so it can avoid many errors.

The javascript/JSON api has a [lot of surface area](http://www.chartjs.org/docs/).
Currently, only the options that I use are provided. More can and will be added as I need
them (or via pull-requests).

There is currently a small amount of code to simplify creating charts.

data to be plotted by chartjs has to meet this interface.

```Go
  type Values interface {
      // X-axis values. If only these are specified then it must be a Bar plot.
      Xs() []float64
      // Optional Y values.
      Ys() []float64
      // Rs are used for chartType `Bubble`
      Rs() []float64
  }
```

Example
-------

This examples shows common use of the library and demonstrates saving to an HTML file.

```Go
package main

import (
	"math"
	"os"

	chartjs "github.com/brentp/go-chartjs"
)

type xy struct {
	x []float64
	y []float64
	r []float64
}

func (v xy) Xs() []float64 {
	return v.x
}
func (v xy) Ys() []float64 {
	return v.y
}
func (v xy) Rs() []float64 {
	return v.r
}

func main() {
	var xys xy
	for i := float64(0); i < 9; i += 0.05 {
		xys.x = append(xys.x, float64(i))
		xys.y = append(xys.y, math.Sin(float64(i)))
		xys.r = append(xys.r, float64(i+1))
	}

	d := chartjs.Dataset{Data: xys, BackgroundColor: &chartjs.RGBA{0, 102, 255, 200}, Label: "sin(x)"}

	chart := chartjs.Chart{Type: chartjs.Bubble, Label: "sin chart"}
	chart.AddDataset(d)
	chart.AddXAxis(chartjs.Axis{Type: chartjs.Linear, Position: chartjs.Bottom})
	chart.AddYAxis(chartjs.Axis{Type: chartjs.Linear, Position: chartjs.Right})
	chart.Options.Responsive = chartjs.False

	wtr, _ := os.Create("test-chartjs.html")
	chart.SaveHTML(wtr, nil)
	wtr.Close()

}
```
