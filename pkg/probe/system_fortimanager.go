package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type SystemFortimanagerResults struct {
	Mode            string `json:"mode"`
	Status_ID       int    `json:"fortimanager_status_id"`
	Registration_ID int    `json:"fortimanager_registration_status_id"`
}

type SystemFortimanagerStatus struct {
	Results SystemFortimanagerResults `json:"results"`
	VDOM    string                    `json:"vdom"`
}

func probeSystemFortimanagerStatus(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		FortimanStat_id = prometheus.NewDesc(
			"fortigate_fortimanager_connection_status",
			"Fortimanager status ID",
			[]string{"vdom", "mode", "status"}, nil,
		)
		FortimanReg_id = prometheus.NewDesc(
			"fortigate_fortimanager_registration_status",
			"Fortimanager registration status ID",
			[]string{"vdom", "mode", "status"}, nil,
		)
	)

	var res []SystemFortimanagerStatus
	if err := c.Get("api/v2/monitor/system/fortimanager/status", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res {
		StatusDown, StatusHandshake, StatusUp := 0.0, 0.0, 0.0
		switch r.Results.Status_ID {
		case 0:
			// No management Tunnel
			StatusDown = 1.0
			break
		case 1:
			// Management tunnel establishment in progress
			StatusHandshake = 1.0
			break
		case 2:
			// Management tunnel is establised
			StatusUp = 1.0
			break
		}

		RegistrationUnknown, RegistrationInProgress, RegistrationRegistered, RegistrationUnregistered := 0.0, 0.0, 0.0, 0.0
		switch r.Results.Registration_ID {
		case 0:
			// FMG does not know about the device
			RegistrationUnknown = 1.0
			break
		case 1:
			// FMG does know the device, but it is not yet fully saved in the list of unregistered devices
			RegistrationInProgress = 1.0
			break
		case 2:
			// FMG does know the device, and device is authorized
			RegistrationRegistered = 1.0
			break
		case 3:
			// FMG does know the device, but it is not yet authorized
			RegistrationUnregistered = 1.0
			break
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
