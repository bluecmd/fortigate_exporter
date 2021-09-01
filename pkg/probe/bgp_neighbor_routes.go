package probe

import (
	"fmt"
	"log"

	"github.com/bluecmd/fortigate_exporter/internal/config"
	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type BGPPath struct {
	LearnedFrom string `json:"learned_from"`
	IsBest      bool   `json:"is_best"`
}

type BGPPaths struct {
	Results []BGPPath `json:"results"`
	VDOM    string    `json:"vdom"`
	Version string    `json:"version"`
}

type PathCount struct {
	Source string
	VDOM   string
}

func probeBGPNeighborPathsIPv4(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	savedConfig := config.GetConfig()
	MaxBGPPaths := savedConfig.MaxBGPPaths

	if MaxBGPPaths == 0 {
		return nil, true
	}

	if meta.VersionMajor < 7 {
		// not supported version. Before 7.0.0 the requested endpoint doesn't exist
		return nil, true
	}
	var (
		BGPNeighborPaths = prometheus.NewDesc(
			"fortigate_bgp_neighbor_ipv4_paths",
			"Count of paths received from an BGP neighbor",
			[]string{"vdom", "neighbor_ip"}, nil,
		)
		BGPNeighborBestPaths = prometheus.NewDesc(
			"fortigate_bgp_neighbor_ipv4_best_paths",
			"Count of best paths for an BGP neighbor",
			[]string{"vdom", "neighbor_ip"}, nil,
		)
	)

	var rs []BGPPaths

	if err := c.Get("api/v2/monitor/router/bgp/paths", fmt.Sprintf("vdom=*&count=%d", MaxBGPPaths), &rs); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}

	srMap := make(map[PathCount]int)
	sr2Map := make(map[PathCount]int)
	for _, r := range rs {

		if len(r.Results) > MaxBGPPaths {
			log.Printf("Error: Received more BGP Paths than maximum (%d > %d) allowed, ignoring metric ...", len(r.Results), MaxBGPPaths)
			return nil, false
		}
		for _, route := range r.Results {
			sr := PathCount{
				Source: route.LearnedFrom,
				VDOM:   r.VDOM,
			}
			srMap[sr] += 1
			if route.IsBest {
				sr2 := PathCount{
					Source: route.LearnedFrom,
					VDOM:   r.VDOM,
				}
				sr2Map[sr2] += 1
			}
		}
	}

	for neighbor, count := range srMap {
		m = append(m, prometheus.MustNewConstMetric(BGPNeighborPaths, prometheus.GaugeValue, float64(count), neighbor.VDOM, neighbor.Source))
	}
	for neighbor, count := range sr2Map {
		m = append(m, prometheus.MustNewConstMetric(BGPNeighborBestPaths, prometheus.GaugeValue, float64(count), neighbor.VDOM, neighbor.Source))
	}

	return m, true
}

func probeBGPNeighborPathsIPv6(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	savedConfig := config.GetConfig()
	MaxBGPPaths := savedConfig.MaxBGPPaths

	if MaxBGPPaths == 0 {
		return nil, true
	}

	if meta.VersionMajor < 7 {
		// not supported version. Before 7.0.0 the requested endpoint doesn't exist
		return nil, true
	}
	var (
		BGPNeighborPaths = prometheus.NewDesc(
			"fortigate_bgp_neighbor_ipv6_paths",
			"Count of paths received from an BGP neighbor",
			[]string{"vdom", "neighbor_ip"}, nil,
		)
		BGPNeighborBestPaths = prometheus.NewDesc(
			"fortigate_bgp_neighbor_ipv6_best_paths",
			"Count of best paths for an BGP neighbor",
			[]string{"vdom", "neighbor_ip"}, nil,
		)
	)

	var rs []BGPPaths

	if err := c.Get("api/v2/monitor/router/bgp/paths6", fmt.Sprintf("vdom=*&count=%d", MaxBGPPaths), &rs); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}

	srMap := make(map[PathCount]int)
	sr2Map := make(map[PathCount]int)
	for _, r := range rs {

		if len(r.Results) > MaxBGPPaths {
			log.Printf("Error: Received more BGP Paths than maximum (%d > %d) allowed, ignoring metric ...", len(r.Results), MaxBGPPaths)
			return nil, false
		}
		for _, route := range r.Results {
			sr := PathCount{
				Source: route.LearnedFrom,
				VDOM:   r.VDOM,
			}
			srMap[sr] += 1
			if route.IsBest {
				sr2 := PathCount{
					Source: route.LearnedFrom,
					VDOM:   r.VDOM,
				}
				sr2Map[sr2] += 1
			}
		}
	}

	for neighbor, count := range srMap {
		m = append(m, prometheus.MustNewConstMetric(BGPNeighborPaths, prometheus.GaugeValue, float64(count), neighbor.VDOM, neighbor.Source))
	}
	for neighbor, count := range sr2Map {
		m = append(m, prometheus.MustNewConstMetric(BGPNeighborBestPaths, prometheus.GaugeValue, float64(count), neighbor.VDOM, neighbor.Source))
	}

	return m, true
}
