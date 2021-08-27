package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type UserFssoResults struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

type UserFsso struct {
	Results []UserFssoResults `json:"results"`
	VDOM    string            `json:"vdom"`
}

func probeUserFsso(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		FssoUsers = prometheus.NewDesc(
			"fortigate_user_fsso_info",
			"Info on Fsso defined connectors",
			[]string{"vdom", "name", "type", "status"}, nil,
		)
	)

	var res []UserFsso
	if err := c.Get("api/v2/monitor/user/fsso", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, r := range res {
		for _, fssoCon := range r.Results {
			m = append(m, prometheus.MustNewConstMetric(FssoUsers, prometheus.GaugeValue, float64(1), r.VDOM, fssoCon.Name, fssoCon.Type, fssoCon.Status))
		}
	}

	return m, true
}
