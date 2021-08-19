package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestVPNSslStats(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/vpn/ssl/stats", "testdata/vpn-stats.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeVPNSslStats, c, r) {
		t.Errorf("probeVPNSslStats() returned non-success")
	}

	em := `
	# HELP fortigate_vpn_ssl_connections Number of current SSL VPN connections
	# TYPE fortigate_vpn_ssl_connections gauge
	fortigate_vpn_ssl_connections{vdom="root"} 2
	# HELP fortigate_vpn_ssl_tunnels Number of current SSL VPN tunnels
	# TYPE fortigate_vpn_ssl_tunnels gauge
	fortigate_vpn_ssl_tunnels{vdom="root"} 2
	# HELP fortigate_vpn_ssl_users Number of current SSL VPN users
	# TYPE fortigate_vpn_ssl_users gauge
	fortigate_vpn_ssl_users{vdom="root"} 3
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
