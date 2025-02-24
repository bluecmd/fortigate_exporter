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

func probeVPNIPSec(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		status = prometheus.NewDesc(
			"fortigate_ipsec_tunnel_up",
			"Status of IPsec tunnel (0 - Down, 1 - Up)",
			[]string{"vdom", "name", "p2serial", "parent"}, nil,
		)
		transmitted = prometheus.NewDesc(
			"fortigate_ipsec_tunnel_transmit_bytes_total",
			"Total number of bytes transmitted over the IPsec tunnel",
			[]string{"vdom", "name", "p2serial", "parent"}, nil,
		)
		received = prometheus.NewDesc(
			"fortigate_ipsec_tunnel_receive_bytes_total",
			"Total number of bytes received over the IPsec tunnel",
			[]string{"vdom", "name", "p2serial", "parent"}, nil,
		)
	)

	type proxyid struct {
		Name     string  `json:"p2name"`
		P2serial int     `json:"p2serial"`
		Status   string  `json:"status"`
		Incoming float64 `json:"incoming_bytes"`
		Outgoing float64 `json:"outgoing_bytes"`
	}
	type tunnel struct {
		Name    string    `json:"name"`
		Type    string    `json:"type"`
		ProxyID []proxyid `json:"proxyid"`
	}
	type ipsecResult struct {
		Results []tunnel `json:"results"`
		VDOM    string
	}
	var res []ipsecResult
	if err := c.Get("api/v2/monitor/vpn/ipsec", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, v := range res {
		for _, i := range v.Results {
			/*
			  type 'dialup' seems to be client vpn.
			  Not sure exactly what the difference is between probeVPNSsl
			*/
			if i.Type == "dialup" {
				continue
			}
			for _, t := range i.ProxyID {
				s := 0.0
				if t.Status == "up" {
					s = 1.0
				}
				m = append(m, prometheus.MustNewConstMetric(status, prometheus.GaugeValue, s, v.VDOM, t.Name, strconv.Itoa(t.P2serial), i.Name))
				m = append(m, prometheus.MustNewConstMetric(transmitted, prometheus.CounterValue, t.Outgoing, v.VDOM, t.Name, strconv.Itoa(t.P2serial), i.Name))
				m = append(m, prometheus.MustNewConstMetric(received, prometheus.CounterValue, t.Incoming, v.VDOM, t.Name, strconv.Itoa(t.P2serial), i.Name))
			}
		}
	}
	return m, true
}
