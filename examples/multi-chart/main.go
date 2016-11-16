package main

import (
	"log"
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
	var xys1 xy
	var xys2 xy

	for i := float64(0); i < 9; i += 0.1 {
		xys1.x = append(xys1.x, float64(i))
		xys2.x = append(xys2.x, float64(i))

		xys1.y = append(xys1.y, math.Sin(float64(i)))

		xys2.y = append(xys2.y, 2*math.Cos(float64(i)))

	}

	// a set of colors to work with.
	colors := []*chartjs.RGBA{
		&chartjs.RGBA{102, 194, 165, 220},
		&chartjs.RGBA{250, 141, 98, 220},
		&chartjs.RGBA{141, 159, 202, 220},
		&chartjs.RGBA{230, 138, 195, 220},
	}

	d1 := chartjs.Dataset{Data: xys1, BorderColor: colors[0], Label: "sin(x)", Fill: chartjs.False,
		PointRadius: 10, PointBorderWidth: 4, BackgroundColor: colors[1]}

	d2 := chartjs.Dataset{Data: xys2, BorderWidth: 8, BorderColor: colors[2], Label: "2 * cos(x)", Fill: chartjs.False}

	chart := chartjs.Chart{Type: chartjs.Line, Label: "test-chart"}
	chart.AddXAxis(chartjs.Axis{Type: chartjs.Linear, Position: chartjs.Bottom, ScaleLabel: &chartjs.ScaleLabel{FontSize: 22, LabelString: "X", Display: chartjs.True}})
	d1.YAxisID = chart.AddYAxis(chartjs.Axis{Type: chartjs.Linear, Position: chartjs.Left,
		ScaleLabel: &chartjs.ScaleLabel{LabelString: "sin(x)", Display: chartjs.True}})
	d2.YAxisID = chart.AddYAxis(chartjs.Axis{Type: chartjs.Linear, Position: chartjs.Right,
		ScaleLabel: &chartjs.ScaleLabel{LabelString: "2 * cos(x)", Display: chartjs.True}})

	chart.AddDataset(d2)
	chart.AddDataset(d1)

	chart.Options.Responsive = chartjs.False

	wtr, err := os.Create("test-chartjs-multi.html")
	if err != nil {
	}
	if err := chart.SaveHTML(wtr, nil); err != nil {
		log.Fatal(err)
	}
	wtr.Close()
}
