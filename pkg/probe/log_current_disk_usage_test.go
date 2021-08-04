package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestLogCurrentDiskUsage(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/log/current-disk-usage", "testdata/log-current-disk-usage.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeLogCurrentDiskUsage, c, r) {
		t.Errorf("probeSystemStatus() returned non-success")
	}

	em := `
	# HELP fortigate_log_free_bytes Current free bytes for log
	# TYPE fortigate_log_free_bytes gauge
	fortigate_log_free_bytes{vdom="root"} 2.93e+10
	# HELP fortigate_log_total_bytes Current total bytes for log
	# TYPE fortigate_log_total_bytes gauge
	fortigate_log_total_bytes{vdom="root"} 3e+10
	# HELP fortigate_log_used_bytes Current used bytes for log
	# TYPE fortigate_log_used_bytes gauge
	fortigate_log_used_bytes{vdom="root"} 7e+08
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
