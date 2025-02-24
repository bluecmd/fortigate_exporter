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

func TestFirewallLoadBalance(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/firewall/load-balance?vdom=*&start=0&count=1000", "testdata/fw-loadbalancers.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeFirewallLoadBalance, c, r) {
		t.Errorf("testLoadBalanceServers() returned non-success")
	}

	em := `
	# HELP fortigate_lb_real_server_active_sessions Number of sessions active on this real server
	# TYPE fortigate_lb_real_server_active_sessions gauge
	fortigate_lb_real_server_active_sessions{id="1",vdom="root",virtual_server="LB-EXAMPLE"} 999
	fortigate_lb_real_server_active_sessions{id="2",vdom="root",virtual_server="LB-EXAMPLE"} 3
	fortigate_lb_real_server_active_sessions{id="3",vdom="root",virtual_server="LB-EXAMPLE"} 0
	fortigate_lb_real_server_active_sessions{id="4",vdom="root",virtual_server="LB-EXAMPLE"} 0
	# HELP fortigate_lb_real_server_info Info metric regarding real servers
	# TYPE fortigate_lb_real_server_info gauge
	fortigate_lb_real_server_info{id="1",ip="10.10.0.1",port="8080",vdom="root",virtual_server="LB-EXAMPLE"} 1
	fortigate_lb_real_server_info{id="2",ip="10.10.0.2",port="8080",vdom="root",virtual_server="LB-EXAMPLE"} 1
	fortigate_lb_real_server_info{id="3",ip="10.10.0.3",port="8080",vdom="root",virtual_server="LB-EXAMPLE"} 1
	fortigate_lb_real_server_info{id="4",ip="10.10.0.4",port="8080",vdom="root",virtual_server="LB-EXAMPLE"} 1
	# HELP fortigate_lb_real_server_mode Mode of this real server: active, standby or disabled
	# TYPE fortigate_lb_real_server_mode gauge
	fortigate_lb_real_server_mode{id="1",mode="active",vdom="root",virtual_server="LB-EXAMPLE"} 1
	fortigate_lb_real_server_mode{id="1",mode="disabled",vdom="root",virtual_server="LB-EXAMPLE"} 0
	fortigate_lb_real_server_mode{id="1",mode="standby",vdom="root",virtual_server="LB-EXAMPLE"} 0
	fortigate_lb_real_server_mode{id="2",mode="active",vdom="root",virtual_server="LB-EXAMPLE"} 0
	fortigate_lb_real_server_mode{id="2",mode="disabled",vdom="root",virtual_server="LB-EXAMPLE"} 0
	fortigate_lb_real_server_mode{id="2",mode="standby",vdom="root",virtual_server="LB-EXAMPLE"} 1
	fortigate_lb_real_server_mode{id="3",mode="active",vdom="root",virtual_server="LB-EXAMPLE"} 0
	fortigate_lb_real_server_mode{id="3",mode="disabled",vdom="root",virtual_server="LB-EXAMPLE"} 1
	fortigate_lb_real_server_mode{id="3",mode="standby",vdom="root",virtual_server="LB-EXAMPLE"} 0
	fortigate_lb_real_server_mode{id="4",mode="active",vdom="root",virtual_server="LB-EXAMPLE"} 0
	fortigate_lb_real_server_mode{id="4",mode="disabled",vdom="root",virtual_server="LB-EXAMPLE"} 1
	fortigate_lb_real_server_mode{id="4",mode="standby",vdom="root",virtual_server="LB-EXAMPLE"} 0
	# HELP fortigate_lb_real_server_processed_bytes_total Number of bytes processed by this real server
	# TYPE fortigate_lb_real_server_processed_bytes_total counter
	fortigate_lb_real_server_processed_bytes_total{id="1",vdom="root",virtual_server="LB-EXAMPLE"} 38260
	fortigate_lb_real_server_processed_bytes_total{id="2",vdom="root",virtual_server="LB-EXAMPLE"} 0
	fortigate_lb_real_server_processed_bytes_total{id="3",vdom="root",virtual_server="LB-EXAMPLE"} 0
	fortigate_lb_real_server_processed_bytes_total{id="4",vdom="root",virtual_server="LB-EXAMPLE"} 0
	# HELP fortigate_lb_real_server_rtt_seconds Round Trip Time (RTT) for this real server. A RTT of 1 ms or less is reported as 1 ms (0.001 s). A RTT of -1 indicates a parsing error.
	# TYPE fortigate_lb_real_server_rtt_seconds gauge
	fortigate_lb_real_server_rtt_seconds{id="1",vdom="root",virtual_server="LB-EXAMPLE"} 0.001
	fortigate_lb_real_server_rtt_seconds{id="2",vdom="root",virtual_server="LB-EXAMPLE"} 0.357
	fortigate_lb_real_server_rtt_seconds{id="3",vdom="root",virtual_server="LB-EXAMPLE"} NaN
	fortigate_lb_real_server_rtt_seconds{id="4",vdom="root",virtual_server="LB-EXAMPLE"} NaN
	# HELP fortigate_lb_real_server_status Status of this real server: up, down or unknown
	# TYPE fortigate_lb_real_server_status gauge
	fortigate_lb_real_server_status{id="1",state="down",vdom="root",virtual_server="LB-EXAMPLE"} 0
	fortigate_lb_real_server_status{id="1",state="unknown",vdom="root",virtual_server="LB-EXAMPLE"} 0
	fortigate_lb_real_server_status{id="1",state="up",vdom="root",virtual_server="LB-EXAMPLE"} 1
	fortigate_lb_real_server_status{id="2",state="down",vdom="root",virtual_server="LB-EXAMPLE"} 1
	fortigate_lb_real_server_status{id="2",state="unknown",vdom="root",virtual_server="LB-EXAMPLE"} 0
	fortigate_lb_real_server_status{id="2",state="up",vdom="root",virtual_server="LB-EXAMPLE"} 0
	fortigate_lb_real_server_status{id="3",state="down",vdom="root",virtual_server="LB-EXAMPLE"} 0
	fortigate_lb_real_server_status{id="3",state="unknown",vdom="root",virtual_server="LB-EXAMPLE"} 1
	fortigate_lb_real_server_status{id="3",state="up",vdom="root",virtual_server="LB-EXAMPLE"} 0
	fortigate_lb_real_server_status{id="4",state="down",vdom="root",virtual_server="LB-EXAMPLE"} 0
	fortigate_lb_real_server_status{id="4",state="unknown",vdom="root",virtual_server="LB-EXAMPLE"} 1
	fortigate_lb_real_server_status{id="4",state="up",vdom="root",virtual_server="LB-EXAMPLE"} 0
	# HELP fortigate_lb_virtual_server_info Info metric regarding virtual servers
	# TYPE fortigate_lb_virtual_server_info gauge
	fortigate_lb_virtual_server_info{ip="169.254.1.1",name="LB-EXAMPLE",port="80",type="http",vdom="root"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}

func TestLoadBalanceServers_6_0_5(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/firewall/load-balance?vdom=*&start=0&count=1000", "testdata/fw-loadbalancers_6_0_5.jsonnet")
	r := prometheus.NewPedanticRegistry()
	meta := &TargetMetadata{
		VersionMajor: 6,
		VersionMinor: 0,
	}
	if !testProbeWithMetadata(probeFirewallLoadBalance, c, meta, r) {
		t.Errorf("TestLoadBalanceServers_6_0_5() failed, but should have succeeded")
	}
}
