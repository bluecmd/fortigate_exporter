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

func TestLogAnalyzerQueue(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/log/fortianalyzer-queue", "testdata/log-fortianalyzer-queue.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeLogAnalyzerQueue, c, r) {
		t.Errorf("probeLogAnalyzerQueue() returned non-success")
	}

	em := `
	# HELP fortigate_log_fortianalyzer_queue_connections Fortianalyzer queue connected state
	# TYPE fortigate_log_fortianalyzer_queue_connections gauge
	fortigate_log_fortianalyzer_queue_connections{vdom="root"} 1
	# HELP fortigate_log_fortianalyzer_queue_logs State of logs in the queue
	# TYPE fortigate_log_fortianalyzer_queue_logs gauge
	fortigate_log_fortianalyzer_queue_logs{state="cached",vdom="root"} 0
	fortigate_log_fortianalyzer_queue_logs{state="failed",vdom="root"} 0
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
