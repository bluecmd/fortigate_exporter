package probe

import (
	"strings"
	"testing"
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/prometheus/common/expfmt"
)

func TestProbeManagedSwitch(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/switch-controller/managed-switch/status", "testdata/managed-switch-status.jsonnet")
	c.prepare("api/v2/monitor/switch-controller/managed-switch/port-stats", "testdata/managed-switch-port-stats.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeManagedSwitch, c, r) {
		t.Errorf("probeManagedSwitchStatus() returned non-success")
	}

	em := `
		# HELP fortigate_managed_switch_collisions_total Total number of collisions
		# TYPE fortigate_managed_switch_collisions_total counter
		fortigate_managed_switch_collisions_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_collisions_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_crc_alignments_total Total number of crc alignments
		# TYPE fortigate_managed_switch_crc_alignments_total counter
		fortigate_managed_switch_crc_alignments_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_crc_alignments_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 1
		# HELP fortigate_managed_switch_fragments_total Total number of fragments
		# TYPE fortigate_managed_switch_fragments_total counter
		fortigate_managed_switch_fragments_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_fragments_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_info Infos about a managed switch
		# TYPE fortigate_managed_switch_info counter
		fortigate_managed_switch_info{os_version="S124EF-v6.4.6-build470,210211 (GA)",serial="S124EF5920010260",state="Authorized",status="Connected",switch_name="FOO-SW-01",vdom="root"} 1
		# HELP fortigate_managed_switch_jabbers_total Total number of jabbers
		# TYPE fortigate_managed_switch_jabbers_total counter
		fortigate_managed_switch_jabbers_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_jabbers_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_l3_packets_total Total number of l3 packets
		# TYPE fortigate_managed_switch_l3_packets_total counter
		fortigate_managed_switch_l3_packets_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 166728
		fortigate_managed_switch_l3_packets_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 3.681687e+06
		fortigate_managed_switch_l3_packets_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 2.451258e+06
		fortigate_managed_switch_l3_packets_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 9
		fortigate_managed_switch_l3_packets_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 5.774658e+06
		fortigate_managed_switch_l3_packets_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_max_poe_budget_watt Max poe budget watt
		# TYPE fortigate_managed_switch_max_poe_budget_watt counter
		fortigate_managed_switch_max_poe_budget_watt{switch_name="FOO-SW-01",vdom="root"} 370
		# HELP fortigate_managed_switch_port_info Infos about a switch port
		# TYPE fortigate_managed_switch_port_info gauge
		fortigate_managed_switch_port_info{duplex="",poe_capable="false",poe_status="",port="port25",status="down",switch_name="FOO-SW-01",vdom="root",vlan="default"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="false",poe_status="",port="port26",status="down",switch_name="FOO-SW-01",vdom="root",vlan="default"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="false",poe_status="",port="port27",status="down",switch_name="FOO-SW-01",vdom="root",vlan="default"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="false",poe_status="",port="port28",status="down",switch_name="FOO-SW-01",vdom="root",vlan="default"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port10",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port11",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port12",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port13",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port14",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port15",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port16",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port17",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port18",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port19",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port20",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port21",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port22",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port23",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port3",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port4",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port5",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port6",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port7",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port8",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="",poe_capable="true",poe_status="enabled",port="port9",status="down",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-101"} 1
		fortigate_managed_switch_port_info{duplex="full",poe_capable="true",poe_status="enabled",port="port1",status="up",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-4001"} 1
		fortigate_managed_switch_port_info{duplex="full",poe_capable="true",poe_status="enabled",port="port2",status="up",switch_name="FOO-SW-01",vdom="root",vlan="VLAN-4001"} 1
		fortigate_managed_switch_port_info{duplex="full",poe_capable="true",poe_status="enabled",port="port24",status="up",switch_name="FOO-SW-01",vdom="root",vlan="default"} 1
		# HELP fortigate_managed_switch_port_power_status Port power status
		# TYPE fortigate_managed_switch_port_power_status gauge
		fortigate_managed_switch_port_power_status{port="port1",switch_name="FOO-SW-01",vdom="root"} 2
		fortigate_managed_switch_port_power_status{port="port10",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port11",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port12",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port13",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port14",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port15",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port16",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port17",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port18",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port19",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port2",switch_name="FOO-SW-01",vdom="root"} 2
		fortigate_managed_switch_port_power_status{port="port20",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port21",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port22",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port23",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port24",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_status{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_status{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_status{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_status{port="port3",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port4",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port5",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port6",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port7",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port8",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_power_status{port="port9",switch_name="FOO-SW-01",vdom="root"} 1
		# HELP fortigate_managed_switch_port_power_watt Port power in watt
		# TYPE fortigate_managed_switch_port_power_watt gauge
		fortigate_managed_switch_port_power_watt{port="port1",switch_name="FOO-SW-01",vdom="root"} 6
		fortigate_managed_switch_port_power_watt{port="port10",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port11",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port17",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port19",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port2",switch_name="FOO-SW-01",vdom="root"} 6.099999904632568
		fortigate_managed_switch_port_power_watt{port="port20",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port21",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port22",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port24",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_power_watt{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_port_status Port status up=1 down=0
		# TYPE fortigate_managed_switch_port_status gauge
		fortigate_managed_switch_port_status{port="port1",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_status{port="port10",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port11",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port17",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port19",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port2",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_status{port="port20",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port21",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port22",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port24",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_port_status{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_port_status{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_rx_bcast_packets_total Total number of received broadcast packets
		# TYPE fortigate_managed_switch_rx_bcast_packets_total counter
		fortigate_managed_switch_rx_bcast_packets_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 386468
		fortigate_managed_switch_rx_bcast_packets_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 218
		fortigate_managed_switch_rx_bcast_packets_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_rx_bcast_packets_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 12
		fortigate_managed_switch_rx_bcast_packets_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 18
		fortigate_managed_switch_rx_bcast_packets_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 6867
		fortigate_managed_switch_rx_bcast_packets_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 6889
		fortigate_managed_switch_rx_bcast_packets_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 524
		fortigate_managed_switch_rx_bcast_packets_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 400003
		fortigate_managed_switch_rx_bcast_packets_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 333423
		fortigate_managed_switch_rx_bcast_packets_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 557
		fortigate_managed_switch_rx_bcast_packets_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_rx_bytes_total Total number of received bytes
		# TYPE fortigate_managed_switch_rx_bytes_total counter
		fortigate_managed_switch_rx_bytes_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 3.609650574e+09
		fortigate_managed_switch_rx_bytes_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 5.05250508e+08
		fortigate_managed_switch_rx_bytes_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 13952
		fortigate_managed_switch_rx_bytes_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 9.0023496e+07
		fortigate_managed_switch_rx_bytes_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 3.3288567e+08
		fortigate_managed_switch_rx_bytes_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 1.164988782e+09
		fortigate_managed_switch_rx_bytes_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 7.844701e+06
		fortigate_managed_switch_rx_bytes_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 9.2797595e+07
		fortigate_managed_switch_rx_bytes_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 1.248215979e+09
		fortigate_managed_switch_rx_bytes_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 1.283515639e+09
		fortigate_managed_switch_rx_bytes_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 1.643049095e+09
		fortigate_managed_switch_rx_bytes_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 1.54309355e+09
		fortigate_managed_switch_rx_bytes_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 5.56071078e+08
		fortigate_managed_switch_rx_bytes_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_rx_drops_total Total number of received drops
		# TYPE fortigate_managed_switch_rx_drops_total counter
		fortigate_managed_switch_rx_drops_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 4
		fortigate_managed_switch_rx_drops_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 5
		fortigate_managed_switch_rx_drops_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 32
		fortigate_managed_switch_rx_drops_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 1.264621e+06
		fortigate_managed_switch_rx_drops_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 410
		fortigate_managed_switch_rx_drops_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_rx_errors_total Total number of received errors
		# TYPE fortigate_managed_switch_rx_errors_total counter
		fortigate_managed_switch_rx_errors_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_errors_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_rx_mcast_packets_total Total number of received multicast packets
		# TYPE fortigate_managed_switch_rx_mcast_packets_total counter
		fortigate_managed_switch_rx_mcast_packets_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 454590
		fortigate_managed_switch_rx_mcast_packets_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 27961
		fortigate_managed_switch_rx_mcast_packets_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 27971
		fortigate_managed_switch_rx_mcast_packets_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 27973
		fortigate_managed_switch_rx_mcast_packets_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 27973
		fortigate_managed_switch_rx_mcast_packets_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 56777
		fortigate_managed_switch_rx_mcast_packets_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 56780
		fortigate_managed_switch_rx_mcast_packets_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 403979
		fortigate_managed_switch_rx_mcast_packets_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 404022
		fortigate_managed_switch_rx_mcast_packets_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 11494
		fortigate_managed_switch_rx_mcast_packets_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_rx_oversize_total Total number of received oversize
		# TYPE fortigate_managed_switch_rx_oversize_total counter
		fortigate_managed_switch_rx_oversize_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_oversize_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_rx_packets_total Total number of received packets
		# TYPE fortigate_managed_switch_rx_packets_total counter
		fortigate_managed_switch_rx_packets_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 1.6097764e+07
		fortigate_managed_switch_rx_packets_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 1.612797e+06
		fortigate_managed_switch_rx_packets_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 218
		fortigate_managed_switch_rx_packets_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 789978
		fortigate_managed_switch_rx_packets_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 670380
		fortigate_managed_switch_rx_packets_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 5.099612e+06
		fortigate_managed_switch_rx_packets_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 35671
		fortigate_managed_switch_rx_packets_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 286308
		fortigate_managed_switch_rx_packets_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 7.76534e+06
		fortigate_managed_switch_rx_packets_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 7.888974e+06
		fortigate_managed_switch_rx_packets_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 4.450257e+06
		fortigate_managed_switch_rx_packets_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 9.964549e+06
		fortigate_managed_switch_rx_packets_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 2.121701e+06
		fortigate_managed_switch_rx_packets_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_rx_ucast_packets_total Total number of received unicast packets
		# TYPE fortigate_managed_switch_rx_ucast_packets_total counter
		fortigate_managed_switch_rx_ucast_packets_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 1.6097764e+07
		fortigate_managed_switch_rx_ucast_packets_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 771739
		fortigate_managed_switch_rx_ucast_packets_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 789977
		fortigate_managed_switch_rx_ucast_packets_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 642407
		fortigate_managed_switch_rx_ucast_packets_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 5.071623e+06
		fortigate_managed_switch_rx_ucast_packets_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 831
		fortigate_managed_switch_rx_ucast_packets_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 251446
		fortigate_managed_switch_rx_ucast_packets_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 7.708039e+06
		fortigate_managed_switch_rx_ucast_packets_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 7.832194e+06
		fortigate_managed_switch_rx_ucast_packets_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 3.646275e+06
		fortigate_managed_switch_rx_ucast_packets_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 9.227104e+06
		fortigate_managed_switch_rx_ucast_packets_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 2.10965e+06
		fortigate_managed_switch_rx_ucast_packets_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_tx_bcast_packets_total Total number of transmitted broadcast packets
		# TYPE fortigate_managed_switch_tx_bcast_packets_total counter
		fortigate_managed_switch_tx_bcast_packets_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 329887
		fortigate_managed_switch_tx_bcast_packets_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 98099
		fortigate_managed_switch_tx_bcast_packets_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 218769
		fortigate_managed_switch_tx_bcast_packets_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 486468
		fortigate_managed_switch_tx_bcast_packets_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 567570
		fortigate_managed_switch_tx_bcast_packets_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 380235
		fortigate_managed_switch_tx_bcast_packets_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 364291
		fortigate_managed_switch_tx_bcast_packets_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 401120
		fortigate_managed_switch_tx_bcast_packets_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 401615
		fortigate_managed_switch_tx_bcast_packets_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 459245
		fortigate_managed_switch_tx_bcast_packets_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 11021
		fortigate_managed_switch_tx_bcast_packets_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 252981
		fortigate_managed_switch_tx_bcast_packets_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_tx_bytes_total Total number of transmitted bytes
		# TYPE fortigate_managed_switch_tx_bytes_total counter
		fortigate_managed_switch_tx_bytes_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 4.094254483e+09
		fortigate_managed_switch_tx_bytes_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 1.057457469e+09
		fortigate_managed_switch_tx_bytes_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 1.20054895e+08
		fortigate_managed_switch_tx_bytes_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 3.3026203e+08
		fortigate_managed_switch_tx_bytes_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 9.95087031e+08
		fortigate_managed_switch_tx_bytes_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 4.98405588e+08
		fortigate_managed_switch_tx_bytes_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 3.56603008e+08
		fortigate_managed_switch_tx_bytes_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 5.37463983e+08
		fortigate_managed_switch_tx_bytes_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 8.244358564e+09
		fortigate_managed_switch_tx_bytes_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 5.925235826e+09
		fortigate_managed_switch_tx_bytes_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 2.713564863e+09
		fortigate_managed_switch_tx_bytes_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 3.19096093e+08
		fortigate_managed_switch_tx_bytes_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 1.695841429e+09
		fortigate_managed_switch_tx_bytes_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_tx_drops_total Total number of transmitted drops
		# TYPE fortigate_managed_switch_tx_drops_total counter
		fortigate_managed_switch_tx_drops_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_tx_drops_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 22
		fortigate_managed_switch_tx_drops_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_tx_errors_total Total number of transmitted errors
		# TYPE fortigate_managed_switch_tx_errors_total counter
		fortigate_managed_switch_tx_errors_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_errors_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_tx_mcast_packets_total Total number of transmitted multicast packets
		# TYPE fortigate_managed_switch_tx_mcast_packets_total counter
		fortigate_managed_switch_tx_mcast_packets_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 811352
		fortigate_managed_switch_tx_mcast_packets_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 270481
		fortigate_managed_switch_tx_mcast_packets_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 269219
		fortigate_managed_switch_tx_mcast_packets_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 657839
		fortigate_managed_switch_tx_mcast_packets_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 667598
		fortigate_managed_switch_tx_mcast_packets_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 667600
		fortigate_managed_switch_tx_mcast_packets_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 667598
		fortigate_managed_switch_tx_mcast_packets_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 315334
		fortigate_managed_switch_tx_mcast_packets_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 315331
		fortigate_managed_switch_tx_mcast_packets_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 402472
		fortigate_managed_switch_tx_mcast_packets_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 410227
		fortigate_managed_switch_tx_mcast_packets_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 114683
		fortigate_managed_switch_tx_mcast_packets_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_tx_oversize_total Total number of transmitted oversize
		# TYPE fortigate_managed_switch_tx_oversize_total counter
		fortigate_managed_switch_tx_oversize_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_oversize_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_tx_packets_total Total number of transmitted packets
		# TYPE fortigate_managed_switch_tx_packets_total counter
		fortigate_managed_switch_tx_packets_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 2.3838574e+07
		fortigate_managed_switch_tx_packets_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 2.352591e+06
		fortigate_managed_switch_tx_packets_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 368615
		fortigate_managed_switch_tx_packets_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 1.283318e+06
		fortigate_managed_switch_tx_packets_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 5.659466e+06
		fortigate_managed_switch_tx_packets_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 2.626036e+06
		fortigate_managed_switch_tx_packets_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 1.135396e+06
		fortigate_managed_switch_tx_packets_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 1.437148e+06
		fortigate_managed_switch_tx_packets_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 1.0005028e+07
		fortigate_managed_switch_tx_packets_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 8.308575e+06
		fortigate_managed_switch_tx_packets_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 9.317269e+06
		fortigate_managed_switch_tx_packets_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 2.203125e+06
		fortigate_managed_switch_tx_packets_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 6.113198e+06
		fortigate_managed_switch_tx_packets_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_tx_ucast_packets_total Total number of transmitted unicast packets
		# TYPE fortigate_managed_switch_tx_ucast_packets_total counter
		fortigate_managed_switch_tx_ucast_packets_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 2.3838574e+07
		fortigate_managed_switch_tx_ucast_packets_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 1.211352e+06
		fortigate_managed_switch_tx_ucast_packets_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 35
		fortigate_managed_switch_tx_ucast_packets_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 795330
		fortigate_managed_switch_tx_ucast_packets_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 4.515159e+06
		fortigate_managed_switch_tx_ucast_packets_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 1.390868e+06
		fortigate_managed_switch_tx_ucast_packets_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 87561
		fortigate_managed_switch_tx_ucast_packets_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 405259
		fortigate_managed_switch_tx_ucast_packets_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 9.288574e+06
		fortigate_managed_switch_tx_ucast_packets_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 7.591629e+06
		fortigate_managed_switch_tx_ucast_packets_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 8.455552e+06
		fortigate_managed_switch_tx_ucast_packets_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 1.781877e+06
		fortigate_managed_switch_tx_ucast_packets_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 5.745534e+06
		fortigate_managed_switch_tx_ucast_packets_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		# HELP fortigate_managed_switch_under_size_total Total number of under size
		# TYPE fortigate_managed_switch_under_size_total counter
		fortigate_managed_switch_under_size_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 0
		`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}

}

// Use this to get the string output of all the metrics in the registry
func printMetrics(registry *prometheus.Registry) {
    // Gather the metrics
    metrics, err := registry.Gather()
    if err != nil {
        fmt.Println("Error gathering metrics:", err)
        return
    }

    // Create a new encoder
    encoder := expfmt.NewEncoder(os.Stdout, expfmt.FmtText)

    // Encode the gathered metrics
    for _, metricFamily := range metrics {
        if err := encoder.Encode(metricFamily); err != nil {
            fmt.Println("Error encoding metric family:", err)
            return
        }
    }
}