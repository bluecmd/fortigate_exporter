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

type IpPool struct {
	Name      string  `json:"name"`
	IPTotal   int     `json:"natip_total"`
	IPInUse   int     `json:"natip_in_use"`
	Clients   int     `json:"clients"`
	Available float64 `json:"available"`
	Used      int     `json:"used"`
	Total     int     `json:"total"`
	PbaPerIp  int     `json:"pba_per_ip"`
}

type IpPoolResponse struct {
	Results map[string]IpPool `json:"results"`
	VDOM    string            `json:"vdom"`
	Version string            `json:"version"`
}

func probeFirewallIpPool(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		mAvailable = prometheus.NewDesc(
			"fortigate_ippool_available_ratio",
			"Percentage available in ippool (0 - 1.0)",
			[]string{"vdom", "name"}, nil,
		)
	)
	var (
		mIpUsed = prometheus.NewDesc(
			"fortigate_ippool_used_ips",
			"Ip addresses in use in ippool",
			[]string{"vdom", "name"}, nil,
		)
	)
	var (
		mIpTotal = prometheus.NewDesc(
			"fortigate_ippool_total_ips",
			"Ip addresses total in ippool",
			[]string{"vdom", "name"}, nil,
		)
	)
	var (
		mClients = prometheus.NewDesc(
			"fortigate_ippool_clients",
			"Amount of clients using ippool",
			[]string{"vdom", "name"}, nil,
		)
	)
	var (
		mUsed = prometheus.NewDesc(
			"fortigate_ippool_used_items",
			"Amount of items used in ippool",
			[]string{"vdom", "name"}, nil,
		)
	)
	var (
		mTotal = prometheus.NewDesc(
			"fortigate_ippool_total_items",
			"Amount of items total in ippool",
			[]string{"vdom", "name"}, nil,
		)
	)

	var (
		mPbaPerIp = prometheus.NewDesc(
			"fortigate_ippool_pba_per_ip",
			"Amount of available port block allocations per ip",
			[]string{"vdom", "name"}, nil,
		)
	)

	var rs []IpPoolResponse

	if err := c.Get("api/v2/monitor/firewall/ippool", "vdom=*", &rs); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}

	for _, r := range rs {
		for _, ippool := range r.Results {
			m = append(m, prometheus.MustNewConstMetric(mAvailable, prometheus.GaugeValue, ippool.Available/100, r.VDOM, ippool.Name))
			m = append(m, prometheus.MustNewConstMetric(mIpUsed, prometheus.GaugeValue, float64(ippool.IPInUse), r.VDOM, ippool.Name))
			m = append(m, prometheus.MustNewConstMetric(mIpTotal, prometheus.GaugeValue, float64(ippool.IPTotal), r.VDOM, ippool.Name))
			m = append(m, prometheus.MustNewConstMetric(mClients, prometheus.GaugeValue, float64(ippool.Clients), r.VDOM, ippool.Name))
			m = append(m, prometheus.MustNewConstMetric(mUsed, prometheus.GaugeValue, float64(ippool.Used), r.VDOM, ippool.Name))
			m = append(m, prometheus.MustNewConstMetric(mTotal, prometheus.GaugeValue, float64(ippool.Total), r.VDOM, ippool.Name))
			m = append(m, prometheus.MustNewConstMetric(mPbaPerIp, prometheus.GaugeValue, float64(ippool.PbaPerIp), r.VDOM, ippool.Name))
		}
	}

	return m, true
}
