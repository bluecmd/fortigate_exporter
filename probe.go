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
				Name: "fortigate_system_version_info",
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
				Name: "fortigate_system_cpu_usage_ratio",
				Help: "Current resource usage ratio of system CPU, per core",
			},
			[]string{"processor"},
		)
		mResMemory = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_system_memory_usage_ratio",
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
				Name: "fortigate_vdom_system_cpu_usage_ratio",
				Help: "Current resource usage ratio of CPU, per VDOM",
			},
			[]string{"vdom"},
		)
		mResMemory = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_vdom_system_memory_usage_ratio",
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
			probeSystemVDOMResources(c, registry)

	// TODO(bluecmd): log/current-disk-usage
	return success, nil
}
