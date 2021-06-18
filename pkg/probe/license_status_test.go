package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestLicenseStatus(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/license/status/select", "testdata/license-status.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeLicenseStatus, c, r) {
		t.Errorf("probeLicenseStatus() returned non-success")
	}

	em := `
        # HELP fortigate_license_vdom_usage The amount of VDOM licenses currently used
        # TYPE fortigate_license_vdom_usage gauge
        fortigate_license_vdom_usage 114
        # HELP fortigate_license_vdom_max The total amount of VDOM licenses available
        # TYPE fortigate_license_vdom_max gauge
        fortigate_license_vdom_max 125
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
