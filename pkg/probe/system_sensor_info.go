package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
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
