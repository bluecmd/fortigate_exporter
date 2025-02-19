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
	"strconv"

	"github.com/prometheus-community/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type BGPNeighbor struct {
	NeighborIP  string `json:"neighbor_ip"`
	LocalIP     string `json:"local_ip"`
	RemoteAS    int    `json:"remote_as"`
	AdminStatus bool   `json:"admin_status"`
	State       string `json:"state"`
}

type BGPNeighborResponse struct {
	Results []BGPNeighbor `json:"results"`
	VDOM    string        `json:"vdom"`
	Version string        `json:"version"`
}

func probeBGPNeighborsIPv4(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	if meta.VersionMajor < 7 {
		// not supported version. Before 7.0.0 the requested endpoint doesn't exist
		return nil, true
	}
	var (
		mBGPNeighbor = prometheus.NewDesc(
			"fortigate_bgp_neighbor_ipv4_info",
			"Configured bgp neighbor over ipv4, return state as value (1 - Idle, 2 - Connect, 3 - Active, 4 - Open sent, 5 - Open confirm, 6 - Established)",
			[]string{"vdom", "remote_as", "state", "admin_status", "local_ip", "neighbor_ip"}, nil,
		)
	)

	var rs []BGPNeighborResponse

	if err := c.Get("api/v2/monitor/router/bgp/neighbors", "vdom=*", &rs); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}

	for _, r := range rs {
		for _, peer := range r.Results {
			m = append(m, prometheus.MustNewConstMetric(mBGPNeighbor, prometheus.GaugeValue, bgpStateToNumber(peer.State), r.VDOM, strconv.Itoa(peer.RemoteAS), peer.State, strconv.FormatBool(peer.AdminStatus), peer.LocalIP, peer.NeighborIP))
		}
	}

	return m, true
}

func probeBGPNeighborsIPv6(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	if meta.VersionMajor < 7 {
		// not supported version. Before 7.0.0 the requested endpoint doesn't exist
		return nil, true
	}

	var (
		mBGPNeighbor = prometheus.NewDesc(
			"fortigate_bgp_neighbor_ipv6_info",
			"Configured bgp neighbor over ipv6, return state as value (1 - Idle, 2 - Connect, 3 - Active, 4 - Open sent, 5 - Open confirm, 6 - Established)",
			[]string{"vdom", "remote_as", "state", "admin_status", "local_ip", "neighbor_ip"}, nil,
		)
	)

	var rs []BGPNeighborResponse

	if err := c.Get("api/v2/monitor/router/bgp/neighbors6", "vdom=*", &rs); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}

	for _, r := range rs {
		for _, peer := range r.Results {
			m = append(m, prometheus.MustNewConstMetric(mBGPNeighbor, prometheus.GaugeValue, bgpStateToNumber(peer.State), r.VDOM, strconv.Itoa(peer.RemoteAS), peer.State, strconv.FormatBool(peer.AdminStatus), peer.LocalIP, peer.NeighborIP))
		}
	}

	return m, true
}

func bgpStateToNumber(bgpState string) float64 {
	switch bgpState {
	case "Idle":
		return 1
	case "Connect":
		return 2
	case "Active":
		return 3
	case "Open sent":
		return 4
	case "Open confirm":
		return 5
	case "Established":
		return 6
	default:
		return 0
	}
}
