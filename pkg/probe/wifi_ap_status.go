package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

func probeWifiAPStatus(c http.FortiHTTP) ([]prometheus.Metric, bool) {
	var (
		wtpCount = prometheus.NewDesc(
			"fortigate_wifi_access_points",
			"Number of connected access points by status",
			[]string{"vdom", "status"}, nil,
		)
		wtpClientCount = prometheus.NewDesc(
			"fortigate_wifi_access_point_client_count",
			"Number of connected clients",
			[]string{"vdom"}, nil,
		)
		wtpMaxClientCount = prometheus.NewDesc(
			"fortigate_wifi_access_point_allowed_client_count",
			"Maximum number of clients which are allowed to connect",
			[]string{"vdom"}, nil,
		)
	)

	type ApiStatusResponse struct {
		Results struct {
			WtpSessionCount float64 `json:"wtp_session_count"`
			WtpActive       float64 `json:"wtp_active"`
			WtpDown         float64 `json:"wtp_down"`
			WtpRebooted     float64 `json:"wtp_rebooted"`
			ClientCount     float64 `json:"client_count"`
			ClientCountMax  float64 `json:"client_count_max"`
		} `json:"results"`
		VDOM string `json:"vdom"`
	}

	// Consider implementing pagination to remove this limit of 1000 entries
	var rs ApiStatusResponse
	if err := c.Get("api/v2/monitor/wifi/ap_status", "vdom=*&start=0&count=1000", &rs); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	var m []prometheus.Metric

	m = append(m, prometheus.MustNewConstMetric(wtpCount, prometheus.GaugeValue, rs.Results.WtpActive, rs.VDOM, "active"))
	m = append(m, prometheus.MustNewConstMetric(wtpCount, prometheus.GaugeValue, rs.Results.WtpDown, rs.VDOM, "down"))
	m = append(m, prometheus.MustNewConstMetric(wtpCount, prometheus.GaugeValue, rs.Results.WtpRebooted, rs.VDOM, "rebooting"))
	m = append(m, prometheus.MustNewConstMetric(wtpClientCount, prometheus.GaugeValue, rs.Results.ClientCount, rs.VDOM))
	m = append(m, prometheus.MustNewConstMetric(wtpMaxClientCount, prometheus.GaugeValue, rs.Results.ClientCountMax, rs.VDOM))

	return m, true
}
