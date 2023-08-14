package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestSystemInterfaces(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/interface/select", "testdata/interface.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemInterface, c, r) {
		t.Errorf("probeSystemInterface() returned non-success")
	}

	em := `
	# HELP fortigate_interface_link_up Whether the link is up or not (not taking into account admin status)
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
	fortigate_interface_transmit_bytes_total{alias="",name="internal5",parent="",vdom="root"} 6.4687125982e+17
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
