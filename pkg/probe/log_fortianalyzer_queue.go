package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type LogAnaQueueResults struct {
	Connected  int `json:"connected"`
	FailedLogs int `json:"failed_logs"`
	CachedLogs int `json:"cached_logs"`
}

type LogAnaQueue struct {
	Results LogAnaQueueResults `json:"results"`
	VDOM    string             `json:"vdom"`
	Version string             `json:"version"`
}

func probeLogAnalyzerQueue(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		logAnaConn = prometheus.NewDesc(
			"fortigate_log_fortianalyzer_queue_connected",
			"Fortianalyzer queue connected state",
			[]string{"vdom"}, nil,
		)
		logAnaFail = prometheus.NewDesc(
			"fortigate_log_fortianalyzer_queue_failed",
			"Failed logs in fortianalyzer queue",
			[]string{"vdom"}, nil,
		)
		logAnaCach = prometheus.NewDesc(
			"fortigate_log_fortianalyzer_queue_cached",
			"Cached logs in fortianalyzer queue",
			[]string{"vdom"}, nil,
		)
	)

	var res []LogAnaQueue
	if err := c.Get("api/v2/monitor/log/fortianalyzer-queue", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res {
		m = append(m, prometheus.MustNewConstMetric(logAnaConn, prometheus.GaugeValue, float64(r.Results.Connected), r.VDOM))
		m = append(m, prometheus.MustNewConstMetric(logAnaFail, prometheus.GaugeValue, float64(r.Results.FailedLogs), r.VDOM))
		m = append(m, prometheus.MustNewConstMetric(logAnaCach, prometheus.GaugeValue, float64(r.Results.CachedLogs), r.VDOM))
	}

	return m, true
}
