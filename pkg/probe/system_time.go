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

func probeSystemTime(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		mTime = prometheus.NewDesc(
			"fortigate_time_seconds",
			"System epoch time in seconds",
			nil, nil,
		)
	)

	type SystemTimeVal struct {
		Time float64 `json:"time"`
	}

	type systemTime struct {
		Results SystemTimeVal
	}

	var stime systemTime

	if err := c.Get("api/v2/monitor/system/time", "vdom=root", &stime); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{
		prometheus.MustNewConstMetric(mTime, prometheus.GaugeValue, stime.Results.Time),
	}
	return m, true
}
