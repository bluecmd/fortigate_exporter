package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"

)

func probeSwitchHealth(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		mSumCPU = prometheus.NewDesc(
			"fortiswitch_health_summary_cpu",
			"Summary CPU health",
			[]string{"rating", "fortiswitch", "VDOM"}, nil,
		)
		mSumMem = prometheus.NewDesc(
			"fortiswitch_health_summary_mem",
			"Summary MEM health",
			[]string{"rating", "fortiswitch", "VDOM"}, nil,
		)
		mSumUpTime = prometheus.NewDesc(
			"fortiswitch_health_summary_uptime",
			"Summary Uptime",
			[]string{"rating", "fortiswitch", "VDOM"}, nil,
		)
		mSumTemp = prometheus.NewDesc(
			"fortiswitch_health_summary_temp",
			"Summary Temperature health",
			[]string{"rating", "fortiswitch", "VDOM"}, nil,
		)
		mTemp = prometheus.NewDesc(
			"fortiswitch_health_temperature",
			"Temperature per switch sensor",
			[]string{"unit", "module", "fortiswitch", "VDOM"}, nil,
		)
		mCpuUser = prometheus.NewDesc(
			"fortiswitch_health_performance_stats_cpu_user",
			"Fortiswitch CPU user usage",
			[]string{"unit", "fortiswitch", "VDOM"}, nil,
		)
		mCpuSystem = prometheus.NewDesc(
			"fortiswitch_health_performance_stats_cpu_system",
			"Fortiswitch CPU system usage",
			[]string{"unit", "fortiswitch", "VDOM"}, nil,
		)
		mCpuIdle = prometheus.NewDesc(
			"fortiswitch_health_performance_stats_cpu_idle",
			"Fortiswitch CPU idle",
			[]string{"unit", "fortiswitch", "VDOM"}, nil,
		)
		mCpuNice = prometheus.NewDesc(
			"fortiswitch_health_performance_stats_cpu_nice",
			"Fortiswitch CPU nice usage",
			[]string{"unit", "fortiswitch", "VDOM"}, nil,
		)
	)
	type Sum struct {
		Value  	float64	 `json:"value"`
		Rating	 string	 `json:"rating"`
	}
	type Status struct {
		Value float64 `json:"value"`
		Unit  string  `json:"unit"`
	}	
	type Uptime struct {
		Days    Status   `json:"days"`
		Hours   Status	 `json:"hours"`
		Minutes Status	 `json:"minutes"`
	}
	type Network struct {
		In1Min  Status `json:"in-1min"`
		In10Min Status `json:"in-10min"`
		In30Min Status `json:"in-30min"`
	}
	type Memory struct {
		Used Status `json:"used"`
	}
	type CPU struct {
		User   Status `json:"user"`
		System Status `json:"system"`
		Nice   Status `json:"nice"`
		Idle   Status `json:"idle"`
	}
	type PerformanceStatus struct {
		CPU     CPU     `json:"cpu"`
		Memory  Memory  `json:"memory"`
		Network Network `json:"network"`
		Uptime  Uptime  `json:"uptime"`
	}
	type Temperature struct {
		Module	string
		Status 	Status
	}
	type Summary struct {
		Overall	string	 `json:"overall"`
		CPU	Sum
		Memory	Sum
		Uptime	Sum
		Temperature	Sum
	}
	type Poe struct {
		Value    int    `json:"value"`
		MaxValue int    `json:"max_value"`
		Unit     string `json:"unit"`
	}
	type Results struct {
		PerformanceStatus PerformanceStatus `json:"performance-status"`
		Temperature       []Temperature     `json:"temperature"`
		Summary           Summary           `json:"summary"`
		Poe               Poe               `json:"poe"`
	}
	
	type swResponse struct {
		Results		map[string]Results	`json:"results"`
		VDOM 		string
	}
	var r swResponse
	//var r map[string]swResponse
	//var r []swResponse

	if err := c.Get("api/v2/monitor/switch-controller/managed-switch/health", "vdom=root", &r); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}
	m := []prometheus.Metric{}
	//for _, sw := range r {
		for fswitch, hr := range r.Results {
			
			m = append(m, prometheus.MustNewConstMetric(mSumCPU, prometheus.GaugeValue, hr.Summary.CPU.Value, hr.Summary.CPU.Rating, fswitch, r.VDOM))
			m = append(m, prometheus.MustNewConstMetric(mSumMem, prometheus.GaugeValue, hr.Summary.Memory.Value, hr.Summary.Memory.Rating, fswitch, r.VDOM))
			m = append(m, prometheus.MustNewConstMetric(mSumUpTime, prometheus.GaugeValue, hr.Summary.Uptime.Value, hr.Summary.Uptime.Rating, fswitch, r.VDOM))
			m = append(m, prometheus.MustNewConstMetric(mSumTemp, prometheus.GaugeValue, hr.Summary.Temperature.Value, hr.Summary.Temperature.Rating, fswitch, r.VDOM))

			for _, ts := range hr.Temperature {
				m = append(m, prometheus.MustNewConstMetric(mTemp, prometheus.GaugeValue, ts.Status.Value, ts.Status.Unit, ts.Module, fswitch, r.VDOM))
                        }
		
			CpuUnit := hr.PerformanceStatus.CPU.System.Unit
			/*if CpuUnit == "%" {
				CpuUnit = "%%"
			}*/

			m = append(m, prometheus.MustNewConstMetric(mCpuUser, prometheus.GaugeValue, hr.PerformanceStatus.CPU.User.Value, CpuUnit, fswitch, r.VDOM))
			m = append(m, prometheus.MustNewConstMetric(mCpuNice, prometheus.GaugeValue, hr.PerformanceStatus.CPU.Nice.Value, CpuUnit, fswitch, r.VDOM))
			m = append(m, prometheus.MustNewConstMetric(mCpuSystem, prometheus.GaugeValue, hr.PerformanceStatus.CPU.System.Value, CpuUnit, fswitch, r.VDOM))
			m = append(m, prometheus.MustNewConstMetric(mCpuIdle, prometheus.GaugeValue, hr.PerformanceStatus.CPU.Idle.Value, CpuUnit, fswitch, r.VDOM))
		}
	//}
	return m, true
}
