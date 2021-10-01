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
	# HELP fortigate_fortimanager_connection_status Fortimanager status ID
	# TYPE fortigate_fortimanager_connection_status gauge
	fortigate_fortimanager_connection_status{mode="normal",status="down",vdom="VDOM1"} 0
	fortigate_fortimanager_connection_status{mode="normal",status="down",vdom="root"} 0
	fortigate_fortimanager_connection_status{mode="normal",status="handshake",vdom="VDOM1"} 0
	fortigate_fortimanager_connection_status{mode="normal",status="handshake",vdom="root"} 0
	fortigate_fortimanager_connection_status{mode="normal",status="up",vdom="VDOM1"} 1
	fortigate_fortimanager_connection_status{mode="normal",status="up",vdom="root"} 1
	# HELP fortigate_fortimanager_registration_status Fortimanager registration status ID
	# TYPE fortigate_fortimanager_registration_status gauge
	fortigate_fortimanager_registration_status{mode="normal",status="inprogress",vdom="VDOM1"} 0
	fortigate_fortimanager_registration_status{mode="normal",status="inprogress",vdom="root"} 0
	fortigate_fortimanager_registration_status{mode="normal",status="registered",vdom="VDOM1"} 1
	fortigate_fortimanager_registration_status{mode="normal",status="registered",vdom="root"} 1
	fortigate_fortimanager_registration_status{mode="normal",status="unknown",vdom="VDOM1"} 0
	fortigate_fortimanager_registration_status{mode="normal",status="unknown",vdom="root"} 0
	fortigate_fortimanager_registration_status{mode="normal",status="unregistered",vdom="VDOM1"} 0
	fortigate_fortimanager_registration_status{mode="normal",status="unregistered",vdom="root"} 0
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
