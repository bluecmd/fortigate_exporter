package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

func probeSystemLinkMonitor(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		linkStatus = prometheus.NewDesc(
			"fortigate_link_status",
			"Signals the status of the link. 1 means that this state is present in every other case the value is 0",
			[]string{"vdom", "monitor", "link", "state"}, nil,
		)
		linkLatency = prometheus.NewDesc(
			"fortigate_link_latency_seconds",
			"Average latency of this link based on the last 30 probes in seconds",
			[]string{"vdom", "monitor", "link"}, nil,
		)
		linkJitter = prometheus.NewDesc(
			"fortigate_link_latency_jitter_seconds",
			"Average of the latency jitter  on this link based on the last 30 probes in seconds",
			[]string{"vdom", "monitor", "link"}, nil,
		)
		linkPacketLoss = prometheus.NewDesc(
			"fortigate_link_packet_loss_ratio",
			"Percentage of packets lost relative to  all sent based on the last 30 probes",
			[]string{"vdom", "monitor", "link"}, nil,
		)
		linkPacketSent = prometheus.NewDesc(
			"fortigate_link_packet_sent_total",
			"Number of packets sent on this link",
			[]string{"vdom", "monitor", "link"}, nil,
		)
		linkPacketReceived = prometheus.NewDesc(
			"fortigate_link_packet_received_total",
			"Number of packets received on this link",
			[]string{"vdom", "monitor", "link"}, nil,
		)
		linkSessions = prometheus.NewDesc(
			"fortigate_link_active_sessions",
			"Number of sessions active on this link",
			[]string{"vdom", "monitor", "link"}, nil,
		)
		linkBandwidthTx = prometheus.NewDesc(
			"fortigate_link_bandwidth_tx_byte_per_second",
			"Bandwidth available on this link for sending",
			[]string{"vdom", "monitor", "link"}, nil,
		)
		linkBandwidthRx = prometheus.NewDesc(
			"fortigate_link_bandwidth_rx_byte_per_second",
			"Bandwidth available on this link for sending",
			[]string{"vdom", "monitor", "link"}, nil,
		)
		linkStatusChanged = prometheus.NewDesc(
			"fortigate_link_status_change_time_seconds",
			"Unix timestamp describing the time when the last status change has occurred",
			[]string{"vdom", "monitor", "link"}, nil,
		)
	)

	type LinkMonitor struct {
		Status         string  `json:"status"`
		Latency        float64 `json:"latency"`
		Jitter         float64 `json:"jitter"`
		PacketLoss     float64 `json:"packet_loss"`
		PacketSent     float64 `json:"packet_sent"`
		PacketReceived float64 `json:"packet_received"`
		Session        float64 `json:"session"`
		TxBandwidth    float64 `json:"tx_bandwidth"`
		RxBandwidth    float64 `json:"rx_bandwidth"`
		StateChanged   float64 `json:"state_changed"`
	}

	type LinkGroup map[string]LinkMonitor

	type linkMonitorResponse struct {
		HTTPMethod string               `json:"http_method"`
		Results    map[string]LinkGroup `json:"results"`
		VDOM       string               `json:"vdom"`
		Path       string               `json:"path"`
		Name       string               `json:"name"`
		Status     string               `json:"status"`
		Serial     string               `json:"serial"`
		Version    string               `json:"version"`
		Build      int64                `json:"build"`
	}

	var rs []linkMonitorResponse

	if err := c.Get("api/v2/monitor/system/link-monitor", "vdom=*", &rs); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}

	for _, r := range rs {
		for linkGroupName, linkGroup := range r.Results {
			for linkName, link := range linkGroup {
				wanStatusUp, wanStatusDown, wanStatusError, wanStatusUnknown := 0.0, 0.0, 0.0, 0.0
				switch link.Status {
				case "up":
					wanStatusUp = 1.0
					break
				case "down":
					wanStatusDown = 1.0
					break
				case "error":
					wanStatusError = 1.0
					break
				default:
					wanStatusUnknown = 1.0
				}

				m = append(m, prometheus.MustNewConstMetric(linkStatus, prometheus.GaugeValue, wanStatusUp, r.VDOM, linkGroupName, linkName, "up"))
				m = append(m, prometheus.MustNewConstMetric(linkStatus, prometheus.GaugeValue, wanStatusDown, r.VDOM, linkGroupName, linkName, "down"))
				m = append(m, prometheus.MustNewConstMetric(linkStatus, prometheus.GaugeValue, wanStatusError, r.VDOM, linkGroupName, linkName, "error"))
				m = append(m, prometheus.MustNewConstMetric(linkStatus, prometheus.GaugeValue, wanStatusUnknown, r.VDOM, linkGroupName, linkName, "unknown"))

				// if no error or unknown status is reported, export the metrics
				if wanStatusError == 0 && wanStatusUnknown == 0 {
					m = append(m, prometheus.MustNewConstMetric(linkLatency, prometheus.GaugeValue, link.Latency/1000, r.VDOM, linkGroupName, linkName))
					m = append(m, prometheus.MustNewConstMetric(linkJitter, prometheus.GaugeValue, link.Jitter/1000, r.VDOM, linkGroupName, linkName))
					m = append(m, prometheus.MustNewConstMetric(linkPacketLoss, prometheus.GaugeValue, link.PacketLoss/100, r.VDOM, linkGroupName, linkName))
					m = append(m, prometheus.MustNewConstMetric(linkPacketSent, prometheus.CounterValue, link.PacketSent, r.VDOM, linkGroupName, linkName))
					m = append(m, prometheus.MustNewConstMetric(linkPacketReceived, prometheus.CounterValue, link.PacketReceived, r.VDOM, linkGroupName, linkName))
					m = append(m, prometheus.MustNewConstMetric(linkSessions, prometheus.GaugeValue, link.Session, r.VDOM, linkGroupName, linkName))
					m = append(m, prometheus.MustNewConstMetric(linkBandwidthTx, prometheus.GaugeValue, link.TxBandwidth/8, r.VDOM, linkGroupName, linkName))
					m = append(m, prometheus.MustNewConstMetric(linkBandwidthRx, prometheus.GaugeValue, link.RxBandwidth/8, r.VDOM, linkGroupName, linkName))
					m = append(m, prometheus.MustNewConstMetric(linkStatusChanged, prometheus.GaugeValue, link.StateChanged, r.VDOM, linkGroupName, linkName))
				}
			}
		}
	}
	return m, true
}
