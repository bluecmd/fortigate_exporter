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
		t.Errorf("probeLogCurrentDiskUsage() returned non-success")
	}

	em := `
	# HELP fortigate_log_disk_total_bytes Disk total bytes for log
	# TYPE fortigate_log_disk_total_bytes gauge
	fortigate_log_disk_total_bytes{vdom="root"} 3e+10
	# HELP fortigate_log_disk_used_bytes Disk used bytes for log
	# TYPE fortigate_log_disk_used_bytes gauge
	fortigate_log_disk_used_bytes{vdom="root"} 7e+08
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
