package chartjs

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"testing"
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

func TestLine(t *testing.T) {

	//var axes Axes
	//axes.AddY(Axis{Type: Linear, Position: Bottom})
	//axes.AddY(Axis{Type: Linear, Position: Left})

	var xys xy
	for i := 0; i < 10; i++ {
		xys.x = append(xys.x, float64(i))
		xys.y = append(xys.y, float64(i))
		xys.r = append(xys.r, float64(i))
	}

	d := Dataset{Data: xys, BackgroundColor: &RGBA{0, 255, 0, 200}, Label: "HHIHIHI"}

	//options := Options{Scales: axes}
	//options.Responsive = true
	//options.MaintainAspectRatio = false

	chart := Chart{Type: Bubble, Label: "test-chart"}
	//chart.Data = Data{Datasets: []Dataset{d}}
	chart.AddDataset(d)
	chart.AddXAxis(Axis{Type: Linear, Position: Bottom})
	chart.AddYAxis(Axis{Type: Linear, Position: Right})

	b, err := json.Marshal(chart)
	if err != nil {
		t.Fatalf("error marshaling chart: %+v", err)
	}
	fmt.Println(string(b))
}

func TestBar(t *testing.T) {

	//var axes Axes
	//axes.AddY(Axis{Type: Linear, Position: Bottom})
	//axes.AddY(Axis{Type: Linear, Position: Left})

	var xs xy
	var labels []string
	for i := 0; i < 10; i++ {
		xs.x = append(xs.x, float64(i))
		labels = append(labels, strconv.Itoa(i))
	}
	d := Dataset{Data: xs, BackgroundColor: &RGBA{0, 255, 0, 200}}

	//options := Options{Scales: axes}
	//options.Responsive = true
	//options.MaintainAspectRatio = false

	chart := Chart{Type: Bar, Label: "test-chart"}
	//chart.Data = Data{Datasets: []Dataset{d}}
	chart.AddDataset(d)
	chart.Data.Labels = labels

	b, err := json.Marshal(chart)
	if err != nil {
		t.Fatalf("error marshaling chart: %+v", err)
	}
	fmt.Println(string(b))
}

func TestHTML(t *testing.T) {

	var xys xy
	for i := float64(0); i < 9; i += 0.05 {
		xys.x = append(xys.x, float64(i))
		xys.y = append(xys.y, math.Sin(float64(i)))
		xys.r = append(xys.r, float64(i))
	}
	fmt.Println(len(xys.x))

	d := Dataset{Data: xys, BackgroundColor: &RGBA{0, 255, 0, 200}, Label: "sin(x)"}

	//options := Options{Scales: axes}
	//options.Responsive = true
	//options.MaintainAspectRatio = false

	chart := Chart{Type: Bubble, Label: "test-chart"}
	//chart.Data = Data{Datasets: []Dataset{d}}
	chart.AddDataset(d)
	chart.AddXAxis(Axis{Type: Linear, Position: Bottom})
	chart.AddYAxis(Axis{Type: Linear, Position: Right})
	chart.Options.Responsive = False

	wtr, err := os.Create("test-chartjs.html")
	if err != nil {
		t.Fatalf("error opening file: %+v", err)
	}
	if err := chart.SaveHTML(wtr, nil); err != nil {
		t.Fatalf("error saving chart: %+v", err)
	}
	wtr.Close()
}
