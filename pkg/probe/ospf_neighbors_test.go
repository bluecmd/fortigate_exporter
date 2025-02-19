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
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestOSPFNeighborsIPv4(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/router/ospf/neighbors", "testdata/router-ospf-neighbors.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeOSPFNeighbors, c, r) {
		t.Errorf("probeOSPFNeighbors() returned non-success")
	}

	em := `
    # HELP fortigate_ospf_neighbor_info List all discovered OSPF neighbors, return state as value (1 - Down, 2 - Attempt, 3 - Init, 4 - Two way, 5 - Exchange start, 6 - Exchange, 7 - Loading, 8 - Full)
    # TYPE fortigate_ospf_neighbor_info gauge
    fortigate_ospf_neighbor_info{neighbor_ip="10.0.0.1",priority="3",router_id="12341",state="Down",vdom="root"} 1
    fortigate_ospf_neighbor_info{neighbor_ip="10.0.0.2",priority="3",router_id="12342",state="Attempt",vdom="root"} 2
    fortigate_ospf_neighbor_info{neighbor_ip="10.0.0.3",priority="3",router_id="12343",state="Init",vdom="root"} 3
    fortigate_ospf_neighbor_info{neighbor_ip="10.0.0.4",priority="3",router_id="12344",state="Two way",vdom="root"} 4
    fortigate_ospf_neighbor_info{neighbor_ip="10.0.0.5",priority="3",router_id="12345",state="Exchange start",vdom="root"} 5
    fortigate_ospf_neighbor_info{neighbor_ip="10.0.0.6",priority="3",router_id="12346",state="Exchange",vdom="root"} 6
    fortigate_ospf_neighbor_info{neighbor_ip="10.0.0.7",priority="3",router_id="12347",state="Loading",vdom="root"} 7
    fortigate_ospf_neighbor_info{neighbor_ip="10.0.0.8",priority="3",router_id="12348",state="Full",vdom="root"} 8
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
