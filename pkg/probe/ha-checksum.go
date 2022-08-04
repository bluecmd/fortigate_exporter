package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type HAChecksumResults struct {
	IsManageMaster int    `json:"is_manage_master"`
	IsRootMaster   int    `json:"is_root_master"`
	SerialNo       string `json:"serial_no"`
}

type HAChecksum struct {
	Results []HAChecksumResults `json:"results"`
}

func probeSystemHAChecksum(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		IsMaster = prometheus.NewDesc(
			"fortigate_ha_role",
			"Master/Slave information",
			[]string{"name", "serial_no"}, nil,
		)
	)

	var res HAChecksum
	if err := c.Get("api/v2/monitor/system/ha-checksums", "scope=global", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, response := range res.Results {
		m = append(m, prometheus.MustNewConstMetric(IsMaster, prometheus.GaugeValue, float64(response.IsManageMaster), "is_manage_master", response.SerialNo))
		m = append(m, prometheus.MustNewConstMetric(IsMaster, prometheus.GaugeValue, float64(response.IsRootMaster), "is_root_master", response.SerialNo))
	}

	return m, true
}
