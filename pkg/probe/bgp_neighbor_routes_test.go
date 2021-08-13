package probe

import (
	"strings"
	"testing"

	"github.com/bluecmd/fortigate_exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestBGPNeighborPathsIPv4(t *testing.T) {
	config.MustReInit()
	c := newFakeClient()
	c.prepare("api/v2/monitor/router/bgp/paths", "testdata/router-bgp-paths-v4.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeBGPNeighborPathsIPv4, c, r) {
		t.Errorf("probeBGPNeighborPathsIPv4() returned non-success")
	}

	em := `
    # HELP fortigate_bgp_neighbor_ipv4_best_paths Count of best paths for an BGP neighbor
    # TYPE fortigate_bgp_neighbor_ipv4_best_paths gauge
    fortigate_bgp_neighbor_ipv4_best_paths{neighbor_ip="10.0.0.1",vdom="root"} 1
    fortigate_bgp_neighbor_ipv4_best_paths{neighbor_ip="10.0.0.2",vdom="root"} 1
    # HELP fortigate_bgp_neighbor_ipv4_paths Count of paths received from an BGP neighbor
    # TYPE fortigate_bgp_neighbor_ipv4_paths gauge
    fortigate_bgp_neighbor_ipv4_paths{neighbor_ip="10.0.0.1",vdom="root"} 1
    fortigate_bgp_neighbor_ipv4_paths{neighbor_ip="10.0.0.2",vdom="root"} 2
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}

func TestBGPNeighborPathsIPv6(t *testing.T) {

	if err := config.Init(); err != nil {
		t.Fatalf("config.Init failed: %+v", err)
	}

	c := newFakeClient()
	c.prepare("api/v2/monitor/router/bgp/paths6", "testdata/router-bgp-paths-v6.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeBGPNeighborPathsIPv6, c, r) {
		t.Errorf("probeBGPNeighborPathsIPv6() returned non-success")
	}

	em := `
    # HELP fortigate_bgp_neighbor_ipv6_best_paths Count of best paths for an BGP neighbor
    # TYPE fortigate_bgp_neighbor_ipv6_best_paths gauge
    fortigate_bgp_neighbor_ipv6_best_paths{neighbor_ip="::",vdom="root"} 1
    fortigate_bgp_neighbor_ipv6_best_paths{neighbor_ip="fd00::1",vdom="root"} 2
    # HELP fortigate_bgp_neighbor_ipv6_paths Count of paths received from an BGP neighbor
    # TYPE fortigate_bgp_neighbor_ipv6_paths gauge
    fortigate_bgp_neighbor_ipv6_paths{neighbor_ip="::",vdom="root"} 1
    fortigate_bgp_neighbor_ipv6_paths{neighbor_ip="fd00::1",vdom="root"} 3
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
