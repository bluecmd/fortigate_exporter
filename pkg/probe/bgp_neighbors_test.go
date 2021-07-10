
package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestBgpPeers(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/router/bgp/neighbors", "testdata/bpg-peers-v4.jsonnet")
	c.prepare("api/v2/monitor/router/bgp/neighbors6", "testdata/bpg-peers-v6.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeBgpNeighbors, c, r) {
		t.Errorf("probeBgpNeighbors() returned non-success")
	}

	em := `
	# HELP fortigate_bgp_neighbor Recived bgp neighbor
    # TYPE fortigate_bgp_neighbor gauge
    fortigate_bgp_neighbor{admin_status="true",local_ip="10.0.0.0",neighbor_ip="10.0.0.1",remote_as="1337",state="Established",type="ipv4",vdom="root"} 1
    fortigate_bgp_neighbor{admin_status="true",local_ip="fd00::1",neighbor_ip="fd00::2",remote_as="1337",state="Established",type="ipv6",vdom="root"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}