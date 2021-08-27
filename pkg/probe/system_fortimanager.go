package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type SystemFortimanagerResults struct {
	Mode            string  `json:"mode"`
	Status_ID       float64 `json:"fortimanager_status_id"`
	Registration_ID float64 `json:"fortimanager_registration_status_id"`
}

type SystemFortimanagerStatus struct {
	Results SystemFortimanagerResults `json:"results"`
	VDOM    string                    `json:"vdom"`
}

func probeSystemFortimanagerStatus(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		FortimanStat_id = prometheus.NewDesc(
			"fortigate_fortimanager_connection_status",
			"Fortimanager status ID",
			[]string{"vdom", "mode"}, nil,
		)
		FortimanReg_id = prometheus.NewDesc(
			"fortigate_fortimanager_registration_status",
			"Fortimanager registration status ID",
			[]string{"vdom", "mode"}, nil,
		)
	)

	var res SystemFortimanagerStatus
	if err := c.Get("api/v2/monitor/system/fortimanager/status", "vdom=root", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	m = append(m, prometheus.MustNewConstMetric(FortimanStat_id, prometheus.GaugeValue, res.Results.Status_ID, res.VDOM, res.Results.Mode))
	m = append(m, prometheus.MustNewConstMetric(FortimanReg_id, prometheus.GaugeValue, res.Results.Registration_ID, res.VDOM, res.Results.Mode))

	return m, true
}
