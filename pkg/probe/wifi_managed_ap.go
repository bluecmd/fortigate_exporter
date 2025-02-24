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
	"strconv"

	"github.com/prometheus-community/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

func probeWifiManagedAP(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		managedAPInfo = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_info",
			"Infos about a managed access point",
			[]string{"vdom", "ap_name", "ap_profile", "os_version", "serial"}, nil,
		)
		managedApJoinTime = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_join_time_seconds",
			"Unix time when the managed access point has joined the mesh",
			[]string{"vdom", "ap_name"}, nil,
		)
		managedAPCPUUsage = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_cpu_usage_ratio",
			"CPU usage of the access point",
			[]string{"vdom", "ap_name"}, nil,
		)
		managedAPMemFree = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_memory_free_bytes",
			"Free memory of the managed access point",
			[]string{"vdom", "ap_name"}, nil,
		)
		managedAPMemTotal = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_memory_bytes_total",
			"Total memory of the managed access point",
			[]string{"vdom", "ap_name"}, nil,
		)

		radioClientInfo = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_radio_info",
			"Informations about radios on managed access points",
			[]string{"vdom", "ap_name", "radio_id", "operating_channel"}, nil,
		)
		radioClientCount = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_radio_client_count",
			"Number of clients that are connected using this radio",
			[]string{"vdom", "ap_name", "radio_id"}, nil,
		)
		radioClientOperTxPower = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_radio_operating_tx_power_ratio",
			"Power usage on the operating channel in percent",
			[]string{"vdom", "ap_name", "radio_id"}, nil,
		)
		radioClientChannelUtilization = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_radio_operating_channel_utilization_ratio",
			"Utilization on the operating channel of the radio",
			[]string{"vdom", "ap_name", "radio_id"}, nil,
		)
		radioClientRadioBandwidthRx = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_radio_bandwidth_rx_bps",
			"Bandwidth of this radio for receiving",
			[]string{"vdom", "ap_name", "radio_id"}, nil,
		)
		radioClientRadioBandwidthTx = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_radio_bandwidth_tx_bps",
			"Bandwidth of this radio for transmitting",
			[]string{"vdom", "ap_name", "radio_id"}, nil,
		)
		radioClientRadioRxBytes = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_radio_rx_bytes_total",
			"Total number of received bytes",
			[]string{"vdom", "ap_name", "radio_id"}, nil,
		)
		radioClientRadioTxBytes = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_radio_tx_bytes_total",
			"Total number of transferred bytes",
			[]string{"vdom", "ap_name", "radio_id"}, nil,
		)
		radioInterferingAps = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_radio_interfering_aps",
			"Number of interfering access points",
			[]string{"vdom", "ap_name", "radio_id"}, nil,
		)
		radioTxPower = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_radio_tx_power_ratio",
			"Set Wifi power for the radio in percent",
			[]string{"vdom", "ap_name", "radio_id"}, nil,
		)
		radioTxDiscardPercentage = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_radio_tx_discard_ratio",
			"Percentage of discarded packets",
			[]string{"vdom", "ap_name", "radio_id"}, nil,
		)
		radioTxRetryPercentage = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_radio_tx_retries_ratio",
			"Percentage of retried connection to all connection attempts",
			[]string{"vdom", "ap_name", "radio_id"}, nil,
		)

		interfaceBytesRx = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_interface_rx_bytes_total",
			"total number of bytes received on this interface",
			[]string{"vdom", "ap_name", "interface"}, nil,
		)
		interfaceBytesTx = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_interface_tx_bytes_total",
			"total number of bytes transferred on this interface",
			[]string{"vdom", "ap_name", "interface"}, nil,
		)
		interfacePackagesRx = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_interface_rx_packets_total",
			"total number of packets received on this interface",
			[]string{"vdom", "ap_name", "interface"}, nil,
		)
		interfacePackagesTx = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_interface_tx_packets_total",
			"total number of packets transferred on this interface",
			[]string{"vdom", "ap_name", "interface"}, nil,
		)
		interfaceErrorsRx = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_interface_rx_errors_total",
			"total number of errors received on this interface",
			[]string{"vdom", "ap_name", "interface"}, nil,
		)
		interfaceErrorsTx = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_interface_tx_errors_total",
			"total number of errors transferred on this interface",
			[]string{"vdom", "ap_name", "interface"}, nil,
		)
		interfaceDroppedRx = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_interface_rx_dropped_packets_total",
			"total number of dropped packets received on this interface",
			[]string{"vdom", "ap_name", "interface"}, nil,
		)
		interfaceDroppedTx = prometheus.NewDesc(
			"fortigate_wifi_managed_ap_interface_tx_dropped_packets_total",
			"total number of dropped packets transferred on this interface",
			[]string{"vdom", "ap_name", "interface"}, nil,
		)
	)

	type Radio struct {
		RadioID                   int     `json:"radio_id"`
		ClientCount               float64 `json:"client_count,omitempty"`
		OperChan                  int     `json:"oper_chan,omitempty"`
		OperTxpower               float64 `json:"oper_txpower,omitempty"`
		ChannelUtilizationPercent float64 `json:"channel_utilization_percent,omitempty"`
		BandwidthRx               float64 `json:"bandwidth_rx,omitempty"`
		BandwidthTx               float64 `json:"bandwidth_tx,omitempty"`
		BytesRx                   float64 `json:"bytes_rx,omitempty"`
		BytesTx                   float64 `json:"bytes_tx,omitempty"`
		InterferingAps            float64 `json:"interfering_aps,omitempty"`
		Txpower                   float64 `json:"txpower,omitempty"`
		TxRetriesPercent          float64 `json:"tx_retries_percent,omitempty"`
		TxDiscardPercentage       float64 `json:"tx_discard_percentage,omitempty"`
	}

	type Wired struct {
		Interface  string  `json:"interface"`
		BytesRx    float64 `json:"bytes_rx"`
		BytesTx    float64 `json:"bytes_tx"`
		PacketsRx  float64 `json:"packets_rx"`
		PacketsTx  float64 `json:"packets_tx"`
		ErrorsRx   float64 `json:"errors_rx"`
		ErrorsTx   float64 `json:"errors_tx"`
		DroppedRx  float64 `json:"dropped_rx"`
		DroppedTx  float64 `json:"dropped_tx"`
		Collisions float64 `json:"collisions"`
	}
	type WANStatus struct {
		Interface     string `json:"interface"`
		LinkSpeedMbps int    `json:"link_speed_mbps"`
		CarrierLink   bool   `json:"carrier_link"`
		FullDuplex    bool   `json:"full_duplex"`
	}

	type Results struct {
		Name        string      `json:"name"`
		VDOM        string      `json:"vdom"`
		Serial      string      `json:"serial"`
		APProfile   string      `json:"ap_profile"`
		OSVersion   string      `json:"os_version"`
		JoinTimeRaw float64     `json:"join_time_raw"`
		CPUUsage    float64     `json:"cpu_usage"`
		MemFree     float64     `json:"mem_free"`
		MemTotal    float64     `json:"mem_total"`
		Radio       []Radio     `json:"radio"`
		Wired       []Wired     `json:"wired"`
		WANStatus   []WANStatus `json:"wan_status"`
	}

	type managedAPResponse []struct {
		Results []Results `json:"results"`
	}

	// Consider implementing pagination to remove this limit of 1000 entries
	var response managedAPResponse
	if err := c.Get("api/v2/monitor/wifi/managed_ap", "vdom=*&start=0&count=1000", &response); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	var m []prometheus.Metric
	for _, rs := range response {
		for _, result := range rs.Results {
			m = append(m, prometheus.MustNewConstMetric(managedAPInfo, prometheus.CounterValue, 1, result.VDOM, result.Name, result.APProfile, result.OSVersion, result.Serial))
			m = append(m, prometheus.MustNewConstMetric(managedApJoinTime, prometheus.CounterValue, result.JoinTimeRaw, result.VDOM, result.Name))
			m = append(m, prometheus.MustNewConstMetric(managedAPCPUUsage, prometheus.GaugeValue, result.CPUUsage/100, result.VDOM, result.Name))
			m = append(m, prometheus.MustNewConstMetric(managedAPMemFree, prometheus.GaugeValue, result.MemFree, result.VDOM, result.Name))
			m = append(m, prometheus.MustNewConstMetric(managedAPMemTotal, prometheus.GaugeValue, result.MemTotal, result.VDOM, result.Name))

			for _, radio := range result.Radio {
				radioId := strconv.Itoa(radio.RadioID)
				m = append(m, prometheus.MustNewConstMetric(radioClientInfo, prometheus.CounterValue, 1, result.VDOM, result.Name, radioId, strconv.Itoa(radio.OperChan)))
				m = append(m, prometheus.MustNewConstMetric(radioClientCount, prometheus.GaugeValue, radio.ClientCount, result.VDOM, result.Name, radioId))
				m = append(m, prometheus.MustNewConstMetric(radioClientOperTxPower, prometheus.GaugeValue, radio.OperTxpower/100, result.VDOM, result.Name, radioId))
				m = append(m, prometheus.MustNewConstMetric(radioClientChannelUtilization, prometheus.GaugeValue, radio.ChannelUtilizationPercent/100, result.VDOM, result.Name, radioId))
				m = append(m, prometheus.MustNewConstMetric(radioClientRadioBandwidthRx, prometheus.GaugeValue, radio.BandwidthRx, result.VDOM, result.Name, radioId))
				m = append(m, prometheus.MustNewConstMetric(radioClientRadioBandwidthTx, prometheus.GaugeValue, radio.BandwidthTx, result.VDOM, result.Name, radioId))
				m = append(m, prometheus.MustNewConstMetric(radioClientRadioRxBytes, prometheus.GaugeValue, radio.BytesRx, result.VDOM, result.Name, radioId))
				m = append(m, prometheus.MustNewConstMetric(radioClientRadioTxBytes, prometheus.GaugeValue, radio.BytesTx, result.VDOM, result.Name, radioId))
				m = append(m, prometheus.MustNewConstMetric(radioInterferingAps, prometheus.GaugeValue, radio.InterferingAps, result.VDOM, result.Name, radioId))
				m = append(m, prometheus.MustNewConstMetric(radioTxPower, prometheus.GaugeValue, radio.Txpower/100, result.VDOM, result.Name, radioId))
				m = append(m, prometheus.MustNewConstMetric(radioTxRetryPercentage, prometheus.GaugeValue, radio.TxRetriesPercent/100, result.VDOM, result.Name, radioId))
				m = append(m, prometheus.MustNewConstMetric(radioTxDiscardPercentage, prometheus.GaugeValue, radio.TxDiscardPercentage/100, result.VDOM, result.Name, radioId))
			}

			for _, wired := range result.Wired {
				m = append(m, prometheus.MustNewConstMetric(interfaceBytesRx, prometheus.GaugeValue, wired.BytesRx, result.VDOM, result.Name, wired.Interface))
				m = append(m, prometheus.MustNewConstMetric(interfaceBytesTx, prometheus.GaugeValue, wired.BytesTx, result.VDOM, result.Name, wired.Interface))
				m = append(m, prometheus.MustNewConstMetric(interfacePackagesRx, prometheus.GaugeValue, wired.PacketsRx, result.VDOM, result.Name, wired.Interface))
				m = append(m, prometheus.MustNewConstMetric(interfacePackagesTx, prometheus.GaugeValue, wired.PacketsTx, result.VDOM, result.Name, wired.Interface))
				m = append(m, prometheus.MustNewConstMetric(interfaceErrorsRx, prometheus.GaugeValue, wired.ErrorsRx, result.VDOM, result.Name, wired.Interface))
				m = append(m, prometheus.MustNewConstMetric(interfaceErrorsTx, prometheus.GaugeValue, wired.ErrorsTx, result.VDOM, result.Name, wired.Interface))
				m = append(m, prometheus.MustNewConstMetric(interfaceDroppedRx, prometheus.GaugeValue, wired.DroppedRx, result.VDOM, result.Name, wired.Interface))
				m = append(m, prometheus.MustNewConstMetric(interfaceDroppedTx, prometheus.GaugeValue, wired.DroppedTx, result.VDOM, result.Name, wired.Interface))
			}
		}
	}

	return m, true
}
