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
		vpncon = prometheus.NewDesc(
			"fortigate_vpn_connections",
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
				m = append(m, prometheus.MustNewConstMetric(transmitted, prometheus.CounterValue, float64(t.Outgoing), v.VDOM, t.Name, i.Name))
				m = append(m, prometheus.MustNewConstMetric(received, prometheus.CounterValue, float64(t.Incoming), v.VDOM, t.Name, i.Name))
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
			prometheus.MustNewConstMetric(mHitCount, prometheus.CounterValue, float64(s.HitCount), ps.VDOM, proto, name, s.UUID, id),
			prometheus.MustNewConstMetric(mBytes, prometheus.CounterValue, float64(s.Bytes), ps.VDOM, proto, name, s.UUID, id),
			prometheus.MustNewConstMetric(mPackets, prometheus.CounterValue, float64(s.Packets), ps.VDOM, proto, name, s.UUID, id),
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

func probeLicenseStatus(c FortiHTTP) ([]prometheus.Metric, bool) {
	var (
		vdomUsed = prometheus.NewDesc(
			"fortigate_license_vdom_usage",
			"The amount of VDOM licenses currently used",
			[]string{}, nil,
		)
		vdomMax = prometheus.NewDesc(
			"fortigate_license_vdom_max",
			"The total amount of VDOM licenses available",
			[]string{}, nil,
		)
	)

	type LicenseStatus struct {
		VDOM struct {
			Type       string `json:"type"`
			CanUpgrade bool   `json:"can_upgrade"`
			Used       int    `json:"used"`
			Max        int    `json:"max"`
		} `json:"vdom"`
	}

	type LicenseResponse struct {
		Results LicenseStatus `json:"results"`
	}
	var r LicenseResponse

	if err := c.Get("api/v2/monitor/license/status/select", "", &r); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{
		prometheus.MustNewConstMetric(vdomUsed, prometheus.GaugeValue, float64(r.Results.VDOM.Used)),
		prometheus.MustNewConstMetric(vdomMax, prometheus.GaugeValue, float64(r.Results.VDOM.Max)),
	}

	return m, true
}

func probeLinkMonitor(c FortiHTTP) ([]prometheus.Metric, bool) {
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
			"Percentage of packages lost relative to  all sent based on the last 30 probes",
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
		Build      int                  `json:"build"`
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

func probeVirtualWanPerf(c FortiHTTP) ([]prometheus.Metric, bool) {
	var (
		mLink = prometheus.NewDesc(
			"fortigate_virtual_wan_healthcheck_status",
			"Status of the health check. If the SD-WAN interface is disabled, disable will be returned. If the interface does not participate in the health check, error will be returned.",
			[]string{"vdom","sla", "interface", "state"}, nil,
		)
		mLatency = prometheus.NewDesc(
			"fortigate_virtual_wan_healthcheck_latency_seconds",
			"Measured latency for this health check",
			[]string{"vdom", "sla", "interface"}, nil,
		)
		mJitter = prometheus.NewDesc(
			"fortigate_virtual_wan_healthcheck_jitter_seconds",
			"Measured latency jitter for this health check",
			[]string{"vdom", "sla", "interface"}, nil,
		)
		mPacketLoss = prometheus.NewDesc(
			"fortigate_virtual_wan_healthcheck_packetloss_ratio",
			"Measured packet loss in percentage for this health check",
			[]string{"vdom", "sla", "interface"}, nil,
		)
		mPacketSent = prometheus.NewDesc(
			"fortigate_virtual_wan_healthcheck_packetsent_total",
			"Number of packets sent for this health check",
			[]string{"vdom", "sla", "interface"}, nil,
		)
		mPacketReceived = prometheus.NewDesc(
			"fortigate_virtual_wan_healthcheck_packetreceived_total",
			"Number of packets received for this health check",
			[]string{"vdom", "sla", "interface"}, nil,
		)
		mSession = prometheus.NewDesc(
			"fortigate_virtual_wan_healthcheck_active_sessions",
			"Active Session count for the health check interface",
			[]string{"vdom", "sla", "interface"}, nil,
		)
		mTXBandwidth = prometheus.NewDesc(
			"fortigate_virtual_wan_healthcheck_bandwidth_tx_byte_per_second",
			"Upload bandwidth of the health check interface",
			[]string{"vdom", "sla", "interface"}, nil,
		)
		mRXBandwidth = prometheus.NewDesc(
			"fortigate_virtual_wan_healthcheck_bandwidth_rx_byte_per_second",
			"Download bandwidth of the health check interface",
			[]string{"vdom", "sla", "interface"}, nil,
		)
		mStateChanged = prometheus.NewDesc(
			"fortigate_virtual_wan_healthcheck_status_change_time_seconds",
			"Unix timestamp describing the time when the last status change has occurred",
			[]string{"vdom", "sla", "interface"}, nil,
		)
	)

	type SLAMember struct {
		Status         string  `json:"status"`
		Latency        float64 `json:"latency"`
		Jitter         float64 `json:"jitter"`
		PacketLoss     float64 `json:"packet_loss"`
		PacketSent     float64 `json:"packet_sent"`
		PacketReceived float64 `json:"packet_received"`
		//todo add slatargetmet
		SLAtargetmet   []string  `json:"sla_targets_met"`
		Session        float64 `json:"session"`
		TxBandwidth    float64 `json:"tx_bandwidth"`
		RxBandwidth    float64 `json:"rx_bandwidth"`
		StateChanged   float64 `json:"state_changed"`
	}

	type VirtualWanSLA map[string]SLAMember

	type VirtualWanMonitorResponse struct {
		HTTPMethod string               `json:"http_method"`
		Results    map[string]VirtualWanSLA `json:"results"`
		VDOM       string               `json:"vdom"`
		Path       string               `json:"path"`
		Name       string               `json:"name"`
		Status     string               `json:"status"`
		Serial     string               `json:"serial"`
		Version    string               `json:"version"`
		Build      int                  `json:"build"`
	}


	var rs []VirtualWanMonitorResponse

	if err := c.Get("api/v2/monitor/virtual-wan/health-check","vdom=*", &rs); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}
	m := []prometheus.Metric{}
	for _, r := range rs {
		for VirtualWanSLAName, VirtualWanSLA := range r.Results {
			for MemberName, Member := range VirtualWanSLA {
				MemberStatusUp, MemberStatusDown, MemberStatusError, MemberStatusDisable, MemberStatusUnknown := 0.0, 0.0, 0.0, 0.0, 0.0
				switch Member.Status {
				case "up":
					MemberStatusUp = 1.0
					break
				case "down":
					MemberStatusDown = 1.0
					break
				case "error":
					MemberStatusError = 1.0
					break
				case "disable":
					MemberStatusDisable = 1.0
					break
				default:
					MemberStatusUnknown = 1.0
				}

				m = append(m, prometheus.MustNewConstMetric(mLink, prometheus.GaugeValue, MemberStatusUp, r.VDOM, VirtualWanSLAName, MemberName, "up"))
				m = append(m, prometheus.MustNewConstMetric(mLink, prometheus.GaugeValue, MemberStatusDown, r.VDOM, VirtualWanSLAName, MemberName, "down"))
				m = append(m, prometheus.MustNewConstMetric(mLink, prometheus.GaugeValue, MemberStatusError, r.VDOM, VirtualWanSLAName, MemberName, "error"))
				m = append(m, prometheus.MustNewConstMetric(mLink, prometheus.GaugeValue, MemberStatusDisable, r.VDOM, VirtualWanSLAName, MemberName, "disable"))
				m = append(m, prometheus.MustNewConstMetric(mLink, prometheus.GaugeValue, MemberStatusUnknown, r.VDOM, VirtualWanSLAName, MemberName, "unknown"))
				// if no error or unknown status is reported, export the metrics
				if MemberStatusUp == 1 {
					m = append(m, prometheus.MustNewConstMetric(mLatency, prometheus.GaugeValue, Member.Latency/1000, r.VDOM, VirtualWanSLAName, MemberName,))
					m = append(m, prometheus.MustNewConstMetric(mJitter, prometheus.GaugeValue, Member.Jitter/1000, r.VDOM, VirtualWanSLAName, MemberName,))
					m = append(m, prometheus.MustNewConstMetric(mPacketLoss, prometheus.GaugeValue, Member.PacketLoss/100, r.VDOM, VirtualWanSLAName, MemberName,))
					m = append(m, prometheus.MustNewConstMetric(mPacketSent, prometheus.GaugeValue, Member.PacketSent, r.VDOM, VirtualWanSLAName, MemberName,))
					m = append(m, prometheus.MustNewConstMetric(mPacketReceived, prometheus.GaugeValue, Member.PacketReceived, r.VDOM, VirtualWanSLAName, MemberName,))
					m = append(m, prometheus.MustNewConstMetric(mSession, prometheus.GaugeValue, Member.Session, VirtualWanSLAName, r.VDOM, MemberName,))
					m = append(m, prometheus.MustNewConstMetric(mTXBandwidth, prometheus.GaugeValue, Member.TxBandwidth/8, r.VDOM, VirtualWanSLAName, MemberName,))
					m = append(m, prometheus.MustNewConstMetric(mRXBandwidth, prometheus.GaugeValue, Member.RxBandwidth/8, r.VDOM, VirtualWanSLAName, MemberName,))
					m = append(m, prometheus.MustNewConstMetric(mStateChanged, prometheus.GaugeValue, Member.StateChanged, r.VDOM, VirtualWanSLAName, MemberName,))
				}
			}
		}
	}
	return m, true
}

func probeCertificates(c FortiHTTP) ([]prometheus.Metric, bool) {
	var (
		certificateInfo = prometheus.NewDesc(
			"fortigate_certificate_info",
			"Info metric containing meta information about the certificate",
			[]string{"name", "source", "scope", "vdom", "status", "type"}, nil,
		)
		certificateValidFrom = prometheus.NewDesc(
			"fortigate_certificate_valid_from_seconds",
			"Unix timestamp from which this certificate is valid",
			[]string{"name", "source", "scope", "vdom"}, nil,
		)
		certificateValidTo = prometheus.NewDesc(
			"fortigate_certificate_valid_to_seconds",
			"Unix timestamp till which this certificate is valid",
			[]string{"name", "source", "scope", "vdom"}, nil,
		)
		certificateCMDBReferences = prometheus.NewDesc(
			"fortigate_certificate_cmdb_references",
			"Number of times the certificate is referenced",
			[]string{"name", "source", "scope", "vdom"}, nil,
		)
	)

	type Results struct {
		Name      string  `json:"name"`
		Source    string  `json:"source"`
		Type      string  `json:"type"`
		Status    string  `json:"status"`
		ValidFrom float64 `json:"valid_from"`
		ValidTo   float64 `json:"valid_to"`
		QRef      float64 `json:"q_ref"`
	}

	type Response struct {
		Results []Results `json:"results"`
		VDOM    string    `json:"vdom"`
		Status  string    `json:"status"`
		Scope   string
	}

	var globalResponse Response
	if err := c.Get("api/v2/monitor/system/available-certificates", "scope=global", &globalResponse); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}
	globalResponse.Scope = "global"

	var vdomResponses []Response

	if err := c.Get("api/v2/monitor/system/available-certificates", "vdom=*", &vdomResponses); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}
	for i, _ := range vdomResponses {
		vdomResponses[i].Scope = "vdom"
	}

	// combine responses
	combinedResponses := make([]Response, 0)
	combinedResponses = append(combinedResponses, vdomResponses...)
	combinedResponses = append(combinedResponses, globalResponse)

	m := []prometheus.Metric{}

	for _, response := range combinedResponses {
		for _, result := range response.Results {
			m = append(m, prometheus.MustNewConstMetric(certificateInfo, prometheus.GaugeValue, 1, result.Name, result.Source, response.Scope, response.VDOM, result.Status, result.Type))
			m = append(m, prometheus.MustNewConstMetric(certificateValidFrom, prometheus.GaugeValue, result.ValidFrom, result.Name, result.Source, response.Scope, response.VDOM))
			m = append(m, prometheus.MustNewConstMetric(certificateValidTo, prometheus.GaugeValue, result.ValidTo, result.Name, result.Source, response.Scope, response.VDOM))
			m = append(m, prometheus.MustNewConstMetric(certificateCMDBReferences, prometheus.GaugeValue, result.QRef, result.Name, result.Source, response.Scope, response.VDOM))
		}
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
		probeLicenseStatus,
		probeLinkMonitor,
		probeVirtualWanPerf,
		probeCertificates,
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
