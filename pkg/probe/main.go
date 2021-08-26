package probe

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bluecmd/fortigate_exporter/internal/logging"

	"github.com/bluecmd/fortigate_exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ProbeHandler(w http.ResponseWriter, r *http.Request) {
	savedConfig := config.GetConfig()

	params := r.URL.Query()
	target := params.Get("target")
	if target == "" {
		http.Error(w, "Target parameter missing or empty", http.StatusBadRequest)
		return
	}

	log := logging.GetSugar().With("target", target)

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
	success, err := pc.Probe(ctx, target, &http.Client{}, savedConfig, log)
	if err != nil {
		log.Errorf("Probe request rejected; error is: %v", err)
		http.Error(w, fmt.Sprintf("probe: %v", err), http.StatusBadRequest)
		return
	}
	duration := time.Since(start)
	probeDurationGauge.Set(duration.Seconds())
	if success {
		probeSuccessGauge.Set(1)
		log.Debugw("Probe succeeded", "duration", duration)
	} else {
		// probeSuccessGauge default is 0
		log.Errorw("Probe failed", "duration", duration)
	}
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}
