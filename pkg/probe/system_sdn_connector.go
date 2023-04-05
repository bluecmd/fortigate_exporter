package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type SystemSdnConnectorResults struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Updating bool `json:"updating"`
	LastUpdate int `json:"last_update"`
}

type SystemSdnConnector struct {
	Results []SystemSdnConnectorResults `json:"results"`
	VDOM    string            `json:"vdom"`
}

func probeSystemSdnConnector(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		SdnConnectorsStatus = prometheus.NewDesc(
			"fortigate_system_sdn_connector_status",
			"Status of SDN connectors (0=Disabled, 1=Down, 2=Unknown, 3=Up, 4=Updating)",
			[]string{"vdom", "name", "type"}, nil,
		)
		SdnConnectorsLastUpdate = prometheus.NewDesc(
			"fortigate_system_sdn_connector_last_update",
			"Last update time for SDN connectors",
			[]string{"vdom", "name", "type"}, nil,
		)
	)

	var res []SystemSdnConnector
	if err := c.Get("api/v2/monitor/system/sdn-connector/status", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res {
		for _, sdnConn := range r.Results {
			if sdnConn.Status == "Disabled" {
				m = append(m, prometheus.MustNewConstMetric(SdnConnectorsStatus, prometheus.GaugeValue, float64(0), r.VDOM, sdnConn.Name, sdnConn.Type))
			} else if sdnConn.Status == "Down" {
				m = append(m, prometheus.MustNewConstMetric(SdnConnectorsStatus, prometheus.GaugeValue, float64(1), r.VDOM, sdnConn.Name, sdnConn.Type))
			} else if sdnConn.Status == "Unknown" {
				m = append(m, prometheus.MustNewConstMetric(SdnConnectorsStatus, prometheus.GaugeValue, float64(2), r.VDOM, sdnConn.Name, sdnConn.Type))
			} else if sdnConn.Status == "Up" {
				m = append(m, prometheus.MustNewConstMetric(SdnConnectorsStatus, prometheus.GaugeValue, float64(3), r.VDOM, sdnConn.Name, sdnConn.Type))
			} else if sdnConn.Status == "Updating" {
				m = append(m, prometheus.MustNewConstMetric(SdnConnectorsStatus, prometheus.GaugeValue, float64(4), r.VDOM, sdnConn.Name, sdnConn.Type))
			}
			m = append(m, prometheus.MustNewConstMetric(SdnConnectorsLastUpdate, prometheus.GaugeValue, float64(sdnConn.LastUpdate), r.VDOM, sdnConn.Name, sdnConn.Type))
		}
	}

	return m, true
}
