package probe

import (
	"log"
	"math"
	"strconv"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

func probeFirewallLoadBalance(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	if meta.VersionMajor < 6 || (meta.VersionMajor == 6 && meta.VersionMinor < 4) {
		// not supported version. Before 6.4.0 there is no real_server_id and therefore this will fail
		return nil, true
	}

	var (
		virtualServerInfo = prometheus.NewDesc(
			"fortigate_lb_virtual_server_info",
			"Info metric regarding virtual servers",
			[]string{"vdom", "name", "ip", "port", "type"}, nil,
		)
		realServerInfo = prometheus.NewDesc(
			"fortigate_lb_real_server_info",
			"Info metric regarding real servers",
			[]string{"vdom", "virtual_server", "id", "ip", "port"}, nil,
		)
		realServerMode = prometheus.NewDesc(
			"fortigate_lb_real_server_mode",
			"Mode of this real server: active, standby or disabled",
			[]string{"vdom", "virtual_server", "id", "mode"}, nil,
		)
		realServerStatus = prometheus.NewDesc(
			"fortigate_lb_real_server_status",
			"Status of this real server: up, down or unknown",
			[]string{"vdom", "virtual_server", "id", "state"}, nil,
		)
		realServerSessions = prometheus.NewDesc(
			"fortigate_lb_real_server_active_sessions",
			"Number of sessions active on this real server",
			[]string{"vdom", "virtual_server", "id"}, nil,
		)
		realServerRTT = prometheus.NewDesc(
			"fortigate_lb_real_server_rtt_seconds",
			"Round Trip Time (RTT) for this real server. A RTT of 1 ms or less is reported as 1 ms (0.001 s). A RTT of -1 indicates a parsing error.",
			[]string{"vdom", "virtual_server", "id"}, nil,
		)
		realServerBytesProcessed = prometheus.NewDesc(
			"fortigate_lb_real_server_processed_bytes_total",
			"Number of bytes processed by this real server",
			[]string{"vdom", "virtual_server", "id"}, nil,
		)
	)

	type RealServer struct {
		IP     string `json:"real_server_ip"`
		Port   int    `json:"real_server_port"`
		ID     int    `json:"real_server_id"`
		Mode   string `json:"mode"`
		Status string `json:"status"`
		// MonitorEvents  float64    `json:"monitor_events"`
		ActiveSessions float64 `json:"active_sessions"`
		RTT            string  `json:"RTT"`
		BytesProcessed float64 `json:"bytes_processed"`
	}

	type VirtualServer struct {
		Name        string       `json:"virtual_server_name"`
		IP          string       `json:"virtual_server_ip"`
		Port        int          `json:"virtual_server_port"`
		Type        string       `json:"virtual_server_type"`
		RealServers []RealServer `json:"list"`
	}

	type LoadBalanceResponse struct {
		HTTPMethod string          `json:"http_method"`
		Results    []VirtualServer `json:"results"`
		VDOM       string          `json:"vdom"`
		Path       string          `json:"path"`
		Name       string          `json:"name"`
		Status     string          `json:"status"`
		Serial     string          `json:"serial"`
		Version    string          `json:"version"`
		Build      int64           `json:"build"`
	}

	// Consider implementing pagination to remove this limit of 1000 entries
	var rs []LoadBalanceResponse
	if err := c.Get("api/v2/monitor/firewall/load-balance", "vdom=*&start=0&count=1000", &rs); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}

	for _, r := range rs {
		for _, virtualServer := range r.Results {
			m = append(m, prometheus.MustNewConstMetric(virtualServerInfo, prometheus.GaugeValue, 1, r.VDOM, virtualServer.Name, virtualServer.IP, strconv.Itoa(virtualServer.Port), virtualServer.Type))

			for _, realServer := range virtualServer.RealServers {
				realServerModeActive, realServerModeStandby, realServerModeDisabled := 0.0, 0.0, 0.0
				switch realServer.Mode {
				case "active":
					realServerModeActive = 1.0
					break
				case "standby":
					realServerModeStandby = 1.0
					break
				case "disabled":
					realServerModeDisabled = 1.0
					break
				}

				realServerStatusUp, realServerStatusDown, realServerStatusUnknown := 0.0, 0.0, 0.0
				switch realServer.Status {
				case "up":
					realServerStatusUp = 1.0
					break
				case "down":
					realServerStatusDown = 1.0
					break
				default:
					realServerStatusUnknown = 1.0
				}

				realServerRTTValue := math.NaN()
				if "<1" == realServer.RTT {
					realServerRTTValue = 0.001
				} else if "" == realServer.RTT {
					// NaN
				} else {
					if realServerRTTValueInMs, err := strconv.ParseFloat(realServer.RTT, 64); err != nil {
						log.Printf("Failed to parse RTT value: %v", err)
					} else {
						realServerRTTValue = realServerRTTValueInMs / 1000
					}
				}

				m = append(m, prometheus.MustNewConstMetric(realServerInfo, prometheus.GaugeValue, 1, r.VDOM, virtualServer.Name, strconv.Itoa(realServer.ID), realServer.IP, strconv.Itoa(realServer.Port)))
				m = append(m, prometheus.MustNewConstMetric(realServerMode, prometheus.GaugeValue, realServerModeActive, r.VDOM, virtualServer.Name, strconv.Itoa(realServer.ID), "active"))
				m = append(m, prometheus.MustNewConstMetric(realServerMode, prometheus.GaugeValue, realServerModeStandby, r.VDOM, virtualServer.Name, strconv.Itoa(realServer.ID), "standby"))
				m = append(m, prometheus.MustNewConstMetric(realServerMode, prometheus.GaugeValue, realServerModeDisabled, r.VDOM, virtualServer.Name, strconv.Itoa(realServer.ID), "disabled"))
				m = append(m, prometheus.MustNewConstMetric(realServerStatus, prometheus.GaugeValue, realServerStatusUp, r.VDOM, virtualServer.Name, strconv.Itoa(realServer.ID), "up"))
				m = append(m, prometheus.MustNewConstMetric(realServerStatus, prometheus.GaugeValue, realServerStatusDown, r.VDOM, virtualServer.Name, strconv.Itoa(realServer.ID), "down"))
				m = append(m, prometheus.MustNewConstMetric(realServerStatus, prometheus.GaugeValue, realServerStatusUnknown, r.VDOM, virtualServer.Name, strconv.Itoa(realServer.ID), "unknown"))
				m = append(m, prometheus.MustNewConstMetric(realServerSessions, prometheus.GaugeValue, realServer.ActiveSessions, r.VDOM, virtualServer.Name, strconv.Itoa(realServer.ID)))
				m = append(m, prometheus.MustNewConstMetric(realServerRTT, prometheus.GaugeValue, realServerRTTValue, r.VDOM, virtualServer.Name, strconv.Itoa(realServer.ID)))
				m = append(m, prometheus.MustNewConstMetric(realServerBytesProcessed, prometheus.CounterValue, realServer.BytesProcessed, r.VDOM, virtualServer.Name, strconv.Itoa(realServer.ID)))
			}
		}
	}
	return m, true
}
