package chartjs

import (
	"encoding/json"
	"html/template"
	"io"
)

// this file implements some syntactic sugar for creating charts

// JQuery holds the path to hosted JQuery
var JQuery = "https://code.jquery.com/jquery-2.2.4.min.js"

// ChartJS holds the path to hosted ChartJS
var ChartJS = "https://cdnjs.cloudflare.com/ajax/libs/Chart.js/2.3.0/Chart.bundle.js"

const tmpl = `<!DOCTYPE html>
<html>
    <head>
		<script src="{{ index . "JQuery" }}"></script>
		<script src="{{ index . "ChartJS" }}"></script>
    </head>
    <body>
        <canvas id="canvas" height="{{ index . "height" }}" width="{{ index . "width" }}"></canvas>
    </body>
    <script>
    var ctx = document.getElementById("canvas").getContext("2d");
	new Chart(ctx, {{ index . "json" }})
    </script>
</html>`

// SaveHTML writes the chart and minimal HTML to an io.Writer.
func (c Chart) SaveHTML(w io.Writer, tmap map[string]interface{}) error {
	if tmap == nil {
		tmap = make(map[string]interface{})
	}

	if _, ok := tmap["height"]; !ok {
		tmap["height"] = 400
	}
	if _, ok := tmap["width"]; !ok {
		tmap["width"] = 400
	}
	cjson, err := json.Marshal(c)
	if err != nil {
		return err
	}
	tmap["json"] = template.JS(cjson)
	if _, ok := tmap["JQuery"]; !ok {
		tmap["JQuery"] = JQuery
	}
	if _, ok := tmap["ChartJS"]; !ok {
		tmap["ChartJS"] = ChartJS
	}

	t, err := template.New("chartjs").Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(w, tmap)
}
