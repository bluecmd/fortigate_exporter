package probe

import (
	"log"
	"strconv"

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
		FortimanStat_info = prometheus.NewDesc(
			"fortigate_fortimanager_info",
			"Fortimanager infos",
			[]string{"vdom", "mode", "connection_status", "registration_status"}, nil,
		)
	)

	var res []SystemFortimanagerStatus
	if err := c.Get("api/v2/monitor/system/fortimanager/status", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res {
		m = append(m, prometheus.MustNewConstMetric(FortimanStat_info, prometheus.GaugeValue, float64(1), r.VDOM, r.Results.Mode, strconv.Itoa(r.Results.Status_ID), strconv.Itoa(r.Results.Registration_ID)))
	}

	return m, true
}
