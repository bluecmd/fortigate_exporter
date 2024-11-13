package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestHAPeer(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/ha-peer", "testdata/ha-peer.jsonnet")
	r := prometheus.NewPedanticRegistry()
	meta := &TargetMetadata{
		VersionMajor: 7,
		VersionMinor: 0,
	}
	if !testProbeWithMetadata(probeSystemHAPeer, c, meta, r) {
		t.Errorf("probeSystemHAPeer() returned non-success")
	}

	em := `
	# HELP fortigate_ha_peer_info Info metrics regarding cluster HA peers
	# TYPE fortigate_ha_peer_info gauge
	fortigate_ha_peer_info{hostname="member-name-1",primary="Unsupported",priority="200",serial="FGT61E4QXXXXXXXX1",vcluster_id="0",vdom="root"} 1
	fortigate_ha_peer_info{hostname="member-name-2",primary="Unsupported",priority="100",serial="FGT61E4QXXXXXXXX2",vcluster_id="0",vdom="root"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}

func TestHAPeer74(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/ha-peer", "testdata/ha-peer-74+.jsonnet")
	r := prometheus.NewPedanticRegistry()
	meta := &TargetMetadata{
		VersionMajor: 7,
		VersionMinor: 4,
	}
	if !testProbeWithMetadata(probeSystemHAPeer, c, meta, r) {
		t.Errorf("probeSystemHAPeer() returned non-success")
	}

	em := `
	# HELP fortigate_ha_peer_info Info metrics regarding cluster HA peers
	# TYPE fortigate_ha_peer_info gauge
	fortigate_ha_peer_info{hostname="member-name-1",primary="true",priority="200",serial="FGT61E4QXXXXXXXX1",vcluster_id="0",vdom="root"} 1
	fortigate_ha_peer_info{hostname="member-name-2",primary="false",priority="100",serial="FGT61E4QXXXXXXXX2",vcluster_id="0",vdom="root"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
