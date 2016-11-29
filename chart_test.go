package chartjs

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"testing"

	"github.com/brentp/go-chartjs/annotation"
	"github.com/brentp/go-chartjs/types"
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

	d := Dataset{Data: xys, BackgroundColor: &types.RGBA{0, 255, 0, 200}, Label: "HHIHIHI"}

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
	d := Dataset{Data: xs, BackgroundColor: &types.RGBA{0, 255, 0, 200}}

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

	d := Dataset{Data: xys, BackgroundColor: &types.RGBA{0, 255, 0, 200}, Label: "sin(x)"}

	//options := Options{Scales: axes}
	//options.Responsive = true
	//options.MaintainAspectRatio = false

	chart := Chart{Type: Bubble, Label: "test-chart"}
	//chart.Data = Data{Datasets: []Dataset{d}}
	chart.AddDataset(d)
	chart.AddXAxis(Axis{Type: Linear, Position: Bottom})
	chart.AddYAxis(Axis{Type: Linear, Position: Right})
	chart.Options.Responsive = types.False

	wtr, err := os.Create("test-chartjs.html")
	if err != nil {
		t.Fatalf("error opening file: %+v", err)
	}
	if err := chart.SaveHTML(wtr, nil); err != nil {
		t.Fatalf("error saving chart: %+v", err)
	}
	wtr.Close()
}

func TestMultipleCharts(t *testing.T) {
	var xys1 xy
	var xys2 xy

	for i := float64(0); i < 9; i += 0.1 {
		xys1.x = append(xys1.x, float64(i))
		xys2.x = append(xys2.x, float64(i))

		xys1.y = append(xys1.y, math.Sin(float64(i)))

		xys2.y = append(xys2.y, 2*math.Cos(float64(i)))

	}

	// a set of colors to work with.
	colors := []*types.RGBA{
		&types.RGBA{102, 194, 165, 220},
		&types.RGBA{250, 141, 98, 220},
		&types.RGBA{141, 159, 202, 220},
		&types.RGBA{230, 138, 195, 220},
	}

	d1 := Dataset{Data: xys1, BorderColor: colors[0], Label: "sin(x)", Fill: types.False,
		PointRadius: 10, PointBorderWidth: 4, BackgroundColor: colors[1]}

	d2 := Dataset{Data: xys2, BorderWidth: 8, BorderColor: colors[2], Label: "2 * cos(x)", Fill: types.False}

	chart := Chart{Type: Line, Label: "test-chart"}
	chart.AddXAxis(Axis{Type: Linear, Position: Bottom, ScaleLabel: &ScaleLabel{FontSize: 22, LabelString: "X", Display: types.True}})
	var err error
	d1.YAxisID, err = chart.AddYAxis(Axis{Type: Linear, Position: Left, ScaleLabel: &ScaleLabel{LabelString: "sin(x)", Display: types.True}})
	if err != nil {
		t.Fatalf("error adding axis: %s", err)
	}
	d2.YAxisID, err = chart.AddYAxis(Axis{Type: Linear, Position: Right, ScaleLabel: &ScaleLabel{LabelString: "2 * cos(x)", Display: types.True}})
	if err != nil {
		t.Fatalf("error adding axis: %s", err)
	}

	chart.AddDataset(d2)
	chart.AddDataset(d1)

	chart.Options.Responsive = types.False

	wtr, err := os.Create("test-chartjs-multi.html")
	if err != nil {
		t.Fatalf("error opening file: %+v", err)
	}
	if err := chart.SaveHTML(wtr, nil); err != nil {
		t.Fatalf("error saving chart: %+v", err)
	}
	wtr.Close()
}

func TestAnno(t *testing.T) {

	var xys xy
	for i := float64(0); i < 9; i += 0.05 {
		xys.x = append(xys.x, float64(i))
		xys.y = append(xys.y, math.Sin(float64(i)))
		xys.r = append(xys.r, float64(i))
	}

	d := Dataset{Data: xys, BackgroundColor: &types.RGBA{0, 255, 0, 200}, Label: "sin(x)"}

	//options := Options{Scales: axes}
	//options.Responsive = true
	//options.MaintainAspectRatio = false

	chart := Chart{Type: Bubble, Label: "test-chart"}
	//chart.Data = Data{Datasets: []Dataset{d}}
	chart.AddDataset(d)
	xid, _ := chart.AddXAxis(Axis{Type: Linear, Position: Bottom})
	yid, _ := chart.AddYAxis(Axis{Type: Linear, Position: Right})
	chart.Options.Responsive = types.False

	chart.Options.Annotation = Annotation{[]annotation.Annotation{annotation.Annotation{Type: annotation.Line, AxisID: xid, Mode: annotation.Vertical,
		Value: 4, EndValue: 6, BorderColor: &types.RGBA{255, 0, 0, 200}, BorderWidth: 4, Label: &annotation.Label{FontSize: 12, FontFamily: "sans-serif",
			Position: "center", BackgroundColor: &types.RGBA{0, 0, 0, 200}, Enabled: types.True, Content: "hello"}}}}

	chart.Options.Annotation.Annotations = append(chart.Options.Annotation.Annotations,
		annotation.Annotation{Type: annotation.Box, XAxisID: xid, YAxisID: yid,
			BackgroundColor: &types.RGBA{0, 0, 0, 200},
			XMin:            1, XMax: 4, YMin: -0.6, YMax: 0.77, BorderColor: &types.RGBA{200, 200, 0, 200}, BorderWidth: 3, Label: &annotation.Label{FontSize: 12, FontFamily: "sans-serif",
				Position:  "center",
				FontColor: &types.RGBA{255, 0, 0, 200},
				Enabled:   types.True, Content: "boxed"}})

	wtr, err := os.Create("test-chart-anno.html")
	if err != nil {
		t.Fatalf("error opening file: %+v", err)
	}
	if err := chart.SaveHTML(wtr, nil); err != nil {
		t.Fatalf("error saving chart: %+v", err)
	}
	wtr.Close()
}
