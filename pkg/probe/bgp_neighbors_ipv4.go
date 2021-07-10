package probe

import (
	"log"
	"strconv"
    

	"github.com/prometheus/client_golang/prometheus"
	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/bluecmd/fortigate_exporter/internal/version"
)

func probeBGPNeighborsIPv4(c http.FortiHTTP) ([]prometheus.Metric, bool) {

	var (
		bpgNeighbor = prometheus.NewDesc(
			"fortigate_bgp_neighbors_ipv4",
			"State of configured IPv4 BGP peers",
			[]string{"vdom", "remote_as", "state", "admin_status", "local_ip", "neighbor_ip", "type"}, nil,
		)
	)

	type BgpNeighbor struct {
		NeighborIP        string       `json:"neighbor_ip"`
		LocalIP           string       `json:"local_ip"`
		RemoteAS          int          `json:"remote_as"`
		AdminStatus       bool         `json:"admin_status"`
		State             string       `json:"state"`
		Type              string       `json:"type"`
	}

	type BpgNeighborResponse struct {
		HTTPMethod string                 `json:"http_method"`
		Results    []BgpNeighbor          `json:"results"`
		VDOM       string                 `json:"vdom"`
		Path       string                 `json:"path"`
		Name       string                 `json:"name"`
		Status     string                 `json:"status"`
		Serial     string                 `json:"serial"`
		Version    string                 `json:"version"`
		Build      int64                  `json:"build"`
	}

	var rs []BpgNeighborResponse

	if err := c.Get("api/v2/monitor/router/bgp/neighbors", "vdom=*", &rs); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}
	
	m := []prometheus.Metric{}

	for _, r := range rs {
		major, _, ok := version.ParseVersion(r.Version)
		if !ok || major < 7 {
			// not supported version. Before 7.0.0 the requested endpoint doesn't exist
			return nil, false
		}
		for _, peer := range r.Results {
			m = append(m, prometheus.MustNewConstMetric(bpgNeighbor, prometheus.GaugeValue, 1, r.VDOM, strconv.Itoa(peer.RemoteAS), peer.State, strconv.FormatBool(peer.AdminStatus), peer.LocalIP, peer.NeighborIP, peer.Type))
		}
	}

	return m, true
}
