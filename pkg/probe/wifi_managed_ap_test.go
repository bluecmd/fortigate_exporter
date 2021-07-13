package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestProbeWifiManagedAP(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/wifi/managed_ap", "testdata/wifi-managed-ap.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeWifiManagedAP, c, r) {
		t.Errorf("probeWifiAPStatus() returned non-success")
	}

	em := `
        # HELP fortigate_wifi_managed_ap_cpu_usage_ratio CPU usage of the access point
        # TYPE fortigate_wifi_managed_ap_cpu_usage_ratio gauge
        fortigate_wifi_managed_ap_cpu_usage_ratio{ap_name="1st Floor",vdom="root"} 0.09
        fortigate_wifi_managed_ap_cpu_usage_ratio{ap_name="2nd Floor",vdom="root"} 0.08
        fortigate_wifi_managed_ap_cpu_usage_ratio{ap_name="3rd Floor",vdom="root"} 0.08
        # HELP fortigate_wifi_managed_ap_info Infos about a managed access point
        # TYPE fortigate_wifi_managed_ap_info counter
        fortigate_wifi_managed_ap_info{ap_name="1st Floor",ap_profile="athome",os_version="FP221E-v6.4-build0460",serial="FP221E0000000000",vdom="root"} 1
        fortigate_wifi_managed_ap_info{ap_name="2nd Floor",ap_profile="athome",os_version="FP221E-v6.4-build0460",serial="FP221E0000000000",vdom="root"} 1
        fortigate_wifi_managed_ap_info{ap_name="3rd Floor",ap_profile="athome",os_version="FP221E-v6.4-build0460",serial="FP221E0000000000",vdom="root"} 1
        # HELP fortigate_wifi_managed_ap_interface_rx_bytes_total total number of bytes received on this interface
        # TYPE fortigate_wifi_managed_ap_interface_rx_bytes_total gauge
        fortigate_wifi_managed_ap_interface_rx_bytes_total{ap_name="1st Floor",interface="lan1",vdom="root"} 2.796197197e+09
        fortigate_wifi_managed_ap_interface_rx_bytes_total{ap_name="2nd Floor",interface="lan1",vdom="root"} 5.90530263e+09
        fortigate_wifi_managed_ap_interface_rx_bytes_total{ap_name="3rd Floor",interface="lan1",vdom="root"} 1.03099754e+09
        # HELP fortigate_wifi_managed_ap_interface_rx_dropped_packets_total total number of dropped packets received on this interface
        # TYPE fortigate_wifi_managed_ap_interface_rx_dropped_packets_total gauge
        fortigate_wifi_managed_ap_interface_rx_dropped_packets_total{ap_name="1st Floor",interface="lan1",vdom="root"} 5918
        fortigate_wifi_managed_ap_interface_rx_dropped_packets_total{ap_name="2nd Floor",interface="lan1",vdom="root"} 5920
        fortigate_wifi_managed_ap_interface_rx_dropped_packets_total{ap_name="3rd Floor",interface="lan1",vdom="root"} 5920
        # HELP fortigate_wifi_managed_ap_interface_rx_errors_total total number of errors received on this interface
        # TYPE fortigate_wifi_managed_ap_interface_rx_errors_total gauge
        fortigate_wifi_managed_ap_interface_rx_errors_total{ap_name="1st Floor",interface="lan1",vdom="root"} 0
        fortigate_wifi_managed_ap_interface_rx_errors_total{ap_name="2nd Floor",interface="lan1",vdom="root"} 0
        fortigate_wifi_managed_ap_interface_rx_errors_total{ap_name="3rd Floor",interface="lan1",vdom="root"} 0
        # HELP fortigate_wifi_managed_ap_interface_rx_packets_total total number of packets received on this interface
        # TYPE fortigate_wifi_managed_ap_interface_rx_packets_total gauge
        fortigate_wifi_managed_ap_interface_rx_packets_total{ap_name="1st Floor",interface="lan1",vdom="root"} 1.5144121e+07
        fortigate_wifi_managed_ap_interface_rx_packets_total{ap_name="2nd Floor",interface="lan1",vdom="root"} 6.463931e+06
        fortigate_wifi_managed_ap_interface_rx_packets_total{ap_name="3rd Floor",interface="lan1",vdom="root"} 1.157464e+06
        # HELP fortigate_wifi_managed_ap_interface_tx_bytes_total total number of bytes transferred on this interface
        # TYPE fortigate_wifi_managed_ap_interface_tx_bytes_total gauge
        fortigate_wifi_managed_ap_interface_tx_bytes_total{ap_name="1st Floor",interface="lan1",vdom="root"} 3.1582484823e+10
        fortigate_wifi_managed_ap_interface_tx_bytes_total{ap_name="2nd Floor",interface="lan1",vdom="root"} 1.757457061e+09
        fortigate_wifi_managed_ap_interface_tx_bytes_total{ap_name="3rd Floor",interface="lan1",vdom="root"} 3.49823037e+08
        # HELP fortigate_wifi_managed_ap_interface_tx_dropped_packets_total total number of dropped packets transferred on this interface
        # TYPE fortigate_wifi_managed_ap_interface_tx_dropped_packets_total gauge
        fortigate_wifi_managed_ap_interface_tx_dropped_packets_total{ap_name="1st Floor",interface="lan1",vdom="root"} 0
        fortigate_wifi_managed_ap_interface_tx_dropped_packets_total{ap_name="2nd Floor",interface="lan1",vdom="root"} 0
        fortigate_wifi_managed_ap_interface_tx_dropped_packets_total{ap_name="3rd Floor",interface="lan1",vdom="root"} 0
        # HELP fortigate_wifi_managed_ap_interface_tx_errors_total total number of errors transferred on this interface
        # TYPE fortigate_wifi_managed_ap_interface_tx_errors_total gauge
        fortigate_wifi_managed_ap_interface_tx_errors_total{ap_name="1st Floor",interface="lan1",vdom="root"} 0
        fortigate_wifi_managed_ap_interface_tx_errors_total{ap_name="2nd Floor",interface="lan1",vdom="root"} 0
        fortigate_wifi_managed_ap_interface_tx_errors_total{ap_name="3rd Floor",interface="lan1",vdom="root"} 0
        # HELP fortigate_wifi_managed_ap_interface_tx_packets_total total number of packets transferred on this interface
        # TYPE fortigate_wifi_managed_ap_interface_tx_packets_total gauge
        fortigate_wifi_managed_ap_interface_tx_packets_total{ap_name="1st Floor",interface="lan1",vdom="root"} 2.5875808e+07
        fortigate_wifi_managed_ap_interface_tx_packets_total{ap_name="2nd Floor",interface="lan1",vdom="root"} 3.867749e+06
        fortigate_wifi_managed_ap_interface_tx_packets_total{ap_name="3rd Floor",interface="lan1",vdom="root"} 1.21672e+06
        # HELP fortigate_wifi_managed_ap_join_time_seconds Unix time when the managed access point has joined the mesh
        # TYPE fortigate_wifi_managed_ap_join_time_seconds counter
        fortigate_wifi_managed_ap_join_time_seconds{ap_name="1st Floor",vdom="root"} 1.61883368e+09
        fortigate_wifi_managed_ap_join_time_seconds{ap_name="2nd Floor",vdom="root"} 1.618833935e+09
        fortigate_wifi_managed_ap_join_time_seconds{ap_name="3rd Floor",vdom="root"} 1.618833476e+09
        # HELP fortigate_wifi_managed_ap_memory_bytes_total Total memory of the managed access point
        # TYPE fortigate_wifi_managed_ap_memory_bytes_total gauge
        fortigate_wifi_managed_ap_memory_bytes_total{ap_name="1st Floor",vdom="root"} 235216
        fortigate_wifi_managed_ap_memory_bytes_total{ap_name="2nd Floor",vdom="root"} 235216
        fortigate_wifi_managed_ap_memory_bytes_total{ap_name="3rd Floor",vdom="root"} 235216
        # HELP fortigate_wifi_managed_ap_memory_free_bytes Free memory of the managed access point
        # TYPE fortigate_wifi_managed_ap_memory_free_bytes gauge
        fortigate_wifi_managed_ap_memory_free_bytes{ap_name="1st Floor",vdom="root"} 80500
        fortigate_wifi_managed_ap_memory_free_bytes{ap_name="2nd Floor",vdom="root"} 80268
        fortigate_wifi_managed_ap_memory_free_bytes{ap_name="3rd Floor",vdom="root"} 80272
        # HELP fortigate_wifi_managed_ap_radio_bandwidth_rx_bps Bandwidth of this radio for receiving
        # TYPE fortigate_wifi_managed_ap_radio_bandwidth_rx_bps gauge
        fortigate_wifi_managed_ap_radio_bandwidth_rx_bps{ap_name="1st Floor",radio_id="1",vdom="root"} 91186
        fortigate_wifi_managed_ap_radio_bandwidth_rx_bps{ap_name="1st Floor",radio_id="2",vdom="root"} 114700
        fortigate_wifi_managed_ap_radio_bandwidth_rx_bps{ap_name="1st Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_bandwidth_rx_bps{ap_name="1st Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_bandwidth_rx_bps{ap_name="1st Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_bandwidth_rx_bps{ap_name="2nd Floor",radio_id="1",vdom="root"} 97341
        fortigate_wifi_managed_ap_radio_bandwidth_rx_bps{ap_name="2nd Floor",radio_id="2",vdom="root"} 252022
        fortigate_wifi_managed_ap_radio_bandwidth_rx_bps{ap_name="2nd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_bandwidth_rx_bps{ap_name="2nd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_bandwidth_rx_bps{ap_name="2nd Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_bandwidth_rx_bps{ap_name="3rd Floor",radio_id="1",vdom="root"} 116875
        fortigate_wifi_managed_ap_radio_bandwidth_rx_bps{ap_name="3rd Floor",radio_id="2",vdom="root"} 22
        fortigate_wifi_managed_ap_radio_bandwidth_rx_bps{ap_name="3rd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_bandwidth_rx_bps{ap_name="3rd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_bandwidth_rx_bps{ap_name="3rd Floor",radio_id="5",vdom="root"} 0
        # HELP fortigate_wifi_managed_ap_radio_bandwidth_tx_bps Bandwidth of this radio for transmitting
        # TYPE fortigate_wifi_managed_ap_radio_bandwidth_tx_bps gauge
        fortigate_wifi_managed_ap_radio_bandwidth_tx_bps{ap_name="1st Floor",radio_id="1",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_bandwidth_tx_bps{ap_name="1st Floor",radio_id="2",vdom="root"} 65554
        fortigate_wifi_managed_ap_radio_bandwidth_tx_bps{ap_name="1st Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_bandwidth_tx_bps{ap_name="1st Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_bandwidth_tx_bps{ap_name="1st Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_bandwidth_tx_bps{ap_name="2nd Floor",radio_id="1",vdom="root"} 16378
        fortigate_wifi_managed_ap_radio_bandwidth_tx_bps{ap_name="2nd Floor",radio_id="2",vdom="root"} 2708
        fortigate_wifi_managed_ap_radio_bandwidth_tx_bps{ap_name="2nd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_bandwidth_tx_bps{ap_name="2nd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_bandwidth_tx_bps{ap_name="2nd Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_bandwidth_tx_bps{ap_name="3rd Floor",radio_id="1",vdom="root"} 2314
        fortigate_wifi_managed_ap_radio_bandwidth_tx_bps{ap_name="3rd Floor",radio_id="2",vdom="root"} 2068
        fortigate_wifi_managed_ap_radio_bandwidth_tx_bps{ap_name="3rd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_bandwidth_tx_bps{ap_name="3rd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_bandwidth_tx_bps{ap_name="3rd Floor",radio_id="5",vdom="root"} 0
        # HELP fortigate_wifi_managed_ap_radio_client_count Number of clients that are connected using this radio
        # TYPE fortigate_wifi_managed_ap_radio_client_count gauge
        fortigate_wifi_managed_ap_radio_client_count{ap_name="1st Floor",radio_id="1",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_client_count{ap_name="1st Floor",radio_id="2",vdom="root"} 4
        fortigate_wifi_managed_ap_radio_client_count{ap_name="1st Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_client_count{ap_name="1st Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_client_count{ap_name="1st Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_client_count{ap_name="2nd Floor",radio_id="1",vdom="root"} 6
        fortigate_wifi_managed_ap_radio_client_count{ap_name="2nd Floor",radio_id="2",vdom="root"} 2
        fortigate_wifi_managed_ap_radio_client_count{ap_name="2nd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_client_count{ap_name="2nd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_client_count{ap_name="2nd Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_client_count{ap_name="3rd Floor",radio_id="1",vdom="root"} 3
        fortigate_wifi_managed_ap_radio_client_count{ap_name="3rd Floor",radio_id="2",vdom="root"} 2
        fortigate_wifi_managed_ap_radio_client_count{ap_name="3rd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_client_count{ap_name="3rd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_client_count{ap_name="3rd Floor",radio_id="5",vdom="root"} 0
        # HELP fortigate_wifi_managed_ap_radio_info Informations about radios on managed access points
        # TYPE fortigate_wifi_managed_ap_radio_info counter
        fortigate_wifi_managed_ap_radio_info{ap_name="1st Floor",operating_channel="0",radio_id="3",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_info{ap_name="1st Floor",operating_channel="0",radio_id="4",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_info{ap_name="1st Floor",operating_channel="0",radio_id="5",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_info{ap_name="1st Floor",operating_channel="11",radio_id="1",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_info{ap_name="1st Floor",operating_channel="48",radio_id="2",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_info{ap_name="2nd Floor",operating_channel="0",radio_id="3",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_info{ap_name="2nd Floor",operating_channel="0",radio_id="4",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_info{ap_name="2nd Floor",operating_channel="0",radio_id="5",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_info{ap_name="2nd Floor",operating_channel="48",radio_id="2",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_info{ap_name="2nd Floor",operating_channel="6",radio_id="1",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_info{ap_name="3rd Floor",operating_channel="0",radio_id="3",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_info{ap_name="3rd Floor",operating_channel="0",radio_id="4",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_info{ap_name="3rd Floor",operating_channel="0",radio_id="5",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_info{ap_name="3rd Floor",operating_channel="11",radio_id="1",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_info{ap_name="3rd Floor",operating_channel="44",radio_id="2",vdom="root"} 1
        # HELP fortigate_wifi_managed_ap_radio_interfering_aps Number of interfering access points
        # TYPE fortigate_wifi_managed_ap_radio_interfering_aps gauge
        fortigate_wifi_managed_ap_radio_interfering_aps{ap_name="1st Floor",radio_id="1",vdom="root"} 5
        fortigate_wifi_managed_ap_radio_interfering_aps{ap_name="1st Floor",radio_id="2",vdom="root"} 4
        fortigate_wifi_managed_ap_radio_interfering_aps{ap_name="1st Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_interfering_aps{ap_name="1st Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_interfering_aps{ap_name="1st Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_interfering_aps{ap_name="2nd Floor",radio_id="1",vdom="root"} 3
        fortigate_wifi_managed_ap_radio_interfering_aps{ap_name="2nd Floor",radio_id="2",vdom="root"} 5
        fortigate_wifi_managed_ap_radio_interfering_aps{ap_name="2nd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_interfering_aps{ap_name="2nd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_interfering_aps{ap_name="2nd Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_interfering_aps{ap_name="3rd Floor",radio_id="1",vdom="root"} 5
        fortigate_wifi_managed_ap_radio_interfering_aps{ap_name="3rd Floor",radio_id="2",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_interfering_aps{ap_name="3rd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_interfering_aps{ap_name="3rd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_interfering_aps{ap_name="3rd Floor",radio_id="5",vdom="root"} 0
        # HELP fortigate_wifi_managed_ap_radio_operating_channel_utilization_ratio Utilization on the operating channel of the radio
        # TYPE fortigate_wifi_managed_ap_radio_operating_channel_utilization_ratio gauge
        fortigate_wifi_managed_ap_radio_operating_channel_utilization_ratio{ap_name="1st Floor",radio_id="1",vdom="root"} 0.05
        fortigate_wifi_managed_ap_radio_operating_channel_utilization_ratio{ap_name="1st Floor",radio_id="2",vdom="root"} 0.14
        fortigate_wifi_managed_ap_radio_operating_channel_utilization_ratio{ap_name="1st Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_operating_channel_utilization_ratio{ap_name="1st Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_operating_channel_utilization_ratio{ap_name="1st Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_operating_channel_utilization_ratio{ap_name="2nd Floor",radio_id="1",vdom="root"} 0.12
        fortigate_wifi_managed_ap_radio_operating_channel_utilization_ratio{ap_name="2nd Floor",radio_id="2",vdom="root"} 0.13
        fortigate_wifi_managed_ap_radio_operating_channel_utilization_ratio{ap_name="2nd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_operating_channel_utilization_ratio{ap_name="2nd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_operating_channel_utilization_ratio{ap_name="2nd Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_operating_channel_utilization_ratio{ap_name="3rd Floor",radio_id="1",vdom="root"} 0.1
        fortigate_wifi_managed_ap_radio_operating_channel_utilization_ratio{ap_name="3rd Floor",radio_id="2",vdom="root"} 0.01
        fortigate_wifi_managed_ap_radio_operating_channel_utilization_ratio{ap_name="3rd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_operating_channel_utilization_ratio{ap_name="3rd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_operating_channel_utilization_ratio{ap_name="3rd Floor",radio_id="5",vdom="root"} 0
        # HELP fortigate_wifi_managed_ap_radio_operating_tx_power_ratio Power usage on the operating channel in percent
        # TYPE fortigate_wifi_managed_ap_radio_operating_tx_power_ratio gauge
        fortigate_wifi_managed_ap_radio_operating_tx_power_ratio{ap_name="1st Floor",radio_id="1",vdom="root"} 0.17
        fortigate_wifi_managed_ap_radio_operating_tx_power_ratio{ap_name="1st Floor",radio_id="2",vdom="root"} 0.2
        fortigate_wifi_managed_ap_radio_operating_tx_power_ratio{ap_name="1st Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_operating_tx_power_ratio{ap_name="1st Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_operating_tx_power_ratio{ap_name="1st Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_operating_tx_power_ratio{ap_name="2nd Floor",radio_id="1",vdom="root"} 0.17
        fortigate_wifi_managed_ap_radio_operating_tx_power_ratio{ap_name="2nd Floor",radio_id="2",vdom="root"} 0.2
        fortigate_wifi_managed_ap_radio_operating_tx_power_ratio{ap_name="2nd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_operating_tx_power_ratio{ap_name="2nd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_operating_tx_power_ratio{ap_name="2nd Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_operating_tx_power_ratio{ap_name="3rd Floor",radio_id="1",vdom="root"} 0.17
        fortigate_wifi_managed_ap_radio_operating_tx_power_ratio{ap_name="3rd Floor",radio_id="2",vdom="root"} 0.2
        fortigate_wifi_managed_ap_radio_operating_tx_power_ratio{ap_name="3rd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_operating_tx_power_ratio{ap_name="3rd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_operating_tx_power_ratio{ap_name="3rd Floor",radio_id="5",vdom="root"} 0
        # HELP fortigate_wifi_managed_ap_radio_rx_bytes_total Total number of received bytes
        # TYPE fortigate_wifi_managed_ap_radio_rx_bytes_total gauge
        fortigate_wifi_managed_ap_radio_rx_bytes_total{ap_name="1st Floor",radio_id="1",vdom="root"} 5.20036172e+09
        fortigate_wifi_managed_ap_radio_rx_bytes_total{ap_name="1st Floor",radio_id="2",vdom="root"} 4.576042303e+09
        fortigate_wifi_managed_ap_radio_rx_bytes_total{ap_name="1st Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_rx_bytes_total{ap_name="1st Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_rx_bytes_total{ap_name="1st Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_rx_bytes_total{ap_name="2nd Floor",radio_id="1",vdom="root"} 3.085351482e+09
        fortigate_wifi_managed_ap_radio_rx_bytes_total{ap_name="2nd Floor",radio_id="2",vdom="root"} 7.356751369e+09
        fortigate_wifi_managed_ap_radio_rx_bytes_total{ap_name="2nd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_rx_bytes_total{ap_name="2nd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_rx_bytes_total{ap_name="2nd Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_rx_bytes_total{ap_name="3rd Floor",radio_id="1",vdom="root"} 3.711981622e+09
        fortigate_wifi_managed_ap_radio_rx_bytes_total{ap_name="3rd Floor",radio_id="2",vdom="root"} 4.881798311e+09
        fortigate_wifi_managed_ap_radio_rx_bytes_total{ap_name="3rd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_rx_bytes_total{ap_name="3rd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_rx_bytes_total{ap_name="3rd Floor",radio_id="5",vdom="root"} 0
        # HELP fortigate_wifi_managed_ap_radio_tx_bytes_total Total number of transferred bytes
        # TYPE fortigate_wifi_managed_ap_radio_tx_bytes_total gauge
        fortigate_wifi_managed_ap_radio_tx_bytes_total{ap_name="1st Floor",radio_id="1",vdom="root"} 4.34100011e+08
        fortigate_wifi_managed_ap_radio_tx_bytes_total{ap_name="1st Floor",radio_id="2",vdom="root"} 2.63637736e+09
        fortigate_wifi_managed_ap_radio_tx_bytes_total{ap_name="1st Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_bytes_total{ap_name="1st Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_bytes_total{ap_name="1st Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_bytes_total{ap_name="2nd Floor",radio_id="1",vdom="root"} 2.11146417e+08
        fortigate_wifi_managed_ap_radio_tx_bytes_total{ap_name="2nd Floor",radio_id="2",vdom="root"} 9.97893287e+08
        fortigate_wifi_managed_ap_radio_tx_bytes_total{ap_name="2nd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_bytes_total{ap_name="2nd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_bytes_total{ap_name="2nd Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_bytes_total{ap_name="3rd Floor",radio_id="1",vdom="root"} 4.9929301e+07
        fortigate_wifi_managed_ap_radio_tx_bytes_total{ap_name="3rd Floor",radio_id="2",vdom="root"} 6.44022911e+08
        fortigate_wifi_managed_ap_radio_tx_bytes_total{ap_name="3rd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_bytes_total{ap_name="3rd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_bytes_total{ap_name="3rd Floor",radio_id="5",vdom="root"} 0
        # HELP fortigate_wifi_managed_ap_radio_tx_discard_ratio Percentage of discarded packets
        # TYPE fortigate_wifi_managed_ap_radio_tx_discard_ratio gauge
        fortigate_wifi_managed_ap_radio_tx_discard_ratio{ap_name="1st Floor",radio_id="1",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_discard_ratio{ap_name="1st Floor",radio_id="2",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_discard_ratio{ap_name="1st Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_discard_ratio{ap_name="1st Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_discard_ratio{ap_name="1st Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_discard_ratio{ap_name="2nd Floor",radio_id="1",vdom="root"} 0.02
        fortigate_wifi_managed_ap_radio_tx_discard_ratio{ap_name="2nd Floor",radio_id="2",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_discard_ratio{ap_name="2nd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_discard_ratio{ap_name="2nd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_discard_ratio{ap_name="2nd Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_discard_ratio{ap_name="3rd Floor",radio_id="1",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_discard_ratio{ap_name="3rd Floor",radio_id="2",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_discard_ratio{ap_name="3rd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_discard_ratio{ap_name="3rd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_discard_ratio{ap_name="3rd Floor",radio_id="5",vdom="root"} 0
        # HELP fortigate_wifi_managed_ap_radio_tx_power_ratio Set Wifi power for the radio in percent
        # TYPE fortigate_wifi_managed_ap_radio_tx_power_ratio gauge
        fortigate_wifi_managed_ap_radio_tx_power_ratio{ap_name="1st Floor",radio_id="1",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_tx_power_ratio{ap_name="1st Floor",radio_id="2",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_tx_power_ratio{ap_name="1st Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_power_ratio{ap_name="1st Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_power_ratio{ap_name="1st Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_power_ratio{ap_name="2nd Floor",radio_id="1",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_tx_power_ratio{ap_name="2nd Floor",radio_id="2",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_tx_power_ratio{ap_name="2nd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_power_ratio{ap_name="2nd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_power_ratio{ap_name="2nd Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_power_ratio{ap_name="3rd Floor",radio_id="1",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_tx_power_ratio{ap_name="3rd Floor",radio_id="2",vdom="root"} 1
        fortigate_wifi_managed_ap_radio_tx_power_ratio{ap_name="3rd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_power_ratio{ap_name="3rd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_power_ratio{ap_name="3rd Floor",radio_id="5",vdom="root"} 0
        # HELP fortigate_wifi_managed_ap_radio_tx_retries_ratio Percentage of retried connection to all connection attempts
        # TYPE fortigate_wifi_managed_ap_radio_tx_retries_ratio gauge
        fortigate_wifi_managed_ap_radio_tx_retries_ratio{ap_name="1st Floor",radio_id="1",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_retries_ratio{ap_name="1st Floor",radio_id="2",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_retries_ratio{ap_name="1st Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_retries_ratio{ap_name="1st Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_retries_ratio{ap_name="1st Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_retries_ratio{ap_name="2nd Floor",radio_id="1",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_retries_ratio{ap_name="2nd Floor",radio_id="2",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_retries_ratio{ap_name="2nd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_retries_ratio{ap_name="2nd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_retries_ratio{ap_name="2nd Floor",radio_id="5",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_retries_ratio{ap_name="3rd Floor",radio_id="1",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_retries_ratio{ap_name="3rd Floor",radio_id="2",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_retries_ratio{ap_name="3rd Floor",radio_id="3",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_retries_ratio{ap_name="3rd Floor",radio_id="4",vdom="root"} 0
        fortigate_wifi_managed_ap_radio_tx_retries_ratio{ap_name="3rd Floor",radio_id="5",vdom="root"} 0
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
