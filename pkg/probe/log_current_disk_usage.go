package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type LogResults struct {
	UsedBytes  int `json:"used_bytes"`
	FreeBytes  int `json:"free_bytes"`
	TotalBytes int `json:"total_bytes"`
}

type Log struct {
	Results LogResults `json:"results"`
	VDOM    string     `json:"vdom"`
	Version string     `json:"version"`
}

func probeLogCurrentDiskUsage(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		logUsed = prometheus.NewDesc(
			"fortigate_log_used_bytes",
			"Current used bytes for log",
			[]string{"vdom"}, nil,
		)
		logFree = prometheus.NewDesc(
			"fortigate_log_free_bytes",
			"Current free bytes for log",
			[]string{"vdom"}, nil,
		)
		logTotal = prometheus.NewDesc(
			"fortigate_log_total_bytes",
			"Current total bytes for log",
			[]string{"vdom"}, nil,
		)
	)

	var res []Log
	if err := c.Get("api/v2/monitor/log/current-disk-usage", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res {
		m = append(m, prometheus.MustNewConstMetric(logUsed, prometheus.GaugeValue, float64(r.Results.UsedBytes), r.VDOM))
		m = append(m, prometheus.MustNewConstMetric(logFree, prometheus.GaugeValue, float64(r.Results.FreeBytes), r.VDOM))
		m = append(m, prometheus.MustNewConstMetric(logTotal, prometheus.GaugeValue, float64(r.Results.TotalBytes), r.VDOM))
	}

	return m, true
}
