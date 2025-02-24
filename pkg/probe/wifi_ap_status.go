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

func probeWifiAPStatus(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		wtpCount = prometheus.NewDesc(
			"fortigate_wifi_access_points",
			"Number of connected access points by status",
			[]string{"vdom", "status"}, nil,
		)
		wtpClientCount = prometheus.NewDesc(
			"fortigate_wifi_fabric_clients",
			"Number of connected clients",
			[]string{"vdom"}, nil,
		)
		wtpMaxClientCount = prometheus.NewDesc(
			"fortigate_wifi_fabric_max_allowed_clients",
			"Maximum number of clients which are allowed to connect",
			[]string{"vdom"}, nil,
		)
	)

	type ApiStatusResponse []struct {
		Results struct {
			WtpSessionCount float64 `json:"wtp_session_count"`
			WtpActive       float64 `json:"wtp_active"`
			WtpDown         float64 `json:"wtp_down"`
			WtpRebooted     float64 `json:"wtp_rebooted"`
			ClientCount     float64 `json:"client_count"`
			ClientCountMax  float64 `json:"client_count_max"`
		} `json:"results"`
		VDOM string `json:"vdom"`
	}

	var response ApiStatusResponse
	if err := c.Get("api/v2/monitor/wifi/ap_status", "vdom=*", &response); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	var m []prometheus.Metric

	for _, rs := range response {
		m = append(m, prometheus.MustNewConstMetric(wtpCount, prometheus.GaugeValue, rs.Results.WtpActive, rs.VDOM, "active"))
		m = append(m, prometheus.MustNewConstMetric(wtpCount, prometheus.GaugeValue, rs.Results.WtpDown, rs.VDOM, "down"))
		m = append(m, prometheus.MustNewConstMetric(wtpCount, prometheus.GaugeValue, rs.Results.WtpRebooted, rs.VDOM, "rebooting"))
		m = append(m, prometheus.MustNewConstMetric(wtpClientCount, prometheus.GaugeValue, rs.Results.ClientCount, rs.VDOM))
		m = append(m, prometheus.MustNewConstMetric(wtpMaxClientCount, prometheus.GaugeValue, rs.Results.ClientCountMax, rs.VDOM))
	}

	return m, true
}
