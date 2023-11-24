package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestSystemVDOMResourcesAll(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/vdom-resource", "testdata/system-vdom-resources.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemVDOMResourcesAll, c, r) {
		t.Errorf("probeSystemVDOMResourcesAll() returned non-success")
	}

	em := `
	# HELP fortigate_vdom_resources Metrics of current usage of vdom resources as well as both the default and user configured maximum values.
	# TYPE fortigate_vdom_resources gauge
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="0",guaranteed="0",if_name="dialup-tunnel",max_custom_value="0",max_guaranteed_value="0",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="root"} 5
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="0",guaranteed="0",if_name="dialup-tunnel",max_custom_value="0",max_guaranteed_value="0",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="test1"} 5
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="0",guaranteed="0",if_name="ipsec-phase1-interface",max_custom_value="0",max_guaranteed_value="0",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="root"} 3
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="0",guaranteed="0",if_name="ipsec-phase2-interface",max_custom_value="0",max_guaranteed_value="0",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="root"} 4
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="0",guaranteed="0",if_name="log-disk-quota",max_custom_value="0",max_guaranteed_value="0",min_custom_value="100",min_guaranteed_value="0",usage_percent="0",vdom="root"} 17
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="0",guaranteed="0",if_name="log-disk-quota",max_custom_value="0",max_guaranteed_value="0",min_custom_value="100",min_guaranteed_value="0",usage_percent="0",vdom="test1"} 17
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="0",guaranteed="0",if_name="onetime-schedule",max_custom_value="5000",max_guaranteed_value="5000",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="root"} 11
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="0",guaranteed="0",if_name="onetime-schedule",max_custom_value="5000",max_guaranteed_value="5000",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="test1"} 11
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="0",guaranteed="0",if_name="sslvpn",max_custom_value="0",max_guaranteed_value="0",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="root"} 15
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="0",guaranteed="0",if_name="sslvpn",max_custom_value="0",max_guaranteed_value="0",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="test1"} 15
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="0",guaranteed="0",if_name="user",max_custom_value="5000",max_guaranteed_value="5000",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="root"} 13
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="0",guaranteed="0",if_name="user",max_custom_value="5000",max_guaranteed_value="5000",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="test1"} 13
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="0",guaranteed="0",if_name="user-group",max_custom_value="2000",max_guaranteed_value="2000",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="test1"} 14
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="20000",guaranteed="0",if_name="ipsec-phase1",max_custom_value="20000",max_guaranteed_value="20000",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="root"} 1
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="20000",guaranteed="0",if_name="ipsec-phase1",max_custom_value="20000",max_guaranteed_value="20000",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="test1"} 1
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="20000",guaranteed="0",if_name="ipsec-phase2",max_custom_value="20000",max_guaranteed_value="20000",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="root"} 2
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="20000",guaranteed="0",if_name="ipsec-phase2",max_custom_value="20000",max_guaranteed_value="20000",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="test1"} 2
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="201024",guaranteed="0",if_name="firewall-policy",max_custom_value="201024",max_guaranteed_value="201014",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="root"} 6
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="33192",guaranteed="0",if_name="firewall-addrgrp",max_custom_value="33192",max_guaranteed_value="33188",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="root"} 8
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="64000",guaranteed="0",if_name="proxy",max_custom_value="64000",max_guaranteed_value="64000",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="root"} 16
	fortigate_vdom_resources{current_usage="0",custom_max="0",global_max="64000",guaranteed="0",if_name="proxy",max_custom_value="64000",max_guaranteed_value="64000",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="test1"} 16
	fortigate_vdom_resources{current_usage="0",custom_max="200000",global_max="0",guaranteed="150000",if_name="session",max_custom_value="22000000",max_guaranteed_value="200000",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="test1"} 0
	fortigate_vdom_resources{current_usage="1",custom_max="0",global_max="0",guaranteed="0",if_name="ipsec-phase1-interface",max_custom_value="0",max_guaranteed_value="0",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="test1"} 3
	fortigate_vdom_resources{current_usage="1",custom_max="0",global_max="0",guaranteed="0",if_name="ipsec-phase2-interface",max_custom_value="0",max_guaranteed_value="0",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="test1"} 4
	fortigate_vdom_resources{current_usage="1",custom_max="0",global_max="33192",guaranteed="0",if_name="firewall-addrgrp",max_custom_value="33192",max_guaranteed_value="33189",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="test1"} 8
	fortigate_vdom_resources{current_usage="17",custom_max="0",global_max="113192",guaranteed="0",if_name="firewall-address",max_custom_value="113192",max_guaranteed_value="113152",min_custom_value="17",min_guaranteed_value="0",usage_percent="0",vdom="test1"} 7
	fortigate_vdom_resources{current_usage="2",custom_max="0",global_max="0",guaranteed="0",if_name="recurring-schedule",max_custom_value="1024",max_guaranteed_value="1024",min_custom_value="2",min_guaranteed_value="0",usage_percent="0",vdom="root"} 12
	fortigate_vdom_resources{current_usage="2",custom_max="0",global_max="0",guaranteed="0",if_name="recurring-schedule",max_custom_value="1024",max_guaranteed_value="1024",min_custom_value="2",min_guaranteed_value="0",usage_percent="0",vdom="test1"} 12
	fortigate_vdom_resources{current_usage="2",custom_max="0",global_max="0",guaranteed="0",if_name="user-group",max_custom_value="2000",max_guaranteed_value="2000",min_custom_value="2",min_guaranteed_value="0",usage_percent="0",vdom="root"} 14
	fortigate_vdom_resources{current_usage="4",custom_max="0",global_max="0",guaranteed="0",if_name="service-group",max_custom_value="4000",max_guaranteed_value="4000",min_custom_value="4",min_guaranteed_value="0",usage_percent="0",vdom="root"} 10
	fortigate_vdom_resources{current_usage="4",custom_max="0",global_max="0",guaranteed="0",if_name="service-group",max_custom_value="4000",max_guaranteed_value="4000",min_custom_value="4",min_guaranteed_value="0",usage_percent="0",vdom="test1"} 10
	fortigate_vdom_resources{current_usage="52",custom_max="0",global_max="0",guaranteed="0",if_name="session",max_custom_value="22000000",max_guaranteed_value="21450000",min_custom_value="1",min_guaranteed_value="0",usage_percent="0",vdom="root"} 0
	fortigate_vdom_resources{current_usage="7",custom_max="0",global_max="201024",guaranteed="0",if_name="firewall-policy",max_custom_value="201024",max_guaranteed_value="201021",min_custom_value="7",min_guaranteed_value="0",usage_percent="0",vdom="test1"} 6
	fortigate_vdom_resources{current_usage="8",custom_max="0",global_max="113192",guaranteed="0",if_name="firewall-address",max_custom_value="113192",max_guaranteed_value="113143",min_custom_value="8",min_guaranteed_value="0",usage_percent="0",vdom="root"} 7
	fortigate_vdom_resources{current_usage="87",custom_max="0",global_max="0",guaranteed="0",if_name="custom-service",max_custom_value="10240",max_guaranteed_value="10240",min_custom_value="87",min_guaranteed_value="0",usage_percent="0",vdom="root"} 9
	fortigate_vdom_resources{current_usage="87",custom_max="0",global_max="0",guaranteed="0",if_name="custom-service",max_custom_value="10240",max_guaranteed_value="10240",min_custom_value="87",min_guaranteed_value="0",usage_percent="0",vdom="test1"} 9
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
