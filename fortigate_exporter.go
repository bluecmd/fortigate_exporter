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
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
)

var (
	authMapFile    = flag.String("auth-file", "", "file containing the authentication map to use when connecting to a Fortigate device")
	listen         = flag.String("listen", ":9710", "address to listen on")
	timeoutSeconds = flag.Int("scrape-timeout", 30, "max seconds to allow a scrape to take")

	authMap = map[string]Auth{}
)

type Auth struct {
	Token string
}

type FortiHTTP interface {
	Get(path string, query string, obj interface{}) error
}

func newFortiClient(ctx context.Context, tgt url.URL, hc *http.Client) (FortiHTTP, error) {
	auth, ok := authMap[tgt.String()]
	if !ok {
		return nil, fmt.Errorf("No API authentication registered for %q", tgt.String())
	}

	if auth.Token != "" {
		if tgt.Scheme != "https" {
			return nil, fmt.Errorf("FortiOS only supports token for HTTPS connections")
		}
		c, err := newFortiTokenClient(ctx, tgt, hc, auth.Token)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
	return nil, fmt.Errorf("Invalid authentication data for %q", tgt.String())
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

	af, err := ioutil.ReadFile(*authMapFile)
	if err != nil {
		log.Fatalf("Failed to read API authentication map file: %v", err)
	}

	if err := yaml.Unmarshal(af, &authMap); err != nil {
		log.Fatalf("Failed to parse API authentication map file: %v", err)
	}

	log.Printf("Loaded %d API keys", len(authMap))

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/probe", probeHandler)
	go http.ListenAndServe(*listen, nil)
	log.Printf("Fortigate exporter running, listening on %q", *listen)
	select {}
}
