package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

func probeSystemAvailableCertificates(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		certificateInfo = prometheus.NewDesc(
			"fortigate_certificate_info",
			"Info metric containing meta information about the certificate",
			[]string{"name", "source", "scope", "vdom", "status", "type"}, nil,
		)
		certificateValidFrom = prometheus.NewDesc(
			"fortigate_certificate_valid_from_seconds",
			"Unix timestamp from which this certificate is valid",
			[]string{"name", "source", "scope", "vdom"}, nil,
		)
		certificateValidTo = prometheus.NewDesc(
			"fortigate_certificate_valid_to_seconds",
			"Unix timestamp till which this certificate is valid",
			[]string{"name", "source", "scope", "vdom"}, nil,
		)
		certificateCMDBReferences = prometheus.NewDesc(
			"fortigate_certificate_cmdb_references",
			"Number of times the certificate is referenced",
			[]string{"name", "source", "scope", "vdom"}, nil,
		)
	)

	type Results struct {
		Name      string  `json:"name"`
		Source    string  `json:"source"`
		Type      string  `json:"type"`
		Status    string  `json:"status"`
		ValidFrom float64 `json:"valid_from"`
		ValidTo   float64 `json:"valid_to"`
		QRef      float64 `json:"q_ref"`
	}

	type Response struct {
		Results []Results `json:"results"`
		VDOM    string    `json:"vdom"`
		Status  string    `json:"status"`
		Scope   string
	}

	var globalResponse Response
	if err := c.Get("api/v2/monitor/system/available-certificates", "scope=global", &globalResponse); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}
	globalResponse.Scope = "global"

	var vdomResponses []Response

	if err := c.Get("api/v2/monitor/system/available-certificates", "vdom=*", &vdomResponses); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}
	for i := range vdomResponses {
		vdomResponses[i].Scope = "vdom"
	}

	// combine responses
	combinedResponses := make([]Response, 0)
	combinedResponses = append(combinedResponses, vdomResponses...)
	combinedResponses = append(combinedResponses, globalResponse)

	m := []prometheus.Metric{}

	for _, response := range combinedResponses {
		for _, result := range response.Results {
			m = append(m, prometheus.MustNewConstMetric(certificateInfo, prometheus.GaugeValue, 1, result.Name, result.Source, response.Scope, response.VDOM, result.Status, result.Type))
			m = append(m, prometheus.MustNewConstMetric(certificateValidFrom, prometheus.GaugeValue, result.ValidFrom, result.Name, result.Source, response.Scope, response.VDOM))
			m = append(m, prometheus.MustNewConstMetric(certificateValidTo, prometheus.GaugeValue, result.ValidTo, result.Name, result.Source, response.Scope, response.VDOM))
			m = append(m, prometheus.MustNewConstMetric(certificateCMDBReferences, prometheus.GaugeValue, result.QRef, result.Name, result.Source, response.Scope, response.VDOM))
		}
	}
	return m, true
}
