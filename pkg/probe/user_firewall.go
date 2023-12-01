package probe

import (
	"log"
	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type UserFirewallResult struct {
	Type            string `json:"type"`
	ID              int    `json:"id"`
	DurationSecs    int    `json:"duration_secs"`
	AuthType        int    `json:"auth_type"`
	IPAddr          string `json:"ipaddr"`
	SrcType         string `json:"src_type"`
	ExpirySecs      int    `json:"expiry_secs"`
	TrafficVolBytes int64  `json:"traffic_vol_bytes"`
	Method          string `json:"method"`
}

type UserFirewall struct {
    HttpMethod string               `json:"http_method"`
    Results    []UserFirewallResult `json:"results"`
    VDOM       string               `json:"vdom"`
    Path       string               `json:"path"`
    Name       string               `json:"name"`
    Action     string               `json:"action"`
    Status     string               `json:"status"`
    Serial     string               `json:"serial"`
    Version    string               `json:"version"`
    Build      int                  `json:"build"`
}


func probeUserFirewall(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		userFirewallDuration = prometheus.NewDesc(
			"fortigate_user_firewall_duration_seconds",
			"Duration of user firewall activity in seconds",
			[]string{"vdom", "ipaddr", "method", "type"}, nil,
		)
		userFirewallTraffic = prometheus.NewDesc(
			"fortigate_user_firewall_traffic_bytes",
			"Traffic volume in bytes for user firewall activity",
			[]string{"vdom", "ipaddr", "method", "type"}, nil,
		)
	)

    var res []UserFirewall
    if err := c.Get("/api/v2/monitor/user/firewall", "vdom=*", &res); err != nil {
        log.Printf("Error: %v", err)
        return nil, false
    }

    metrics := []prometheus.Metric{}
    for _, fw := range res { 
        for _, r := range fw.Results {
            metrics = append(metrics, prometheus.MustNewConstMetric(
                userFirewallDuration, prometheus.GaugeValue, float64(r.DurationSecs), fw.VDOM, r.IPAddr, r.Method, r.Type,
            ))
            metrics = append(metrics, prometheus.MustNewConstMetric(
                userFirewallTraffic, prometheus.GaugeValue, float64(r.TrafficVolBytes), fw.VDOM, r.IPAddr, r.Method, r.Type,
            ))
        }
    }

	return metrics, true
}
