package probe

import (
	"log"
	"strconv"

	"github.com/bluecmd/fortigate_exporter/internal/config"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type BGPNeighbor struct {
	NeighborIP  string `json:"neighbor_ip"`
	LocalIP     string `json:"local_ip"`
	RemoteAS    int    `json:"remote_as"`
	AdminStatus bool   `json:"admin_status"`
	State       string `json:"state"`
}

type BGPNeighborResponse struct {
	Results []BGPNeighbor `json:"results"`
	VDOM    string        `json:"vdom"`
	Version string        `json:"version"`
}

func probeBGPNeighborsIPv4(c http.FortiHTTP, meta *TargetMetadata, _ config.FortiExporterConfig) ([]prometheus.Metric, bool) {
	if meta.VersionMajor < 7 {
		// not supported version. Before 7.0.0 the requested endpoint doesn't exist
		return nil, true
	}
	var (
		mBGPNeighbor = prometheus.NewDesc(
			"fortigate_bgp_neighbor_ipv4_info",
			"Configured bgp neighbor over ipv4",
			[]string{"vdom", "remote_as", "state", "admin_status", "local_ip", "neighbor_ip"}, nil,
		)
	)

	var rs []BGPNeighborResponse

	if err := c.Get("api/v2/monitor/router/bgp/neighbors", "vdom=*", &rs); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}

	for _, r := range rs {
		for _, peer := range r.Results {
			m = append(m, prometheus.MustNewConstMetric(mBGPNeighbor, prometheus.GaugeValue, 1, r.VDOM, strconv.Itoa(peer.RemoteAS), peer.State, strconv.FormatBool(peer.AdminStatus), peer.LocalIP, peer.NeighborIP))
		}
	}

	return m, true
}

func probeBGPNeighborsIPv6(c http.FortiHTTP, meta *TargetMetadata, _ config.FortiExporterConfig) ([]prometheus.Metric, bool) {
	if meta.VersionMajor < 7 {
		// not supported version. Before 7.0.0 the requested endpoint doesn't exist
		return nil, true
	}

	var (
		mBGPNeighbor = prometheus.NewDesc(
			"fortigate_bgp_neighbor_ipv6_info",
			"Configured bgp neighbor over ipv6",
			[]string{"vdom", "remote_as", "state", "admin_status", "local_ip", "neighbor_ip"}, nil,
		)
	)

	var rs []BGPNeighborResponse

	if err := c.Get("api/v2/monitor/router/bgp/neighbors6", "vdom=*", &rs); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}

	for _, r := range rs {
		for _, peer := range r.Results {
			m = append(m, prometheus.MustNewConstMetric(mBGPNeighbor, prometheus.GaugeValue, 1, r.VDOM, strconv.Itoa(peer.RemoteAS), peer.State, strconv.FormatBool(peer.AdminStatus), peer.LocalIP, peer.NeighborIP))
		}
	}

	return m, true
}
