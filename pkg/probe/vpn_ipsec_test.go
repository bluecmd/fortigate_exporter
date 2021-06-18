package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestVPNIPSec(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/vpn/ipsec", "testdata/ipsec.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeVPNIPSec, c, r) {
		t.Errorf("probeVPNIPSec() returned non-success")
	}

	em := `
	# HELP fortigate_ipsec_tunnel_receive_bytes_total Total number of bytes received over the IPsec tunnel
	# TYPE fortigate_ipsec_tunnel_receive_bytes_total counter
	fortigate_ipsec_tunnel_receive_bytes_total{name="tunnel_1-sub",p2serial="1",parent="tunnel_1",vdom="root"} 1.429824e+07
	fortigate_ipsec_tunnel_receive_bytes_total{name="tunnel_1-sub",p2serial="12",parent="tunnel_1",vdom="root"} 1.429824e+07
	# HELP fortigate_ipsec_tunnel_transmit_bytes_total Total number of bytes transmitted over the IPsec tunnel
	# TYPE fortigate_ipsec_tunnel_transmit_bytes_total counter
	fortigate_ipsec_tunnel_transmit_bytes_total{name="tunnel_1-sub",p2serial="1",parent="tunnel_1",vdom="root"} 1.424856e+07
	fortigate_ipsec_tunnel_transmit_bytes_total{name="tunnel_1-sub",p2serial="12",parent="tunnel_1",vdom="root"} 1.424856e+07
	# HELP fortigate_ipsec_tunnel_up Status of IPsec tunnel
	# TYPE fortigate_ipsec_tunnel_up gauge
	fortigate_ipsec_tunnel_up{name="tunnel_1-sub",p2serial="1",parent="tunnel_1",vdom="root"} 1
	fortigate_ipsec_tunnel_up{name="tunnel_1-sub",p2serial="12",parent="tunnel_1",vdom="root"} 0

	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
