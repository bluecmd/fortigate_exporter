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

type UserFssoResults struct {
	Name   string `json:"name"`
	ID     int    `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

type UserFsso struct {
	Results []UserFssoResults `json:"results"`
	VDOM    string            `json:"vdom"`
}

func probeUserFsso(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		FssoUsers = prometheus.NewDesc(
			"fortigate_user_fsso_info",
			"Info on Fsso defined connectors",
			[]string{"vdom", "name", "id", "type", "status"}, nil,
		)
	)

	var res []UserFsso
	if err := c.Get("api/v2/monitor/user/fsso", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res {
		for _, fssoCon := range r.Results {
			if fssoCon.Type == "fsso" {
				m = append(m, prometheus.MustNewConstMetric(FssoUsers, prometheus.GaugeValue, float64(1), r.VDOM, fssoCon.Name, "", fssoCon.Type, fssoCon.Status))
			} else {
				m = append(m, prometheus.MustNewConstMetric(FssoUsers, prometheus.GaugeValue, float64(1), r.VDOM, "", strconv.Itoa(fssoCon.ID), fssoCon.Type, fssoCon.Status))
			}
		}
	}

	return m, true
}
