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
	c.prepare("api/v2/monitor/system/status", "testdata/status.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !probeSystemStatus(c, r) {
		t.Errorf("probeSystemStatus() returned non-success")
	}

	em := `
	# HELP fortigate_system_version_info System version and build information
	# TYPE fortigate_system_version_info gauge
	fortigate_system_version_info{build="1112",serial="FGVMEVZFNTS3OAC8",version="v6.2.4"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}

func TestProbeSystemResources(t *testing.T) {
	c := newFakeClient()
	c.prepare(
		"api/v2/monitor/system/resource/usage?interval=1-min&scope=global",
		"testdata/usage.jsonnet",
	)
	r := prometheus.NewPedanticRegistry()
	if !probeSystemResources(c, r) {
		t.Errorf("probeSystemResources() returned non-success")
	}

	em := `
	# HELP fortigate_system_cpu_usage_ratio Current resource usage ratio of system CPU, per core
	# TYPE fortigate_system_cpu_usage_ratio gauge
	fortigate_system_cpu_usage_ratio{processor="0"} 0.32
	# HELP fortigate_system_memory_usage_ratio Current resource usage ratio of system memory
	# TYPE fortigate_system_memory_usage_ratio gauge
	fortigate_system_memory_usage_ratio 0.76
	# HELP fortigate_current_sessions Current amount of sessions, per IP version
	# TYPE fortigate_current_sessions gauge
	fortigate_current_sessions{protocol="ipv4"} 5
	fortigate_current_sessions{protocol="ipv6"} 1
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}

func TestProbeSystemVDOMResources(t *testing.T) {
	c := newFakeClient()
	c.prepare(
		"api/v2/monitor/system/resource/usage?interval=1-min&vdom=*",
		"testdata/usage-vdom.jsonnet",
	)

	r := prometheus.NewPedanticRegistry()
	if !probeSystemVDOMResources(c, r) {
		t.Errorf("probeSystemVDOMResources() returned non-success")
	}

	em := `
	# HELP fortigate_vdom_system_cpu_usage_ratio Current resource usage ratio of CPU, per VDOM
	# TYPE fortigate_vdom_system_cpu_usage_ratio gauge
	fortigate_vdom_system_cpu_usage_ratio{vdom="FG-traffic"} 0
	fortigate_vdom_system_cpu_usage_ratio{vdom="root"} 0.01
	# HELP fortigate_vdom_system_memory_usage_ratio Current resource usage ratio of memory, per VDOM
	# TYPE fortigate_vdom_system_memory_usage_ratio gauge
	fortigate_vdom_system_memory_usage_ratio{vdom="FG-traffic"} 0
	fortigate_vdom_system_memory_usage_ratio{vdom="root"} 0.78
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
