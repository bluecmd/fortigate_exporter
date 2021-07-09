package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestProbeClients(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/wifi/client", "testdata/wifi-client.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeWifiClients, c, r) {
		t.Errorf("probeWifiAPStatus() returned non-success")
	}

	em := `
        # HELP fortigate_wifi_client_bandwidth_rx_bit_per_second Bandwidth for receiving traffic
        # TYPE fortigate_wifi_client_bandwidth_rx_bit_per_second gauge
        fortigate_wifi_client_bandwidth_rx_bit_per_second{mac="00:00:00:00:00:00",vdom="root"} 0
        fortigate_wifi_client_bandwidth_rx_bit_per_second{mac="00:00:00:AA:00:00",vdom="root"} 0
        # HELP fortigate_wifi_client_bandwidth_tx_bit_per_second Bandwidth for transmitting traffic
        # TYPE fortigate_wifi_client_bandwidth_tx_bit_per_second gauge
        fortigate_wifi_client_bandwidth_tx_bit_per_second{mac="00:00:00:00:00:00",vdom="root"} 0
        fortigate_wifi_client_bandwidth_tx_bit_per_second{mac="00:00:00:AA:00:00",vdom="root"} 0
        # HELP fortigate_wifi_client_data_rate_bit_per_second Data rate of the client connection
        # TYPE fortigate_wifi_client_data_rate_bit_per_second gauge
        fortigate_wifi_client_data_rate_bit_per_second{mac="00:00:00:00:00:00",vdom="root"} 1e+06
        fortigate_wifi_client_data_rate_bit_per_second{mac="00:00:00:AA:00:00",vdom="root"} 1.3e+08
        # HELP fortigate_wifi_client_info Number of connected access points by status
        # TYPE fortigate_wifi_client_info counter
        fortigate_wifi_client_info{hostname="",mac="00:00:00:AA:00:00",vdom="root",wtp_name="3rd Floor"} 1
        fortigate_wifi_client_info{hostname="wled-WLED",mac="00:00:00:00:00:00",vdom="root",wtp_name="2nd Floor"} 1
        # HELP fortigate_wifi_client_signal_noise_decibel Signal noise on the frequency of the client
        # TYPE fortigate_wifi_client_signal_noise_decibel gauge
        fortigate_wifi_client_signal_noise_decibel{mac="00:00:00:00:00:00",vdom="root"} -95
        fortigate_wifi_client_signal_noise_decibel{mac="00:00:00:AA:00:00",vdom="root"} -95
        # HELP fortigate_wifi_client_signal_strength_decibel Signal strength of the connected client
        # TYPE fortigate_wifi_client_signal_strength_decibel gauge
        fortigate_wifi_client_signal_strength_decibel{mac="00:00:00:00:00:00",vdom="root"} -59
        fortigate_wifi_client_signal_strength_decibel{mac="00:00:00:AA:00:00",vdom="root"} -59
        # HELP fortigate_wifi_client_tx_discard_percentage Percentage of discarded packets
        # TYPE fortigate_wifi_client_tx_discard_percentage gauge
        fortigate_wifi_client_tx_discard_percentage{mac="00:00:00:00:00:00",vdom="root"} 0
        fortigate_wifi_client_tx_discard_percentage{mac="00:00:00:AA:00:00",vdom="root"} 0
        # HELP fortigate_wifi_client_tx_retries_percentage Percentage of retried connection to all connection attempts
        # TYPE fortigate_wifi_client_tx_retries_percentage gauge
        fortigate_wifi_client_tx_retries_percentage{mac="00:00:00:00:00:00",vdom="root"} 0
        fortigate_wifi_client_tx_retries_percentage{mac="00:00:00:AA:00:00",vdom="root"} 0
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
