package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

func probeVPNSsl(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		vpncon = prometheus.NewDesc(
			"fortigate_vpn_connections",
			"Number of VPN connections",
			[]string{"vdom"}, nil,
		)
	)

	type result struct {
		Results []map[string]interface{}
		VDOM    string
	}
	var res []result
	if err := c.Get("api/v2/monitor/vpn/ssl", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, v := range res {
		count := len(v.Results)
		m = append(m, prometheus.MustNewConstMetric(vpncon, prometheus.GaugeValue, float64(count), v.VDOM))
	}

	return m, true
}
