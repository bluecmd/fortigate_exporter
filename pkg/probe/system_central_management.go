package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type SystemCentralManagementResults struct {
	Mode         string `json:"mode"`
	Status       string `json:"status"`
	Registration string `json:"registration_status"`
}

type SystemCentralManagementStatus struct {
	Results SystemCentralManagementResults `json:"results"`
	VDOM    string                         `json:"vdom"`
}

func probeSystemCentralManagementStatus(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	// New API in FortiOS v7.4.0+
	if meta.VersionMajor < 7 || (meta.VersionMajor == 7 && meta.VersionMinor < 4) {
		log.Printf("Unsupported probe System/CentralManagement/Status - requires minimum FortiOS v7.4.0")
		return []prometheus.Metric{}, false
	}
	var (
		FortimanStat_id = prometheus.NewDesc(
			"fortigate_central_management_connection_status",
			"Fortimanager status",
			[]string{"vdom", "mode", "status"}, nil,
		)
		FortimanReg_id = prometheus.NewDesc(
			"fortigate_central_management_registration_status",
			"Fortimanager registration status",
			[]string{"vdom", "mode", "status"}, nil,
		)
	)

	var res []SystemCentralManagementStatus
	if err := c.Get("api/v2/monitor/system/central-management/status", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res {
		StatusDown, StatusHandshake, StatusUp := 0.0, 0.0, 0.0
		switch r.Results.Status {
		case "down":
			// No management Tunnel
			StatusDown = 1.0
		case "handshake":
			// Management tunnel establishment in progress
			StatusHandshake = 1.0
		case "up":
			// Management tunnel is establised
			StatusUp = 1.0
		}

		RegistrationUnknown, RegistrationInProgress, RegistrationRegistered, RegistrationUnregistered := 0.0, 0.0, 0.0, 0.0
		switch r.Results.Registration {
		case "unknown":
			// FMG does not know about the device
			RegistrationUnknown = 1.0
		case "in_progress":
			// FMG does know the device, but it is not yet fully saved in the list of unregistered devices
			RegistrationInProgress = 1.0
		case "registered":
			// FMG does know the device, and device is authorized
			RegistrationRegistered = 1.0
		case "unregistered ":
			// FMG does know the device, but it is not yet authorized
			RegistrationUnregistered = 1.0
		}

		m = append(m, prometheus.MustNewConstMetric(FortimanStat_id, prometheus.GaugeValue, StatusDown, r.VDOM, r.Results.Mode, "down"))
		m = append(m, prometheus.MustNewConstMetric(FortimanStat_id, prometheus.GaugeValue, StatusHandshake, r.VDOM, r.Results.Mode, "handshake"))
		m = append(m, prometheus.MustNewConstMetric(FortimanStat_id, prometheus.GaugeValue, StatusUp, r.VDOM, r.Results.Mode, "up"))

		m = append(m, prometheus.MustNewConstMetric(FortimanReg_id, prometheus.GaugeValue, RegistrationUnknown, r.VDOM, r.Results.Mode, "unknown"))
		m = append(m, prometheus.MustNewConstMetric(FortimanReg_id, prometheus.GaugeValue, RegistrationInProgress, r.VDOM, r.Results.Mode, "inprogress"))
		m = append(m, prometheus.MustNewConstMetric(FortimanReg_id, prometheus.GaugeValue, RegistrationRegistered, r.VDOM, r.Results.Mode, "registered"))
		m = append(m, prometheus.MustNewConstMetric(FortimanReg_id, prometheus.GaugeValue, RegistrationUnregistered, r.VDOM, r.Results.Mode, "unregistered"))
	}

	return m, true
}
