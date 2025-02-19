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

type LogResults struct {
	UsedBytes  float64 `json:"used_bytes"`
	TotalBytes float64 `json:"total_bytes"`
}

type Log struct {
	Results LogResults `json:"results"`
	VDOM    string     `json:"vdom"`
}

func probeLogCurrentDiskUsage(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		logUsed = prometheus.NewDesc(
			"fortigate_log_disk_used_bytes",
			"Disk used bytes for log",
			[]string{"vdom"}, nil,
		)
		logTotal = prometheus.NewDesc(
			"fortigate_log_disk_total_bytes",
			"Disk total bytes for log",
			[]string{"vdom"}, nil,
		)
	)

	var res []Log
	if err := c.Get("api/v2/monitor/log/current-disk-usage", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res {
		m = append(m, prometheus.MustNewConstMetric(logUsed, prometheus.GaugeValue, r.Results.UsedBytes, r.VDOM))
		m = append(m, prometheus.MustNewConstMetric(logTotal, prometheus.GaugeValue, r.Results.TotalBytes, r.VDOM))
	}

	return m, true
}
