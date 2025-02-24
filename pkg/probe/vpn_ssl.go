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

	"github.com/prometheus-community/fortigate_exporter/internal/config"
	"github.com/prometheus-community/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type VPNUser struct {
	UserName string `json:"user_name"`
}

type VPNUsers struct {
	Results []VPNUser `json:"results"`
	VDOM    string    `json:"vdom"`
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
			"Number of VPN users connections",
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
				// Structure for summarizing multi VPN per user
				type VPNUserDesc struct {
					VDOM     string
					UserName string
				}
				userMap := map[VPNUserDesc]float64{}

				for _, result := range r.Results {
					userDesc := VPNUserDesc{r.VDOM, result.UserName}
					userMap[userDesc]++
				}
				for userDesc, counter := range userMap {
					m = append(m, prometheus.MustNewConstMetric(vpnusr, prometheus.GaugeValue, counter, userDesc.VDOM, userDesc.UserName))
				}
			}
		}
	}

	return m, true
}
