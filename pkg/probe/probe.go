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
	"strings"

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
	name     string
	function probeFunc
}

func (p *ProbeCollector) Probe(ctx context.Context, target map[string]string, hc *http.Client, savedConfig config.FortiExporterConfig) (bool, error) {
	tgt, err := url.Parse(target["target"])
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

	if target["token"] != "" && savedConfig.AuthKeys[config.Target(target["target"])].Token == "" {
		// Add the target and its apikey to the savedConfig and use, if exists, a target entry as a template for include/exclude
		// This will only happend the "first" time
		savedConfig.AuthKeys[config.Target(target["target"])] = config.TargetAuth{Token: config.Token(target["token"]),
			Probes: savedConfig.AuthKeys[config.Target(target["profile"])].Probes}
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

	includedProbes := savedConfig.AuthKeys[config.Target(u.String())].Probes.Include
	excludedProbes := savedConfig.AuthKeys[config.Target(u.String())].Probes.Exclude

	// TODO: Make parallel
	success := true
	for _, aProbe := range []probeDetailedFunc{
		// Always keep probeSystemTime on top of the list to have the probe processed first.
		// Therefore time returned is more accurate when integrated in Prometheus because
		// timestamp for the metrics probe, in Prometheus, is obtained from the query time, not the reply time.
		// This is especially important when running all the probes takes many seconds.
		{"System/Time/Clock", probeSystemTime},
		{"BGP/NeighborPaths/IPv4", probeBGPNeighborPathsIPv4},
		{"BGP/NeighborPaths/IPv6", probeBGPNeighborPathsIPv6},
		{"BGP/Neighbors/IPv4", probeBGPNeighborsIPv4},
		{"BGP/Neighbors/IPv6", probeBGPNeighborsIPv6},
		{"Firewall/LoadBalance", probeFirewallLoadBalance},
		{"Firewall/Policies", probeFirewallPolicies},
		{"Firewall/IpPool", probeFirewallIpPool},
		{"License/Status", probeLicenseStatus},
		{"Log/Fortianalyzer/Status", probeLogAnalyzer},
		{"Log/Fortianalyzer/Queue", probeLogAnalyzerQueue},
		{"Log/DiskUsage", probeLogCurrentDiskUsage},
		{"System/AvailableCertificates", probeSystemAvailableCertificates},
		{"System/Fortimanager/Status", probeSystemFortimanagerStatus},
		{"System/HAStatistics", probeSystemHAStatistics},
		{"System/Interface", probeSystemInterface},
		{"System/LinkMonitor", probeSystemLinkMonitor},
		{"System/Resource/Usage", probeSystemResourceUsage},
		{"System/SDNConnector", probeSystemSDNConnector},
		{"System/SensorInfo", probeSystemSensorInfo},
		{"System/Status", probeSystemStatus},
		{"System/VDOMResources", probeSystemVDOMResources},
		{"System/HAChecksum", probeSystemHAChecksum},
		{"User/Fsso", probeUserFsso},
		{"VPN/IPSec", probeVPNIPSec},
		{"VPN/Ssl/Connections", probeVPNSsl},
		{"VPN/Ssl/Stats", probeVPNSslStats},
		{"VirtualWAN/HealthCheck", probeVirtualWANHealthCheck},
		{"WebUI/State", probeWebUIState},
		{"Wifi/APStatus", probeWifiAPStatus},
		{"Wifi/Clients", probeWifiClients},
		{"Wifi/ManagedAP", probeWifiManagedAP},
		{"Switch/ManagedSwitch", probeManagedSwitch},
		{"OSPF/Neighbors", probeOSPFNeighbors},
	} {
		wanted := false

		if len(includedProbes) == 0 {
			wanted = true
		} else {
			for _, wantedProbe := range includedProbes {
				if strings.HasPrefix(aProbe.name, wantedProbe) {
					wanted = true
					break
				}
			}
		}

		if len(excludedProbes) != 0 {
			for _, unwantedProbe := range excludedProbes {
				if strings.HasPrefix(aProbe.name, unwantedProbe) {
					wanted = false
					break
				}
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
