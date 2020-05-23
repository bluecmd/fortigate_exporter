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

func (c *fakeClient) Get(path string, obj interface{}) error {
	return json.Unmarshal(c.data[path], obj)
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
	if !probeSystem(c, r) {
		t.Errorf("probeSystem() returned non-success")
	}

	em := `
	# HELP fortigate_system_version Contains the system version and build information
	# TYPE fortigate_system_version gauge
	fortigate_system_version{build="1337",serial="S/N",version="1234"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
