package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/internal/config"
	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type VPNUser struct {
	UserName string `json:"user_name"`
}

type VPNUsers struct {
	Results []VPNUser `json:"results"`
	VDOM    string    `json:"vdom"`
	Version string    `json:"version"`
}

func probeVPNSsl(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	savedConfig := config.GetConfig()
	MaxVPNUsers := savedConfig.MaxVPNUsers

	var (
		vpncon = prometheus.NewDesc(
			"fortigate_vpn_connections",
			"Number of VPN connections",
			[]string{"vdom"}, nil,
		)
		vpnusr = prometheus.NewDesc(
			"fortigate_vpn_users",
			"Users of VPN connections",
			[]string{"vdom", "user"}, nil,
		)
	)

	var res []VPNUsers
	if err := c.Get("api/v2/monitor/vpn/ssl", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res {
		count := len(r.Results)

		m = append(m, prometheus.MustNewConstMetric(vpncon, prometheus.GaugeValue, float64(count), r.VDOM))

		if MaxVPNUsers != 0 {
			if count > MaxVPNUsers {
				log.Printf("Error: Received more VPN Users than maximum (%d > %d) allowed, ignoring metric ...", count, MaxVPNUsers)
			} else {
				for _, result := range r.Results {
					m = append(m, prometheus.MustNewConstMetric(vpnusr, prometheus.GaugeValue, float64(1), r.VDOM, result.UserName))
				}
			}
		}
	}

	return m, true
}
