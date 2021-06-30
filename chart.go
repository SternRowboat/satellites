package main

import (
	"net/http"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"

	log "github.com/sirupsen/logrus"
)

// generate random data for line chart
func (d Database) getItems(i int) []opts.LineData {
	values := []opts.LineData{}
	for _, p := range d.satellite[i].data.packets {
		value := opts.LineData{}
		value.Value = p.value
		values = append(values, value)
	}
	return values
}

func (d Database) httpserver(w http.ResponseWriter, _ *http.Request) {
	log.Debug("All Data:", d.satellite)
	log.Info("Page Refreshed")
	// create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithTooltipOpts(opts.Tooltip{Show: true, Trigger: "axis"}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start: 0,
			End:   100,
		}),
		// charts.WithTooltipOpts(line.Trigger),
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeInfographic}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Satellite values over time",
			Subtitle: "Refresh to include the latest data",
		}))
	line1 := d.getItems(0)
	line2 := d.getItems(1)

	xData := []string{}
	for _, p := range d.satellite[0].data.packets {
		time := parseTime(p.UnixTimestamp)
		xData = append(xData, time)
	}
	line.SetXAxis(xData).
		AddSeries("Sat 1", line1).
		AddSeries("Sat 2", line2)
	line.Render(w)
}

func parseTime(t int64) string {
	return time.Unix(t, 0).Local().Format("08:00:07")
}

func Chart(d Database) {
	http.HandleFunc("/", d.httpserver)
	http.ListenAndServe(":8081", nil)
	log.Info("Ready to serve, go to http://localhost:8081/")
}
