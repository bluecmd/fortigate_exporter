package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestFirewallIpPool(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/firewall/ippool", "testdata/fw-ippool.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeFirewallIpPool, c, r) {
		t.Errorf("probeFirewallIpPool() returned non-success")
	}

	em := `
	# HELP fortigate_ippool_available_ratio Percentage available in ippool (0 - 1.0)
	# TYPE fortigate_ippool_available_ratio gauge
	fortigate_ippool_available_ratio{name="ippool_name",vdom="FG-traffic"} 1
	# HELP fortigate_ippool_clients Amount of clients using ippool
	# TYPE fortigate_ippool_clients gauge
	fortigate_ippool_clients{name="ippool_name",vdom="FG-traffic"} 0
	# HELP fortigate_ippool_total_ips Ip addresses total in ippool
	# TYPE fortigate_ippool_total_ips gauge
	fortigate_ippool_total_ips{name="ippool_name",vdom="FG-traffic"} 1
	# HELP fortigate_ippool_total_items Amount of items total in ippool
	# TYPE fortigate_ippool_total_items gauge
	fortigate_ippool_total_items{name="ippool_name",vdom="FG-traffic"} 472
	# HELP fortigate_ippool_used_ips Ip addresses in use in ippool
	# TYPE fortigate_ippool_used_ips gauge
	fortigate_ippool_used_ips{name="ippool_name",vdom="FG-traffic"} 0
	# HELP fortigate_ippool_used_items Amount of items used in ippool
	# TYPE fortigate_ippool_used_items gauge
	fortigate_ippool_used_items{name="ippool_name",vdom="FG-traffic"} 0
	# HELP fortigate_ippool_pba_per_ip Amount of available port block allocations per ip
    # TYPE fortigate_ippool_pba_per_ip gauge
    fortigate_ippool_pba_per_ip{name="ippool_name",vdom="FG-traffic"} 472
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
