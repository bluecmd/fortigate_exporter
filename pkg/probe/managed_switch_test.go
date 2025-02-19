// Copyright 2025 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestProbeManagedSwitch(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/switch-controller/managed-switch", "testdata/managed-switch.jsonnet")
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
		fortigate_managed_switch_l3_packets_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 0
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
		fortigate_managed_switch_l3_packets_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_l3_packets_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 0
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
		fortigate_managed_switch_rx_bcast_packets_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 8.9598762e+07
		fortigate_managed_switch_rx_bcast_packets_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 79999
		fortigate_managed_switch_rx_bcast_packets_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 24635
		fortigate_managed_switch_rx_bcast_packets_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 131902
		fortigate_managed_switch_rx_bcast_packets_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 157391
		fortigate_managed_switch_rx_bcast_packets_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 30
		fortigate_managed_switch_rx_bcast_packets_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 25713
		fortigate_managed_switch_rx_bcast_packets_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 6972
		fortigate_managed_switch_rx_bcast_packets_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 21
		fortigate_managed_switch_rx_bcast_packets_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 79872
		fortigate_managed_switch_rx_bcast_packets_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 85672
		fortigate_managed_switch_rx_bcast_packets_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 28151
		fortigate_managed_switch_rx_bcast_packets_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 180123
		fortigate_managed_switch_rx_bcast_packets_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 29334
		fortigate_managed_switch_rx_bcast_packets_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 6.012792e+06
		fortigate_managed_switch_rx_bcast_packets_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bcast_packets_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 21657
		fortigate_managed_switch_rx_bcast_packets_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 21
		fortigate_managed_switch_rx_bcast_packets_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 12661
		fortigate_managed_switch_rx_bcast_packets_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 67595
		fortigate_managed_switch_rx_bcast_packets_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 51230
		fortigate_managed_switch_rx_bcast_packets_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 72
		# HELP fortigate_managed_switch_rx_bytes_total Total number of received bytes
		# TYPE fortigate_managed_switch_rx_bytes_total counter
		fortigate_managed_switch_rx_bytes_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 1.34930414247e+11
		fortigate_managed_switch_rx_bytes_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 1.33440202045e+11
		fortigate_managed_switch_rx_bytes_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 5.961056557e+09
		fortigate_managed_switch_rx_bytes_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 2.2679215185e+10
		fortigate_managed_switch_rx_bytes_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 3.0008494845e+10
		fortigate_managed_switch_rx_bytes_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 37874
		fortigate_managed_switch_rx_bytes_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 4.3476254701e+10
		fortigate_managed_switch_rx_bytes_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 2.06123042e+10
		fortigate_managed_switch_rx_bytes_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 49529
		fortigate_managed_switch_rx_bytes_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 8.127628504e+09
		fortigate_managed_switch_rx_bytes_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 7.546860589e+10
		fortigate_managed_switch_rx_bytes_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 3.0651113732e+10
		fortigate_managed_switch_rx_bytes_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 1.553331059e+09
		fortigate_managed_switch_rx_bytes_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 2.21170508e+08
		fortigate_managed_switch_rx_bytes_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 1.52616703893e+12
		fortigate_managed_switch_rx_bytes_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_bytes_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 1.547577e+08
		fortigate_managed_switch_rx_bytes_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 34237
		fortigate_managed_switch_rx_bytes_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 1.93704991e+09
		fortigate_managed_switch_rx_bytes_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 3.81008558e+08
		fortigate_managed_switch_rx_bytes_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 8.304739935e+09
		fortigate_managed_switch_rx_bytes_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 5.9554837e+07
		# HELP fortigate_managed_switch_rx_drops_total Total number of received drops
		# TYPE fortigate_managed_switch_rx_drops_total counter
		fortigate_managed_switch_rx_drops_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 1766
		fortigate_managed_switch_rx_drops_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 13
		fortigate_managed_switch_rx_drops_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 1460
		fortigate_managed_switch_rx_drops_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 5842
		fortigate_managed_switch_rx_drops_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 7949
		fortigate_managed_switch_rx_drops_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 36
		fortigate_managed_switch_rx_drops_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 1918
		fortigate_managed_switch_rx_drops_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 1036
		fortigate_managed_switch_rx_drops_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 46
		fortigate_managed_switch_rx_drops_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 3274
		fortigate_managed_switch_rx_drops_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 44
		fortigate_managed_switch_rx_drops_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 593
		fortigate_managed_switch_rx_drops_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 7
		fortigate_managed_switch_rx_drops_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 132
		fortigate_managed_switch_rx_drops_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 9.502165e+06
		fortigate_managed_switch_rx_drops_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_drops_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 152
		fortigate_managed_switch_rx_drops_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 45
		fortigate_managed_switch_rx_drops_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 38
		fortigate_managed_switch_rx_drops_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 60
		fortigate_managed_switch_rx_drops_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 3779
		fortigate_managed_switch_rx_drops_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 74
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
		fortigate_managed_switch_rx_mcast_packets_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 2.0284271e+08
		fortigate_managed_switch_rx_mcast_packets_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 3.712478e+06
		fortigate_managed_switch_rx_mcast_packets_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 124724
		fortigate_managed_switch_rx_mcast_packets_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 552526
		fortigate_managed_switch_rx_mcast_packets_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 543593
		fortigate_managed_switch_rx_mcast_packets_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 44
		fortigate_managed_switch_rx_mcast_packets_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 107638
		fortigate_managed_switch_rx_mcast_packets_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 34664
		fortigate_managed_switch_rx_mcast_packets_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 58
		fortigate_managed_switch_rx_mcast_packets_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 228891
		fortigate_managed_switch_rx_mcast_packets_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 2.478241e+06
		fortigate_managed_switch_rx_mcast_packets_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 95464
		fortigate_managed_switch_rx_mcast_packets_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 93216
		fortigate_managed_switch_rx_mcast_packets_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 339331
		fortigate_managed_switch_rx_mcast_packets_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 3.9493597e+07
		fortigate_managed_switch_rx_mcast_packets_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_mcast_packets_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 238933
		fortigate_managed_switch_rx_mcast_packets_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 58
		fortigate_managed_switch_rx_mcast_packets_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 91546
		fortigate_managed_switch_rx_mcast_packets_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 110379
		fortigate_managed_switch_rx_mcast_packets_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 221841
		fortigate_managed_switch_rx_mcast_packets_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 88
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
		fortigate_managed_switch_rx_packets_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 3.31743835e+08
		fortigate_managed_switch_rx_packets_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 4.33947186e+08
		fortigate_managed_switch_rx_packets_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 3.0300639e+07
		fortigate_managed_switch_rx_packets_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 6.8944093e+07
		fortigate_managed_switch_rx_packets_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 9.8296351e+07
		fortigate_managed_switch_rx_packets_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 235
		fortigate_managed_switch_rx_packets_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 5.572529e+07
		fortigate_managed_switch_rx_packets_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 2.4175761e+07
		fortigate_managed_switch_rx_packets_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 262
		fortigate_managed_switch_rx_packets_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 2.8622586e+07
		fortigate_managed_switch_rx_packets_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 1.76679625e+08
		fortigate_managed_switch_rx_packets_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 4.0420467e+07
		fortigate_managed_switch_rx_packets_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 7.760252e+06
		fortigate_managed_switch_rx_packets_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 2.123154e+06
		fortigate_managed_switch_rx_packets_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 1.461469556e+09
		fortigate_managed_switch_rx_packets_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_packets_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 1.185479e+06
		fortigate_managed_switch_rx_packets_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 257
		fortigate_managed_switch_rx_packets_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 8.844873e+06
		fortigate_managed_switch_rx_packets_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 1.004016e+06
		fortigate_managed_switch_rx_packets_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 3.2368241e+07
		fortigate_managed_switch_rx_packets_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 671591
		# HELP fortigate_managed_switch_rx_ucast_packets_total Total number of received unicast packets
		# TYPE fortigate_managed_switch_rx_ucast_packets_total counter
		fortigate_managed_switch_rx_ucast_packets_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 3.9302363e+07
		fortigate_managed_switch_rx_ucast_packets_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 4.30154709e+08
		fortigate_managed_switch_rx_ucast_packets_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 3.015128e+07
		fortigate_managed_switch_rx_ucast_packets_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 6.8259665e+07
		fortigate_managed_switch_rx_ucast_packets_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 9.7595367e+07
		fortigate_managed_switch_rx_ucast_packets_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 161
		fortigate_managed_switch_rx_ucast_packets_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 5.5591939e+07
		fortigate_managed_switch_rx_ucast_packets_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 2.4134125e+07
		fortigate_managed_switch_rx_ucast_packets_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 183
		fortigate_managed_switch_rx_ucast_packets_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 2.8313823e+07
		fortigate_managed_switch_rx_ucast_packets_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 1.74115712e+08
		fortigate_managed_switch_rx_ucast_packets_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 4.0296852e+07
		fortigate_managed_switch_rx_ucast_packets_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 7.486913e+06
		fortigate_managed_switch_rx_ucast_packets_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 1.754489e+06
		fortigate_managed_switch_rx_ucast_packets_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 1.415963167e+09
		fortigate_managed_switch_rx_ucast_packets_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_rx_ucast_packets_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 924889
		fortigate_managed_switch_rx_ucast_packets_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 178
		fortigate_managed_switch_rx_ucast_packets_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 8.740666e+06
		fortigate_managed_switch_rx_ucast_packets_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 826042
		fortigate_managed_switch_rx_ucast_packets_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 3.209517e+07
		fortigate_managed_switch_rx_ucast_packets_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 671431
		# HELP fortigate_managed_switch_tx_bcast_packets_total Total number of transmitted broadcast packets
		# TYPE fortigate_managed_switch_tx_bcast_packets_total counter
		fortigate_managed_switch_tx_bcast_packets_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 110109
		fortigate_managed_switch_tx_bcast_packets_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 2.2913751e+07
		fortigate_managed_switch_tx_bcast_packets_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 1.121818e+06
		fortigate_managed_switch_tx_bcast_packets_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 1.1829478e+07
		fortigate_managed_switch_tx_bcast_packets_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 7.782054e+06
		fortigate_managed_switch_tx_bcast_packets_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 5
		fortigate_managed_switch_tx_bcast_packets_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 1.086198e+06
		fortigate_managed_switch_tx_bcast_packets_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 242528
		fortigate_managed_switch_tx_bcast_packets_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 6
		fortigate_managed_switch_tx_bcast_packets_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 5.629362e+06
		fortigate_managed_switch_tx_bcast_packets_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 2.0913412e+07
		fortigate_managed_switch_tx_bcast_packets_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 1.206344e+06
		fortigate_managed_switch_tx_bcast_packets_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 9.430355e+06
		fortigate_managed_switch_tx_bcast_packets_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 1.7169614e+07
		fortigate_managed_switch_tx_bcast_packets_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 1.059981e+06
		fortigate_managed_switch_tx_bcast_packets_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bcast_packets_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 1.247639e+07
		fortigate_managed_switch_tx_bcast_packets_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 5
		fortigate_managed_switch_tx_bcast_packets_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 1.1202087e+07
		fortigate_managed_switch_tx_bcast_packets_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 7.154671e+06
		fortigate_managed_switch_tx_bcast_packets_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 3.104095e+06
		fortigate_managed_switch_tx_bcast_packets_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 8.095255e+06
		# HELP fortigate_managed_switch_tx_bytes_total Total number of transmitted bytes
		# TYPE fortigate_managed_switch_tx_bytes_total counter
		fortigate_managed_switch_tx_bytes_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 1.0884719344e+10
		fortigate_managed_switch_tx_bytes_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 8.34905285534e+11
		fortigate_managed_switch_tx_bytes_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 1.18132462128e+11
		fortigate_managed_switch_tx_bytes_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 8.9042475131e+10
		fortigate_managed_switch_tx_bytes_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 1.47636141701e+11
		fortigate_managed_switch_tx_bytes_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 69093
		fortigate_managed_switch_tx_bytes_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 6.0519307071e+10
		fortigate_managed_switch_tx_bytes_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 2.1693872232e+10
		fortigate_managed_switch_tx_bytes_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 76804
		fortigate_managed_switch_tx_bytes_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 3.5669898069e+10
		fortigate_managed_switch_tx_bytes_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 1.94290370163e+11
		fortigate_managed_switch_tx_bytes_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 3.5080851166e+10
		fortigate_managed_switch_tx_bytes_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 1.2193940704e+10
		fortigate_managed_switch_tx_bytes_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 1.6826662742e+10
		fortigate_managed_switch_tx_bytes_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 4.02394112529e+11
		fortigate_managed_switch_tx_bytes_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_bytes_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 1.3117367928e+10
		fortigate_managed_switch_tx_bytes_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 73034
		fortigate_managed_switch_tx_bytes_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 1.3093340484e+10
		fortigate_managed_switch_tx_bytes_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 7.859708084e+09
		fortigate_managed_switch_tx_bytes_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 5.0520676898e+10
		fortigate_managed_switch_tx_bytes_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 8.656136671e+09
		# HELP fortigate_managed_switch_tx_drops_total Total number of transmitted drops
		# TYPE fortigate_managed_switch_tx_drops_total counter
		fortigate_managed_switch_tx_drops_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_drops_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 1
		fortigate_managed_switch_tx_drops_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 0
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
		fortigate_managed_switch_tx_drops_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
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
		fortigate_managed_switch_tx_mcast_packets_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 2.5744161e+07
		fortigate_managed_switch_tx_mcast_packets_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 4.223098e+07
		fortigate_managed_switch_tx_mcast_packets_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 3.482425e+06
		fortigate_managed_switch_tx_mcast_packets_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 2.0664998e+07
		fortigate_managed_switch_tx_mcast_packets_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 1.5550817e+07
		fortigate_managed_switch_tx_mcast_packets_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 9
		fortigate_managed_switch_tx_mcast_packets_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 2.55728e+06
		fortigate_managed_switch_tx_mcast_packets_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 628108
		fortigate_managed_switch_tx_mcast_packets_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 7
		fortigate_managed_switch_tx_mcast_packets_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 1.6710185e+07
		fortigate_managed_switch_tx_mcast_packets_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 4.2340767e+07
		fortigate_managed_switch_tx_mcast_packets_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 3.269723e+06
		fortigate_managed_switch_tx_mcast_packets_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 3.0923127e+07
		fortigate_managed_switch_tx_mcast_packets_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 3.0690631e+07
		fortigate_managed_switch_tx_mcast_packets_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 4.4133987e+07
		fortigate_managed_switch_tx_mcast_packets_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_mcast_packets_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 3.0787572e+07
		fortigate_managed_switch_tx_mcast_packets_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 7
		fortigate_managed_switch_tx_mcast_packets_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 3.0922857e+07
		fortigate_managed_switch_tx_mcast_packets_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 2.0645025e+07
		fortigate_managed_switch_tx_mcast_packets_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 9.267043e+06
		fortigate_managed_switch_tx_mcast_packets_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 2.536097e+07
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
		fortigate_managed_switch_tx_packets_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 6.186197e+07
		fortigate_managed_switch_tx_packets_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 6.9038986e+08
		fortigate_managed_switch_tx_packets_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 9.6717251e+07
		fortigate_managed_switch_tx_packets_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 1.34325577e+08
		fortigate_managed_switch_tx_packets_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 1.99318286e+08
		fortigate_managed_switch_tx_packets_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 189
		fortigate_managed_switch_tx_packets_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 7.2474488e+07
		fortigate_managed_switch_tx_packets_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 2.6579199e+07
		fortigate_managed_switch_tx_packets_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 195
		fortigate_managed_switch_tx_packets_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 6.114537e+07
		fortigate_managed_switch_tx_packets_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 2.00133854e+08
		fortigate_managed_switch_tx_packets_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 4.9880064e+07
		fortigate_managed_switch_tx_packets_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 4.9132301e+07
		fortigate_managed_switch_tx_packets_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 4.9247686e+07
		fortigate_managed_switch_tx_packets_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 1.060966881e+09
		fortigate_managed_switch_tx_packets_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_packets_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 4.4330569e+07
		fortigate_managed_switch_tx_packets_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 203
		fortigate_managed_switch_tx_packets_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 4.9672181e+07
		fortigate_managed_switch_tx_packets_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 2.8807986e+07
		fortigate_managed_switch_tx_packets_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 7.6660806e+07
		fortigate_managed_switch_tx_packets_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 3.3915497e+07
		# HELP fortigate_managed_switch_tx_ucast_packets_total Total number of transmitted unicast packets
		# TYPE fortigate_managed_switch_tx_ucast_packets_total counter
		fortigate_managed_switch_tx_ucast_packets_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 3.60077e+07
		fortigate_managed_switch_tx_ucast_packets_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 6.25245129e+08
		fortigate_managed_switch_tx_ucast_packets_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 9.2113008e+07
		fortigate_managed_switch_tx_ucast_packets_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 1.01831101e+08
		fortigate_managed_switch_tx_ucast_packets_total{port="port12",switch_name="FOO-SW-01",vdom="root"} 1.75985415e+08
		fortigate_managed_switch_tx_ucast_packets_total{port="port13",switch_name="FOO-SW-01",vdom="root"} 175
		fortigate_managed_switch_tx_ucast_packets_total{port="port14",switch_name="FOO-SW-01",vdom="root"} 6.883101e+07
		fortigate_managed_switch_tx_ucast_packets_total{port="port15",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port16",switch_name="FOO-SW-01",vdom="root"} 2.5708563e+07
		fortigate_managed_switch_tx_ucast_packets_total{port="port17",switch_name="FOO-SW-01",vdom="root"} 182
		fortigate_managed_switch_tx_ucast_packets_total{port="port18",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port19",switch_name="FOO-SW-01",vdom="root"} 3.8805823e+07
		fortigate_managed_switch_tx_ucast_packets_total{port="port2",switch_name="FOO-SW-01",vdom="root"} 1.36879675e+08
		fortigate_managed_switch_tx_ucast_packets_total{port="port20",switch_name="FOO-SW-01",vdom="root"} 4.5403997e+07
		fortigate_managed_switch_tx_ucast_packets_total{port="port21",switch_name="FOO-SW-01",vdom="root"} 8.778819e+06
		fortigate_managed_switch_tx_ucast_packets_total{port="port22",switch_name="FOO-SW-01",vdom="root"} 1.387441e+06
		fortigate_managed_switch_tx_ucast_packets_total{port="port23",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port24",switch_name="FOO-SW-01",vdom="root"} 1.015772913e+09
		fortigate_managed_switch_tx_ucast_packets_total{port="port25",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port26",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port27",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port28",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port3",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_tx_ucast_packets_total{port="port4",switch_name="FOO-SW-01",vdom="root"} 1.066607e+06
		fortigate_managed_switch_tx_ucast_packets_total{port="port5",switch_name="FOO-SW-01",vdom="root"} 191
		fortigate_managed_switch_tx_ucast_packets_total{port="port6",switch_name="FOO-SW-01",vdom="root"} 7.547237e+06
		fortigate_managed_switch_tx_ucast_packets_total{port="port7",switch_name="FOO-SW-01",vdom="root"} 1.00829e+06
		fortigate_managed_switch_tx_ucast_packets_total{port="port8",switch_name="FOO-SW-01",vdom="root"} 6.4289668e+07
		fortigate_managed_switch_tx_ucast_packets_total{port="port9",switch_name="FOO-SW-01",vdom="root"} 459272
		# HELP fortigate_managed_switch_under_size_total Total number of under size
		# TYPE fortigate_managed_switch_under_size_total counter
		fortigate_managed_switch_under_size_total{port="internal",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port1",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port10",switch_name="FOO-SW-01",vdom="root"} 0
		fortigate_managed_switch_under_size_total{port="port11",switch_name="FOO-SW-01",vdom="root"} 1
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
