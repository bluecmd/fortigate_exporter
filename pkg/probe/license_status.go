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

func probeLicenseStatus(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		vdomUsed = prometheus.NewDesc(
			"fortigate_license_vdom_usage",
			"The amount of VDOM licenses currently used",
			[]string{}, nil,
		)
		vdomMax = prometheus.NewDesc(
			"fortigate_license_vdom_max",
			"The total amount of VDOM licenses available",
			[]string{}, nil,
		)
	)

	type LicenseStatus struct {
		VDOM struct {
			Type       string  `json:"type"`
			CanUpgrade bool    `json:"can_upgrade"`
			Used       float64 `json:"used"`
			Max        float64 `json:"max"`
		} `json:"vdom"`
	}

	type LicenseResponse struct {
		Results LicenseStatus `json:"results"`
	}
	var r LicenseResponse

	if err := c.Get("api/v2/monitor/license/status/select", "", &r); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{
		prometheus.MustNewConstMetric(vdomUsed, prometheus.GaugeValue, float64(r.Results.VDOM.Used)),
		prometheus.MustNewConstMetric(vdomMax, prometheus.GaugeValue, float64(r.Results.VDOM.Max)),
	}

	return m, true
}
