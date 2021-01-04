// All currently supported probes
//
// Copyright (C) 2020  Christian Svensson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/prometheus/client_golang/prometheus"
)

type ProbeCollector struct {
	metrics []prometheus.Metric
}

type probeFunc func(FortiHTTP) ([]prometheus.Metric, bool)

func probeSystemStatus(c FortiHTTP) ([]prometheus.Metric, bool) {
	var (
		mVersion = prometheus.NewDesc(
			"fortigate_version_info",
			"System version and build information",
			[]string{"serial", "version", "build"}, nil,
		)
	)

	type systemStatus struct {
		Status  string
		Serial  string
		Version string
		Build   int
	}
	var st systemStatus

	if err := c.Get("api/v2/monitor/system/status", "", &st); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{
		prometheus.MustNewConstMetric(mVersion, prometheus.GaugeValue, 1.0, st.Serial, st.Version, fmt.Sprintf("%d", st.Build)),
	}
	return m, true
}

func probeSystemResources(c FortiHTTP) ([]prometheus.Metric, bool) {
	var (
		mResCPU = prometheus.NewDesc(
			"fortigate_cpu_usage_ratio",
			"Current resource usage ratio of system CPU, per core",
			[]string{"processor"}, nil,
		)
		mResMemory = prometheus.NewDesc(
			"fortigate_memory_usage_ratio",
			"Current resource usage ratio of system memory",
			[]string{}, nil,
		)
		mResSession = prometheus.NewDesc(
			"fortigate_current_sessions",
			"Current amount of sessions, per IP version",
			[]string{"protocol"}, nil,
		)
	)

	type resUsage struct {
		Current int
	}
	type resContainer struct {
		CPU []resUsage
		Mem []resUsage
		// Ignore "disk", we get that from log/current-disk-usage instead with better resolution
		Session  []resUsage
		Session6 []resUsage
		// TODO(bluecmd): These are TODO
		// Setuprate []resUsage
		// Setuprate6 []resUsage
		// NpuSession []resUsage `json:"npu_session"`
		// NpuSession6 []resUsage `json:"npu_session6"`
		// NturboSession []resUsage  `json:"nturbo_session"`
		// NturboSession6 []resUsage `json:"nturbo_session6"`
		// DiskLograte []resUsage `json:"disk_lograte"`
		// FazLograte []resUsage `json:"faz_lograte"`
		// ForticloudLograte []resUsage `json:"forticloud_lograte"`
	}
	type systemResourceUsage struct {
		Results resContainer
		VDOM    string
	}
	var sr systemResourceUsage

	if err := c.Get("api/v2/monitor/system/resource/usage", "interval=1-min&scope=global", &sr); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	// CPU[0] is the average over all cores, ignore it
	m := []prometheus.Metric{}
	for i, cpu := range sr.Results.CPU[1:] {
		m = append(m, prometheus.MustNewConstMetric(
			mResCPU, prometheus.GaugeValue, float64(cpu.Current)/100.0, fmt.Sprintf("%d", i)))
	}
	m = append(m, prometheus.MustNewConstMetric(mResMemory, prometheus.GaugeValue, float64(sr.Results.Mem[0].Current)/100.0))
	m = append(m, prometheus.MustNewConstMetric(mResSession, prometheus.GaugeValue, float64(sr.Results.Session[0].Current), "ipv4"))
	m = append(m, prometheus.MustNewConstMetric(mResSession, prometheus.GaugeValue, float64(sr.Results.Session6[0].Current), "ipv6"))
	return m, true
}

func probeSystemVDOMResources(c FortiHTTP) ([]prometheus.Metric, bool) {
	var (
		mResCPU = prometheus.NewDesc(
			"fortigate_vdom_cpu_usage_ratio",
			"Current resource usage ratio of CPU, per VDOM",
			[]string{"vdom"}, nil,
		)
		mResMemory = prometheus.NewDesc(
			"fortigate_vdom_memory_usage_ratio",
			"Current resource usage ratio of memory, per VDOM",
			[]string{"vdom"}, nil,
		)
		mResSession = prometheus.NewDesc(
			"fortigate_vdom_current_sessions",
			"Current amount of sessions, per VDOM and IP version",
			[]string{"vdom", "protocol"}, nil,
		)
	)

	type resUsage struct {
		Current int
	}
	type resContainer struct {
		CPU []resUsage
		Mem []resUsage
		// Ignore "disk", we get that from log/current-disk-usage instead with better resolution
		Session  []resUsage
		Session6 []resUsage
		// TODO(bluecmd): These are TODO
		// Setuprate []resUsage
		// Setuprate6 []resUsage
		// DiskLograte []resUsage `json:"disk_lograte"`
		// FazLograte []resUsage `json:"faz_lograte"`
		// ForticloudLograte []resUsage `json:"forticloud_lograte"`
	}
	type systemResourceUsage struct {
		Results resContainer
		VDOM    string
	}
	var sr []systemResourceUsage

	if err := c.Get("api/v2/monitor/system/resource/usage", "interval=1-min&vdom=*", &sr); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}
	m := []prometheus.Metric{}
	for _, s := range sr {
		m = append(m, prometheus.MustNewConstMetric(mResCPU, prometheus.GaugeValue, float64(s.Results.CPU[0].Current)/100.0, s.VDOM))
		m = append(m, prometheus.MustNewConstMetric(mResMemory, prometheus.GaugeValue, float64(s.Results.Mem[0].Current)/100.0, s.VDOM))
		m = append(m, prometheus.MustNewConstMetric(mResSession, prometheus.GaugeValue, float64(s.Results.Session[0].Current), s.VDOM, "ipv4"))
		m = append(m, prometheus.MustNewConstMetric(mResSession, prometheus.GaugeValue, float64(s.Results.Session6[0].Current), s.VDOM, "ipv6"))
	}
	return m, true
}

func probeVPNStatistics(c FortiHTTP) ([]prometheus.Metric, bool) {
	var (
		// TODO(bluecmd): Rename as it is a gauge, issue #5
		vpncon = prometheus.NewDesc(
			"fortigate_vpn_connections_count_total",
			"Number of VPN connections",
			[]string{"vdom"}, nil,
		)
	)

	type result struct {
		Results []map[string]interface{}
		VDOM    string
	}
	var res []result
	if err := c.Get("api/v2/monitor/vpn/ssl", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, v := range res {
		count := len(v.Results)
		m = append(m, prometheus.MustNewConstMetric(vpncon, prometheus.GaugeValue, float64(count), v.VDOM))
	}

	return m, true

}
func probeIPSec(c FortiHTTP) ([]prometheus.Metric, bool) {
	var (
		status = prometheus.NewDesc(
			"fortigate_ipsec_tunnel_up",
			"Status of IPsec tunnel",
			[]string{"vdom", "name", "parent"}, nil,
		)
		transmitted = prometheus.NewDesc(
			"fortigate_ipsec_tunnel_transmit_bytes_total",
			"Total number of bytes transmitted over the IPsec tunnel",
			[]string{"vdom", "name", "parent"}, nil,
		)
		received = prometheus.NewDesc(
			"fortigate_ipsec_tunnel_receive_bytes_total",
			"Total number of bytes received over the IPsec tunnel",
			[]string{"vdom", "name", "parent"}, nil,
		)
	)

	type proxyid struct {
		Name     string `json:"p2name"`
		Status   string `json:"status"`
		Incoming int    `json:"incoming_bytes"`
		Outgoing int    `json:"outgoing_bytes"`
	}
	type tunnel struct {
		Name    string    `json:"name"`
		Type    string    `json:"type"`
		ProxyID []proxyid `json:"proxyid"`
	}
	type ipsecResult struct {
		Results []tunnel `json:"results"`
		VDOM    string
	}
	var res []ipsecResult
	if err := c.Get("api/v2/monitor/vpn/ipsec", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, v := range res {
		for _, i := range v.Results {
			/*
			  type 'dialup' seems to be client vpn.
			  Not sure exactly what the difference is between probeVPNStatistics
			*/
			if i.Type == "dialup" {
				continue
			}
			for _, t := range i.ProxyID {
				s := 0.0
				if t.Status == "up" {
					s = 1.0
				}
				m = append(m, prometheus.MustNewConstMetric(status, prometheus.GaugeValue, s, v.VDOM, t.Name, i.Name))
				m = append(m, prometheus.MustNewConstMetric(transmitted, prometheus.GaugeValue, float64(t.Outgoing), v.VDOM, t.Name, i.Name))
				m = append(m, prometheus.MustNewConstMetric(received, prometheus.GaugeValue, float64(t.Incoming), v.VDOM, t.Name, i.Name))
			}
		}
	}
	return m, true

}

func probeFirewallPolicies(c FortiHTTP) ([]prometheus.Metric, bool) {
	var (
		mHitCount = prometheus.NewDesc(
			"fortigate_policy_hit_count_total",
			"Number of times a policy has been hit",
			[]string{"vdom", "protocol", "name", "uuid", "id"}, nil,
		)
		mBytes = prometheus.NewDesc(
			"fortigate_policy_bytes_total",
			"Number of bytes that has passed through a policy",
			[]string{"vdom", "protocol", "name", "uuid", "id"}, nil,
		)
		mPackets = prometheus.NewDesc(
			"fortigate_policy_packets_total",
			"Number of packets that has passed through a policy",
			[]string{"vdom", "protocol", "name", "uuid", "id"}, nil,
		)
		mActiveSessions = prometheus.NewDesc(
			"fortigate_policy_active_sessions",
			"Number of active sessions for a policy",
			[]string{"vdom", "protocol", "name", "uuid", "id"}, nil,
		)
	)

	type pStats struct {
		ID               int `json:"policyid"`
		UUID             string
		ActiveSessions   int `json:"active_sessions"`
		Bytes            int
		Packets          int
		SoftwareBytes    int `json:"software_bytes"`
		SoftwarePackets  int `json:"software_packets"`
		ASICBytes        int `json:"asic_bytes"`
		ASICPackets      int `json:"asic_packets"`
		NTurboBytes      int `json:"nturbo_bytes"`
		NTurboPackets    int `json:"nturbo_packets"`
		HitCount         int `json:"hit_count"`
		SessionCount     int `json:"session_count"`
		SessionLastUsed  int `json:"session_last_used"`
		SessionFirstUsed int `json:"session_first_used"`
		LastUsed         int `json:"last_used"`
		FirstUsed        int `json:"first_used"`
	}
	type policyStats struct {
		Results []pStats
		VDOM    string
		Version string
	}
	var ps4 []policyStats
	var ps6 []policyStats

	// NOTE: ip_version=ipv4 is a no-op if combined policies are not active
	if err := c.Get("api/v2/monitor/firewall/policy/select", "vdom=*&ip_version=ipv4", &ps4); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	combined := false
	maj, min, ok := ParseVersion(ps4[0].Version)
	if !ok {
		log.Printf("Could not parse version number %q", ps4[0].Version)
		return nil, false
	}
	// If we are at 6.4 or later we use combined policies
	if maj > 6 || (maj == 6 && min >= 4) {
		combined = true
	}

	if !combined {
		if err := c.Get("api/v2/monitor/firewall/policy6/select", "vdom=*", &ps6); err != nil {
			log.Printf("Error: %v", err)
			return nil, false
		}
	} else {
		if err := c.Get("api/v2/monitor/firewall/policy/select", "vdom=*&ip_version=ipv6", &ps6); err != nil {
			log.Printf("Error: %v", err)
			return nil, false
		}
	}

	type pConfig struct {
		ID     int `json:"policyid"`
		Name   string
		UUID   string
		Action string
		Status string
	}

	type policyConfig struct {
		Results []pConfig
		VDOM    string
	}
	var pc []policyConfig
	var pc6 []policyConfig

	query := "vdom=*&policyid|name|uuid|action|status"

	if err := c.Get("api/v2/cmdb/firewall/policy", query, &pc); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}
	if !combined {
		if err := c.Get("api/v2/cmdb/firewall/policy6", query, &pc6); err != nil {
			log.Printf("Error: %v", err)
			return nil, false
		}
	}

	pc4Map := map[string]*pConfig{}
	pc6Map := map[string]*pConfig{}
	for _, pc := range pc {
		for i, c := range pc.Results {
			pc4Map[c.UUID] = &pc.Results[i]
		}
	}
	if !combined {
		for _, pc := range pc6 {
			for i, c := range pc.Results {
				pc6Map[c.UUID] = &pc.Results[i]
			}
		}
	} else {
		pc6Map = pc4Map
	}

	process := func(ps *policyStats, s *pStats, pcMap map[string]*pConfig, proto string) []prometheus.Metric {
		id := fmt.Sprintf("%d", s.ID)
		name := "Implicit Deny"
		if s.ID > 0 {
			c, ok := pcMap[s.UUID]
			if !ok {
				log.Printf("Warning: Failed to map %q to policy config - this should not happen", s.UUID)
				name = "<UNKNOWN>"
			} else {
				name = c.Name
			}
		}
		m := []prometheus.Metric{
			prometheus.MustNewConstMetric(mHitCount, prometheus.GaugeValue, float64(s.HitCount), ps.VDOM, proto, name, s.UUID, id),
			prometheus.MustNewConstMetric(mBytes, prometheus.GaugeValue, float64(s.Bytes), ps.VDOM, proto, name, s.UUID, id),
			prometheus.MustNewConstMetric(mPackets, prometheus.GaugeValue, float64(s.Packets), ps.VDOM, proto, name, s.UUID, id),
			prometheus.MustNewConstMetric(mActiveSessions, prometheus.GaugeValue, float64(s.ActiveSessions), ps.VDOM, proto, name, s.UUID, id),
		}
		return m
	}

	m := []prometheus.Metric{}
	for _, ps := range ps4 {
		for _, s := range ps.Results {
			m = append(m, process(&ps, &s, pc4Map, "ipv4")...)
		}
	}

	for _, ps := range ps6 {
		for _, s := range ps.Results {
			m = append(m, process(&ps, &s, pc6Map, "ipv6")...)
		}
	}

	return m, true
}

func probeInterfaces(c FortiHTTP) ([]prometheus.Metric, bool) {
	var (
		mLink = prometheus.NewDesc(
			"fortigate_interface_link_up",
			"Whether the link is up or not",
			[]string{"vdom", "name", "alias", "parent"}, nil,
		)
		mSpeed = prometheus.NewDesc(
			"fortigate_interface_speed_bps",
			"Speed negotiated on the port in bits/s",
			[]string{"vdom", "name", "alias", "parent"}, nil,
		)
		mTxPkts = prometheus.NewDesc(
			"fortigate_interface_transmit_packets_total",
			"Number of packets transmitted on the interface",
			[]string{"vdom", "name", "alias", "parent"}, nil,
		)
		mRxPkts = prometheus.NewDesc(
			"fortigate_interface_receive_packets_total",
			"Number of packets received on the interface",
			[]string{"vdom", "name", "alias", "parent"}, nil,
		)
		mTxB = prometheus.NewDesc(
			"fortigate_interface_transmit_bytes_total",
			"Number of bytes transmitted on the interface",
			[]string{"vdom", "name", "alias", "parent"}, nil,
		)
		mRxB = prometheus.NewDesc(
			"fortigate_interface_receive_bytes_total",
			"Number of bytes received on the interface",
			[]string{"vdom", "name", "alias", "parent"}, nil,
		)
		mTxErr = prometheus.NewDesc(
			"fortigate_interface_transmit_errors_total",
			"Number of transmission errors detected on the interface",
			[]string{"vdom", "name", "alias", "parent"}, nil,
		)
		mRxErr = prometheus.NewDesc(
			"fortigate_interface_receive_errors_total",
			"Number of reception errors detected on the interface",
			[]string{"vdom", "name", "alias", "parent"}, nil,
		)
	)

	type ifResult struct {
		Id        string
		Name      string
		Alias     string
		Link      bool
		Speed     float64
		Duplex    int
		TxPackets int64 `json:"tx_packets"`
		RxPackets int64 `json:"rx_packets"`
		TxBytes   int64 `json:"tx_bytes"`
		RxBytes   int64 `json:"rx_bytes"`
		TxErrors  int64 `json:"tx_errors"`
		RxErrors  int64 `json:"rx_errors"`
		Interface string
	}
	type ifResponse struct {
		Results map[string]ifResult
		VDOM    string
	}
	var r []ifResponse

	if err := c.Get("api/v2/monitor/system/interface/select", "vdom=*&include_vlan=true&include_aggregate=true", &r); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}
	m := []prometheus.Metric{}
	for _, v := range r {
		for _, ir := range v.Results {
			linkf := 0.0
			if ir.Link {
				linkf = 1.0
			}
			m = append(m, prometheus.MustNewConstMetric(mLink, prometheus.GaugeValue, linkf, v.VDOM, ir.Name, ir.Alias, ir.Interface))
			m = append(m, prometheus.MustNewConstMetric(mSpeed, prometheus.GaugeValue, ir.Speed*1000*1000, v.VDOM, ir.Name, ir.Alias, ir.Interface))
			m = append(m, prometheus.MustNewConstMetric(mTxPkts, prometheus.CounterValue, float64(ir.TxPackets), v.VDOM, ir.Name, ir.Alias, ir.Interface))
			m = append(m, prometheus.MustNewConstMetric(mRxPkts, prometheus.CounterValue, float64(ir.RxPackets), v.VDOM, ir.Name, ir.Alias, ir.Interface))
			m = append(m, prometheus.MustNewConstMetric(mTxB, prometheus.CounterValue, float64(ir.TxBytes), v.VDOM, ir.Name, ir.Alias, ir.Interface))
			m = append(m, prometheus.MustNewConstMetric(mRxB, prometheus.CounterValue, float64(ir.RxBytes), v.VDOM, ir.Name, ir.Alias, ir.Interface))
			m = append(m, prometheus.MustNewConstMetric(mTxErr, prometheus.CounterValue, float64(ir.TxErrors), v.VDOM, ir.Name, ir.Alias, ir.Interface))
			m = append(m, prometheus.MustNewConstMetric(mRxErr, prometheus.CounterValue, float64(ir.RxErrors), v.VDOM, ir.Name, ir.Alias, ir.Interface))
		}
	}
	return m, true
}

func probeHAStatistics(c FortiHTTP) ([]prometheus.Metric, bool) {
	var (
		memberInfo = prometheus.NewDesc(
			"fortigate_ha_member_info",
			"Info metric regarding cluster members",
			[]string{"vdom", "hostname", "serial"}, nil,
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

	m := []prometheus.Metric{}
	for _, result := range r.Results {
		m = append(m, prometheus.MustNewConstMetric(memberInfo, prometheus.GaugeValue, 1, r.VDOM, result.Hostname, result.SerialNo))
		m = append(m, prometheus.MustNewConstMetric(memberSessions, prometheus.GaugeValue, result.Sessions, r.VDOM, result.Hostname))
		m = append(m, prometheus.MustNewConstMetric(memberPackets, prometheus.GaugeValue, result.Tpacket, r.VDOM, result.Hostname))
		m = append(m, prometheus.MustNewConstMetric(memberVirusEvents, prometheus.CounterValue, result.VirEvents, r.VDOM, result.Hostname))
		m = append(m, prometheus.MustNewConstMetric(memberNetworkUsage, prometheus.GaugeValue, result.NetUsage/100, r.VDOM, result.Hostname))
		m = append(m, prometheus.MustNewConstMetric(memberBytesTotal, prometheus.CounterValue, result.TransferredBytes, r.VDOM, result.Hostname))
		m = append(m, prometheus.MustNewConstMetric(memberIPSEvents, prometheus.CounterValue, result.IPSEvents, r.VDOM, result.Hostname))
		m = append(m, prometheus.MustNewConstMetric(memberCpuUsage, prometheus.GaugeValue, result.CpuUsage/100, r.VDOM, result.Hostname))
		m = append(m, prometheus.MustNewConstMetric(memberMemoryUsage, prometheus.GaugeValue, result.MemUsage/100, r.VDOM, result.Hostname))
	}
	return m, true
}

func (p *ProbeCollector) Probe(ctx context.Context, target string, hc *http.Client) (bool, error) {
	tgt, err := url.Parse(target)
	if err != nil {
		return false, fmt.Errorf("url.Parse failed: %v", err)
	}

	if tgt.Scheme != "https" && tgt.Scheme != "http" {
		return false, fmt.Errorf("Unsupported scheme %q", tgt.Scheme)
	}

	// Filter anything else than scheme and hostname
	u := url.URL{
		Scheme: tgt.Scheme,
		Host:   tgt.Host,
	}
	c, err := newFortiClient(ctx, u, hc)
	if err != nil {
		return false, err
	}

	// TODO: Make parallel
	success := true
	for _, f := range []probeFunc{
		probeSystemStatus,
		probeSystemResources,
		probeSystemVDOMResources,
		probeFirewallPolicies,
		probeInterfaces,
		probeVPNStatistics,
		probeIPSec,
		probeHAStatistics,
	} {
		m, ok := f(c)
		if !ok {
			success = false
		}
		p.metrics = append(p.metrics, m...)
	}

	return success, nil
}

func (p *ProbeCollector) Collect(c chan<- prometheus.Metric) {
	// Collect result of new probe functions
	for _, m := range p.metrics {
		c <- m
	}
}

func (p *ProbeCollector) Describe(c chan<- *prometheus.Desc) {
}
