package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

func probeSystemInterfaceTransceivers(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		mVersion = prometheus.NewDesc(
			"fortigate_interface_transceivers_info",
			"List of transceivers being used by the FortiGate",
			[]string{"name", "type", "vendor", "partnumber", "serialnumber", "description"}, nil,
		)
	)

	type ifResult struct {
		Description    string
		Interface      string
		Type           string
		Vendor         string
		VendorPartNr   string `json:"vendor_part_number"`
		VendorSerialNr string `json:"vendor_serial_number"`
	}
	type ifResponse struct {
		Results []ifResult
	}
	var r ifResponse

	if err := c.Get("api/v2/monitor/system/interface/transceivers", "scope=global", &r); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, result := range r.Results {
		m = append(m, prometheus.MustNewConstMetric(mVersion, prometheus.GaugeValue, 1.0, result.Interface, result.Type, result.Vendor, result.VendorPartNr, result.VendorSerialNr, result.Description))
	}
	return m, true
}
