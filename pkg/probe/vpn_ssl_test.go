package probe

import (
	"flag"
	"strings"
	"testing"

	"github.com/bluecmd/fortigate_exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestVPNSsl(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/vpn/ssl", "testdata/vpn.jsonnet")
	r := prometheus.NewPedanticRegistry()
	flag.Set("max-vpn-users", "10")
	config.MustReInit()
	if !testProbe(probeVPNSsl, c, r) {
		t.Errorf("probeSystemStatus() returned non-success")
	}

	em := `
	# HELP fortigate_vpn_connections Number of VPN connections
	# TYPE fortigate_vpn_connections gauge
	fortigate_vpn_connections{vdom="root"} 1
	# HELP fortigate_vpn_users Users of VPN connections
	# TYPE fortigate_vpn_users gauge
	fortigate_vpn_users{user="user1",vdom="root"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
