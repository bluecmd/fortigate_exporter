package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestFirewallPoliciesPre64(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/firewall/policy/select", "testdata/fw-policy-pre64.jsonnet")
	c.prepare("api/v2/monitor/firewall/policy6/select", "testdata/fw-policy6-pre64.jsonnet")
	c.prepare("api/v2/cmdb/firewall/policy", "testdata/fw-policy-config-pre64.jsonnet")
	c.prepare("api/v2/cmdb/firewall/policy6", "testdata/fw-policy6-config-pre64.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeFirewallPolicies, c, r) {
		t.Errorf("probeFirewallPolicies() returned non-success")
	}

	em := `
	# HELP fortigate_policy_active_sessions Number of active sessions for a policy
	# TYPE fortigate_policy_active_sessions gauge
	fortigate_policy_active_sessions{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="FG-traffic"} 0
	fortigate_policy_active_sessions{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="root"} 0
	fortigate_policy_active_sessions{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="FG-traffic"} 0
	fortigate_policy_active_sessions{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="root"} 0
	fortigate_policy_active_sessions{id="1",name="",protocol="ipv4",uuid="078f184c-9e9d-51ea-9fbb-66c20957b9c0",vdom="FG-traffic"} 2
	fortigate_policy_active_sessions{id="1",name="ipv6 policy",protocol="ipv6",uuid="4a2e2fe4-9e9d-51ea-75b1-b5b486b12192",vdom="FG-traffic"} 0
	fortigate_policy_active_sessions{id="2",name="ping",protocol="ipv4",uuid="24843c52-9e9d-51ea-b838-3500a9e54b2e",vdom="FG-traffic"} 0
	# HELP fortigate_policy_bytes_total Number of bytes that has passed through a policy
	# TYPE fortigate_policy_bytes_total counter
	fortigate_policy_bytes_total{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="FG-traffic"} 64687125982
	fortigate_policy_bytes_total{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="root"} 0
	fortigate_policy_bytes_total{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="FG-traffic"} 0
	fortigate_policy_bytes_total{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="root"} 432
	fortigate_policy_bytes_total{id="1",name="",protocol="ipv4",uuid="078f184c-9e9d-51ea-9fbb-66c20957b9c0",vdom="FG-traffic"} 5.34459022e+08
	fortigate_policy_bytes_total{id="1",name="ipv6 policy",protocol="ipv6",uuid="4a2e2fe4-9e9d-51ea-75b1-b5b486b12192",vdom="FG-traffic"} 0
	fortigate_policy_bytes_total{id="2",name="ping",protocol="ipv4",uuid="24843c52-9e9d-51ea-b838-3500a9e54b2e",vdom="FG-traffic"} 0
	# HELP fortigate_policy_hit_count_total Number of times a policy has been hit
	# TYPE fortigate_policy_hit_count_total counter
	fortigate_policy_hit_count_total{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="FG-traffic"} 0
	fortigate_policy_hit_count_total{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="root"} 0
	fortigate_policy_hit_count_total{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="FG-traffic"} 0
	fortigate_policy_hit_count_total{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="root"} 8
	fortigate_policy_hit_count_total{id="1",name="",protocol="ipv4",uuid="078f184c-9e9d-51ea-9fbb-66c20957b9c0",vdom="FG-traffic"} 4662
	fortigate_policy_hit_count_total{id="1",name="ipv6 policy",protocol="ipv6",uuid="4a2e2fe4-9e9d-51ea-75b1-b5b486b12192",vdom="FG-traffic"} 0
	fortigate_policy_hit_count_total{id="2",name="ping",protocol="ipv4",uuid="24843c52-9e9d-51ea-b838-3500a9e54b2e",vdom="FG-traffic"} 0
	# HELP fortigate_policy_packets_total Number of packets that has passed through a policy
	# TYPE fortigate_policy_packets_total counter
	fortigate_policy_packets_total{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="FG-traffic"} 0
	fortigate_policy_packets_total{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="root"} 0
	fortigate_policy_packets_total{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="FG-traffic"} 0
	fortigate_policy_packets_total{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="root"} 6
	fortigate_policy_packets_total{id="1",name="",protocol="ipv4",uuid="078f184c-9e9d-51ea-9fbb-66c20957b9c0",vdom="FG-traffic"} 792806
	fortigate_policy_packets_total{id="1",name="ipv6 policy",protocol="ipv6",uuid="4a2e2fe4-9e9d-51ea-75b1-b5b486b12192",vdom="FG-traffic"} 0
	fortigate_policy_packets_total{id="2",name="ping",protocol="ipv4",uuid="24843c52-9e9d-51ea-b838-3500a9e54b2e",vdom="FG-traffic"} 0
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}

func TestFirewallPolicies(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/firewall/policy/select?ip_version=ipv4", "testdata/fw-policy-v4.jsonnet")
	c.prepare("api/v2/monitor/firewall/policy/select?ip_version=ipv6", "testdata/fw-policy-v6.jsonnet")
	c.prepare("api/v2/cmdb/firewall/policy", "testdata/fw-policy-config.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeFirewallPolicies, c, r) {
		t.Errorf("probeFirewallPolicies() returned non-success")
	}

	em := `
	# HELP fortigate_policy_active_sessions Number of active sessions for a policy
	# TYPE fortigate_policy_active_sessions gauge
	fortigate_policy_active_sessions{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="FG-traffic"} 0
	fortigate_policy_active_sessions{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="root"} 0
	fortigate_policy_active_sessions{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="FG-traffic"} 1
	fortigate_policy_active_sessions{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="root"} 0
	fortigate_policy_active_sessions{id="1",name="",protocol="ipv4",uuid="078f184c-9e9d-51ea-9fbb-66c20957b9c0",vdom="FG-traffic"} 2
	fortigate_policy_active_sessions{id="1",name="",protocol="ipv6",uuid="078f184c-9e9d-51ea-9fbb-66c20957b9c0",vdom="FG-traffic"} 10
	fortigate_policy_active_sessions{id="2",name="ping",protocol="ipv4",uuid="24843c52-9e9d-51ea-b838-3500a9e54b2e",vdom="FG-traffic"} 0
	fortigate_policy_active_sessions{id="2",name="ping",protocol="ipv6",uuid="24843c52-9e9d-51ea-b838-3500a9e54b2e",vdom="FG-traffic"} 1
	# HELP fortigate_policy_bytes_total Number of bytes that has passed through a policy
	# TYPE fortigate_policy_bytes_total counter
	fortigate_policy_bytes_total{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="FG-traffic"} 0
	fortigate_policy_bytes_total{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="root"} 0
	fortigate_policy_bytes_total{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="FG-traffic"} 1
	fortigate_policy_bytes_total{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="root"} 0
	fortigate_policy_bytes_total{id="1",name="",protocol="ipv4",uuid="078f184c-9e9d-51ea-9fbb-66c20957b9c0",vdom="FG-traffic"} 5.34459022e+08
	fortigate_policy_bytes_total{id="1",name="",protocol="ipv6",uuid="078f184c-9e9d-51ea-9fbb-66c20957b9c0",vdom="FG-traffic"} 1000
	fortigate_policy_bytes_total{id="2",name="ping",protocol="ipv4",uuid="24843c52-9e9d-51ea-b838-3500a9e54b2e",vdom="FG-traffic"} 0
	fortigate_policy_bytes_total{id="2",name="ping",protocol="ipv6",uuid="24843c52-9e9d-51ea-b838-3500a9e54b2e",vdom="FG-traffic"} 2
	# HELP fortigate_policy_hit_count_total Number of times a policy has been hit
	# TYPE fortigate_policy_hit_count_total counter
	fortigate_policy_hit_count_total{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="FG-traffic"} 0
	fortigate_policy_hit_count_total{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="root"} 0
	fortigate_policy_hit_count_total{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="FG-traffic"} 1
	fortigate_policy_hit_count_total{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="root"} 0
	fortigate_policy_hit_count_total{id="1",name="",protocol="ipv4",uuid="078f184c-9e9d-51ea-9fbb-66c20957b9c0",vdom="FG-traffic"} 4662
	fortigate_policy_hit_count_total{id="1",name="",protocol="ipv6",uuid="078f184c-9e9d-51ea-9fbb-66c20957b9c0",vdom="FG-traffic"} 11000
	fortigate_policy_hit_count_total{id="2",name="ping",protocol="ipv4",uuid="24843c52-9e9d-51ea-b838-3500a9e54b2e",vdom="FG-traffic"} 0
	fortigate_policy_hit_count_total{id="2",name="ping",protocol="ipv6",uuid="24843c52-9e9d-51ea-b838-3500a9e54b2e",vdom="FG-traffic"} 0
	# HELP fortigate_policy_packets_total Number of packets that has passed through a policy
	# TYPE fortigate_policy_packets_total counter
	fortigate_policy_packets_total{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="FG-traffic"} 0
	fortigate_policy_packets_total{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="root"} 0
	fortigate_policy_packets_total{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="FG-traffic"} 1
	fortigate_policy_packets_total{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="root"} 0
	fortigate_policy_packets_total{id="1",name="",protocol="ipv4",uuid="078f184c-9e9d-51ea-9fbb-66c20957b9c0",vdom="FG-traffic"} 792806
	fortigate_policy_packets_total{id="1",name="",protocol="ipv6",uuid="078f184c-9e9d-51ea-9fbb-66c20957b9c0",vdom="FG-traffic"} 2000
	fortigate_policy_packets_total{id="2",name="ping",protocol="ipv4",uuid="24843c52-9e9d-51ea-b838-3500a9e54b2e",vdom="FG-traffic"} 0
	fortigate_policy_packets_total{id="2",name="ping",protocol="ipv6",uuid="24843c52-9e9d-51ea-b838-3500a9e54b2e",vdom="FG-traffic"} 3
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
