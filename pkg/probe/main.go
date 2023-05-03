package probe

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bluecmd/fortigate_exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ProbeHandler(w http.ResponseWriter, r *http.Request) {
	savedConfig := config.GetConfig()

	params := r.URL.Query()
	paramMap := make(map[string]string)
	target := params.Get("target")
	paramMap["target"] = params.Get("target")
	if params.Get("token") != "" {
		paramMap["token"] = params.Get("token")
	}
	if params.Get("profile") != "" {
		paramMap["profile"] = params.Get("profile")
	}

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
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(savedConfig.ScrapeTimeout)*time.Second)
	defer cancel()
	registry := prometheus.NewRegistry()
	registry.MustRegister(probeSuccessGauge)
	registry.MustRegister(probeDurationGauge)
	start := time.Now()
	pc := &ProbeCollector{}
	registry.MustRegister(pc)
	success, err := pc.Probe(ctx, paramMap, &http.Client{}, savedConfig)
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
