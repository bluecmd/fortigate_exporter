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

func probeVirtualWANHealthCheck(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		mLink = prometheus.NewDesc(
			"fortigate_virtual_wan_status",
			"Status of the Interface. If the SD-WAN interface is disabled, disable will be returned. If the interface does not participate in the health check, error will be returned.",
			[]string{"vdom", "sla", "interface", "state"}, nil,
		)
		mLatency = prometheus.NewDesc(
			"fortigate_virtual_wan_latency_seconds",
			"Measured latency for this Health check",
			[]string{"vdom", "sla", "interface"}, nil,
		)
		mJitter = prometheus.NewDesc(
			"fortigate_virtual_wan_latency_jitter_seconds",
			"Measured latency jitter for this Health check",
			[]string{"vdom", "sla", "interface"}, nil,
		)
		mPacketLoss = prometheus.NewDesc(
			"fortigate_virtual_wan_packet_loss_ratio",
			"Measured packet loss in percentage for this Health check",
			[]string{"vdom", "sla", "interface"}, nil,
		)
		mPacketSent = prometheus.NewDesc(
			"fortigate_virtual_wan_packet_sent_total",
			"Number of packets sent for this Health check",
			[]string{"vdom", "sla", "interface"}, nil,
		)
		mPacketReceived = prometheus.NewDesc(
			"fortigate_virtual_wan_packet_received_total",
			"Number of packets received for this Health check",
			[]string{"vdom", "sla", "interface"}, nil,
		)
		mSession = prometheus.NewDesc(
			"fortigate_virtual_wan_active_sessions",
			"Active Session count for the health check interface",
			[]string{"vdom", "sla", "interface"}, nil,
		)
		mTXBandwidth = prometheus.NewDesc(
			"fortigate_virtual_wan_bandwidth_tx_byte_per_second",
			"Upload bandwidth of the health check interface",
			[]string{"vdom", "sla", "interface"}, nil,
		)
		mRXBandwidth = prometheus.NewDesc(
			"fortigate_virtual_wan_bandwidth_rx_byte_per_second",
			"Download bandwidth of the health check interface",
			[]string{"vdom", "sla", "interface"}, nil,
		)
		mStateChanged = prometheus.NewDesc(
			"fortigate_virtual_wan_status_change_time_seconds",
			"Unix timestamp describing the time when the last status change has occurred",
			[]string{"vdom", "sla", "interface"}, nil,
		)
	)

	type SLAMember struct {
		Status         string  `json:"status"`
		Latency        float64 `json:"latency"`
		Jitter         float64 `json:"jitter"`
		PacketLoss     float64 `json:"packet_loss"`
		PacketSent     float64 `json:"packet_sent"`
		PacketReceived float64 `json:"packet_received"`
		//todo add slatargetmet
		SLATargetMet []float64 `json:"sla_targets_met"`
		Session      float64   `json:"session"`
		TxBandwidth  float64   `json:"tx_bandwidth"`
		RxBandwidth  float64   `json:"rx_bandwidth"`
		StateChanged float64   `json:"state_changed"`
	}

	type VirtualWanSLA map[string]SLAMember

	type VirtualWanMonitorResponse struct {
		HTTPMethod string                   `json:"http_method"`
		Results    map[string]VirtualWanSLA `json:"results"`
		VDOM       string                   `json:"vdom"`
		Path       string                   `json:"path"`
		Name       string                   `json:"name"`
		Status     string                   `json:"status"`
		Serial     string                   `json:"serial"`
		Version    string                   `json:"version"`
		Build      int64                    `json:"build"`
	}

	var rs []VirtualWanMonitorResponse

	if err := c.Get("api/v2/monitor/virtual-wan/health-check", "vdom=*", &rs); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}
	m := []prometheus.Metric{}
	for _, r := range rs {
		for VirtualWanSLAName, VirtualWanSLA := range r.Results {
			for MemberName, Member := range VirtualWanSLA {
				MemberStatusUp, MemberStatusDown, MemberStatusError, MemberStatusDisable, MemberStatusUnknown := 0.0, 0.0, 0.0, 0.0, 0.0
				switch Member.Status {
				case "up":
					MemberStatusUp = 1.0
				case "down":
					MemberStatusDown = 1.0
				case "error":
					MemberStatusError = 1.0
				case "disable":
					MemberStatusDisable = 1.0
				default:
					MemberStatusUnknown = 1.0
				}

				m = append(m, prometheus.MustNewConstMetric(mLink, prometheus.GaugeValue, MemberStatusUp, r.VDOM, VirtualWanSLAName, MemberName, "up"))
				m = append(m, prometheus.MustNewConstMetric(mLink, prometheus.GaugeValue, MemberStatusDown, r.VDOM, VirtualWanSLAName, MemberName, "down"))
				m = append(m, prometheus.MustNewConstMetric(mLink, prometheus.GaugeValue, MemberStatusError, r.VDOM, VirtualWanSLAName, MemberName, "error"))
				m = append(m, prometheus.MustNewConstMetric(mLink, prometheus.GaugeValue, MemberStatusDisable, r.VDOM, VirtualWanSLAName, MemberName, "disable"))
				m = append(m, prometheus.MustNewConstMetric(mLink, prometheus.GaugeValue, MemberStatusUnknown, r.VDOM, VirtualWanSLAName, MemberName, "unknown"))
				// if no error or unknown status is reported, export the metrics
				if MemberStatusUp == 1 {
					m = append(m, prometheus.MustNewConstMetric(mLatency, prometheus.GaugeValue, Member.Latency/1000, r.VDOM, VirtualWanSLAName, MemberName))
					m = append(m, prometheus.MustNewConstMetric(mJitter, prometheus.GaugeValue, Member.Jitter/1000, r.VDOM, VirtualWanSLAName, MemberName))
					m = append(m, prometheus.MustNewConstMetric(mPacketLoss, prometheus.GaugeValue, Member.PacketLoss/100, r.VDOM, VirtualWanSLAName, MemberName))
					m = append(m, prometheus.MustNewConstMetric(mPacketSent, prometheus.GaugeValue, Member.PacketSent, r.VDOM, VirtualWanSLAName, MemberName))
					m = append(m, prometheus.MustNewConstMetric(mPacketReceived, prometheus.GaugeValue, Member.PacketReceived, r.VDOM, VirtualWanSLAName, MemberName))
					m = append(m, prometheus.MustNewConstMetric(mSession, prometheus.GaugeValue, Member.Session, r.VDOM, VirtualWanSLAName, MemberName))
					m = append(m, prometheus.MustNewConstMetric(mTXBandwidth, prometheus.GaugeValue, Member.TxBandwidth/8, r.VDOM, VirtualWanSLAName, MemberName))
					m = append(m, prometheus.MustNewConstMetric(mRXBandwidth, prometheus.GaugeValue, Member.RxBandwidth/8, r.VDOM, VirtualWanSLAName, MemberName))
					m = append(m, prometheus.MustNewConstMetric(mStateChanged, prometheus.GaugeValue, Member.StateChanged, r.VDOM, VirtualWanSLAName, MemberName))
				}
			}
		}
	}
	return m, true
}
