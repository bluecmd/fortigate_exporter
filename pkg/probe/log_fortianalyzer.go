package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type LogAnaResults struct {
	Registration string `json:"registration"`
	Connection   string `json:"connection"`
	Received     int    `json:"received"`
}

type LogAna struct {
	Results LogAnaResults `json:"results"`
	VDOM    string        `json:"vdom"`
	Version string        `json:"version"`
}

func probeLogAnalyzer(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		logAnaReg = prometheus.NewDesc(
			"fortigate_log_fortianalyzer_registration",
			"Fortianalyzer registration state",
			[]string{"vdom", "registration"}, nil,
		)
		logAnaCon = prometheus.NewDesc(
			"fortigate_log_fortianalyzer_connection",
			"Fortianalyzer connection state",
			[]string{"vdom", "connection"}, nil,
		)
		logAnaRcv = prometheus.NewDesc(
			"fortigate_log_fortianalyzer_received",
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
		m = append(m, prometheus.MustNewConstMetric(logAnaReg, prometheus.GaugeValue, float64(1), r.VDOM, r.Results.Registration))
		m = append(m, prometheus.MustNewConstMetric(logAnaCon, prometheus.GaugeValue, float64(1), r.VDOM, r.Results.Connection))
		m = append(m, prometheus.MustNewConstMetric(logAnaRcv, prometheus.GaugeValue, float64(r.Results.Received), r.VDOM))
	}

	return m, true
}
