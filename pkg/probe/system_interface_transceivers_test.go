package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestSystemInterfaceTransceivers(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/interface/transceivers", "testdata/interface-transceivers.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemInterfaceTransceivers, c, r) {
		t.Errorf("probeSystemInterfaceTransceivers() returned non-success")
	}

	em := `
	# HELP fortigate_interface_transceivers List of transceivers being used by the FortiGate
	# TYPE fortigate_interface_transceivers gauge
	fortigate_interface_transceivers{description="",name="ha1",partnumber="FTLX8574D3BCLFTN",type="SFP/SFP+/SFP28",vendor="FORTINET"} 1
	fortigate_interface_transceivers{description="",name="ha2",partnumber="FTLX8574D3BCLFTN",type="SFP/SFP+/SFP28",vendor="FORTINET"} 1
	fortigate_interface_transceivers{description="",name="port33",partnumber="FTL410QE4CFTN",type="QSFP/QSFP+",vendor="FORTINET"} 1
	fortigate_interface_transceivers{description="",name="port34",partnumber="FTL410QE4CFTN",type="QSFP/QSFP+",vendor="FORTINET"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
