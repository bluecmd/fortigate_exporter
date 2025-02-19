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

func TestProbeWifiAPStatus(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/wifi/ap_status", "testdata/wifi-ap-status.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeWifiAPStatus, c, r) {
		t.Errorf("probeWifiAPStatus() returned non-success")
	}

	em := `
        # HELP fortigate_wifi_fabric_max_allowed_clients Maximum number of clients which are allowed to connect
        # TYPE fortigate_wifi_fabric_max_allowed_clients gauge
        fortigate_wifi_fabric_max_allowed_clients{vdom="root"} 0
        # HELP fortigate_wifi_fabric_clients Number of connected clients
        # TYPE fortigate_wifi_fabric_clients gauge
        fortigate_wifi_fabric_clients{vdom="root"} 17
        # HELP fortigate_wifi_access_points Number of connected access points by status
        # TYPE fortigate_wifi_access_points gauge
        fortigate_wifi_access_points{status="active",vdom="root"} 3
        fortigate_wifi_access_points{status="down",vdom="root"} 0
        fortigate_wifi_access_points{status="rebooting",vdom="root"} 0
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
