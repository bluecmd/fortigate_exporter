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

	success :=
		probeSystemStatus(c, registry) &&
			probeSystemResources(c, registry) &&
			probeSystemVDOMResources(c, registry) &&
			probeFirewallPolicies(c, registry)

	// TODO(bluecmd): log/current-disk-usage
	return success, nil
}
