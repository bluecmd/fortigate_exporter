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

func probeSystemStatus(c FortiHTTP, registry *prometheus.Registry) bool {
	var (
		mVersion = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_version_info",
				Help: "System version and build information",
			},
			[]string{"serial", "version", "build"},
		)
	)

	registry.MustRegister(mVersion)

	type systemStatus struct {
		Status  string
		Serial  string
		Version string
		Build   int
	}
	var st systemStatus

	if err := c.Get("api/v2/monitor/system/status", "", &st); err != nil {
		log.Printf("Error: %v", err)
		return false
	}

	mVersion.WithLabelValues(st.Serial, st.Version, fmt.Sprintf("%d", st.Build)).Set(1)
	return true
}

func probeSystemResources(c FortiHTTP, registry *prometheus.Registry) bool {
	var (
		mResCPU = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_cpu_usage_ratio",
				Help: "Current resource usage ratio of system CPU, per core",
			},
			[]string{"processor"},
		)
		mResMemory = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_memory_usage_ratio",
				Help: "Current resource usage ratio of system memory",
			},
			[]string{},
		)
		mResSession = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_current_sessions",
				Help: "Current amount of sessions, per IP version",
			},
			[]string{"protocol"},
		)
	)

	registry.MustRegister(mResCPU)
	registry.MustRegister(mResMemory)
	registry.MustRegister(mResSession)

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
		return false
	}

	// CPU[0] is the average over all cores, ignore it
	for i, cpu := range sr.Results.CPU[1:] {
		mResCPU.WithLabelValues(fmt.Sprintf("%d", i)).Set(float64(cpu.Current) / 100.0)
	}
	mResMemory.WithLabelValues().Set(float64(sr.Results.Mem[0].Current) / 100.0)
	mResSession.WithLabelValues("ipv4").Set(float64(sr.Results.Session[0].Current))
	mResSession.WithLabelValues("ipv6").Set(float64(sr.Results.Session6[0].Current))
	return true
}

func probeSystemVDOMResources(c FortiHTTP, registry *prometheus.Registry) bool {
	var (
		mResCPU = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_vdom_cpu_usage_ratio",
				Help: "Current resource usage ratio of CPU, per VDOM",
			},
			[]string{"vdom"},
		)
		mResMemory = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_vdom_memory_usage_ratio",
				Help: "Current resource usage ratio of memory, per VDOM",
			},
			[]string{"vdom"},
		)
		mResSession = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_vdom_current_sessions",
				Help: "Current amount of sessions, per VDOM and IP version",
			},
			[]string{"vdom", "protocol"},
		)
	)

	registry.MustRegister(mResCPU)
	registry.MustRegister(mResMemory)
	registry.MustRegister(mResSession)

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
		return false
	}
	for _, s := range sr {
		mResCPU.WithLabelValues(s.VDOM).Set(float64(s.Results.CPU[0].Current) / 100.0)
		mResMemory.WithLabelValues(s.VDOM).Set(float64(s.Results.Mem[0].Current) / 100.0)
		mResSession.WithLabelValues(s.VDOM, "ipv4").Set(float64(s.Results.Session[0].Current))
		mResSession.WithLabelValues(s.VDOM, "ipv6").Set(float64(s.Results.Session6[0].Current))
	}
	return true
}

func probeVPNStatistics(c FortiHTTP, registry *prometheus.Registry) bool {
	var (
		vpncon = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_vpn_connections_count_total",
				Help: "Number of VPN connections",
			},
			[]string{"vdom"},
		)
	)
	registry.MustRegister(vpncon)
	type result struct {
		Results []map[string]interface{}
		VDOM    string
	}
	var res []result
	if err := c.Get("api/v2/monitor/vpn/ssl", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return false
	}

	for _, v := range res {
		count := len(v.Results)
		vpncon.WithLabelValues(v.VDOM).Set(float64(count))
	}

	return true

}
func probeIPSec(c FortiHTTP, registry *prometheus.Registry) bool {
	var (
		status = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_ipsec_tunnel_up",
				Help: "Status of Ipsec tunnel",
			}, []string{"vdom", "name", "parent"},
		)
		transmitted = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_ipsec_tunnel_transmit_bytes_total",
				Help: "Status of Ipsec tunnel",
			},
			[]string{"vdom", "name", "parent"},
		)
		received = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_ipsec_tunnel_receive_bytes_total",
				Help: "Status of Ipsec tunnel",
			},
			[]string{"vdom", "name", "parent"},
		)
	)

	registry.MustRegister(status)
	registry.MustRegister(transmitted)
	registry.MustRegister(received)

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
		return false
	}
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
				status.WithLabelValues(v.VDOM, t.Name, i.Name).Set(s)
				transmitted.WithLabelValues(v.VDOM, t.Name, i.Name).Set(float64(t.Outgoing))
				received.WithLabelValues(v.VDOM, t.Name, i.Name).Set(float64(t.Incoming))
			}
		}
	}
	return true

}

func probeFirewallPolicies(c FortiHTTP, registry *prometheus.Registry) bool {
	var (
		mHitCount = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_policy_hit_count_total",
				Help: "Number of times a policy has been hit",
			},
			[]string{"vdom", "protocol", "name", "uuid", "id"},
		)
		mBytes = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_policy_bytes_total",
				Help: "Number of bytes that has passed through a policy",
			},
			[]string{"vdom", "protocol", "name", "uuid", "id"},
		)
		mPackets = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_policy_packets_total",
				Help: "Number of packets that has passed through a policy",
			},
			[]string{"vdom", "protocol", "name", "uuid", "id"},
		)
		mActiveSessions = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_policy_active_sessions",
				Help: "Number of active sessions for a policy",
			},
			[]string{"vdom", "protocol", "name", "uuid", "id"},
		)
	)

	registry.MustRegister(mHitCount)
	registry.MustRegister(mBytes)
	registry.MustRegister(mPackets)
	registry.MustRegister(mActiveSessions)

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
	}
	var ps4 []policyStats
	var ps6 []policyStats

	if err := c.Get("api/v2/monitor/firewall/policy/select", "vdom=*", &ps4); err != nil {
		log.Printf("Error: %v", err)
		return false
	}
	if err := c.Get("api/v2/monitor/firewall/policy6/select", "vdom=*", &ps6); err != nil {
		log.Printf("Error: %v", err)
		return false
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
	var pc4 []policyConfig
	var pc6 []policyConfig

	query := "vdom=*&policyid|name|uuid|action|status"

	if err := c.Get("api/v2/cmdb/firewall/policy", query, &pc4); err != nil {
		log.Printf("Error: %v", err)
		return false
	}
	if err := c.Get("api/v2/cmdb/firewall/policy6", query, &pc6); err != nil {
		log.Printf("Error: %v", err)
		return false
	}

	pc4Map := map[string]*pConfig{}
	pc6Map := map[string]*pConfig{}
	for _, pc := range pc4 {
		for i, c := range pc.Results {
			pc4Map[c.UUID] = &pc.Results[i]
		}
	}
	for _, pc := range pc6 {
		for i, c := range pc.Results {
			pc6Map[c.UUID] = &pc.Results[i]
		}
	}

	process := func(ps *policyStats, s *pStats, pcMap map[string]*pConfig, proto string) {
		id := fmt.Sprintf("%d", s.ID)
		name := "Implicit Deny"
		if s.ID > 0 {
			c, ok := pcMap[s.UUID]
			if !ok {
				log.Printf("Warning: Failed to map %q to policy name - this should not happen", s.UUID)
				name = "<UNKNOWN>"
			} else {
				name = c.Name
			}
		}
		mHitCount.WithLabelValues(ps.VDOM, proto, name, s.UUID, id).Set(float64(s.HitCount))
		mBytes.WithLabelValues(ps.VDOM, proto, name, s.UUID, id).Set(float64(s.Bytes))
		mPackets.WithLabelValues(ps.VDOM, proto, name, s.UUID, id).Set(float64(s.Packets))
		mActiveSessions.WithLabelValues(ps.VDOM, proto, name, s.UUID, id).Set(float64(s.ActiveSessions))
	}

	for _, ps := range ps4 {
		for _, s := range ps.Results {
			process(&ps, &s, pc4Map, "ipv4")
		}
	}

	for _, ps := range ps6 {
		for _, s := range ps.Results {
			process(&ps, &s, pc6Map, "ipv6")
		}
	}

	return true
}

func probeInterfaces(c FortiHTTP, registry *prometheus.Registry) bool {
	var (
		mLink = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_interface_link_up",
				Help: "Whether the link is up or not",
			},
			[]string{"vdom", "name", "alias", "parent"},
		)
		mSpeed = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_interface_speed_bps",
				Help: "Speed negotiated on the port in bits/s",
			},
			[]string{"vdom", "name", "alias", "parent"},
		)
		mTxPkts = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "fortigate_interface_transmit_packets_total",
				Help: "Number of packets transmitted on the interface",
			},
			[]string{"vdom", "name", "alias", "parent"},
		)
		mRxPkts = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "fortigate_interface_receive_packets_total",
				Help: "Number of packets received on the interface",
			},
			[]string{"vdom", "name", "alias", "parent"},
		)
		mTxB = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "fortigate_interface_transmit_bytes_total",
				Help: "Number of bytes transmitted on the interface",
			},
			[]string{"vdom", "name", "alias", "parent"},
		)
		mRxB = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "fortigate_interface_receive_bytes_total",
				Help: "Number of bytes received on the interface",
			},
			[]string{"vdom", "name", "alias", "parent"},
		)
		mTxErr = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "fortigate_interface_transmit_errors_total",
				Help: "Number of transmission errors detected on the interface",
			},
			[]string{"vdom", "name", "alias", "parent"},
		)
		mRxErr = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "fortigate_interface_receive_errors_total",
				Help: "Number of reception errors detected on the interface",
			},
			[]string{"vdom", "name", "alias", "parent"},
		)
	)

	registry.MustRegister(mLink)
	registry.MustRegister(mSpeed)
	registry.MustRegister(mTxPkts)
	registry.MustRegister(mRxPkts)
	registry.MustRegister(mTxB)
	registry.MustRegister(mRxB)
	registry.MustRegister(mTxErr)
	registry.MustRegister(mRxErr)

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
		return false
	}
	for _, v := range r {
		for _, ir := range v.Results {
			linkf := 0.0
			if ir.Link {
				linkf = 1.0
			}
			mLink.WithLabelValues(v.VDOM, ir.Name, ir.Alias, ir.Interface).Set(linkf)
			mSpeed.WithLabelValues(v.VDOM, ir.Name, ir.Alias, ir.Interface).Set(ir.Speed * 1000 * 1000)
			mTxPkts.WithLabelValues(v.VDOM, ir.Name, ir.Alias, ir.Interface).Add(float64(ir.TxPackets))
			mRxPkts.WithLabelValues(v.VDOM, ir.Name, ir.Alias, ir.Interface).Add(float64(ir.RxPackets))
			mTxB.WithLabelValues(v.VDOM, ir.Name, ir.Alias, ir.Interface).Add(float64(ir.TxBytes))
			mRxB.WithLabelValues(v.VDOM, ir.Name, ir.Alias, ir.Interface).Add(float64(ir.RxBytes))
			mTxErr.WithLabelValues(v.VDOM, ir.Name, ir.Alias, ir.Interface).Add(float64(ir.TxErrors))
			mRxErr.WithLabelValues(v.VDOM, ir.Name, ir.Alias, ir.Interface).Add(float64(ir.RxErrors))
		}
	}
	return true
}

func probe(ctx context.Context, target string, registry *prometheus.Registry, hc *http.Client) (bool, error) {
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
	success :=
		probeSystemStatus(c, registry) &&
			probeSystemResources(c, registry) &&
			probeSystemVDOMResources(c, registry) &&
			probeFirewallPolicies(c, registry) &&
			probeInterfaces(c, registry) &&
			probeVPNStatistics(c, registry) &&
			probeIPSec(c, registry) &&
			probeHaStatistics(c, registry)

	// TODO(bluecmd): log/current-disk-usage
	return success, nil
}

func probeHaStatistics(c FortiHTTP, registry *prometheus.Registry) bool {
	var (
		memberInfo = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_ha_member_info",
				Help: "Info metric regarding cluster members",
			},
			[]string{"vdom", "hostname", "serial"},
		)
		memberSessions = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_ha_member_sessions",
				Help: "Sessions which are handled by this HA member",
			},
			[]string{"vdom", "hostname"},
		)
		memberPackets = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_ha_member_packets_total",
				Help: "Packets which are handled by this HA member",
			},
			[]string{"vdom", "hostname"},
		)
		memberVirusEvents = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "fortigate_ha_member_virus_events_total",
				Help: "Virus events which are detected by this HA member",
			},
			[]string{"vdom", "hostname"},
		)
		memberNetworkUsage = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_ha_member_network_usage_ratio",
				Help: "Network usage by HA member",
			},
			[]string{"vdom", "hostname"},
		)
		memberBytesTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "fortigate_ha_member_bytes_total",
				Help: "Bytes transferred by HA member",
			},
			[]string{"vdom", "hostname"},
		)
		memberIpsEvents = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "fortigate_ha_member_ips_events_total",
				Help: "IPS events processed by HA member",
			},
			[]string{"vdom", "hostname"},
		)
		memberCpuUsage = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_ha_member_cpu_usage_ratio",
				Help: "CPU usage by HA member",
			},
			[]string{"vdom", "hostname"},
		)
		memberMemoryUsage = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_ha_member_memory_usage_ratio",
				Help: "Memory usage by HA member",
			},
			[]string{"vdom", "hostname"},
		)
	)

	registry.MustRegister(memberInfo)
	registry.MustRegister(memberSessions)
	registry.MustRegister(memberPackets)
	registry.MustRegister(memberVirusEvents)
	registry.MustRegister(memberNetworkUsage)
	registry.MustRegister(memberBytesTotal)
	registry.MustRegister(memberIpsEvents)
	registry.MustRegister(memberCpuUsage)
	registry.MustRegister(memberMemoryUsage)

	type HaResults struct {
		Hostname         string  `json:"hostname"`
		SerialNo         string  `json:"serial_no"`
		Tnow             float64 `json:"tnow"`
		Sessions         float64 `json:"sessions"`
		Tpacket          float64 `json:"tpacket"`
		VirEvents        float64 `json:"vir_usage"`
		NetUsage         float64 `json:"net_usage"`
		TransferredBytes float64 `json:"tbyte"`
		IpsEvents        float64 `json:"intr_usage"`
		CpuUsage         float64 `json:"cpu_usage"`
		MemUsage         float64 `json:"mem_usage"`
	}

	type HaResponse struct {
		HTTPMethod string      `json:"http_method"`
		Results    []HaResults `json:"results"`
		VDOM       string      `json:"vdom"`
		Path       string      `json:"path"`
		Name       string      `json:"name"`
		Status     string      `json:"status"`
		Serial     string      `json:"serial"`
		Version    string      `json:"version"`
		Build      int64       `json:"build"`
	}
	var r HaResponse

	if err := c.Get("api/v2/monitor/system/ha-statistics", "", &r); err != nil {
		log.Printf("Error: %v", err)
		return false
	}
	for _, result := range r.Results {
		memberInfo.WithLabelValues(r.VDOM, result.Hostname, result.SerialNo).Set(1)
		memberSessions.WithLabelValues(r.VDOM, result.Hostname).Set(result.Sessions)
		memberPackets.WithLabelValues(r.VDOM, result.Hostname).Set(result.Tpacket)
		memberVirusEvents.WithLabelValues(r.VDOM, result.Hostname).Add(result.VirEvents)
		memberNetworkUsage.WithLabelValues(r.VDOM, result.Hostname).Add(result.NetUsage / 100)
		memberBytesTotal.WithLabelValues(r.VDOM, result.Hostname).Add(result.TransferredBytes)
		memberIpsEvents.WithLabelValues(r.VDOM, result.Hostname).Add(result.IpsEvents)
		memberCpuUsage.WithLabelValues(r.VDOM, result.Hostname).Add(result.CpuUsage / 100)
		memberMemoryUsage.WithLabelValues(r.VDOM, result.Hostname).Set(result.MemUsage / 100)
	}
	return true
}
