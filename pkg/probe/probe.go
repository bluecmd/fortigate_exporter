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
	"net/http"
	"net/url"

	"github.com/bluecmd/fortigate_exporter/internal/config"
	fortiHTTP "github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type ProbeCollector struct {
	metrics []prometheus.Metric
}

type probeFunc func(fortiHTTP.FortiHTTP) ([]prometheus.Metric, bool)

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

	// TODO: Make parallel
	success := true
	for _, f := range []probeFunc{
		probeSystemStatus,
		probeSystemResourceUsage,
		probeSystemVDOMResources,
		probeFirewallPolicies,
		probeSystemInterface,
		probeVPNSsl,
		probeVPNIPSec,
		probeSystemHAStatistics,
		probeLicenseStatus,
		probeSystemLinkMonitor,
		probeVirtualWANHealthCheck,
		probeSystemAvailableCertificates,
		probeFirewallLoadBalance,
		probeBgpPeers,
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
