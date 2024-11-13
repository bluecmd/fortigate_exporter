package probe

import (
	"log"
	"strconv"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

func probeSystemHAPeer(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		memberInfo = prometheus.NewDesc(
			"fortigate_ha_peer_info",
			"Info metrics regarding cluster HA peers",
			[]string{"vdom", "hostname", "serial", "priority", "vcluster_id", "primary"}, nil,
		)
	)

	type HAResults struct {
		Hostname   string `json:"hostname"`
		SerialNo   string `json:"serial_no"`
		VclusterID int    `json:"vcluster_id"`
		Priority   int    `json:"priority"`
		Primary    bool   `json:"primary"` // Set to true if primary node in FortiOS 7.4+
	}

	type HAResponse struct {
		HTTPMethod string      `json:"http_method"`
		Results    []HAResults `json:"results"`
		VDOM       string      `json:"vdom"`
		Path       string      `json:"path"`
		Name       string      `json:"name"`
		Status     string      `json:"status"`
		Serial     string      `json:"serial"`
		Version    string      `json:"version"`
		Build      int64       `json:"build"`
	}
	var r HAResponse

	if err := c.Get("api/v2/monitor/system/ha-peer", "", &r); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, result := range r.Results {
		strPrimary := strconv.FormatBool(result.Primary)
		if meta.VersionMajor < 7 || (meta.VersionMajor == 7 && meta.VersionMinor < 4) {
			strPrimary = "Unsupported"
		}
		m = append(m, prometheus.MustNewConstMetric(memberInfo, prometheus.GaugeValue, 1, r.VDOM, result.Hostname, result.SerialNo, strconv.Itoa(result.Priority), strconv.Itoa(result.VclusterID), strPrimary))
	}
	return m, true
}
