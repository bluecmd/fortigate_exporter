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

func TestSystemResourceUsage(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/resource/usage", "testdata/usage.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemResourceUsage, c, r) {
		t.Errorf("probeSystemResourceUsage() returned non-success")
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
	if !testProbe(probeSystemVDOMResources, c, r) {
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
