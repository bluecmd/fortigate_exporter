package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestSystemFortimanagerStatus(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/fortimanager/status", "testdata/system-fortimanager-status.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemFortimanagerStatus, c, r) {
		t.Errorf("probeSystemFortimanagerStatus() returned non-success")
	}

	em := `
	# HELP fortigate_fortimanager_info Fortimanager infos
	# TYPE fortigate_fortimanager_info gauge
	fortigate_fortimanager_info{connection_status="2",mode="normal",registration_status="2",vdom="VDOM1"} 1
	fortigate_fortimanager_info{connection_status="2",mode="normal",registration_status="2",vdom="root"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
