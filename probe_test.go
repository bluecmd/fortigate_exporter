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
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type fakeClient struct {
	data map[string][]byte
}

func (c *fakeClient) prepare(path string, data string) {
	c.data[path] = []byte(data)
}

func (c *fakeClient) Get(path string, query string, obj interface{}) error {
	if query != "" {
		query = "?" + query
	}
	return json.Unmarshal(c.data[path+query], obj)
}

func newFakeClient() *fakeClient {
	return &fakeClient{data: map[string][]byte{}}
}

func TestProbeSystemStatus(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/status", `{
		"serial": "S/N",
		"version": "1234",
		"build": 1337
	}`)
	r := prometheus.NewPedanticRegistry()
	if !probeSystemStatus(c, r) {
		t.Errorf("probeSystemStatus() returned non-success")
	}

	em := `
	# HELP fortigate_system_version_info System version and build information
	# TYPE fortigate_system_version_info gauge
	fortigate_system_version_info{build="1337",serial="S/N",version="1234"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}

func TestProbeSystemResources(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/resource/usage?interval=1-min&scope=global", `{
		"http_method":"GET",
		"results":{
			"cpu": [{"current":0, "historical":{}}, {"current":1, "historical":{}}, {"current":2, "historical":{}}],
			"mem": [{"current": 45}],
			"session": [{"current": 100}],
			"session6": [{"current": 50}]
		}
	}`)
	r := prometheus.NewPedanticRegistry()
	if !probeSystemResources(c, r) {
		t.Errorf("probeSystemResources() returned non-success")
	}

	em := `
	# HELP fortigate_system_cpu_usage_ratio Current resource usage ratio of system CPU, per core
	# TYPE fortigate_system_cpu_usage_ratio gauge
	fortigate_system_cpu_usage_ratio{processor="0"} 0.01
	fortigate_system_cpu_usage_ratio{processor="1"} 0.02
	# HELP fortigate_system_memory_usage_ratio Current resource usage ratio of system memory
	# TYPE fortigate_system_memory_usage_ratio gauge
	fortigate_system_memory_usage_ratio 0.45
	# HELP fortigate_system_sessions_total Current amount of system sessions, per IP version
	# TYPE fortigate_system_sessions_total gauge
	fortigate_system_sessions_total{protocol="ipv4"} 100
	fortigate_system_sessions_total{protocol="ipv6"} 50
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}

func TestProbeSystemVDOMResources(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/resource/usage?interval=1-min&vdom=*", `[{
		"http_method":"GET",
		"results":{
			"cpu": [{"current":0, "historical":{}}],
			"mem": [{"current": 45}],
			"session": [{"current": 100}],
			"session6": [{"current": 50}]
		},
		"vdom": "test-1"
	},{
		"http_method":"GET",
		"results":{
			"cpu": [{"current":1, "historical":{}}],
			"mem": [{"current": 46}],
			"session": [{"current": 101}],
			"session6": [{"current": 51}]
		},
		"vdom": "test-2"
	}]`)
	r := prometheus.NewPedanticRegistry()
	if !probeSystemVDOMResources(c, r) {
		t.Errorf("probeSystemVDOMResources() returned non-success")
	}

	em := `
	# HELP fortigate_system_vdom_cpu_usage_ratio Current resource usage ratio of CPU, per VDOM
	# TYPE fortigate_system_vdom_cpu_usage_ratio gauge
	fortigate_system_vdom_cpu_usage_ratio{vdom="test-1"} 0
	fortigate_system_vdom_cpu_usage_ratio{vdom="test-2"} 0.01
	# HELP fortigate_system_vdom_memory_usage_ratio Current resource usage ratio of memory, per VDOM
	# TYPE fortigate_system_vdom_memory_usage_ratio gauge
	fortigate_system_vdom_memory_usage_ratio{vdom="test-1"} 0.45
	fortigate_system_vdom_memory_usage_ratio{vdom="test-2"} 0.46
	# HELP fortigate_system_vdom_sessions_total Current amount of sessions, per VDOM and IP version
	# TYPE fortigate_system_vdom_sessions_total gauge
	fortigate_system_vdom_sessions_total{protocol="ipv4",vdom="test-1"} 100
	fortigate_system_vdom_sessions_total{protocol="ipv4",vdom="test-2"} 101
	fortigate_system_vdom_sessions_total{protocol="ipv6",vdom="test-1"} 50
	fortigate_system_vdom_sessions_total{protocol="ipv6",vdom="test-2"} 51
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
