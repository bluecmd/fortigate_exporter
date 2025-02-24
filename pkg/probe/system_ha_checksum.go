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

type HAChecksumResults struct {
	IsManageMaster int    `json:"is_manage_master"`
	IsRootMaster   int    `json:"is_root_master"`
	SerialNo       string `json:"serial_no"`
}

type HAChecksum struct {
	Results []HAChecksumResults `json:"results"`
}

func probeSystemHAChecksum(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		IsMaster = prometheus.NewDesc(
			"fortigate_ha_member_has_role",
			"Master/Slave information",
			[]string{"role", "serial"}, nil,
		)
	)

	var res HAChecksum
	if err := c.Get("api/v2/monitor/system/ha-checksums", "scope=global", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, response := range res.Results {
		m = append(m, prometheus.MustNewConstMetric(IsMaster, prometheus.GaugeValue, float64(response.IsManageMaster), "manage_master", response.SerialNo))
		m = append(m, prometheus.MustNewConstMetric(IsMaster, prometheus.GaugeValue, float64(response.IsRootMaster), "root_master", response.SerialNo))
	}

	return m, true
}
