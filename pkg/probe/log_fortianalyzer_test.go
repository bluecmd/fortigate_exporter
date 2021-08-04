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
	# HELP fortigate_log_fortianalyzer_connection Fortianalyzer connection state
	# TYPE fortigate_log_fortianalyzer_connection gauge
	fortigate_log_fortianalyzer_connection{connection="allow",vdom="root"} 1
	# HELP fortigate_log_fortianalyzer_received Received logs in fortianalyzer
	# TYPE fortigate_log_fortianalyzer_received gauge
	fortigate_log_fortianalyzer_received{vdom="root"} 999
	# HELP fortigate_log_fortianalyzer_registration Fortianalyzer registration state
	# TYPE fortigate_log_fortianalyzer_registration gauge
	fortigate_log_fortianalyzer_registration{registration="registered",vdom="root"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
