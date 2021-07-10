
package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestBgpNeighborsIPv4(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/router/bgp/neighbors", "testdata/bpg-neighbors-v4.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeBgpNeighborsIPv4, c, r) {
		t.Errorf("probeBgpNeighborsIPv4() returned non-success")
	}

	em := `
	# HELP fortigate_bgp_neighbors_ipv4 Confiured bgp neighbors over ipv4
    # TYPE fortigate_bgp_neighbors_ipv4 gauge
    fortigate_bgp_neighbors_ipv4{admin_status="true",local_ip="10.0.0.0",neighbor_ip="10.0.0.1",remote_as="1337",state="Established",type="ipv4",vdom="root"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}