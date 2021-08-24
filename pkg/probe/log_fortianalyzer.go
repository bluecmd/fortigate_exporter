package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type LogAnaResults struct {
	Registration string  `json:"registration"`
	Connection   string  `json:"connection"`
	Received     float64 `json:"received"`
}

type LogAna struct {
	Results LogAnaResults `json:"results"`
	VDOM    string        `json:"vdom"`
}

func probeLogAnalyzer(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		logAnaInfo = prometheus.NewDesc(
			"fortigate_log_fortianalyzer_registration_info",
			"Fortianalyzer state info",
			[]string{"vdom", "registration", "connection"}, nil,
		)
		logAnaRcv = prometheus.NewDesc(
			"fortigate_log_fortianalyzer_logs_received",
			"Received logs in fortianalyzer",
			[]string{"vdom"}, nil,
		)
	)

	var res []LogAna
	if err := c.Get("api/v2/monitor/log/fortianalyzer", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res {
		m = append(m, prometheus.MustNewConstMetric(logAnaInfo, prometheus.GaugeValue, float64(1), r.VDOM, r.Results.Registration, r.Results.Connection))
		m = append(m, prometheus.MustNewConstMetric(logAnaRcv, prometheus.GaugeValue, r.Results.Received, r.VDOM))
	}

	return m, true
}
