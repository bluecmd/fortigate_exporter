package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestLinkStatus(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/link-monitor", "testdata/link-monitor.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemLinkMonitor, c, r) {
		t.Errorf("probeSystemLinkMonitor() returned non-success")
	}

	em := `
        # HELP fortigate_link_active_sessions Number of sessions active on this link
        # TYPE fortigate_link_active_sessions gauge
        fortigate_link_active_sessions{link="wan1",monitor="wan-mon",vdom="root"} 77
        # HELP fortigate_link_bandwidth_rx_byte_per_second Bandwidth available on this link for sending
        # TYPE fortigate_link_bandwidth_rx_byte_per_second gauge
        fortigate_link_bandwidth_rx_byte_per_second{link="wan1",monitor="wan-mon",vdom="root"} 4194.625
        # HELP fortigate_link_bandwidth_tx_byte_per_second Bandwidth available on this link for sending
        # TYPE fortigate_link_bandwidth_tx_byte_per_second gauge
        fortigate_link_bandwidth_tx_byte_per_second{link="wan1",monitor="wan-mon",vdom="root"} 8582.125
        # HELP fortigate_link_latency_jitter_seconds Average of the latency jitter  on this link based on the last 30 probes in seconds
        # TYPE fortigate_link_latency_jitter_seconds gauge
        fortigate_link_latency_jitter_seconds{link="wan1",monitor="wan-mon",vdom="root"} 0.0011268666982650758
        # HELP fortigate_link_latency_seconds Average latency of this link based on the last 30 probes in seconds
        # TYPE fortigate_link_latency_seconds gauge
        fortigate_link_latency_seconds{link="wan1",monitor="wan-mon",vdom="root"} 0.006810200214385986
        # HELP fortigate_link_packet_loss_ratio Percentage of packets lost relative to  all sent based on the last 30 probes
        # TYPE fortigate_link_packet_loss_ratio gauge
        fortigate_link_packet_loss_ratio{link="wan1",monitor="wan-mon",vdom="root"} 0
        # HELP fortigate_link_packet_received_total Number of packets received on this link
        # TYPE fortigate_link_packet_received_total counter
        fortigate_link_packet_received_total{link="wan1",monitor="wan-mon",vdom="root"} 278807
        # HELP fortigate_link_packet_sent_total Number of packets sent on this link
        # TYPE fortigate_link_packet_sent_total counter
        fortigate_link_packet_sent_total{link="wan1",monitor="wan-mon",vdom="root"} 278878
        # HELP fortigate_link_status Signals the status of the link. 1 means that this state is present in every other case the value is 0
        # TYPE fortigate_link_status gauge
        fortigate_link_status{link="wan1",monitor="wan-mon",state="down",vdom="root"} 0
        fortigate_link_status{link="wan1",monitor="wan-mon",state="error",vdom="root"} 0
        fortigate_link_status{link="wan1",monitor="wan-mon",state="unknown",vdom="root"} 0
        fortigate_link_status{link="wan1",monitor="wan-mon",state="up",vdom="root"} 1
        # HELP fortigate_link_status_change_time_seconds Unix timestamp describing the time when the last status change has occurred
        # TYPE fortigate_link_status_change_time_seconds gauge
        fortigate_link_status_change_time_seconds{link="wan1",monitor="wan-mon",vdom="root"} 1.61291602e+09
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}

// testing status error and empty results
func TestLinkStatusFailure(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/link-monitor", "testdata/link-monitor-error.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemLinkMonitor, c, r) {
		t.Errorf("probeSystemLinkMonitor() returned non-success")
	}

	em := `
        # HELP fortigate_link_status Signals the status of the link. 1 means that this state is present in every other case the value is 0
        # TYPE fortigate_link_status gauge
        fortigate_link_status{link="port3",monitor="google-dns-v4",state="down",vdom="bluecmd"} 0
        fortigate_link_status{link="port3",monitor="google-dns-v4",state="error",vdom="bluecmd"} 1
        fortigate_link_status{link="port3",monitor="google-dns-v4",state="unknown",vdom="bluecmd"} 0
        fortigate_link_status{link="port3",monitor="google-dns-v4",state="up",vdom="bluecmd"} 0
        fortigate_link_status{link="port3",monitor="google-dns-v6",state="down",vdom="bluecmd"} 0
        fortigate_link_status{link="port3",monitor="google-dns-v6",state="error",vdom="bluecmd"} 1
        fortigate_link_status{link="port3",monitor="google-dns-v6",state="unknown",vdom="bluecmd"} 0
        fortigate_link_status{link="port3",monitor="google-dns-v6",state="up",vdom="bluecmd"} 0
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}

func TestLinkStatusUnknown(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/link-monitor", "testdata/link-monitor-unknown.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemLinkMonitor, c, r) {
		t.Errorf("probeSystemLinkMonitor() returned non-success")
	}

	em := `
        # HELP fortigate_link_status Signals the status of the link. 1 means that this state is present in every other case the value is 0
        # TYPE fortigate_link_status gauge
        fortigate_link_status{link="port3",monitor="google-dns-v4",state="down",vdom="bluecmd"} 0
        fortigate_link_status{link="port3",monitor="google-dns-v4",state="error",vdom="bluecmd"} 0
        fortigate_link_status{link="port3",monitor="google-dns-v4",state="unknown",vdom="bluecmd"} 1
        fortigate_link_status{link="port3",monitor="google-dns-v4",state="up",vdom="bluecmd"} 0
        fortigate_link_status{link="port3",monitor="google-dns-v6",state="down",vdom="bluecmd"} 0
        fortigate_link_status{link="port3",monitor="google-dns-v6",state="error",vdom="bluecmd"} 0
        fortigate_link_status{link="port3",monitor="google-dns-v6",state="unknown",vdom="bluecmd"} 1
        fortigate_link_status{link="port3",monitor="google-dns-v6",state="up",vdom="bluecmd"} 0
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
