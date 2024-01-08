package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestSystemArpTable(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/network/arp/select", "testdata/arp.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemArpTable, c, r) {
		t.Errorf("probeSystemArpTable() returned non-success")
	}

	em := `
	# HELP fortigate_arp_entry_age_seconds Age of ARP table entry in seconds
	# TYPE fortigate_arp_entry_age_seconds gauge
	fortigate_arp_entry_age_seconds{ip="192.168.1.23",mac="aa:05:f3:a6:33:2a",interface="port1"} 0
	fortigate_arp_entry_age_seconds{ip="192.168.1.1",mac="80:16:05:f8:71:10",interface="port1"} 35
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
