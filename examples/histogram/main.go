package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"

	chartjs "github.com/brentp/go-chartjs"
	"github.com/brentp/go-chartjs/types"
	"github.com/gonum/floats"
	"github.com/pkg/browser"
)

type hister struct {
	vals     []float64
	bins     []int
	binNames []string
	N        int
	Min      float64
	Max      float64
}

func (h *hister) hist() {
	if h.N <= 0 {
		h.N = int(0.5 + math.Pow(floats.Sum(h.vals), 0.33))
	}
	log.Println(h.N)
	if h.N < 1 {
		h.N = 1
	}
	if h.Min == 0 {
		h.Min = floats.Min(h.vals)
	}
	if h.Max == 0 {
		h.Max = floats.Max(h.vals)
	}

	w := (h.Max - h.Min) / float64(h.N)
	bins := make([]int, h.N)
	h.binNames = make([]string, h.N)
	for i := range h.binNames {
		v := h.Min + float64(i)*w
		if w < 2 {
			h.binNames[i] = fmt.Sprintf("%.2f-%.2f", v, v+w)
		} else {
			h.binNames[i] = fmt.Sprintf("%.0f-%.0f", v, v+w)
		}
	}

	for _, v := range h.vals {
		if v < h.Min || v > h.Max {
			continue
		}
		b := int((v - h.Min) / w)
		if v == h.Max {
			b = h.N - 1
		}
		bins[b]++
	}
	h.bins = bins
}

func (h hister) Xs() []float64 {
	bins := make([]float64, len(h.bins))
	for i, b := range h.bins {
		bins[i] = float64(b)
	}
	return bins
}

func (h hister) Ys() []float64 {
	return nil
}
func (h hister) Rs() []float64 {
	return nil
}

// IsStdin checks if we are getting data from stdin.
func isStdin() bool {
	// http://stackoverflow.com/a/26567513
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func main() {
	if !isStdin() {
		fmt.Fprintln(os.Stderr, "expecting values on stdin")
		os.Exit(0)
	}

	vals := make([]float64, 0, 100)
	stdin := bufio.NewReader(os.Stdin)

	for {
		line, err := stdin.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		v, err := strconv.ParseFloat(line[:len(line)-1], 64)
		if err != nil {
			panic(err)
		}
		vals = append(vals, v)
	}

	h := hister{vals: vals}
	h.N = 16
	h.hist()

	chart := &chartjs.Chart{}
	d := chartjs.Dataset{Data: h, Type: chartjs.Bar}
	d.BackgroundColor = &types.RGBA{102, 194, 165, 220}
	d.BorderWidth = 2
	d.Fill = types.True

	yax := chartjs.Axis{Type: chartjs.Linear, Position: chartjs.Left}
	yax.ScaleLabel = &chartjs.ScaleLabel{Display: types.True, LabelString: "Count"}
	_, err := chart.AddYAxis(yax)
	check(err)
	chart.Options.Scales.YAxes[0].Tick = &chartjs.Tick{BeginAtZero: types.True}

	chart.Data.Labels = h.binNames
	chart.Type = chartjs.Bar
	chart.AddDataset(d)

	chart.Options.Responsive = chartjs.False
	chart.Options.Legend = &chartjs.Legend{Display: chartjs.False}

	wtr, err := os.Create("hist.html")

	check(err)
	if err := chart.SaveHTML(wtr, nil); err != nil {
		log.Fatal(err)
	}
	wtr.Close()
	browser.OpenFile("hist.html")

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
