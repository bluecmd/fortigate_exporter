package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestSwitchHealth(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/switch-controller/managed-switch/health", "testdata/fsw-health.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSwitchHealth, c, r) {
		t.Errorf("probeSwitchHealth() returned non-success")
	}
        
	em := `
	# HELP fortiswitch_health_performance_stats_cpu_idle Fortiswitch CPU idle
	# TYPE fortiswitch_health_performance_stats_cpu_idle gauge
	fortiswitch_health_performance_stats_cpu_idle{VDOM="root",fortiswitch="FS00000000000024",unit="%%"} 100
	fortiswitch_health_performance_stats_cpu_idle{VDOM="root",fortiswitch="FS00000000000027",unit="%%"} 100
	fortiswitch_health_performance_stats_cpu_idle{VDOM="root",fortiswitch="FS00000000000030",unit="%%"} 100
	fortiswitch_health_performance_stats_cpu_idle{VDOM="root",fortiswitch="FS00000000000038",unit="%%"} 100
	# HELP fortiswitch_health_performance_stats_cpu_nice Fortiswitch CPU nice usage
	# TYPE fortiswitch_health_performance_stats_cpu_nice gauge
	fortiswitch_health_performance_stats_cpu_nice{VDOM="root",fortiswitch="FS00000000000024",unit="%%"} 0
	fortiswitch_health_performance_stats_cpu_nice{VDOM="root",fortiswitch="FS00000000000027",unit="%%"} 0
	fortiswitch_health_performance_stats_cpu_nice{VDOM="root",fortiswitch="FS00000000000030",unit="%%"} 0
	fortiswitch_health_performance_stats_cpu_nice{VDOM="root",fortiswitch="FS00000000000038",unit="%%"} 0
	# HELP fortiswitch_health_performance_stats_cpu_system Fortiswitch CPU system usage
	# TYPE fortiswitch_health_performance_stats_cpu_system gauge
	fortiswitch_health_performance_stats_cpu_system{VDOM="root",fortiswitch="FS00000000000024",unit="%%"} 0
	fortiswitch_health_performance_stats_cpu_system{VDOM="root",fortiswitch="FS00000000000027",unit="%%"} 0
	fortiswitch_health_performance_stats_cpu_system{VDOM="root",fortiswitch="FS00000000000030",unit="%%"} 0
	fortiswitch_health_performance_stats_cpu_system{VDOM="root",fortiswitch="FS00000000000038",unit="%%"} 0
	# HELP fortiswitch_health_performance_stats_cpu_user Fortiswitch CPU user usage
	# TYPE fortiswitch_health_performance_stats_cpu_user gauge
	fortiswitch_health_performance_stats_cpu_user{VDOM="root",fortiswitch="FS00000000000024",unit="%%"} 0
	fortiswitch_health_performance_stats_cpu_user{VDOM="root",fortiswitch="FS00000000000027",unit="%%"} 0
	fortiswitch_health_performance_stats_cpu_user{VDOM="root",fortiswitch="FS00000000000030",unit="%%"} 0
	fortiswitch_health_performance_stats_cpu_user{VDOM="root",fortiswitch="FS00000000000038",unit="%%"} 0
	# HELP fortiswitch_health_summary_cpu Summary CPU health
	# TYPE fortiswitch_health_summary_cpu gauge
	fortiswitch_health_summary_cpu{VDOM="root",fortiswitch="FS00000000000024",rating="good"} 0
	fortiswitch_health_summary_cpu{VDOM="root",fortiswitch="FS00000000000027",rating="good"} 0
	fortiswitch_health_summary_cpu{VDOM="root",fortiswitch="FS00000000000030",rating="good"} 0
	fortiswitch_health_summary_cpu{VDOM="root",fortiswitch="FS00000000000038",rating="good"} 0
	# HELP fortiswitch_health_summary_mem Summary MEM health
	# TYPE fortiswitch_health_summary_mem gauge
	fortiswitch_health_summary_mem{VDOM="root",fortiswitch="FS00000000000024",rating="good"} 10
	fortiswitch_health_summary_mem{VDOM="root",fortiswitch="FS00000000000027",rating="good"} 15
	fortiswitch_health_summary_mem{VDOM="root",fortiswitch="FS00000000000030",rating="good"} 50
	fortiswitch_health_summary_mem{VDOM="root",fortiswitch="FS00000000000038",rating="good"} 32
	# HELP fortiswitch_health_summary_temp Summary Temperature health
	# TYPE fortiswitch_health_summary_temp gauge
	fortiswitch_health_summary_temp{VDOM="root",fortiswitch="FS00000000000024",rating="good"} 48.952749999999995
	fortiswitch_health_summary_temp{VDOM="root",fortiswitch="FS00000000000027",rating="good"} 46.156000000000006
	fortiswitch_health_summary_temp{VDOM="root",fortiswitch="FS00000000000030",rating="good"} 39.71875
	fortiswitch_health_summary_temp{VDOM="root",fortiswitch="FS00000000000038",rating="good"} 41.624750000000006
	# HELP fortiswitch_health_summary_uptime Summary Uptime
	# TYPE fortiswitch_health_summary_uptime gauge
	fortiswitch_health_summary_uptime{VDOM="root",fortiswitch="FS00000000000024",rating="good"} 3.928968e+07
	fortiswitch_health_summary_uptime{VDOM="root",fortiswitch="FS00000000000027",rating="good"} 3.928974e+07
	fortiswitch_health_summary_uptime{VDOM="root",fortiswitch="FS00000000000030",rating="good"} 2.661288e+07
	fortiswitch_health_summary_uptime{VDOM="root",fortiswitch="FS00000000000038",rating="good"} 2.661258e+07
	# HELP fortiswitch_health_temperature Temperature per switch sensor
	# TYPE fortiswitch_health_temperature gauge
	fortiswitch_health_temperature{VDOM="root",fortiswitch="FS00000000000024",module="sensor1(CPU  Board Temp)",unit="celsius"} 41.937
	fortiswitch_health_temperature{VDOM="root",fortiswitch="FS00000000000024",module="sensor2(MAIN Board Temp1)",unit="celsius"} 63.875
	fortiswitch_health_temperature{VDOM="root",fortiswitch="FS00000000000024",module="sensor3(MAIN Board Temp2)",unit="celsius"} 51.312
	fortiswitch_health_temperature{VDOM="root",fortiswitch="FS00000000000024",module="sensor4(MAIN Board Temp3)",unit="celsius"} 38.687
	fortiswitch_health_temperature{VDOM="root",fortiswitch="FS00000000000027",module="sensor1(CPU  Board Temp)",unit="celsius"} 39
	fortiswitch_health_temperature{VDOM="root",fortiswitch="FS00000000000027",module="sensor2(MAIN Board Temp1)",unit="celsius"} 60.625
	fortiswitch_health_temperature{VDOM="root",fortiswitch="FS00000000000027",module="sensor3(MAIN Board Temp2)",unit="celsius"} 48.937
	fortiswitch_health_temperature{VDOM="root",fortiswitch="FS00000000000027",module="sensor4(MAIN Board Temp3)",unit="celsius"} 36.062
	fortiswitch_health_temperature{VDOM="root",fortiswitch="FS00000000000030",module="sensor1(CPU  Board Temp)",unit="celsius"} 33.875
	fortiswitch_health_temperature{VDOM="root",fortiswitch="FS00000000000030",module="sensor2(MAIN Board Temp1)",unit="celsius"} 53.75
	fortiswitch_health_temperature{VDOM="root",fortiswitch="FS00000000000030",module="sensor3(MAIN Board Temp2)",unit="celsius"} 41
	fortiswitch_health_temperature{VDOM="root",fortiswitch="FS00000000000030",module="sensor4(MAIN Board Temp3)",unit="celsius"} 30.25
	fortiswitch_health_temperature{VDOM="root",fortiswitch="FS00000000000038",module="sensor1(CPU  Board Temp)",unit="celsius"} 35.437
	fortiswitch_health_temperature{VDOM="root",fortiswitch="FS00000000000038",module="sensor2(MAIN Board Temp1)",unit="celsius"} 55.625
	fortiswitch_health_temperature{VDOM="root",fortiswitch="FS00000000000038",module="sensor3(MAIN Board Temp2)",unit="celsius"} 43.125
	fortiswitch_health_temperature{VDOM="root",fortiswitch="FS00000000000038",module="sensor4(MAIN Board Temp3)",unit="celsius"} 32.312
        `
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
