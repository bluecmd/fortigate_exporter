// Copyright 2025 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package probe

import (
	"log"

	"github.com/prometheus-community/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type VPNCurrentResults struct {
	Users       int `json:"users"`
	Tunnels     int `json:"tunnels"`
	Connections int `json:"connections"`
}

type VPNResults struct {
	Current VPNCurrentResults `json:"current"`
}

type VPNStats struct {
	Results VPNResults `json:"results"`
	VDOM    string     `json:"vdom"`
	Version string     `json:"version"`
}

func probeVPNSslStats(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		vpnCurUsr = prometheus.NewDesc(
			"fortigate_vpn_ssl_users",
			"Number of current SSL VPN users",
			[]string{"vdom"}, nil,
		)
		vpnCurTun = prometheus.NewDesc(
			"fortigate_vpn_ssl_tunnels",
			"Number of current SSL VPN tunnels",
			[]string{"vdom"}, nil,
		)
		vpnCurCon = prometheus.NewDesc(
			"fortigate_vpn_ssl_connections",
			"Number of current SSL VPN connections",
			[]string{"vdom"}, nil,
		)
	)

	var res []VPNStats
	if err := c.Get("api/v2/monitor/vpn/ssl/stats", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res {
		m = append(m, prometheus.MustNewConstMetric(vpnCurUsr, prometheus.GaugeValue, float64(r.Results.Current.Users), r.VDOM))
		m = append(m, prometheus.MustNewConstMetric(vpnCurTun, prometheus.GaugeValue, float64(r.Results.Current.Tunnels), r.VDOM))
		m = append(m, prometheus.MustNewConstMetric(vpnCurCon, prometheus.GaugeValue, float64(r.Results.Current.Connections), r.VDOM))
	}

	return m, true
}
