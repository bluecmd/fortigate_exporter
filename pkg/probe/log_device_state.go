package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type LogDevSttResFlags struct {
	IsAvailable bool `json:"is_available"`
	IsEnabled   bool `json:"is_enabled"`
}

type LogDevSttRes struct {
	Memory        LogDevSttResFlags `json:"memory"`
	Disk          LogDevSttResFlags `json:"disk"`
	FortiAnalyzer LogDevSttResFlags `json:"fortianalyzer"`
	FortiCloud    LogDevSttResFlags `json:"forticloud"`
}

type LogDevStt struct {
	Results LogDevSttRes `json:"results"`
	VDOM    string       `json:"vdom"`
	Version string       `json:"version"`
}

func btof(b bool) float64 {
	if b {
		return float64(1)
	}
	return float64(0)
}

func probeLogDeviceState(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		logMemAva = prometheus.NewDesc(
			"fortigate_log_memory_available",
			"Is memory available for log?",
			[]string{"vdom"}, nil,
		)
		logMemEna = prometheus.NewDesc(
			"fortigate_log_memory_enabled",
			"Is memory enabled for log?",
			[]string{"vdom"}, nil,
		)
		logDskAva = prometheus.NewDesc(
			"fortigate_log_disk_available",
			"Is disk available for log?",
			[]string{"vdom"}, nil,
		)
		logDskEna = prometheus.NewDesc(
			"fortigate_log_disk_enabled",
			"Is disk enabled for log?",
			[]string{"vdom"}, nil,
		)
		logAnaAva = prometheus.NewDesc(
			"fortigate_log_fortianalyzer_available",
			"Is fortianalyzer available for log?",
			[]string{"vdom"}, nil,
		)
		logAnaEna = prometheus.NewDesc(
			"fortigate_log_fortianalyzer_enabled",
			"Is fortianalyzer enabled for log?",
			[]string{"vdom"}, nil,
		)
		logCldAva = prometheus.NewDesc(
			"fortigate_log_forticloud_available",
			"Is forticloud available for log?",
			[]string{"vdom"}, nil,
		)
		logCldEna = prometheus.NewDesc(
			"fortigate_log_forticloud_enabled",
			"Is forticloud enabled for log?",
			[]string{"vdom"}, nil,
		)
	)

	var res []LogDevStt
	if err := c.Get("api/v2/monitor/log/device/state", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res {
		m = append(m, prometheus.MustNewConstMetric(logMemAva, prometheus.GaugeValue, btof(r.Results.Memory.IsAvailable), r.VDOM))
		m = append(m, prometheus.MustNewConstMetric(logMemEna, prometheus.GaugeValue, btof(r.Results.Memory.IsEnabled), r.VDOM))
		m = append(m, prometheus.MustNewConstMetric(logDskAva, prometheus.GaugeValue, btof(r.Results.Disk.IsAvailable), r.VDOM))
		m = append(m, prometheus.MustNewConstMetric(logDskEna, prometheus.GaugeValue, btof(r.Results.Disk.IsEnabled), r.VDOM))
		m = append(m, prometheus.MustNewConstMetric(logAnaAva, prometheus.GaugeValue, btof(r.Results.FortiAnalyzer.IsAvailable), r.VDOM))
		m = append(m, prometheus.MustNewConstMetric(logAnaEna, prometheus.GaugeValue, btof(r.Results.FortiAnalyzer.IsEnabled), r.VDOM))
		m = append(m, prometheus.MustNewConstMetric(logCldAva, prometheus.GaugeValue, btof(r.Results.FortiCloud.IsAvailable), r.VDOM))
		m = append(m, prometheus.MustNewConstMetric(logCldEna, prometheus.GaugeValue, btof(r.Results.FortiCloud.IsEnabled), r.VDOM))
	}

	return m, true
}
