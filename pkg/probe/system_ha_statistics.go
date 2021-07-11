package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

func probeSystemHAStatistics(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		memberInfo = prometheus.NewDesc(
			"fortigate_ha_member_info",
			"Info metric regarding cluster members",
			[]string{"vdom", "hostname", "serial", "group"}, nil,
		)
		memberSessions = prometheus.NewDesc(
			"fortigate_ha_member_sessions",
			"Sessions which are handled by this HA member",
			[]string{"vdom", "hostname"}, nil,
		)
		memberPackets = prometheus.NewDesc(
			"fortigate_ha_member_packets_total",
			"Packets which are handled by this HA member",
			[]string{"vdom", "hostname"}, nil,
		)
		memberVirusEvents = prometheus.NewDesc(
			"fortigate_ha_member_virus_events_total",
			"Virus events which are detected by this HA member",
			[]string{"vdom", "hostname"}, nil,
		)
		memberNetworkUsage = prometheus.NewDesc(
			"fortigate_ha_member_network_usage_ratio",
			"Network usage by HA member",
			[]string{"vdom", "hostname"}, nil,
		)
		memberBytesTotal = prometheus.NewDesc(
			"fortigate_ha_member_bytes_total",
			"Bytes transferred by HA member",
			[]string{"vdom", "hostname"}, nil,
		)
		memberIPSEvents = prometheus.NewDesc(
			"fortigate_ha_member_ips_events_total",
			"IPS events processed by HA member",
			[]string{"vdom", "hostname"}, nil,
		)
		memberCpuUsage = prometheus.NewDesc(
			"fortigate_ha_member_cpu_usage_ratio",
			"CPU usage by HA member",
			[]string{"vdom", "hostname"}, nil,
		)
		memberMemoryUsage = prometheus.NewDesc(
			"fortigate_ha_member_memory_usage_ratio",
			"Memory usage by HA member",
			[]string{"vdom", "hostname"}, nil,
		)
	)

	type HAResults struct {
		Hostname         string  `json:"hostname"`
		SerialNo         string  `json:"serial_no"`
		Tnow             float64 `json:"tnow"`
		Sessions         float64 `json:"sessions"`
		Tpacket          float64 `json:"tpacket"`
		VirEvents        float64 `json:"vir_usage"`
		NetUsage         float64 `json:"net_usage"`
		TransferredBytes float64 `json:"tbyte"`
		IPSEvents        float64 `json:"intr_usage"`
		CpuUsage         float64 `json:"cpu_usage"`
		MemUsage         float64 `json:"mem_usage"`
	}

	type HAResponse struct {
		HTTPMethod string      `json:"http_method"`
		Results    []HAResults `json:"results"`
		VDOM       string      `json:"vdom"`
		Path       string      `json:"path"`
		Name       string      `json:"name"`
		Status     string      `json:"status"`
		Serial     string      `json:"serial"`
		Version    string      `json:"version"`
		Build      int64       `json:"build"`
	}
	var r HAResponse

	if err := c.Get("api/v2/monitor/system/ha-statistics", "", &r); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	type HAConfig struct {
		Result struct {
			GroupName string `json:"group-name"`
		} `json:"results"`
	}
	var rc HAConfig

	if err := c.Get("api/v2/cmdb/system/ha", "", &rc); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, result := range r.Results {
		m = append(m, prometheus.MustNewConstMetric(memberInfo, prometheus.GaugeValue, 1, r.VDOM, result.Hostname, result.SerialNo, rc.Result.GroupName))
		m = append(m, prometheus.MustNewConstMetric(memberSessions, prometheus.GaugeValue, result.Sessions, r.VDOM, result.Hostname))
		m = append(m, prometheus.MustNewConstMetric(memberPackets, prometheus.CounterValue, result.Tpacket, r.VDOM, result.Hostname))
		m = append(m, prometheus.MustNewConstMetric(memberVirusEvents, prometheus.CounterValue, result.VirEvents, r.VDOM, result.Hostname))
		m = append(m, prometheus.MustNewConstMetric(memberNetworkUsage, prometheus.GaugeValue, result.NetUsage/100, r.VDOM, result.Hostname))
		m = append(m, prometheus.MustNewConstMetric(memberBytesTotal, prometheus.CounterValue, result.TransferredBytes, r.VDOM, result.Hostname))
		m = append(m, prometheus.MustNewConstMetric(memberIPSEvents, prometheus.CounterValue, result.IPSEvents, r.VDOM, result.Hostname))
		m = append(m, prometheus.MustNewConstMetric(memberCpuUsage, prometheus.GaugeValue, result.CpuUsage/100, r.VDOM, result.Hostname))
		m = append(m, prometheus.MustNewConstMetric(memberMemoryUsage, prometheus.GaugeValue, result.MemUsage/100, r.VDOM, result.Hostname))
	}
	return m, true
}
