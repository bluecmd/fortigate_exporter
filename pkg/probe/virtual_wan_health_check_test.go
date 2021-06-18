package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestVirtualWANHealthCheck(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/virtual-wan/health-check", "testdata/virtual_wan_health_check.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeVirtualWANHealthCheck, c, r) {
		t.Errorf("probeVirtualWANHealthCheck() returned non-success")
	}

	em := `
		# HELP fortigate_virtual_wan_active_sessions Active Session count for the health check interface
		# TYPE fortigate_virtual_wan_active_sessions gauge
		fortigate_virtual_wan_active_sessions{interface="WAN1_VL300",sla="Internet Check",vdom="root"} 710
		# HELP fortigate_virtual_wan_bandwidth_rx_byte_per_second Download bandwidth of the health check interface
		# TYPE fortigate_virtual_wan_bandwidth_rx_byte_per_second gauge
		fortigate_virtual_wan_bandwidth_rx_byte_per_second{interface="WAN1_VL300",sla="Internet Check",vdom="root"} 32125.375
		# HELP fortigate_virtual_wan_bandwidth_tx_byte_per_second Upload bandwidth of the health check interface
		# TYPE fortigate_virtual_wan_bandwidth_tx_byte_per_second gauge
		fortigate_virtual_wan_bandwidth_tx_byte_per_second{interface="WAN1_VL300",sla="Internet Check",vdom="root"} 14662
		# HELP fortigate_virtual_wan_latency_jitter_seconds Measured latency jitter for this Health check
		# TYPE fortigate_virtual_wan_latency_jitter_seconds gauge
		fortigate_virtual_wan_latency_jitter_seconds{interface="WAN1_VL300",sla="Internet Check",vdom="root"} 3.116671182215214e-05
		# HELP fortigate_virtual_wan_latency_seconds Measured latency for this Health check
		# TYPE fortigate_virtual_wan_latency_seconds gauge
		fortigate_virtual_wan_latency_seconds{interface="WAN1_VL300",sla="Internet Check",vdom="root"} 0.005611332893371582
		# HELP fortigate_virtual_wan_packet_loss_ratio Measured packet loss in percentage for this Health check
		# TYPE fortigate_virtual_wan_packet_loss_ratio gauge
		fortigate_virtual_wan_packet_loss_ratio{interface="WAN1_VL300",sla="Internet Check",vdom="root"} 0
		# HELP fortigate_virtual_wan_packet_received_total Number of packets received for this Health check
		# TYPE fortigate_virtual_wan_packet_received_total gauge
		fortigate_virtual_wan_packet_received_total{interface="WAN1_VL300",sla="Internet Check",vdom="root"} 306895
		# HELP fortigate_virtual_wan_packet_sent_total Number of packets sent for this Health check
		# TYPE fortigate_virtual_wan_packet_sent_total gauge
		fortigate_virtual_wan_packet_sent_total{interface="WAN1_VL300",sla="Internet Check",vdom="root"} 306958
		# HELP fortigate_virtual_wan_status Status of the Interface. If the SD-WAN interface is disabled, disable will be returned. If the interface does not participate in the health check, error will be returned.
		# TYPE fortigate_virtual_wan_status gauge
		fortigate_virtual_wan_status{interface="WAN1_VL300",sla="Internet Check",state="disable",vdom="root"} 0
		fortigate_virtual_wan_status{interface="WAN1_VL300",sla="Internet Check",state="down",vdom="root"} 0
		fortigate_virtual_wan_status{interface="WAN1_VL300",sla="Internet Check",state="error",vdom="root"} 0
		fortigate_virtual_wan_status{interface="WAN1_VL300",sla="Internet Check",state="unknown",vdom="root"} 0
		fortigate_virtual_wan_status{interface="WAN1_VL300",sla="Internet Check",state="up",vdom="root"} 1
		fortigate_virtual_wan_status{interface="wan2",sla="Internet Check",state="disable",vdom="root"} 1
		fortigate_virtual_wan_status{interface="wan2",sla="Internet Check",state="down",vdom="root"} 0
		fortigate_virtual_wan_status{interface="wan2",sla="Internet Check",state="error",vdom="root"} 0
		fortigate_virtual_wan_status{interface="wan2",sla="Internet Check",state="unknown",vdom="root"} 0
		fortigate_virtual_wan_status{interface="wan2",sla="Internet Check",state="up",vdom="root"} 0
		# HELP fortigate_virtual_wan_status_change_time_seconds Unix timestamp describing the time when the last status change has occurred
		# TYPE fortigate_virtual_wan_status_change_time_seconds gauge
		fortigate_virtual_wan_status_change_time_seconds{interface="WAN1_VL300",sla="Internet Check",vdom="root"} 1.6141078e+09
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
