package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestHAChecksum(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/ha-checksums", "testdata/ha-checksum.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemHAChecksum, c, r) {
		t.Errorf("probeSystemHAChecksum() returned non-success")
	}

	em := `
	# HELP fortigate_ha_role Master/Slave information
	# TYPE fortigate_ha_role gauge
	fortigate_ha_role{name="is_manage_master", serial_no="SERIAL111111111"} 1
	fortigate_ha_role{name="is_manage_master", serial_no="SERIAL222222222"} 0
	fortigate_ha_role{name="is_root_master", serial_no="SERIAL111111111"} 1
	fortigate_ha_role{name="is_root_master", serial_no="SERIAL222222222"} 0
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
