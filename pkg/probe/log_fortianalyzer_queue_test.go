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
		t.Errorf("probeSystemStatus() returned non-success")
	}

	em := `
	# HELP fortigate_log_fortianalyzer_queue_cached Cached logs in fortianalyzer queue
	# TYPE fortigate_log_fortianalyzer_queue_cached gauge
	fortigate_log_fortianalyzer_queue_cached{vdom="root"} 0
	# HELP fortigate_log_fortianalyzer_queue_connected Fortianalyzer queue connected state
	# TYPE fortigate_log_fortianalyzer_queue_connected gauge
	fortigate_log_fortianalyzer_queue_connected{vdom="root"} 1
	# HELP fortigate_log_fortianalyzer_queue_failed Failed logs in fortianalyzer queue
	# TYPE fortigate_log_fortianalyzer_queue_failed gauge
	fortigate_log_fortianalyzer_queue_failed{vdom="root"} 0
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
