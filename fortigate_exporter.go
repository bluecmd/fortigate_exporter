// Server executable of fortigate_exporter
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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	apiKeyFile     = flag.String("api-key-file", "", "file containing the api key map to use when connecting to a Fortigate device")
	listen         = flag.String("listen", ":9710", "address to listen on")
	timeoutSeconds = flag.Int("scrape-timeout", 30, "max seconds to allow a scrape to take")

	apiKeyMap = map[string]string{}
)

type fortiApiClient struct {
	u   url.URL
	hc  *http.Client
	ctx context.Context
}

func (c *fortiApiClient) String() string {
	u := url.URL{
		Scheme: c.u.Scheme,
		Host:   c.u.Host,
	}
	return u.String()
}

func (c *fortiApiClient) Get(path string, obj interface{}) error {
	u := c.u
	u.Path = path

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	req = req.WithContext(c.ctx)
	resp, err := c.hc.Do(req)
	if err != nil {
		// Strip request details to hide sensitive access tokens
		uerr := err.(*url.Error)
		return uerr.Err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Response code was %d, expected 200", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, obj)
}

func newFortiApiClient(ctx context.Context, tgt *url.URL, hc *http.Client) (*fortiApiClient, error) {
	u := url.URL{
		Scheme: tgt.Scheme,
		Host:   tgt.Host,
	}
	q := u.Query()
	tk := u.String()

	ak, ok := apiKeyMap[tk]
	if !ok {
		return nil, fmt.Errorf("No API key registered for %q", tk)
	}

	q.Add("access_token", ak)
	u.RawQuery = q.Encode()
	return &fortiApiClient{u, hc, ctx}, nil
}

func probeSystem(c *fortiApiClient, registry *prometheus.Registry) bool {
	var (
		probeSystemVersion = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fortigate_system_version",
				Help: "Contains the system version and build information",
			},
			[]string{"serial", "version", "build"},
		)
	)

	registry.MustRegister(probeSystemVersion)

	type systemStatus struct {
		Status  string
		Serial  string
		Version string
		Build   int
	}
	var st systemStatus

	if err := c.Get("api/v2/monitor/system/status", &st); err != nil {
		log.Printf("Probe of %q failed. Metric: system status. Error: %q", c, err)
		return false
	}
	probeSystemVersion.WithLabelValues(st.Serial, st.Version, fmt.Sprintf("%d", st.Build)).Set(1)
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

	c, err := newFortiApiClient(ctx, tgt, hc)
	if err != nil {
		return false, err
	}

	success :=
		probeSystem(c, registry) &&
			true
	return success, nil
}

func probeHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	target := params.Get("target")
	if target == "" {
		http.Error(w, "Target parameter missing or empty", http.StatusBadRequest)
		return
	}
	probeSuccessGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "probe_success",
		Help: "Whether or not the probe succeeded",
	})
	probeDurationGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "probe_duration_seconds",
		Help: "How many seconds the probe took to complete",
	})
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(*timeoutSeconds)*time.Second)
	defer cancel()
	registry := prometheus.NewRegistry()
	registry.MustRegister(probeSuccessGauge)
	registry.MustRegister(probeDurationGauge)
	start := time.Now()
	success, err := probe(ctx, target, registry, &http.Client{})
	if err != nil {
		log.Printf("Probe request rejected; error is: %v", err)
		http.Error(w, fmt.Sprintf("probe: %v", err), http.StatusBadRequest)
		return
	}
	duration := time.Since(start).Seconds()
	probeDurationGauge.Set(duration)
	if success {
		probeSuccessGauge.Set(1)
		log.Printf("Probe of %q succeeded, took %.3f seconds", target, duration)
	} else {
		// probeSuccessGauge default is 0
		log.Printf("Probe of %q failed, took %.3f seconds", target, duration)
	}
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func main() {
	flag.Parse()

	akm, err := ioutil.ReadFile(*apiKeyFile)
	if err != nil {
		log.Fatalf("Failed to read API key file: %v", err)
	}

	for _, line := range strings.Split(string(akm), "\n") {
		kv := strings.SplitN(strings.TrimSpace(line), " ", 2)
		if len(kv) != 2 {
			continue
		}
		apiKeyMap[kv[0]] = strings.TrimSpace(kv[1])
	}
	log.Printf("Loaded %d API keys", len(apiKeyMap))

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/probe", probeHandler)
	go http.ListenAndServe(*listen, nil)
	log.Printf("Fortigate exporter running, listening on %q", *listen)
	select {}
}
