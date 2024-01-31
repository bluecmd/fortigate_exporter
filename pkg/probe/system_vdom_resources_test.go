package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestSystemVDOMResourcesAll(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/vdom-resource", "testdata/system-vdom-resources.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemVDOMResourceUsage, c, r) {
		t.Errorf("probeSystemVDOMResourcesAll() returned non-success")
	}

	em := `
	# HELP fortigate_vdom_current_session_usage_percent Current percent usage of sessions, per VDOM
	# TYPE fortigate_vdom_current_session_usage_percent gauge
	fortigate_vdom_current_session_usage_percent{if_name="session",vdom="root"} 0
	fortigate_vdom_current_session_usage_percent{if_name="session",vdom="test1"} 0
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
