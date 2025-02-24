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

func TestVPNSslStats(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/vpn/ssl/stats", "testdata/vpn-stats.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeVPNSslStats, c, r) {
		t.Errorf("probeVPNSslStats() returned non-success")
	}

	em := `
	# HELP fortigate_vpn_ssl_connections Number of current SSL VPN connections
	# TYPE fortigate_vpn_ssl_connections gauge
	fortigate_vpn_ssl_connections{vdom="root"} 2
	# HELP fortigate_vpn_ssl_tunnels Number of current SSL VPN tunnels
	# TYPE fortigate_vpn_ssl_tunnels gauge
	fortigate_vpn_ssl_tunnels{vdom="root"} 2
	# HELP fortigate_vpn_ssl_users Number of current SSL VPN users
	# TYPE fortigate_vpn_ssl_users gauge
	fortigate_vpn_ssl_users{vdom="root"} 3
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
