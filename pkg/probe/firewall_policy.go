package probe

import (
	"fmt"
	"log"

	"github.com/bluecmd/fortigate_exporter/internal/version"
	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

func probeFirewallPolicies(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
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
		ID               int64 `json:"policyid"`
		UUID             string
		ActiveSessions   float64 `json:"active_sessions"`
		Bytes            float64
		Packets          float64
		SoftwareBytes    float64 `json:"software_bytes"`
		SoftwarePackets  float64 `json:"software_packets"`
		ASICBytes        float64 `json:"asic_bytes"`
		ASICPackets      float64 `json:"asic_packets"`
		NTurboBytes      float64 `json:"nturbo_bytes"`
		NTurboPackets    float64 `json:"nturbo_packets"`
		HitCount         float64 `json:"hit_count"`
		SessionCount     float64 `json:"session_count"`
		SessionLastUsed  float64 `json:"session_last_used"`
		SessionFirstUsed float64 `json:"session_first_used"`
		LastUsed         float64 `json:"last_used"`
		FirstUsed        float64 `json:"first_used"`
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
	maj, min, ok := version.ParseVersion(ps4[0].Version)
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
		ID     int64 `json:"policyid"`
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
			prometheus.MustNewConstMetric(mHitCount, prometheus.CounterValue, s.HitCount, ps.VDOM, proto, name, s.UUID, id),
			prometheus.MustNewConstMetric(mBytes, prometheus.CounterValue, s.Bytes, ps.VDOM, proto, name, s.UUID, id),
			prometheus.MustNewConstMetric(mPackets, prometheus.CounterValue, s.Packets, ps.VDOM, proto, name, s.UUID, id),
			prometheus.MustNewConstMetric(mActiveSessions, prometheus.GaugeValue, s.ActiveSessions, ps.VDOM, proto, name, s.UUID, id),
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
