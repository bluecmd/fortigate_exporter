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

func probeSystemInterface(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		mLink = prometheus.NewDesc(
			"fortigate_interface_link_up",
			"Whether the link is up or not (not taking into account admin status)",
			[]string{"vdom", "name", "alias", "parent"}, nil,
		)
		mSpeed = prometheus.NewDesc(
			"fortigate_interface_speed_bps",
			"Speed negotiated on the port in bits/s",
			[]string{"vdom", "name", "alias", "parent"}, nil,
		)
		mTxPkts = prometheus.NewDesc(
			"fortigate_interface_transmit_packets_total",
			"Number of packets transmitted on the interface",
			[]string{"vdom", "name", "alias", "parent"}, nil,
		)
		mRxPkts = prometheus.NewDesc(
			"fortigate_interface_receive_packets_total",
			"Number of packets received on the interface",
			[]string{"vdom", "name", "alias", "parent"}, nil,
		)
		mTxB = prometheus.NewDesc(
			"fortigate_interface_transmit_bytes_total",
			"Number of bytes transmitted on the interface",
			[]string{"vdom", "name", "alias", "parent"}, nil,
		)
		mRxB = prometheus.NewDesc(
			"fortigate_interface_receive_bytes_total",
			"Number of bytes received on the interface",
			[]string{"vdom", "name", "alias", "parent"}, nil,
		)
		mTxErr = prometheus.NewDesc(
			"fortigate_interface_transmit_errors_total",
			"Number of transmission errors detected on the interface",
			[]string{"vdom", "name", "alias", "parent"}, nil,
		)
		mRxErr = prometheus.NewDesc(
			"fortigate_interface_receive_errors_total",
			"Number of reception errors detected on the interface",
			[]string{"vdom", "name", "alias", "parent"}, nil,
		)
	)

	type ifResult struct {
		Id        string
		Name      string
		Alias     string
		Link      bool
		Speed     float64
		Duplex    float64
		TxPackets float64 `json:"tx_packets"`
		RxPackets float64 `json:"rx_packets"`
		TxBytes   float64 `json:"tx_bytes"`
		RxBytes   float64 `json:"rx_bytes"`
		TxErrors  float64 `json:"tx_errors"`
		RxErrors  float64 `json:"rx_errors"`
		Interface string
	}
	type ifResponse struct {
		Results map[string]ifResult
		VDOM    string
	}
	var r []ifResponse

	if err := c.Get("api/v2/monitor/system/interface/select", "vdom=*&include_vlan=true&include_aggregate=true", &r); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}
	m := []prometheus.Metric{}
	for _, v := range r {
		for _, ir := range v.Results {
			linkf := 0.0
			if ir.Link {
				linkf = 1.0
			}
			m = append(m, prometheus.MustNewConstMetric(mLink, prometheus.GaugeValue, linkf, v.VDOM, ir.Name, ir.Alias, ir.Interface))
			m = append(m, prometheus.MustNewConstMetric(mSpeed, prometheus.GaugeValue, ir.Speed*1000*1000, v.VDOM, ir.Name, ir.Alias, ir.Interface))
			m = append(m, prometheus.MustNewConstMetric(mTxPkts, prometheus.CounterValue, ir.TxPackets, v.VDOM, ir.Name, ir.Alias, ir.Interface))
			m = append(m, prometheus.MustNewConstMetric(mRxPkts, prometheus.CounterValue, ir.RxPackets, v.VDOM, ir.Name, ir.Alias, ir.Interface))
			m = append(m, prometheus.MustNewConstMetric(mTxB, prometheus.CounterValue, ir.TxBytes, v.VDOM, ir.Name, ir.Alias, ir.Interface))
			m = append(m, prometheus.MustNewConstMetric(mRxB, prometheus.CounterValue, ir.RxBytes, v.VDOM, ir.Name, ir.Alias, ir.Interface))
			m = append(m, prometheus.MustNewConstMetric(mTxErr, prometheus.CounterValue, ir.TxErrors, v.VDOM, ir.Name, ir.Alias, ir.Interface))
			m = append(m, prometheus.MustNewConstMetric(mRxErr, prometheus.CounterValue, ir.RxErrors, v.VDOM, ir.Name, ir.Alias, ir.Interface))
		}
	}
	return m, true
}
