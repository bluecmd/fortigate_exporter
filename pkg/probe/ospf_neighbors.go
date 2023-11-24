package probe

import (
	"log"
	"strconv"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type OSPFNeighbor struct {
	NeighborIP string `json:"neighbor_ip"`
	RouterID   string `json:"router_id"`
	State      string `json:"state"`
	Priority   int    `json:"priority"`
}

type OSPFNeighborResponse struct {
	Results []OSPFNeighbor `json:"results"`
	VDOM    string         `json:"vdom"`
	Version string         `json:"version"`
}

func probeOSPFNeighbors(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	if meta.VersionMajor < 7 {
		// not supported version. Before 7.0.0 the requested endpoint doesn't exist
		return nil, true
	}
	var (
		mOSPFNeighbor = prometheus.NewDesc(
			"fortigate_ospf_neighbor_info",
			"List all discovered OSPF neighbors, return state as value (1 - Down, 2 - Attempt, 3 - Init, 4 - Two way, 5 - Exchange start, 6 - Exchange, 7 - Loading, 8 - Full)",
			[]string{"vdom", "state", "priority", "router_id", "neighbor_ip"}, nil,
		)
	)

	var rs []OSPFNeighborResponse

	if err := c.Get("api/v2/monitor/router/ospf/neighbors", "vdom=*", &rs); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}

	for _, r := range rs {
		for _, peer := range r.Results {
			m = append(m, prometheus.MustNewConstMetric(mOSPFNeighbor, prometheus.GaugeValue, ospfStateToNumber(peer.State), r.VDOM, peer.State, strconv.Itoa(peer.Priority), peer.RouterID, peer.NeighborIP))
		}
	}

	return m, true
}

func ospfStateToNumber(ospfState string) float64 {
	switch ospfState {
	case "Attempt":
		return 1
	case "Init":
		return 2
	case "Two way":
		return 3
	case "Exchange start":
		return 4
	case "Exchange":
		return 5
	case "Loading":
		return 6
	case "Full":
		return 7
	default: // Down
		return 0
	}
}
