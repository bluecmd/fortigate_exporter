package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

func probeSystemArpTable(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		mArpEntryAge = prometheus.NewDesc(
			"fortigate_arp_entry_age_seconds",
			"Age of ARP table entry in seconds",
			[]string{"ip", "mac", "interface"}, nil,
		)
	)

	type arpEntry struct {
		IP        string `json:"ip"`
		Age       int    `json:"age"`
		MAC       string `json:"mac"`
		Interface string `json:"interface"`
	}

	type arpResponse struct {
		Results []arpEntry `json:"results"`
		VDOM    string     `json:"vdom"`
	}

	var r arpResponse

	if err := c.Get("api/v2/monitor/network/arp/select", "", &r); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	metrics := []prometheus.Metric{}
	for _, entry := range r.Results {
		metrics = append(metrics, prometheus.MustNewConstMetric(mArpEntryAge, prometheus.GaugeValue, float64(entry.Age), entry.IP, entry.MAC, entry.Interface))
		// Add more metrics if needed
	}

	return metrics, true
}
