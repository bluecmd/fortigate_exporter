package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestLogAnalyzer(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/log/fortianalyzer", "testdata/log-fortianalyzer.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeLogAnalyzer, c, r) {
		t.Errorf("probeSystemStatus() returned non-success")
	}

	em := `
	# HELP fortigate_log_fortianalyzer_logs_received Received logs in fortianalyzer
	# TYPE fortigate_log_fortianalyzer_logs_received gauge
	fortigate_log_fortianalyzer_logs_received{vdom="root"} 999
	# HELP fortigate_log_fortianalyzer_registration_info Fortianalyzer state info
	# TYPE fortigate_log_fortianalyzer_registration_info gauge
	fortigate_log_fortianalyzer_registration_info{connection="allow",registration="registered",vdom="root"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
