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
var ChartJS = "https://cdnjs.cloudflare.com/ajax/libs/Chart.js/2.4.0/Chart.bundle.js"

const tmpl = `<!DOCTYPE html>
<html>
    <head>
		<script src="{{ index . "JQuery" }}"></script>
		<script src="{{ index . "ChartJS" }}"></script>
		<script>
		{{ index . "extra"}}
		</script>
    </head>
    <body>
	{{ $height := index . "height" }}
	{{ $width := index . "width" }}
	{{ range $i, $json := index . "charts" }}
	<canvas id="canvas{{ $i }}" style="height:{{ $height }}px;width:{{ $width }}px"></canvas>
		<hr>
	{{ end }}
	{{ index . "customHTML" }}
    </body>
    <script>
	Chart.defaults.line.cubicInterpolationMode = 'monotone';
	Chart.defaults.global.animation.duration = 0;
	var charts = []
	{{ range $i, $json := index . "charts" }}
		var ctx = document.getElementById("canvas{{ $i }}").getContext("2d");
		var chart = new Chart(ctx, {{ $json }});
		charts.push(chart)
	{{ end }}
	{{ index . "custom" }}
    </script>
</html>`

// SaveCharts writes the charts and the required HTML to an io.Writer
func SaveCharts(w io.Writer, tmap map[string]interface{}, charts ...Chart) error {
	if tmap == nil {
		tmap = make(map[string]interface{})
	}

	if _, ok := tmap["height"]; !ok {
		tmap["height"] = 400
	}
	if _, ok := tmap["width"]; !ok {
		tmap["width"] = 400
	}
	jscharts := make([]template.JS, 0, len(charts))
	for _, c := range charts {
		cjson, err := json.Marshal(c)
		if err != nil {
			return err
		}
		jscharts = append(jscharts, template.JS(cjson))
	}
	for k, v := range tmap {
		if chart, ok := v.(Chart); ok {
			cjson, err := json.Marshal(chart)
			if err != nil {
				return err
			}
			tmap[k] = template.JS(cjson)
		}
	}

	tmap["charts"] = jscharts
	if _, ok := tmap["JQuery"]; !ok {
		tmap["JQuery"] = JQuery
	}
	if _, ok := tmap["ChartJS"]; !ok {
		tmap["ChartJS"] = ChartJS
	}
	if _, ok := tmap["custom"]; !ok {
		tmap["custom"] = ""
	}
	if _, ok := tmap["customHTML"]; !ok {
		tmap["customHTML"] = ""
	}
	if _, ok := tmap["template"]; !ok {
		tmap["template"] = tmpl
	}
	t, err := template.New("chartjs").Parse(tmap["template"].(string))
	if err != nil {
		return err
	}
	return t.Execute(w, tmap)
}

// SaveHTML writes the chart and minimal HTML to an io.Writer.
func (c Chart) SaveHTML(w io.Writer, tmap map[string]interface{}) error {
	return SaveCharts(w, tmap, c)
}
