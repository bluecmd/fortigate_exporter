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

type SystemSensorInfoResults struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

type SystemSensorInfo struct {
	Results []SystemSensorInfoResults `json:"results"`
}

func probeSystemSensorInfo(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		sensorTemperature = prometheus.NewDesc(
			"fortigate_sensor_temperature_celsius",
			"Sensor temperature in degree celsius",
			[]string{"name"}, nil,
		)
		sensorFan = prometheus.NewDesc(
			"fortigate_sensor_fan_rpm",
			"Sensor fan rotation speed in RPM",
			[]string{"name"}, nil,
		)
		sensorVoltage = prometheus.NewDesc(
			"fortigate_sensor_voltage_volts",
			"Sensor voltage in volts",
			[]string{"name"}, nil,
		)
	)

	var res SystemSensorInfo
	if err := c.Get("api/v2/monitor/system/sensor-info", "vdom=root", &res); err != nil {
		log.Printf("Warning: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res.Results {
		switch r.Type {
		case "temperature":
			m = append(m, prometheus.MustNewConstMetric(sensorTemperature, prometheus.GaugeValue, r.Value, r.Name))
		case "fan":
			m = append(m, prometheus.MustNewConstMetric(sensorFan, prometheus.GaugeValue, r.Value, r.Name))
		case "voltage":
			m = append(m, prometheus.MustNewConstMetric(sensorVoltage, prometheus.GaugeValue, r.Value, r.Name))
		}
	}

	return m, true
}
