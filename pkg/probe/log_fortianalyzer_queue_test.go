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
