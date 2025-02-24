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
