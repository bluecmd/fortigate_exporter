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

func TestLicenseStatus(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/license/status/select", "testdata/license-status.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeLicenseStatus, c, r) {
		t.Errorf("probeLicenseStatus() returned non-success")
	}

	em := `
        # HELP fortigate_license_vdom_usage The amount of VDOM licenses currently used
        # TYPE fortigate_license_vdom_usage gauge
        fortigate_license_vdom_usage 114
        # HELP fortigate_license_vdom_max The total amount of VDOM licenses available
        # TYPE fortigate_license_vdom_max gauge
        fortigate_license_vdom_max 125
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
