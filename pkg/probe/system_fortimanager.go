package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type SystemFortimanagerResults struct {
	Mode            string `json:"mode"`
	Status_ID       int    `json:"fortimanager_status_id"`
	Registration_ID int    `json:"fortimanager_registration_status_id"`
}

type SystemFortimanagerStatus struct {
	Results SystemFortimanagerResults `json:"results"`
	VDOM    string                    `json:"vdom"`
}

func probeSystemFortimanagerStatus(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		FortimanStat_id = prometheus.NewDesc(
			"fortigate_fortimanager_connection_status",
			"Fortimanager tunnel status (0 - Down, 1 - Handshake in progress, 2 - Up)",
			[]string{"vdom", "mode"}, nil,
		)
		FortimanReg_id = prometheus.NewDesc(
			"fortigate_fortimanager_registration_status",
			"Fortimanager registration status (0 - Unknown Registration, 1 - In Progress, 2 - Registered, 3 - Registered but Unauthorized)",
			[]string{"vdom", "mode"}, nil,
		)
	)

	var res []SystemFortimanagerStatus
	if err := c.Get("api/v2/monitor/system/fortimanager/status", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res {
		m = append(m, prometheus.MustNewConstMetric(FortimanStat_id, prometheus.GaugeValue, float64(r.Results.Status_ID), r.VDOM, r.Results.Mode))
		m = append(m, prometheus.MustNewConstMetric(FortimanReg_id, prometheus.GaugeValue, float64(r.Results.Registration_ID), r.VDOM, r.Results.Mode))
	}

	return m, true
}
