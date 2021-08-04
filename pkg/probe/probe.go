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

package probe

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/bluecmd/fortigate_exporter/internal/config"
	"github.com/bluecmd/fortigate_exporter/internal/version"
	fortiHTTP "github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type ProbeCollector struct {
	metrics []prometheus.Metric
}

type TargetMetadata struct {
	VersionMajor int
	VersionMinor int
}

type probeFunc func(fortiHTTP.FortiHTTP, *TargetMetadata) ([]prometheus.Metric, bool)

type probeDetailedFunc struct {
	category string
	name     string
	function probeFunc
}

func (p *ProbeCollector) Probe(ctx context.Context, target string, hc *http.Client, savedConfig config.FortiExporterConfig) (bool, error) {
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
	c, err := fortiHTTP.NewFortiClient(ctx, u, hc, savedConfig)
	if err != nil {
		return false, err
	}

	type systemStatus struct {
		Status  string
		Version string
	}
	var st systemStatus

	// Test client connection before we blast all the probes.
	// The "system status" group has access group "any" so it is a good source
	// to test the authentication as well as fetching the OS version.
	if err := c.Get("api/v2/monitor/system/status", "", &st); err != nil {
		log.Printf("Error: API connectivity test failed, %v", err)
		return false, nil
	}

	if st.Status != "success" {
		log.Printf("Error: API connectivity test returned status: %s", st.Status)
		return false, nil
	}

	major, minor, ok := version.ParseVersion(st.Version)
	if !ok {
		log.Printf("Error: Failed to parse OS version: %q", st.Version)
		return false, nil
	}

	meta := &TargetMetadata{
		VersionMajor: major,
		VersionMinor: minor,
	}

	wantedProbes := savedConfig.AuthKeys[config.Target(u.String())].Probes

	// TODO: Make parallel
	success := true
	for _, aProbe := range []probeDetailedFunc{
		{"BGP", "BGPNeighborPathsIPv4", probeBGPNeighborPathsIPv4},
		{"BGP", "BGPNeighborPathsIPv6", probeBGPNeighborPathsIPv6},
		{"BGP", "BGPNeighborsIPv4", probeBGPNeighborsIPv4},
		{"BGP", "BGPNeighborsIPv6", probeBGPNeighborsIPv6},
		{"Firewall", "FirewallLoadBalance", probeFirewallLoadBalance},
		{"Firewall", "FirewallPolicies", probeFirewallPolicies},
		{"Licence", "LicenseStatus", probeLicenseStatus},
		{"System", "SystemAvailableCertificates", probeSystemAvailableCertificates},
		{"System", "SystemHAStatistics", probeSystemHAStatistics},
		{"System", "SystemInterface", probeSystemInterface},
		{"System", "SystemLinkMonitor", probeSystemLinkMonitor},
		{"System", "SystemResourceUsage", probeSystemResourceUsage},
		{"System", "SystemStatus", probeSystemStatus},
		{"System", "SystemVDOMResources", probeSystemVDOMResources},
		{"VPN", "VPNIPSec", probeVPNIPSec},
		{"VPN", "VPNSsl", probeVPNSsl},
		{"Virtual", "VirtualWANHealthCheck", probeVirtualWANHealthCheck},
		{"Wifi", "WifiAPStatus", probeWifiAPStatus},
		{"Wifi", "WifiClients", probeWifiClients},
		{"Wifi", "WifiManagedAP", probeWifiManagedAP},
	} {
		wanted := false

		if len(wantedProbes) == 0 {
			wanted = true
		}

		for _, wantedProbe := range wantedProbes {
			if aProbe.category == wantedProbe || aProbe.name == wantedProbe {
				wanted = true
				break
			}
		}

		if !wanted {
			continue
		}

		m, ok := aProbe.function(c, meta)
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
