package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestSystemCentralManagementStatus(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/central-management/status", "testdata/system-central-management-status.jsonnet")
	r := prometheus.NewPedanticRegistry()
	meta := &TargetMetadata{
		VersionMajor: 7,
		VersionMinor: 4,
	}
	if !testProbeWithMetadata(probeSystemCentralManagementStatus, c, meta, r) {
		t.Errorf("probeSystemCentralManagementStatus() returned non-success")
	}

	em := `
	# HELP fortigate_central_management_connection_status Fortimanager status
	# TYPE fortigate_central_management_connection_status gauge
	fortigate_central_management_connection_status{mode="",status="down",vdom="test1"} 0
	fortigate_central_management_connection_status{mode="",status="down",vdom="root"} 0
	fortigate_central_management_connection_status{mode="",status="handshake",vdom="test1"} 1
	fortigate_central_management_connection_status{mode="",status="handshake",vdom="root"} 0
	fortigate_central_management_connection_status{mode="",status="up",vdom="test1"} 0
	fortigate_central_management_connection_status{mode="",status="up",vdom="root"} 1
	# HELP fortigate_central_management_registration_status Fortimanager registration status
	# TYPE fortigate_central_management_registration_status gauge
	fortigate_central_management_registration_status{mode="",status="inprogress",vdom="test1"} 0
	fortigate_central_management_registration_status{mode="",status="inprogress",vdom="root"} 0
	fortigate_central_management_registration_status{mode="",status="registered",vdom="test1"} 0
	fortigate_central_management_registration_status{mode="",status="registered",vdom="root"} 1
	fortigate_central_management_registration_status{mode="",status="unknown",vdom="test1"} 1
	fortigate_central_management_registration_status{mode="",status="unknown",vdom="root"} 0
	fortigate_central_management_registration_status{mode="",status="unregistered",vdom="test1"} 0
	fortigate_central_management_registration_status{mode="",status="unregistered",vdom="root"} 0
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
