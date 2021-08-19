package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type LogAnaQueueResults struct {
	Connected  float64 `json:"connected"`
	FailedLogs float64 `json:"failed_logs"`
	CachedLogs float64 `json:"cached_logs"`
}

type LogAnaQueue struct {
	Results LogAnaQueueResults `json:"results"`
	VDOM    string             `json:"vdom"`
}

func probeLogAnalyzerQueue(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		logAnaConn = prometheus.NewDesc(
			"fortigate_log_fortianalyzer_queue_connections",
			"Fortianalyzer queue connected state",
			[]string{"vdom"}, nil,
		)
		logAnaLogs = prometheus.NewDesc(
			"fortigate_log_fortianalyzer_queue_logs",
			"State of logs in the queue",
			[]string{"vdom", "state"}, nil,
		)
	)

	var res []LogAnaQueue
	if err := c.Get("api/v2/monitor/log/fortianalyzer-queue", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res {
		m = append(m, prometheus.MustNewConstMetric(logAnaConn, prometheus.GaugeValue, r.Results.Connected, r.VDOM))
		m = append(m, prometheus.MustNewConstMetric(logAnaLogs, prometheus.GaugeValue, r.Results.FailedLogs, r.VDOM, "Failed"))
		m = append(m, prometheus.MustNewConstMetric(logAnaLogs, prometheus.GaugeValue, r.Results.CachedLogs, r.VDOM, "Cached"))
	}

	return m, true
}
