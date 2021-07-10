
package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestBgpNeighborsIPv6(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/router/bgp/neighbors6", "testdata/bpg-neighbors-v6.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeBgpNeighborsIPv6, c, r) {
		t.Errorf("probeBgpNeighborsIPv6() returned non-success")
	}

	em := `
	# HELP fortigate_bgp_neighbors_ipv6 Confiured bgp neighbors over ipv6
    # TYPE fortigate_bgp_neighbors_ipv6 gauge
    fortigate_bgp_neighbors_ipv6{admin_status="true",local_ip="fd00::1",neighbor_ip="fd00::2",remote_as="1337",state="Established",type="ipv6",vdom="root"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}