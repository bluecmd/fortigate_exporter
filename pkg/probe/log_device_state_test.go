package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestLogDeviceState(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/log/device/state", "testdata/log-device-state.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeLogDeviceState, c, r) {
		t.Errorf("probeLogDeviceState() returned non-success")
	}

	em := `
	# HELP fortigate_log_disk_available Is disk available for log?
	# TYPE fortigate_log_disk_available gauge
	fortigate_log_disk_available{vdom="FG-traffic"} 1
	fortigate_log_disk_available{vdom="root"} 1
	# HELP fortigate_log_disk_enabled Is disk enabled for log?
	# TYPE fortigate_log_disk_enabled gauge
	fortigate_log_disk_enabled{vdom="FG-traffic"} 1
	fortigate_log_disk_enabled{vdom="root"} 1
	# HELP fortigate_log_fortianalyzer_available Is fortianalyzer available for log?
	# TYPE fortigate_log_fortianalyzer_available gauge
	fortigate_log_fortianalyzer_available{vdom="FG-traffic"} 1
	fortigate_log_fortianalyzer_available{vdom="root"} 1
	# HELP fortigate_log_fortianalyzer_enabled Is fortianalyzer enabled for log?
	# TYPE fortigate_log_fortianalyzer_enabled gauge
	fortigate_log_fortianalyzer_enabled{vdom="FG-traffic"} 0
	fortigate_log_fortianalyzer_enabled{vdom="root"} 0
	# HELP fortigate_log_forticloud_available Is forticloud available for log?
	# TYPE fortigate_log_forticloud_available gauge
	fortigate_log_forticloud_available{vdom="FG-traffic"} 0
	fortigate_log_forticloud_available{vdom="root"} 0
	# HELP fortigate_log_forticloud_enabled Is forticloud enabled for log?
	# TYPE fortigate_log_forticloud_enabled gauge
	fortigate_log_forticloud_enabled{vdom="FG-traffic"} 0
	fortigate_log_forticloud_enabled{vdom="root"} 0
	# HELP fortigate_log_memory_available Is memory available for log?
	# TYPE fortigate_log_memory_available gauge
	fortigate_log_memory_available{vdom="FG-traffic"} 1
	fortigate_log_memory_available{vdom="root"} 1
	# HELP fortigate_log_memory_enabled Is memory enabled for log?
	# TYPE fortigate_log_memory_enabled gauge
	fortigate_log_memory_enabled{vdom="FG-traffic"} 0
	fortigate_log_memory_enabled{vdom="root"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
