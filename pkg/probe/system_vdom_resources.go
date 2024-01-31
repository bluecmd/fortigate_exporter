package probe

import (
	"encoding/json"
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type HAVDOMResults struct {
	ID                 float64 `json:"id"`
	CustomMax          int     `json:"custom_max"`
	MinCustomValue     int     `json:"min_custom_value"`
	MaxCustomValue     int     `json:"max_custom_value"`
	Guaranteed         int     `json:"guaranteed"`
	MinGuaranteedValue int     `json:"min_guaranteed_value"`
	MaxGuaranteedValue int     `json:"max_guaranteed_value"`
	GlobalMax          int     `json:"global_max"`
	CurrentUsage       int     `json:"current_usage"`
	UsagePercent       int     `json:"usage_percent"`
}

type HAResults struct {
	Cpu         float64 `json:"cpu"`
	Memory      float64 `json:"memory"`
	SetupRate   float64 `json:"setup_rate"`
	IsDeletable bool    `json:"is_deletable"`
	Interfaces  map[string]HAVDOMResults
}

func (o *HAResults) SystemVDOMResourcesUnmarshalJSON(data []byte) error {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	if _, found := m["cpu"]; found {
		if err := json.Unmarshal(m["cpu"], &o.Cpu); err != nil {
			return err
		}
		delete(m, "cpu")
	}

	if _, found := m["memory"]; found {
		if err := json.Unmarshal(m["memory"], &o.Memory); err != nil {
			return err
		}
		delete(m, "memory")
	}

	if _, found := m["setup_rate"]; found {
		if err := json.Unmarshal(m["setup_rate"], &o.SetupRate); err != nil {
			return err
		}
		delete(m, "setup_rate")
	}

	if _, found := m["is_deletable"]; found {
		if err := json.Unmarshal(m["is_deletable"], &o.IsDeletable); err != nil {
			return err
		}
		delete(m, "is_deletable")
	}

	o.Interfaces = make(map[string]HAVDOMResults)
	for k, v := range m {
		var p HAVDOMResults
		if err := json.Unmarshal(v, &p); err != nil {
			return err
		}
		o.Interfaces[k] = p
	}
	return nil
}

func probeSystemVDOMResourceUsage(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		mUsageSession = prometheus.NewDesc(
			"fortigate_vdom_current_session_usage_percent",
			"Current percent usage of sessions, per VDOM",
			[]string{"vdom", "if_name"}, nil,
		)
	)

	type HAResponse struct {
		HTTPMethod string          `json:"http_method"`
		Results    json.RawMessage `json:"results"`
		VDOM       string          `json:"vdom"`
		Path       string          `json:"path"`
		Name       string          `json:"name"`
		Status     string          `json:"status"`
		Serial     string          `json:"serial"`
		Version    string          `json:"version"`
		Build      int64           `json:"build"`
	}
	var res []HAResponse

	if err := c.Get("api/v2/monitor/system/vdom-resource", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res {
		f := &HAResults{}
		if err := f.SystemVDOMResourcesUnmarshalJSON(r.Results); err != nil {
			log.Printf("Error: %v", err)
			continue
		}
		if result, found := f.Interfaces["session"]; found {
			m = append(m, prometheus.MustNewConstMetric(mUsageSession, prometheus.GaugeValue, float64(result.UsagePercent), r.VDOM, "session"))

		}
	}
	return m, true
}
