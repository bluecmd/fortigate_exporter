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
	"log"

	"github.com/prometheus-community/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

func probeWifiClients(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		clientInfo = prometheus.NewDesc(
			"fortigate_wifi_client_info",
			"Number of connected access points by status",
			[]string{"vdom", "mac", "hostname", "wtp_name"}, nil,
		)
		clientDataRate = prometheus.NewDesc(
			"fortigate_wifi_client_data_rate_bps",
			"Data rate of the client connection",
			[]string{"vdom", "mac"}, nil,
		)
		wtpBandwidthRx = prometheus.NewDesc(
			"fortigate_wifi_client_bandwidth_rx_bps",
			"Bandwidth for receiving traffic",
			[]string{"vdom", "mac"}, nil,
		)
		wtpBandwidthTx = prometheus.NewDesc(
			"fortigate_wifi_client_bandwidth_tx_bps",
			"Bandwidth for transmitting traffic",
			[]string{"vdom", "mac"}, nil,
		)
		signalStrength = prometheus.NewDesc(
			"fortigate_wifi_client_signal_strength_dBm",
			"Signal strength of the connected client",
			[]string{"vdom", "mac"}, nil,
		)
		signalNoise = prometheus.NewDesc(
			"fortigate_wifi_client_signal_noise_dBm",
			"Signal noise on the frequency of the client",
			[]string{"vdom", "mac"}, nil,
		)
		txDiscardPercentage = prometheus.NewDesc(
			"fortigate_wifi_client_tx_discard_ratio",
			"Percentage of discarded packets",
			[]string{"vdom", "mac"}, nil,
		)
		txRetryPercentage = prometheus.NewDesc(
			"fortigate_wifi_client_tx_retries_ratio",
			"Percentage of retried connection to all connection attempts",
			[]string{"vdom", "mac"}, nil,
		)
	)

	type Results struct {
		MAC                 string  `json:"mac"`
		DataRateBps         float64 `json:"data_rate_bps"`
		BandwidthTx         float64 `json:"bandwidth_tx"`
		BandwidthRx         float64 `json:"bandwidth_rx"`
		Signal              float64 `json:"signal"`
		Noise               float64 `json:"noise"`
		TxDiscardPercentage float64 `json:"tx_discard_percentage"`
		TxRetryPercentage   float64 `json:"tx_retry_percentage"`
		Hostname            string  `json:"hostname,omitempty"`
		WtpName             string  `json:"wtp_name"`
	}

	type ApiWifiClientResponse []struct {
		Results []Results `json:"results"`
		VDOM    string    `json:"vdom"`
	}

	// Consider implementing pagination to remove this limit of 1000 entries
	var response ApiWifiClientResponse
	if err := c.Get("api/v2/monitor/wifi/client", "vdom=*&start=0&count=1000", &response); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	var m []prometheus.Metric
	for _, rs := range response {
		for _, result := range rs.Results {
			m = append(m, prometheus.MustNewConstMetric(clientInfo, prometheus.CounterValue, 1, rs.VDOM, result.MAC, result.Hostname, result.WtpName))
			m = append(m, prometheus.MustNewConstMetric(clientDataRate, prometheus.GaugeValue, result.DataRateBps, rs.VDOM, result.MAC))
			m = append(m, prometheus.MustNewConstMetric(wtpBandwidthRx, prometheus.GaugeValue, result.BandwidthRx, rs.VDOM, result.MAC))
			m = append(m, prometheus.MustNewConstMetric(wtpBandwidthTx, prometheus.GaugeValue, result.BandwidthTx, rs.VDOM, result.MAC))
			m = append(m, prometheus.MustNewConstMetric(signalStrength, prometheus.GaugeValue, result.Signal, rs.VDOM, result.MAC))
			m = append(m, prometheus.MustNewConstMetric(signalNoise, prometheus.GaugeValue, result.Noise, rs.VDOM, result.MAC))
			m = append(m, prometheus.MustNewConstMetric(txDiscardPercentage, prometheus.GaugeValue, result.TxDiscardPercentage/100, rs.VDOM, result.MAC))
			m = append(m, prometheus.MustNewConstMetric(txRetryPercentage, prometheus.GaugeValue, result.TxRetryPercentage/100, rs.VDOM, result.MAC))
		}
	}
	return m, true
}
