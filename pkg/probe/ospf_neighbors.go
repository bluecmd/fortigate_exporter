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

type OSPFNeighbor struct {
	NeighborIP string `json:"neighbor_ip"`
	RouterID   string `json:"router_id"`
	State      string `json:"state"`
	Priority   int    `json:"priority"`
}

type OSPFNeighborResponse struct {
	Results []OSPFNeighbor `json:"results"`
	VDOM    string         `json:"vdom"`
	Version string         `json:"version"`
}

func probeOSPFNeighbors(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	if meta.VersionMajor < 7 {
		// not supported version. Before 7.0.0 the requested endpoint doesn't exist
		return nil, true
	}
	var (
		mOSPFNeighbor = prometheus.NewDesc(
			"fortigate_ospf_neighbor_info",
			"List all discovered OSPF neighbors, return state as value (1 - Down, 2 - Attempt, 3 - Init, 4 - Two way, 5 - Exchange start, 6 - Exchange, 7 - Loading, 8 - Full)",
			[]string{"vdom", "state", "priority", "router_id", "neighbor_ip"}, nil,
		)
	)

	var rs []OSPFNeighborResponse

	if err := c.Get("api/v2/monitor/router/ospf/neighbors", "vdom=*", &rs); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}

	for _, r := range rs {
		for _, peer := range r.Results {
			m = append(m, prometheus.MustNewConstMetric(mOSPFNeighbor, prometheus.GaugeValue, ospfStateToNumber(peer.State), r.VDOM, peer.State, strconv.Itoa(peer.Priority), peer.RouterID, peer.NeighborIP))
		}
	}

	return m, true
}

func ospfStateToNumber(ospfState string) float64 {
	switch ospfState {
	case "Down":
		return 1
	case "Attempt":
		return 2
	case "Init":
		return 3
	case "Two way":
		return 4
	case "Exchange start":
		return 5
	case "Exchange":
		return 6
	case "Loading":
		return 7
	case "Full":
		return 8
	default: // Down
		return 1
	}
}
