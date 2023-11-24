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
    fortigate_ospf_neighbor_info{neighbor_ip="10.0.0.1",priority="3",router_id="12345",state="Full",vdom="root"} 7
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
