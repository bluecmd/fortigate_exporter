// Copyright 2025 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
        # HELP fortigate_wifi_client_bandwidth_rx_bps Bandwidth for receiving traffic
        # TYPE fortigate_wifi_client_bandwidth_rx_bps gauge
        fortigate_wifi_client_bandwidth_rx_bps{mac="00:00:00:00:00:00",vdom="root"} 0
        fortigate_wifi_client_bandwidth_rx_bps{mac="00:00:00:AA:00:00",vdom="root"} 0
        # HELP fortigate_wifi_client_bandwidth_tx_bps Bandwidth for transmitting traffic
        # TYPE fortigate_wifi_client_bandwidth_tx_bps gauge
        fortigate_wifi_client_bandwidth_tx_bps{mac="00:00:00:00:00:00",vdom="root"} 0
        fortigate_wifi_client_bandwidth_tx_bps{mac="00:00:00:AA:00:00",vdom="root"} 0
        # HELP fortigate_wifi_client_data_rate_bps Data rate of the client connection
        # TYPE fortigate_wifi_client_data_rate_bps gauge
        fortigate_wifi_client_data_rate_bps{mac="00:00:00:00:00:00",vdom="root"} 1e+06
        fortigate_wifi_client_data_rate_bps{mac="00:00:00:AA:00:00",vdom="root"} 1.3e+08
        # HELP fortigate_wifi_client_info Number of connected access points by status
        # TYPE fortigate_wifi_client_info counter
        fortigate_wifi_client_info{hostname="",mac="00:00:00:AA:00:00",vdom="root",wtp_name="3rd Floor"} 1
        fortigate_wifi_client_info{hostname="wled-WLED",mac="00:00:00:00:00:00",vdom="root",wtp_name="2nd Floor"} 1
        # HELP fortigate_wifi_client_signal_noise_dBm Signal noise on the frequency of the client
        # TYPE fortigate_wifi_client_signal_noise_dBm gauge
        fortigate_wifi_client_signal_noise_dBm{mac="00:00:00:00:00:00",vdom="root"} -95
        fortigate_wifi_client_signal_noise_dBm{mac="00:00:00:AA:00:00",vdom="root"} -95
        # HELP fortigate_wifi_client_signal_strength_dBm Signal strength of the connected client
        # TYPE fortigate_wifi_client_signal_strength_dBm gauge
        fortigate_wifi_client_signal_strength_dBm{mac="00:00:00:00:00:00",vdom="root"} -59
        fortigate_wifi_client_signal_strength_dBm{mac="00:00:00:AA:00:00",vdom="root"} -59
        # HELP fortigate_wifi_client_tx_discard_ratio Percentage of discarded packets
        # TYPE fortigate_wifi_client_tx_discard_ratio gauge
        fortigate_wifi_client_tx_discard_ratio{mac="00:00:00:00:00:00",vdom="root"} 0
        fortigate_wifi_client_tx_discard_ratio{mac="00:00:00:AA:00:00",vdom="root"} 0
        # HELP fortigate_wifi_client_tx_retries_ratio Percentage of retried connection to all connection attempts
        # TYPE fortigate_wifi_client_tx_retries_ratio gauge
        fortigate_wifi_client_tx_retries_ratio{mac="00:00:00:00:00:00",vdom="root"} 0
        fortigate_wifi_client_tx_retries_ratio{mac="00:00:00:AA:00:00",vdom="root"} 0
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
