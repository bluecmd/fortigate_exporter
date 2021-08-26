package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type LogResults struct {
	UsedBytes  float64 `json:"used_bytes"`
	TotalBytes float64 `json:"total_bytes"`
}

type Log struct {
	Results LogResults `json:"results"`
	VDOM    string     `json:"vdom"`
}

func probeLogCurrentDiskUsage(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		logUsed = prometheus.NewDesc(
			"fortigate_log_disk_used_bytes",
			"Disk used bytes for log",
			[]string{"vdom"}, nil,
		)
		logTotal = prometheus.NewDesc(
			"fortigate_log_disk_total_bytes",
			"Disk total bytes for log",
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
		m = append(m, prometheus.MustNewConstMetric(logUsed, prometheus.GaugeValue, r.Results.UsedBytes, r.VDOM))
		m = append(m, prometheus.MustNewConstMetric(logTotal, prometheus.GaugeValue, r.Results.TotalBytes, r.VDOM))
	}

	return m, true
}
