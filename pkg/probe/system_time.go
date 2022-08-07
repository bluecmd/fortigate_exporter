package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
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
