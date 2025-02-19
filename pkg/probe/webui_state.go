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

func probeWebUIState(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		mRebootTime = prometheus.NewDesc(
			"fortigate_last_reboot_seconds",
			"Last system reboot epoch time in seconds",
			nil, nil,
		)
		mSnapshotTime = prometheus.NewDesc(
			"fortigate_last_snapshot_seconds",
			"Last snapshot epoch time in seconds",
			nil, nil,
		)
	)

	type WebUIState struct {
		SnapshotUTCTime float64 `json:"snapshot_utc_time"`
		UTCLastReboot   float64 `json:"utc_last_reboot"`
	}

	type webuiState struct {
		Results WebUIState
	}

	var state webuiState

	if err := c.Get("api/v2/monitor/web-ui/state", "", &state); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{
		prometheus.MustNewConstMetric(mRebootTime, prometheus.GaugeValue, state.Results.UTCLastReboot/1000),
		prometheus.MustNewConstMetric(mSnapshotTime, prometheus.GaugeValue, state.Results.SnapshotUTCTime/1000),
	}
	return m, true
}
