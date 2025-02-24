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

type LogAnaQueueResults struct {
	Connected  float64 `json:"connected"`
	FailedLogs float64 `json:"failed_logs"`
	CachedLogs float64 `json:"cached_logs"`
}

type LogAnaQueue struct {
	Results LogAnaQueueResults `json:"results"`
	VDOM    string             `json:"vdom"`
}

func probeLogAnalyzerQueue(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		logAnaConn = prometheus.NewDesc(
			"fortigate_log_fortianalyzer_queue_connections",
			"Fortianalyzer queue connected state",
			[]string{"vdom"}, nil,
		)
		logAnaLogs = prometheus.NewDesc(
			"fortigate_log_fortianalyzer_queue_logs",
			"State of logs in the queue",
			[]string{"vdom", "state"}, nil,
		)
	)

	var res []LogAnaQueue
	if err := c.Get("api/v2/monitor/log/fortianalyzer-queue", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res {
		m = append(m, prometheus.MustNewConstMetric(logAnaConn, prometheus.GaugeValue, r.Results.Connected, r.VDOM))
		// This is assuming the failed and cached are gauges, which without access to API
		// documentation is too hard to conclusively figure out. If there are data available
		// that proves that the failed logs (or cached logs) are counters instead, we need
		// to either just change the metric type - or split these two up.
		m = append(m, prometheus.MustNewConstMetric(logAnaLogs, prometheus.GaugeValue, r.Results.FailedLogs, r.VDOM, "failed"))
		m = append(m, prometheus.MustNewConstMetric(logAnaLogs, prometheus.GaugeValue, r.Results.CachedLogs, r.VDOM, "cached"))
	}

	return m, true
}
