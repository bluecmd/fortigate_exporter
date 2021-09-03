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
