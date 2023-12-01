package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestUserFirewall(t *testing.T) {
	c := newFakeClient()
	c.prepare("/api/v2/monitor/user/firewall", "testdata/user-firewall.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeUserFirewall, c, r) {
		t.Errorf("probeUserFirewall() returned non-success")
	}

	em := `
	# HELP fortigate_user_firewall_duration_seconds Duration of user firewall activity in seconds
	# TYPE fortigate_user_firewall_duration_seconds gauge
	fortigate_user_firewall_duration_seconds{ipaddr="192.168.24.18",method="Firewall",type="auth_logon",vdom="VD_ES-WIFI"} 21476
	fortigate_user_firewall_duration_seconds{ipaddr="192.168.24.25",method="Firewall",type="auth_logon",vdom="VD_ES-WIFI"} 6366
	fortigate_user_firewall_duration_seconds{ipaddr="192.168.27.150",method="Firewall",type="auth_logon",vdom="VD_ES-WIFI"} 39266
	# HELP fortigate_user_firewall_traffic_bytes Traffic volume in bytes for user firewall activity
	# TYPE fortigate_user_firewall_traffic_bytes gauge
	fortigate_user_firewall_traffic_bytes{ipaddr="192.168.24.18",method="Firewall",type="auth_logon",vdom="VD_ES-WIFI"} 908744605
	fortigate_user_firewall_traffic_bytes{ipaddr="192.168.24.25",method="Firewall",type="auth_logon",vdom="VD_ES-WIFI"} 1738875
	fortigate_user_firewall_traffic_bytes{ipaddr="192.168.27.150",method="Firewall",type="auth_logon",vdom="VD_ES-WIFI"} 79707852
	`
	

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
