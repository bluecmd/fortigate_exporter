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

type LogAnaResults struct {
	Registration string  `json:"registration"`
	Connection   string  `json:"connection"`
	Received     float64 `json:"received"`
}

type LogAna struct {
	Results LogAnaResults `json:"results"`
	VDOM    string        `json:"vdom"`
}

func probeLogAnalyzer(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		logAnaInfo = prometheus.NewDesc(
			"fortigate_log_fortianalyzer_registration_info",
			"Fortianalyzer state info",
			[]string{"vdom", "registration", "connection"}, nil,
		)
		logAnaRcv = prometheus.NewDesc(
			"fortigate_log_fortianalyzer_logs_received",
			"Received logs in fortianalyzer",
			[]string{"vdom"}, nil,
		)
	)

	var res []LogAna
	if err := c.Get("api/v2/monitor/log/fortianalyzer", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res {
		m = append(m, prometheus.MustNewConstMetric(logAnaInfo, prometheus.GaugeValue, float64(1), r.VDOM, r.Results.Registration, r.Results.Connection))
		m = append(m, prometheus.MustNewConstMetric(logAnaRcv, prometheus.GaugeValue, r.Results.Received, r.VDOM))
	}

	return m, true
}
