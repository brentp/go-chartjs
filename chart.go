// Package chartjs simplifies making chartjs.org plots in go.
package chartjs

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/brentp/go-chartjs/annotation"
	"github.com/brentp/go-chartjs/types"
)

var True = types.True
var False = types.False

var chartTypes = [...]string{
	"line",
	"bar",
	"bubble",
}

type chartType int

func (c chartType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + chartTypes[c] + `"`), nil
}

const (
	// Line is a "line" plot
	Line chartType = iota
	// Bar is a "bar" plot
	Bar
	// Bubble is a "bubble" plot
	Bubble
)

// FloatFormat determines how many decimal places are sent in the JSON.
var FloatFormat = "%.2f"

// Values dictates the interface of data to be plotted.
type Values interface {
	// X-axis values. If only these are specified then it must be a Bar plot.
	Xs() []float64
	// Optional Y values.
	Ys() []float64
	// Rs are used to size points for chartType `Bubble`
	Rs() []float64
}

func marshalValuesJSON(v Values) ([]byte, error) {
	xs, ys, rs := v.Xs(), v.Ys(), v.Rs()
	if len(xs) == 0 {
		if len(rs) != 0 {
			return nil, fmt.Errorf("chart: bad format of Values data")
		}
		xs = ys[:len(ys)]
		ys = nil
	}
	buf := bytes.NewBuffer(make([]byte, 0, 8*len(xs)))
	buf.WriteRune('[')
	if len(rs) > 0 {
		if len(xs) != len(ys) || len(xs) != len(rs) {
			return nil, fmt.Errorf("chart: bad format of Values. All axes must be of the same length")
		}
		for i, x := range xs {
			if i > 0 {
				buf.WriteRune(',')
			}
			y, r := ys[i], rs[i]
			_, err := buf.WriteString(fmt.Sprintf(("{\"x\":" + FloatFormat + ",\"y\":" + FloatFormat + ",\"r\":" + FloatFormat + "}"), x, y, r))
			if err != nil {
				return nil, err
			}
		}
	} else if len(ys) > 0 {
		if len(xs) != len(ys) {
			return nil, fmt.Errorf("chart: bad format of Values. X and Y must be of the same length")
		}
		for i, x := range xs {
			if i > 0 {
				buf.WriteRune(',')
			}
			y := ys[i]
			_, err := buf.WriteString(fmt.Sprintf(("{\"x\":" + FloatFormat + ",\"y\":" + FloatFormat + "}"), x, y))
			if err != nil {
				return nil, err
			}
		}

	} else {
		for i, x := range xs {
			if i > 0 {
				buf.WriteRune(',')
			}
			_, err := buf.WriteString(fmt.Sprintf(FloatFormat, x))
			if err != nil {
				return nil, err
			}
		}
	}

	buf.WriteRune(']')
	return buf.Bytes(), nil
}

// shape indicates the type of marker used for plotting.
type shape int

var shapes = []string{
	"",
	"circle",
	"triangle",
	"rect",
	"rectRot",
	"cross",
	"crossRot",
	"star",
	"line",
	"dash",
}

const (
	empty = iota
	Circle
	Triangle
	Rect
	RectRot
	Cross
	CrossRot
	Star
	LinePoint
	Dash
)

func (s shape) MarshalJSON() ([]byte, error) {
	return []byte(`"` + shapes[s] + `"`), nil
}

// Dataset wraps the "dataset" JSON
type Dataset struct {
	Data            Values      `json:"-"`
	Type            chartType   `json:"type,omitempty"`
	BackgroundColor *types.RGBA `json:"backgroundColor,omitempty"`
	// BorderColor is the color of the line.
	BorderColor *types.RGBA `json:"borderColor,omitempty"`
	// BorderWidth is the width of the line.
	BorderWidth int `json:"borderWidth,omitempty"`

	// Label indicates the name of the dataset to be shown in the legend.
	Label       string     `json:"label,omitempty"`
	Fill        types.Bool `json:"fill,omitempty"`
	LineTension float64    `json:"lineTension,omitempty"`

	PointBackgroundColor  *types.RGBA `json:"pointBackgroundColor,omitempty"`
	PointBorderColor      *types.RGBA `json:"pointBorderColor,omitempty"`
	PointBorderWidth      int         `json:"pointBorderWidth,omitempty"`
	PointRadius           float64     `json:"pointRadius,omitempty"`
	PointHoverBorderColor *types.RGBA `json:"pointHoverBorderColor,omitempty"`
	PointStyle            shape       `json:"pointStyle,omitempty"`

	ShowLine types.Bool `json:"showLine,omitempty"`
	SpanGaps types.Bool `json:"spanGaps,omitempty"`

	// Axis ID that matches the ID on the Axis where this dataset is to be drawn.
	XAxisID string `json:"xAxisID,omitempty"`
	YAxisID string `json:"yAxisID,omitempty"`
}

// MarshalJSON implements json.Marshaler interface.
func (d Dataset) MarshalJSON() ([]byte, error) {
	o, err := marshalValuesJSON(d.Data)
	// avoid recursion by creating an alias.
	type alias Dataset
	buf, err := json.Marshal(alias(d))
	if err != nil {
		return nil, err
	}
	// replace '}' with ',' to continue struct
	if len(buf) > 0 {
		buf[len(buf)-1] = ','
	}
	buf = append(buf, []byte(`"data":`)...)
	buf = append(buf, o...)
	buf = append(buf, '}')
	return buf, nil
}

// Data wraps the "data" JSON
type Data struct {
	Datasets []Dataset `json:"datasets"`
	Labels   []string  `json:"labels"`
}

type axisType int

var axisTypes = []string{
	"category",
	"linear",
	"logarithmic",
	"time",
	"radialLinear",
}

const (
	// Category is a categorical axis (this is the default),
	// used for bar plots.
	Category axisType = iota
	// Linear axis should be use for scatter plots.
	Linear
	// Log axis
	Log
	// Time axis
	Time
	// Radial axis
	Radial
)

func (t axisType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + axisTypes[t] + "\""), nil
}

type axisPosition int

const (
	// Bottom puts the axis on the bottom (used for Y-axis)
	Bottom axisPosition = iota + 1
	// Top puts the axis on the bottom (used for Y-axis)
	Top
	// Left puts the axis on the bottom (used for X-axis)
	Left
	// Right puts the axis on the bottom (used for X-axis)
	Right
)

var axisPositions = []string{
	"",
	"bottom",
	"top",
	"left",
	"right",
}

func (p axisPosition) MarshalJSON() ([]byte, error) {
	return []byte(`"` + axisPositions[p] + `"`), nil
}

// Axis corresponds to 'scale' in chart.js lingo.
type Axis struct {
	Type      axisType     `json:"type"`
	Position  axisPosition `json:"position,omitempty"`
	Label     string       `json:"label,omitempty"`
	ID        string       `json:"id,omitempty"`
	GridLines types.Bool   `json:"gridLine,omitempty"`
	Stacked   types.Bool   `json:"stacked,omitempty"`

	// need to differentiate between false and empty to use a pointer
	Display    types.Bool  `json:"display,omitempty"`
	ScaleLabel *ScaleLabel `json:"scaleLabel,omitempty"`
}

// ScaleLabel corresponds to scale title.
// Display: True must be specified for this to be shown.
type ScaleLabel struct {
	Display     types.Bool  `json:"display,omitempty"`
	LabelString string      `json:"labelString,omitempty"`
	FontColor   *types.RGBA `json:"fontColor,omitempty"`
	FontFamily  string      `json:"fontFamily,omitempty"`
	FontSize    int         `json:"fontSize,omitempty"`
	FontStyle   string      `json:"fontStyle,omitempty"`
}

// Axes holds the X and Y axies. Its simpler to use Chart.AddXAxis, Chart.AddYAxis.
type Axes struct {
	XAxes []Axis `json:"xAxes,omitempty"`
	YAxes []Axis `json:"yAxes,omitempty"`
}

// AddX adds a X-Axis.
func (a *Axes) AddX(x Axis) {
	a.XAxes = append(a.XAxes, x)
}

// AddY adds a Y-Axis.
func (a *Axes) AddY(y Axis) {
	a.YAxes = append(a.YAxes, y)
}

// Option wraps the chartjs "option"
type Option struct {
	Responsive          types.Bool `json:"responsive,omitempty"`
	MaintainAspectRatio types.Bool `json:"maintainAspectRatio,omitempty"`
}

type Annotation struct {
	Annotations []annotation.Annotation `json:"annotations,omitempty"`
}

// Options wraps the chartjs "options"
type Options struct {
	Option
	Scales     Axes       `json:"scales,omitempty"`
	Annotation Annotation `json:"annotation,omitempty"`
}

// Chart is the top-level type from chartjs.
type Chart struct {
	Type    chartType `json:"type"`
	Label   string    `json:"label,omitempty"`
	Data    Data      `json:"data,omitempty"`
	Options Options   `json:"options,omitempty"`
}

// AddDataset adds a dataset to the chart.
func (c *Chart) AddDataset(d Dataset) {
	c.Data.Datasets = append(c.Data.Datasets, d)
}

// AddXAxis adds an x-axis to the chart and returns the ID of the added axis.
func (c *Chart) AddXAxis(x Axis) (string, error) {
	if x.ID == "" {
		x.ID = fmt.Sprintf("xaxis%d", len(c.Options.Scales.XAxes))
	}
	if x.Position == Left || x.Position == Right {
		return "", fmt.Errorf("chart: added x-axis to left or right")
	}
	c.Options.Scales.XAxes = append(c.Options.Scales.XAxes, x)
	return x.ID, nil
}

// AddYAxis adds an y-axis to the chart and return the ID of the added axis.
func (c *Chart) AddYAxis(y Axis) (string, error) {
	if y.ID == "" {
		y.ID = fmt.Sprintf("yaxis%d", len(c.Options.Scales.YAxes))
	}
	if y.Position == Top || y.Position == Bottom {
		return "", fmt.Errorf("chart: added y-axis to top or bottom")
	}
	c.Options.Scales.YAxes = append(c.Options.Scales.YAxes, y)
	return y.ID, nil
}
