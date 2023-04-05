package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
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
