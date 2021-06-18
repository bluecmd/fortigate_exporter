package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestHAStatistics(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/ha-statistics", "testdata/ha-statistics.jsonnet")
	c.prepare("api/v2/cmdb/system/ha", "testdata/ha-config.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemHAStatistics, c, r) {
		t.Errorf("probeSystemHAStatistics() returned non-success")
	}

	em := `
        # HELP fortigate_ha_member_bytes_total Bytes transferred by HA member
        # TYPE fortigate_ha_member_bytes_total counter
        fortigate_ha_member_bytes_total{hostname="member-name-1",vdom="root"} 2.02844842379e+11
        fortigate_ha_member_bytes_total{hostname="member-name-2",vdom="root"} 40
        # HELP fortigate_ha_member_cpu_usage_ratio CPU usage by HA member
        # TYPE fortigate_ha_member_cpu_usage_ratio gauge
        fortigate_ha_member_cpu_usage_ratio{hostname="member-name-1",vdom="root"} 0.01
        fortigate_ha_member_cpu_usage_ratio{hostname="member-name-2",vdom="root"} 0
        # HELP fortigate_ha_member_info Info metric regarding cluster members
        # TYPE fortigate_ha_member_info gauge
        fortigate_ha_member_info{group="my-cluster",hostname="member-name-1",serial="FGT61E4QXXXXXXXX1",vdom="root"} 1
        fortigate_ha_member_info{group="my-cluster",hostname="member-name-2",serial="FGT61E4QXXXXXXXX2",vdom="root"} 1
        # HELP fortigate_ha_member_ips_events_total IPS events processed by HA member
        # TYPE fortigate_ha_member_ips_events_total counter
        fortigate_ha_member_ips_events_total{hostname="member-name-1",vdom="root"} 0
        fortigate_ha_member_ips_events_total{hostname="member-name-2",vdom="root"} 0
        # HELP fortigate_ha_member_memory_usage_ratio Memory usage by HA member
        # TYPE fortigate_ha_member_memory_usage_ratio gauge
        fortigate_ha_member_memory_usage_ratio{hostname="member-name-1",vdom="root"} 0.67
        fortigate_ha_member_memory_usage_ratio{hostname="member-name-2",vdom="root"} 0.68
        # HELP fortigate_ha_member_network_usage_ratio Network usage by HA member
        # TYPE fortigate_ha_member_network_usage_ratio gauge
        fortigate_ha_member_network_usage_ratio{hostname="member-name-1",vdom="root"} 1.52
        fortigate_ha_member_network_usage_ratio{hostname="member-name-2",vdom="root"} 0.43
        # HELP fortigate_ha_member_packets_total Packets which are handled by this HA member
        # TYPE fortigate_ha_member_packets_total counter
        fortigate_ha_member_packets_total{hostname="member-name-1",vdom="root"} 5.49981862e+08
        fortigate_ha_member_packets_total{hostname="member-name-2",vdom="root"} 1
        # HELP fortigate_ha_member_sessions Sessions which are handled by this HA member
        # TYPE fortigate_ha_member_sessions gauge
        fortigate_ha_member_sessions{hostname="member-name-1",vdom="root"} 148
        fortigate_ha_member_sessions{hostname="member-name-2",vdom="root"} 12
        # HELP fortigate_ha_member_virus_events_total Virus events which are detected by this HA member
        # TYPE fortigate_ha_member_virus_events_total counter
        fortigate_ha_member_virus_events_total{hostname="member-name-1",vdom="root"} 0
        fortigate_ha_member_virus_events_total{hostname="member-name-2",vdom="root"} 0
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}

func TestHAStatisticsNoConfigAccess(t *testing.T) {
	// The only difference here to TestHAStatistics is that the "group" label
	// is empty in fortigate_ha_member_info.
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/ha-statistics", "testdata/ha-statistics.jsonnet")
	c.prepare("api/v2/cmdb/system/ha", "testdata/ha-config-no-access.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemHAStatistics, c, r) {
		t.Errorf("probeSystemHAStatistics() returned non-success")
	}

	em := `
        # HELP fortigate_ha_member_bytes_total Bytes transferred by HA member
        # TYPE fortigate_ha_member_bytes_total counter
        fortigate_ha_member_bytes_total{hostname="member-name-1",vdom="root"} 2.02844842379e+11
        fortigate_ha_member_bytes_total{hostname="member-name-2",vdom="root"} 40
        # HELP fortigate_ha_member_cpu_usage_ratio CPU usage by HA member
        # TYPE fortigate_ha_member_cpu_usage_ratio gauge
        fortigate_ha_member_cpu_usage_ratio{hostname="member-name-1",vdom="root"} 0.01
        fortigate_ha_member_cpu_usage_ratio{hostname="member-name-2",vdom="root"} 0
        # HELP fortigate_ha_member_info Info metric regarding cluster members
        # TYPE fortigate_ha_member_info gauge
        fortigate_ha_member_info{group="",hostname="member-name-1",serial="FGT61E4QXXXXXXXX1",vdom="root"} 1
        fortigate_ha_member_info{group="",hostname="member-name-2",serial="FGT61E4QXXXXXXXX2",vdom="root"} 1
        # HELP fortigate_ha_member_ips_events_total IPS events processed by HA member
        # TYPE fortigate_ha_member_ips_events_total counter
        fortigate_ha_member_ips_events_total{hostname="member-name-1",vdom="root"} 0
        fortigate_ha_member_ips_events_total{hostname="member-name-2",vdom="root"} 0
        # HELP fortigate_ha_member_memory_usage_ratio Memory usage by HA member
        # TYPE fortigate_ha_member_memory_usage_ratio gauge
        fortigate_ha_member_memory_usage_ratio{hostname="member-name-1",vdom="root"} 0.67
        fortigate_ha_member_memory_usage_ratio{hostname="member-name-2",vdom="root"} 0.68
        # HELP fortigate_ha_member_network_usage_ratio Network usage by HA member
        # TYPE fortigate_ha_member_network_usage_ratio gauge
        fortigate_ha_member_network_usage_ratio{hostname="member-name-1",vdom="root"} 1.52
        fortigate_ha_member_network_usage_ratio{hostname="member-name-2",vdom="root"} 0.43
        # HELP fortigate_ha_member_packets_total Packets which are handled by this HA member
        # TYPE fortigate_ha_member_packets_total counter
        fortigate_ha_member_packets_total{hostname="member-name-1",vdom="root"} 5.49981862e+08
        fortigate_ha_member_packets_total{hostname="member-name-2",vdom="root"} 1
        # HELP fortigate_ha_member_sessions Sessions which are handled by this HA member
        # TYPE fortigate_ha_member_sessions gauge
        fortigate_ha_member_sessions{hostname="member-name-1",vdom="root"} 148
        fortigate_ha_member_sessions{hostname="member-name-2",vdom="root"} 12
        # HELP fortigate_ha_member_virus_events_total Virus events which are detected by this HA member
        # TYPE fortigate_ha_member_virus_events_total counter
        fortigate_ha_member_virus_events_total{hostname="member-name-1",vdom="root"} 0
        fortigate_ha_member_virus_events_total{hostname="member-name-2",vdom="root"} 0
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
