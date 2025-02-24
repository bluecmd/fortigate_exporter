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

func TestWebUIState(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/web-ui/state", "testdata/web-ui-state.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeWebUIState, c, r) {
		t.Errorf("probeWebUIState() returned non-success")
	}

	em := `
	# HELP fortigate_last_reboot_seconds Last system reboot epoch time in seconds
	# TYPE fortigate_last_reboot_seconds gauge
	fortigate_last_reboot_seconds 1.657116965e+09
	# HELP fortigate_last_snapshot_seconds Last snapshot epoch time in seconds
	# TYPE fortigate_last_snapshot_seconds gauge
	fortigate_last_snapshot_seconds 1.659857566e+09
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
