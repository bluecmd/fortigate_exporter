// Tests of fortigate_exporter
//
// Copyright (C) 2020  Christian Svensson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/google/go-jsonnet"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type fakeClient struct {
	data map[string][]byte
}

func (c *fakeClient) prepare(path string, jfile string) {
	vm := jsonnet.MakeVM()
	b, err := ioutil.ReadFile(jfile)
	if err != nil {
		log.Fatalf("Failed to read jsonnet %q: %v", jfile, err)
	}
	output, err := vm.EvaluateSnippet(jfile, string(b))
	if err != nil {
		log.Fatalf("Failed to evaluate jsonnet %q: %v", jfile, err)
	}
	c.data[path] = []byte(output)
}

func (c *fakeClient) Get(path string, query string, obj interface{}) error {
	d, ok := c.data[path]
	if !ok {
		log.Fatalf("Tried to get unprepared URL %q", path)
	}
	return json.Unmarshal(d, obj)
}

func newFakeClient() *fakeClient {
	return &fakeClient{data: map[string][]byte{}}
}

func TestSystemStatus(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/status", "testdata/status.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !probeSystemStatus(c, r) {
		t.Errorf("probeSystemStatus() returned non-success")
	}

	em := `
	# HELP fortigate_version_info System version and build information
	# TYPE fortigate_version_info gauge
	fortigate_version_info{build="1112",serial="FGVMEVZFNTS3OAC8",version="v6.2.4"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
func TestVPNConnection(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/vpn/ssl", "testdata/vpn.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !probeVPNStatistics(c, r) {
		t.Errorf("probeSystemStatus() returned non-success")
	}

	em := `
	# HELP fortigate_vpn_connections_count_total Number of VPN connections
	# TYPE fortigate_vpn_connections_count_total gauge
	fortigate_vpn_connections_count_total{vdom="root"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
func TestIPSec(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/vpn/ipsec", "testdata/ipsec.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !probeIPSec(c, r) {
		t.Errorf("probeSystemStatus() returned non-success")
	}

	em := `
	# HELP fortigate_ipsec_tunnel_receive_bytes_total Status of Ipsec tunnel
	# TYPE fortigate_ipsec_tunnel_receive_bytes_total gauge
	fortigate_ipsec_tunnel_receive_bytes_total{name="tunnel_1-sub",parent="tunnel_1",vdom="root"} 1.429824e+07
	# HELP fortigate_ipsec_tunnel_transmit_bytes_total Status of Ipsec tunnel
	# TYPE fortigate_ipsec_tunnel_transmit_bytes_total gauge
	fortigate_ipsec_tunnel_transmit_bytes_total{name="tunnel_1-sub",parent="tunnel_1",vdom="root"} 1.424856e+07
	# HELP fortigate_ipsec_tunnel_up Status of Ipsec tunnel
	# TYPE fortigate_ipsec_tunnel_up gauge
	fortigate_ipsec_tunnel_up{name="tunnel_1-sub",parent="tunnel_1",vdom="root"} 1

	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}

func TestSystemResources(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/resource/usage", "testdata/usage.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !probeSystemResources(c, r) {
		t.Errorf("probeSystemResources() returned non-success")
	}

	em := `
	# HELP fortigate_cpu_usage_ratio Current resource usage ratio of system CPU, per core
	# TYPE fortigate_cpu_usage_ratio gauge
	fortigate_cpu_usage_ratio{processor="0"} 0.32
	# HELP fortigate_memory_usage_ratio Current resource usage ratio of system memory
	# TYPE fortigate_memory_usage_ratio gauge
	fortigate_memory_usage_ratio 0.76
	# HELP fortigate_current_sessions Current amount of sessions, per IP version
	# TYPE fortigate_current_sessions gauge
	fortigate_current_sessions{protocol="ipv4"} 5
	fortigate_current_sessions{protocol="ipv6"} 1
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}

func TestSystemVDOMResources(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/resource/usage", "testdata/usage-vdom.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !probeSystemVDOMResources(c, r) {
		t.Errorf("probeSystemVDOMResources() returned non-success")
	}

	em := `
	# HELP fortigate_vdom_cpu_usage_ratio Current resource usage ratio of CPU, per VDOM
	# TYPE fortigate_vdom_cpu_usage_ratio gauge
	fortigate_vdom_cpu_usage_ratio{vdom="FG-traffic"} 0
	fortigate_vdom_cpu_usage_ratio{vdom="root"} 0.01
	# HELP fortigate_vdom_memory_usage_ratio Current resource usage ratio of memory, per VDOM
	# TYPE fortigate_vdom_memory_usage_ratio gauge
	fortigate_vdom_memory_usage_ratio{vdom="FG-traffic"} 0
	fortigate_vdom_memory_usage_ratio{vdom="root"} 0.78
	# HELP fortigate_vdom_current_sessions Current amount of sessions, per VDOM and IP version
	# TYPE fortigate_vdom_current_sessions gauge
	fortigate_vdom_current_sessions{protocol="ipv4",vdom="FG-traffic"} 0
	fortigate_vdom_current_sessions{protocol="ipv4",vdom="root"} 18
	fortigate_vdom_current_sessions{protocol="ipv6",vdom="FG-traffic"} 7
	fortigate_vdom_current_sessions{protocol="ipv6",vdom="root"} 7
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}

func TestFirewallPolicies(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/firewall/policy/select", "testdata/fw-policy.jsonnet")
	c.prepare("api/v2/monitor/firewall/policy6/select", "testdata/fw-policy6.jsonnet")
	c.prepare("api/v2/cmdb/firewall/policy", "testdata/fw-policy-config.jsonnet")
	c.prepare("api/v2/cmdb/firewall/policy6", "testdata/fw-policy6-config.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !probeFirewallPolicies(c, r) {
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
	# TYPE fortigate_policy_bytes_total gauge
	fortigate_policy_bytes_total{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="FG-traffic"} 0
	fortigate_policy_bytes_total{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="root"} 0
	fortigate_policy_bytes_total{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="FG-traffic"} 0
	fortigate_policy_bytes_total{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="root"} 432
	fortigate_policy_bytes_total{id="1",name="",protocol="ipv4",uuid="078f184c-9e9d-51ea-9fbb-66c20957b9c0",vdom="FG-traffic"} 5.34459022e+08
	fortigate_policy_bytes_total{id="1",name="ipv6 policy",protocol="ipv6",uuid="4a2e2fe4-9e9d-51ea-75b1-b5b486b12192",vdom="FG-traffic"} 0
	fortigate_policy_bytes_total{id="2",name="ping",protocol="ipv4",uuid="24843c52-9e9d-51ea-b838-3500a9e54b2e",vdom="FG-traffic"} 0
	# HELP fortigate_policy_hit_count_total Number of times a policy has been hit
	# TYPE fortigate_policy_hit_count_total gauge
	fortigate_policy_hit_count_total{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="FG-traffic"} 0
	fortigate_policy_hit_count_total{id="0",name="Implicit Deny",protocol="ipv4",uuid="",vdom="root"} 0
	fortigate_policy_hit_count_total{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="FG-traffic"} 0
	fortigate_policy_hit_count_total{id="0",name="Implicit Deny",protocol="ipv6",uuid="",vdom="root"} 8
	fortigate_policy_hit_count_total{id="1",name="",protocol="ipv4",uuid="078f184c-9e9d-51ea-9fbb-66c20957b9c0",vdom="FG-traffic"} 4662
	fortigate_policy_hit_count_total{id="1",name="ipv6 policy",protocol="ipv6",uuid="4a2e2fe4-9e9d-51ea-75b1-b5b486b12192",vdom="FG-traffic"} 0
	fortigate_policy_hit_count_total{id="2",name="ping",protocol="ipv4",uuid="24843c52-9e9d-51ea-b838-3500a9e54b2e",vdom="FG-traffic"} 0
	# HELP fortigate_policy_packets_total Number of packets that has passed through a policy
	# TYPE fortigate_policy_packets_total gauge
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

func TestInterfaces(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/interface/select", "testdata/interface.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !probeInterfaces(c, r) {
		t.Errorf("probeInterfaces() returned non-success")
	}

	em := `
	# HELP fortigate_interface_link_up Whether the link is up or not
	# TYPE fortigate_interface_link_up gauge
	fortigate_interface_link_up{alias="",name="b",parent="",vdom="root"} 0
	fortigate_interface_link_up{alias="",name="internal1",parent="",vdom="infra"} 1
	fortigate_interface_link_up{alias="",name="internal2",parent="",vdom="infra"} 1
	fortigate_interface_link_up{alias="",name="internal3",parent="",vdom="root"} 0
	fortigate_interface_link_up{alias="",name="internal4",parent="",vdom="root"} 0
	fortigate_interface_link_up{alias="",name="internal5",parent="",vdom="root"} 0
	fortigate_interface_link_up{alias="",name="modem",parent="",vdom="root"} 0
	fortigate_interface_link_up{alias="",name="npu0_vlink0",parent="",vdom="root"} 1
	fortigate_interface_link_up{alias="",name="npu0_vlink1",parent="",vdom="root"} 1
	fortigate_interface_link_up{alias="",name="vlan-knx",parent="downlink",vdom="knx"} 1
	fortigate_interface_link_up{alias="",name="vlan-ocp-knx",parent="a",vdom="knx"} 1
	fortigate_interface_link_up{alias="",name="wan1",parent="",vdom="main"} 1
	fortigate_interface_link_up{alias="",name="wan2",parent="",vdom="root"} 0
	fortigate_interface_link_up{alias="(mgmt)",name="dmz",parent="",vdom="root"} 1
	fortigate_interface_link_up{alias="(ocp-mgmt)",name="a",parent="",vdom="main"} 1
	fortigate_interface_link_up{alias="(vlan-infra-mgmt)",name="downlink",parent="",vdom="infra"} 1
	# HELP fortigate_interface_receive_bytes_total Number of bytes received on the interface
	# TYPE fortigate_interface_receive_bytes_total counter
	fortigate_interface_receive_bytes_total{alias="",name="b",parent="",vdom="root"} 0
	fortigate_interface_receive_bytes_total{alias="",name="internal1",parent="",vdom="infra"} 1.7367580312e+10
	fortigate_interface_receive_bytes_total{alias="",name="internal2",parent="",vdom="infra"} 7.929384567e+09
	fortigate_interface_receive_bytes_total{alias="",name="internal3",parent="",vdom="root"} 0
	fortigate_interface_receive_bytes_total{alias="",name="internal4",parent="",vdom="root"} 0
	fortigate_interface_receive_bytes_total{alias="",name="internal5",parent="",vdom="root"} 0
	fortigate_interface_receive_bytes_total{alias="",name="modem",parent="",vdom="root"} 0
	fortigate_interface_receive_bytes_total{alias="",name="npu0_vlink0",parent="",vdom="root"} 90
	fortigate_interface_receive_bytes_total{alias="",name="npu0_vlink1",parent="",vdom="root"} 90
	fortigate_interface_receive_bytes_total{alias="",name="vlan-knx",parent="downlink",vdom="knx"} 964072
	fortigate_interface_receive_bytes_total{alias="",name="vlan-ocp-knx",parent="a",vdom="knx"} 2.3445384e+07
	fortigate_interface_receive_bytes_total{alias="",name="wan1",parent="",vdom="main"} 4.8782482049e+10
	fortigate_interface_receive_bytes_total{alias="",name="wan2",parent="",vdom="root"} 0
	fortigate_interface_receive_bytes_total{alias="(mgmt)",name="dmz",parent="",vdom="root"} 1.310564319e+09
	fortigate_interface_receive_bytes_total{alias="(ocp-mgmt)",name="a",parent="",vdom="main"} 1.4568944108e+10
	fortigate_interface_receive_bytes_total{alias="(vlan-infra-mgmt)",name="downlink",parent="",vdom="infra"} 2.5353784011e+10
	# HELP fortigate_interface_receive_errors_total Number of reception errors detected on the interface
	# TYPE fortigate_interface_receive_errors_total counter
	fortigate_interface_receive_errors_total{alias="",name="b",parent="",vdom="root"} 0
	fortigate_interface_receive_errors_total{alias="",name="internal1",parent="",vdom="infra"} 0
	fortigate_interface_receive_errors_total{alias="",name="internal2",parent="",vdom="infra"} 0
	fortigate_interface_receive_errors_total{alias="",name="internal3",parent="",vdom="root"} 0
	fortigate_interface_receive_errors_total{alias="",name="internal4",parent="",vdom="root"} 0
	fortigate_interface_receive_errors_total{alias="",name="internal5",parent="",vdom="root"} 0
	fortigate_interface_receive_errors_total{alias="",name="modem",parent="",vdom="root"} 0
	fortigate_interface_receive_errors_total{alias="",name="npu0_vlink0",parent="",vdom="root"} 0
	fortigate_interface_receive_errors_total{alias="",name="npu0_vlink1",parent="",vdom="root"} 0
	fortigate_interface_receive_errors_total{alias="",name="vlan-knx",parent="downlink",vdom="knx"} 0
	fortigate_interface_receive_errors_total{alias="",name="vlan-ocp-knx",parent="a",vdom="knx"} 0
	fortigate_interface_receive_errors_total{alias="",name="wan1",parent="",vdom="main"} 0
	fortigate_interface_receive_errors_total{alias="",name="wan2",parent="",vdom="root"} 0
	fortigate_interface_receive_errors_total{alias="(mgmt)",name="dmz",parent="",vdom="root"} 0
	fortigate_interface_receive_errors_total{alias="(ocp-mgmt)",name="a",parent="",vdom="main"} 0
	fortigate_interface_receive_errors_total{alias="(vlan-infra-mgmt)",name="downlink",parent="",vdom="infra"} 0
	# HELP fortigate_interface_receive_packets_total Number of packets received on the interface
	# TYPE fortigate_interface_receive_packets_total counter
	fortigate_interface_receive_packets_total{alias="",name="b",parent="",vdom="root"} 0
	fortigate_interface_receive_packets_total{alias="",name="internal1",parent="",vdom="infra"} 5.278112e+07
	fortigate_interface_receive_packets_total{alias="",name="internal2",parent="",vdom="infra"} 5.495165e+07
	fortigate_interface_receive_packets_total{alias="",name="internal3",parent="",vdom="root"} 0
	fortigate_interface_receive_packets_total{alias="",name="internal4",parent="",vdom="root"} 0
	fortigate_interface_receive_packets_total{alias="",name="internal5",parent="",vdom="root"} 0
	fortigate_interface_receive_packets_total{alias="",name="modem",parent="",vdom="root"} 0
	fortigate_interface_receive_packets_total{alias="",name="npu0_vlink0",parent="",vdom="root"} 1
	fortigate_interface_receive_packets_total{alias="",name="npu0_vlink1",parent="",vdom="root"} 1
	fortigate_interface_receive_packets_total{alias="",name="vlan-knx",parent="downlink",vdom="knx"} 5325
	fortigate_interface_receive_packets_total{alias="",name="vlan-ocp-knx",parent="a",vdom="knx"} 134805
	fortigate_interface_receive_packets_total{alias="",name="wan1",parent="",vdom="main"} 4.0481777e+07
	fortigate_interface_receive_packets_total{alias="",name="wan2",parent="",vdom="root"} 0
	fortigate_interface_receive_packets_total{alias="(mgmt)",name="dmz",parent="",vdom="root"} 6.376122e+06
	fortigate_interface_receive_packets_total{alias="(ocp-mgmt)",name="a",parent="",vdom="main"} 1.43943096e+08
	fortigate_interface_receive_packets_total{alias="(vlan-infra-mgmt)",name="downlink",parent="",vdom="infra"} 1.0787347e+08
	# HELP fortigate_interface_speed_bps Speed negotiated on the port in bits/s
	# TYPE fortigate_interface_speed_bps gauge
	fortigate_interface_speed_bps{alias="",name="b",parent="",vdom="root"} 0
	fortigate_interface_speed_bps{alias="",name="internal1",parent="",vdom="infra"} 1e+09
	fortigate_interface_speed_bps{alias="",name="internal2",parent="",vdom="infra"} 1e+09
	fortigate_interface_speed_bps{alias="",name="internal3",parent="",vdom="root"} 0
	fortigate_interface_speed_bps{alias="",name="internal4",parent="",vdom="root"} 0
	fortigate_interface_speed_bps{alias="",name="internal5",parent="",vdom="root"} 0
	fortigate_interface_speed_bps{alias="",name="modem",parent="",vdom="root"} 0
	fortigate_interface_speed_bps{alias="",name="npu0_vlink0",parent="",vdom="root"} 1e+09
	fortigate_interface_speed_bps{alias="",name="npu0_vlink1",parent="",vdom="root"} 1e+09
	fortigate_interface_speed_bps{alias="",name="vlan-knx",parent="downlink",vdom="knx"} 2e+09
	fortigate_interface_speed_bps{alias="",name="vlan-ocp-knx",parent="a",vdom="knx"} 1e+09
	fortigate_interface_speed_bps{alias="",name="wan1",parent="",vdom="main"} 1e+09
	fortigate_interface_speed_bps{alias="",name="wan2",parent="",vdom="root"} 0
	fortigate_interface_speed_bps{alias="(mgmt)",name="dmz",parent="",vdom="root"} 1e+09
	fortigate_interface_speed_bps{alias="(ocp-mgmt)",name="a",parent="",vdom="main"} 1e+09
	fortigate_interface_speed_bps{alias="(vlan-infra-mgmt)",name="downlink",parent="",vdom="infra"} 2e+09
	# HELP fortigate_interface_transmit_bytes_total Number of bytes transmitted on the interface
	# TYPE fortigate_interface_transmit_bytes_total counter
	fortigate_interface_transmit_bytes_total{alias="",name="b",parent="",vdom="root"} 0
	fortigate_interface_transmit_bytes_total{alias="",name="internal1",parent="",vdom="infra"} 1.3038411253e+10
	fortigate_interface_transmit_bytes_total{alias="",name="internal2",parent="",vdom="infra"} 1.0426856559e+10
	fortigate_interface_transmit_bytes_total{alias="",name="internal3",parent="",vdom="root"} 0
	fortigate_interface_transmit_bytes_total{alias="",name="internal4",parent="",vdom="root"} 0
	fortigate_interface_transmit_bytes_total{alias="",name="internal5",parent="",vdom="root"} 0
	fortigate_interface_transmit_bytes_total{alias="",name="modem",parent="",vdom="root"} 0
	fortigate_interface_transmit_bytes_total{alias="",name="npu0_vlink0",parent="",vdom="root"} 90
	fortigate_interface_transmit_bytes_total{alias="",name="npu0_vlink1",parent="",vdom="root"} 90
	fortigate_interface_transmit_bytes_total{alias="",name="vlan-knx",parent="downlink",vdom="knx"} 2.1098754e+07
	fortigate_interface_transmit_bytes_total{alias="",name="vlan-ocp-knx",parent="a",vdom="knx"} 742101
	fortigate_interface_transmit_bytes_total{alias="",name="wan1",parent="",vdom="main"} 8.056925505e+09
	fortigate_interface_transmit_bytes_total{alias="",name="wan2",parent="",vdom="root"} 0
	fortigate_interface_transmit_bytes_total{alias="(mgmt)",name="dmz",parent="",vdom="root"} 2.489018103e+09
	fortigate_interface_transmit_bytes_total{alias="(ocp-mgmt)",name="a",parent="",vdom="main"} 5.5827906482e+10
	fortigate_interface_transmit_bytes_total{alias="(vlan-infra-mgmt)",name="downlink",parent="",vdom="infra"} 2.3561313347e+10
	# HELP fortigate_interface_transmit_errors_total Number of transmission errors detected on the interface
	# TYPE fortigate_interface_transmit_errors_total counter
	fortigate_interface_transmit_errors_total{alias="",name="b",parent="",vdom="root"} 0
	fortigate_interface_transmit_errors_total{alias="",name="internal1",parent="",vdom="infra"} 0
	fortigate_interface_transmit_errors_total{alias="",name="internal2",parent="",vdom="infra"} 0
	fortigate_interface_transmit_errors_total{alias="",name="internal3",parent="",vdom="root"} 0
	fortigate_interface_transmit_errors_total{alias="",name="internal4",parent="",vdom="root"} 0
	fortigate_interface_transmit_errors_total{alias="",name="internal5",parent="",vdom="root"} 0
	fortigate_interface_transmit_errors_total{alias="",name="modem",parent="",vdom="root"} 0
	fortigate_interface_transmit_errors_total{alias="",name="npu0_vlink0",parent="",vdom="root"} 0
	fortigate_interface_transmit_errors_total{alias="",name="npu0_vlink1",parent="",vdom="root"} 0
	fortigate_interface_transmit_errors_total{alias="",name="vlan-knx",parent="downlink",vdom="knx"} 0
	fortigate_interface_transmit_errors_total{alias="",name="vlan-ocp-knx",parent="a",vdom="knx"} 0
	fortigate_interface_transmit_errors_total{alias="",name="wan1",parent="",vdom="main"} 0
	fortigate_interface_transmit_errors_total{alias="",name="wan2",parent="",vdom="root"} 0
	fortigate_interface_transmit_errors_total{alias="(mgmt)",name="dmz",parent="",vdom="root"} 0
	fortigate_interface_transmit_errors_total{alias="(ocp-mgmt)",name="a",parent="",vdom="main"} 0
	fortigate_interface_transmit_errors_total{alias="(vlan-infra-mgmt)",name="downlink",parent="",vdom="infra"} 0
	# HELP fortigate_interface_transmit_packets_total Number of packets transmitted on the interface
	# TYPE fortigate_interface_transmit_packets_total counter
	fortigate_interface_transmit_packets_total{alias="",name="b",parent="",vdom="root"} 0
	fortigate_interface_transmit_packets_total{alias="",name="internal1",parent="",vdom="infra"} 6.3722042e+07
	fortigate_interface_transmit_packets_total{alias="",name="internal2",parent="",vdom="infra"} 6.0035128e+07
	fortigate_interface_transmit_packets_total{alias="",name="internal3",parent="",vdom="root"} 0
	fortigate_interface_transmit_packets_total{alias="",name="internal4",parent="",vdom="root"} 0
	fortigate_interface_transmit_packets_total{alias="",name="internal5",parent="",vdom="root"} 0
	fortigate_interface_transmit_packets_total{alias="",name="modem",parent="",vdom="root"} 0
	fortigate_interface_transmit_packets_total{alias="",name="npu0_vlink0",parent="",vdom="root"} 1
	fortigate_interface_transmit_packets_total{alias="",name="npu0_vlink1",parent="",vdom="root"} 1
	fortigate_interface_transmit_packets_total{alias="",name="vlan-knx",parent="downlink",vdom="knx"} 119021
	fortigate_interface_transmit_packets_total{alias="",name="vlan-ocp-knx",parent="a",vdom="knx"} 3638
	fortigate_interface_transmit_packets_total{alias="",name="wan1",parent="",vdom="main"} 2.1184365e+07
	fortigate_interface_transmit_packets_total{alias="",name="wan2",parent="",vdom="root"} 0
	fortigate_interface_transmit_packets_total{alias="(mgmt)",name="dmz",parent="",vdom="root"} 6.632568e+06
	fortigate_interface_transmit_packets_total{alias="(ocp-mgmt)",name="a",parent="",vdom="main"} 1.30413416e+08
	fortigate_interface_transmit_packets_total{alias="(vlan-infra-mgmt)",name="downlink",parent="",vdom="infra"} 1.23901895e+08
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}

func TestHaStatistics(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/ha-statistics", "testdata/ha-statistics.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !probeHaStatistics(c, r) {
		t.Errorf("probeHaStatistics() returned non-success")
	}

	em := `
        # HELP fortigate_ha_member_bytes_total Bytes transferred by HA member
        # TYPE fortigate_ha_member_bytes_total counter
        fortigate_ha_member_bytes_total{hostname="member-name-1",vdom="root"} 2.02844842379e+11
        fortigate_ha_member_bytes_total{hostname="member-name-2",vdom="root"} 40
        # HELP fortigate_ha_member_cpu_usage_ratio CPU usage by HA member
        # TYPE fortigate_ha_member_cpu_usage_ratio gauge
        fortigate_ha_member_cpu_usage_ratio{hostname="member-name-1",vdom="root"} 0.01
        fortigate_ha_member_cpu_usage_ratio{hostname="member-name-2",vdom="root"} 0
        # HELP fortigate_ha_member_info Info metric regarding cluster members
        # TYPE fortigate_ha_member_info gauge
        fortigate_ha_member_info{hostname="member-name-1",serial="FGT61E4QXXXXXXXX1",vdom="root"} 1
        fortigate_ha_member_info{hostname="member-name-2",serial="FGT61E4QXXXXXXXX2",vdom="root"} 1
        # HELP fortigate_ha_member_ips_events_total IPS events processed by HA member
        # TYPE fortigate_ha_member_ips_events_total counter
        fortigate_ha_member_ips_events_total{hostname="member-name-1",vdom="root"} 0
        fortigate_ha_member_ips_events_total{hostname="member-name-2",vdom="root"} 0
        # HELP fortigate_ha_member_memory_usage_ratio Memory usage by HA member
        # TYPE fortigate_ha_member_memory_usage_ratio gauge
        fortigate_ha_member_memory_usage_ratio{hostname="member-name-1",vdom="root"} 0.67
        fortigate_ha_member_memory_usage_ratio{hostname="member-name-2",vdom="root"} 0.68
        # HELP fortigate_ha_member_network_usage_ratio Network usage by HA member
        # TYPE fortigate_ha_member_network_usage_ratio gauge
        fortigate_ha_member_network_usage_ratio{hostname="member-name-1",vdom="root"} 1.52
        fortigate_ha_member_network_usage_ratio{hostname="member-name-2",vdom="root"} 0.43
        # HELP fortigate_ha_member_packets_total Packets which are handled by this HA member
        # TYPE fortigate_ha_member_packets_total gauge
        fortigate_ha_member_packets_total{hostname="member-name-1",vdom="root"} 5.49981862e+08
        fortigate_ha_member_packets_total{hostname="member-name-2",vdom="root"} 1
        # HELP fortigate_ha_member_sessions Sessions which are handled by this HA member
        # TYPE fortigate_ha_member_sessions gauge
        fortigate_ha_member_sessions{hostname="member-name-1",vdom="root"} 148
        fortigate_ha_member_sessions{hostname="member-name-2",vdom="root"} 12
        # HELP fortigate_ha_member_virus_events_total Virus events which are detected by this HA member
        # TYPE fortigate_ha_member_virus_events_total counter
        fortigate_ha_member_virus_events_total{hostname="member-name-1",vdom="root"} 0
        fortigate_ha_member_virus_events_total{hostname="member-name-2",vdom="root"} 0
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
