package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestProbeWifiAPStatus(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/wifi/ap_status", "testdata/wifi-ap-status.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeWifiAPStatus, c, r) {
		t.Errorf("probeWifiAPStatus() returned non-success")
	}

	em := `
        # HELP fortigate_wifi_access_point_allowed_client_count Maximum number of clients which are allowed to connect
        # TYPE fortigate_wifi_access_point_allowed_client_count gauge
        fortigate_wifi_access_point_allowed_client_count{vdom="root"} 0
        # HELP fortigate_wifi_access_point_client_count Number of connected clients
        # TYPE fortigate_wifi_access_point_client_count gauge
        fortigate_wifi_access_point_client_count{vdom="root"} 17
        # HELP fortigate_wifi_access_points Number of connected access points by status
        # TYPE fortigate_wifi_access_points gauge
        fortigate_wifi_access_points{status="active",vdom="root"} 3
        fortigate_wifi_access_points{status="down",vdom="root"} 0
        fortigate_wifi_access_points{status="rebooting",vdom="root"} 0
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
