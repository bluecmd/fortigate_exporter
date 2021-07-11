package probe

import (
	"fmt"
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

func probeSystemResourceUsage(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		mResCPU = prometheus.NewDesc(
			"fortigate_cpu_usage_ratio",
			"Current resource usage ratio of system CPU, per core",
			[]string{"processor"}, nil,
		)
		mResMemory = prometheus.NewDesc(
			"fortigate_memory_usage_ratio",
			"Current resource usage ratio of system memory",
			[]string{}, nil,
		)
		mResSession = prometheus.NewDesc(
			"fortigate_current_sessions",
			"Current amount of sessions, per IP version",
			[]string{"protocol"}, nil,
		)
	)

	type resUsage struct {
		Current float64
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
		return nil, false
	}

	// CPU[0] is the average over all cores, ignore it
	m := []prometheus.Metric{}
	for i, cpu := range sr.Results.CPU[1:] {
		m = append(m, prometheus.MustNewConstMetric(
			mResCPU, prometheus.GaugeValue, float64(cpu.Current)/100.0, fmt.Sprintf("%d", i)))
	}
	m = append(m, prometheus.MustNewConstMetric(mResMemory, prometheus.GaugeValue, float64(sr.Results.Mem[0].Current)/100.0))
	m = append(m, prometheus.MustNewConstMetric(mResSession, prometheus.GaugeValue, float64(sr.Results.Session[0].Current), "ipv4"))
	m = append(m, prometheus.MustNewConstMetric(mResSession, prometheus.GaugeValue, float64(sr.Results.Session6[0].Current), "ipv6"))
	return m, true
}

func probeSystemVDOMResources(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		mResCPU = prometheus.NewDesc(
			"fortigate_vdom_cpu_usage_ratio",
			"Current resource usage ratio of CPU, per VDOM",
			[]string{"vdom"}, nil,
		)
		mResMemory = prometheus.NewDesc(
			"fortigate_vdom_memory_usage_ratio",
			"Current resource usage ratio of memory, per VDOM",
			[]string{"vdom"}, nil,
		)
		mResSession = prometheus.NewDesc(
			"fortigate_vdom_current_sessions",
			"Current amount of sessions, per VDOM and IP version",
			[]string{"vdom", "protocol"}, nil,
		)
	)

	type resUsage struct {
		Current float64
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
		return nil, false
	}
	m := []prometheus.Metric{}
	for _, s := range sr {
		m = append(m, prometheus.MustNewConstMetric(mResCPU, prometheus.GaugeValue, float64(s.Results.CPU[0].Current)/100.0, s.VDOM))
		m = append(m, prometheus.MustNewConstMetric(mResMemory, prometheus.GaugeValue, float64(s.Results.Mem[0].Current)/100.0, s.VDOM))
		m = append(m, prometheus.MustNewConstMetric(mResSession, prometheus.GaugeValue, float64(s.Results.Session[0].Current), s.VDOM, "ipv4"))
		m = append(m, prometheus.MustNewConstMetric(mResSession, prometheus.GaugeValue, float64(s.Results.Session6[0].Current), s.VDOM, "ipv6"))
	}
	return m, true
}
