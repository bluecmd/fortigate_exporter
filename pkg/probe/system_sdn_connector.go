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

type SystemSDNConnectorResults struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Status     string `json:"status"`
	Updating   bool   `json:"updating"`
	LastUpdate int    `json:"last_update"`
}

type SystemSDNConnector struct {
	Results []SystemSDNConnectorResults `json:"results"`
	VDOM    string                      `json:"vdom"`
}

func probeSystemSDNConnector(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		SDNConnectorsStatus = prometheus.NewDesc(
			"fortigate_system_sdn_connector_status",
			"Status of SDN connectors (0=Disabled, 1=Down, 2=Unknown, 3=Up, 4=Updating)",
			[]string{"vdom", "name", "type"}, nil,
		)
		SDNConnectorsLastUpdate = prometheus.NewDesc(
			"fortigate_system_sdn_connector_last_update_seconds",
			"Last update time for SDN connectors (in seconds from epoch)",
			[]string{"vdom", "name", "type"}, nil,
		)
	)

	var res []SystemSDNConnector
	if err := c.Get("api/v2/monitor/system/sdn-connector/status", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res {
		for _, sdnConn := range r.Results {
			if sdnConn.Status == "Disabled" {
				m = append(m, prometheus.MustNewConstMetric(SDNConnectorsStatus, prometheus.GaugeValue, float64(0), r.VDOM, sdnConn.Name, sdnConn.Type))
			} else if sdnConn.Status == "Down" {
				m = append(m, prometheus.MustNewConstMetric(SDNConnectorsStatus, prometheus.GaugeValue, float64(1), r.VDOM, sdnConn.Name, sdnConn.Type))
			} else if sdnConn.Status == "Unknown" {
				m = append(m, prometheus.MustNewConstMetric(SDNConnectorsStatus, prometheus.GaugeValue, float64(2), r.VDOM, sdnConn.Name, sdnConn.Type))
			} else if sdnConn.Status == "Up" {
				m = append(m, prometheus.MustNewConstMetric(SDNConnectorsStatus, prometheus.GaugeValue, float64(3), r.VDOM, sdnConn.Name, sdnConn.Type))
			} else if sdnConn.Status == "Updating" {
				m = append(m, prometheus.MustNewConstMetric(SDNConnectorsStatus, prometheus.GaugeValue, float64(4), r.VDOM, sdnConn.Name, sdnConn.Type))
			}
			m = append(m, prometheus.MustNewConstMetric(SDNConnectorsLastUpdate, prometheus.GaugeValue, float64(sdnConn.LastUpdate), r.VDOM, sdnConn.Name, sdnConn.Type))
		}
	}

	return m, true
}
