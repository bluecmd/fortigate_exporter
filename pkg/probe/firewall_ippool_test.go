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
	# HELP fortigate_ippool_available Percentage available in ippool
	# TYPE fortigate_ippool_available gauge
	fortigate_ippool_available{name="ippool_name",vdom="FG-traffic"} 100
	# HELP fortigate_ippool_clients Amount of clients using ippool
	# TYPE fortigate_ippool_clients gauge
	fortigate_ippool_clients{name="ippool_name",vdom="FG-traffic"} 0
	# HELP fortigate_ippool_total_ip Ip addresses total in ippool
	# TYPE fortigate_ippool_total_ip gauge
	fortigate_ippool_total_ip{name="ippool_name",vdom="FG-traffic"} 1
	# HELP fortigate_ippool_total_items Amount of items total in ippool
	# TYPE fortigate_ippool_total_items gauge
	fortigate_ippool_total_items{name="ippool_name",vdom="FG-traffic"} 472
	# HELP fortigate_ippool_used_ip Ip addresses in use in ippool
	# TYPE fortigate_ippool_used_ip gauge
	fortigate_ippool_used_ip{name="ippool_name",vdom="FG-traffic"} 0
	# HELP fortigate_ippool_used_items Amount of items used in ippool
	# TYPE fortigate_ippool_used_items gauge
	fortigate_ippool_used_items{name="ippool_name",vdom="FG-traffic"} 0
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
