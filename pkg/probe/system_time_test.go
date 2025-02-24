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

func TestSystemTime(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/time", "testdata/system-time.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemTime, c, r) {
		t.Errorf("probeSystemTime() returned non-success")
	}

	em := `
	# HELP fortigate_time_seconds System epoch time in seconds
	# TYPE fortigate_time_seconds gauge
	fortigate_time_seconds 1.630313596e+09
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
