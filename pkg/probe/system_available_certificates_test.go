package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestCertificates(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/available-certificates?scope=global", "testdata/available-certificates-scope-global.jsonnet")
	c.prepare("api/v2/monitor/system/available-certificates?vdom=*", "testdata/available-certificates-vdom.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemAvailableCertificates, c, r) {
		t.Errorf("testCertificates() returned non-success")
	}

	em := `
        # HELP fortigate_certificate_cmdb_references Number of times the certificate is referenced
        # TYPE fortigate_certificate_cmdb_references gauge
        fortigate_certificate_cmdb_references{name="Fortinet_CA_SSL",scope="global",source="factory",vdom="root"} 0
        fortigate_certificate_cmdb_references{name="Fortinet_CA_SSL",scope="vdom",source="factory",vdom="root"} 5
        fortigate_certificate_cmdb_references{name="Fortinet_CA_Untrusted",scope="global",source="factory",vdom="root"} 0
        fortigate_certificate_cmdb_references{name="Fortinet_Factory",scope="global",source="factory",vdom="root"} 4
        fortigate_certificate_cmdb_references{name="Fortinet_SSL",scope="global",source="factory",vdom="root"} 0
        fortigate_certificate_cmdb_references{name="Fortinet_SSL_DSA1024",scope="global",source="factory",vdom="root"} 0
        fortigate_certificate_cmdb_references{name="Fortinet_SSL_DSA2048",scope="global",source="factory",vdom="root"} 0
        fortigate_certificate_cmdb_references{name="Fortinet_SSL_ECDSA256",scope="global",source="factory",vdom="root"} 0
        fortigate_certificate_cmdb_references{name="Fortinet_SSL_ECDSA384",scope="global",source="factory",vdom="root"} 0
        fortigate_certificate_cmdb_references{name="Fortinet_SSL_ECDSA521",scope="global",source="factory",vdom="root"} 0
        fortigate_certificate_cmdb_references{name="Fortinet_SSL_ED25519",scope="global",source="factory",vdom="root"} 0
        fortigate_certificate_cmdb_references{name="Fortinet_SSL_ED448",scope="global",source="factory",vdom="root"} 0
        fortigate_certificate_cmdb_references{name="Fortinet_SSL_RSA1024",scope="global",source="factory",vdom="root"} 0
        fortigate_certificate_cmdb_references{name="Fortinet_SSL_RSA2048",scope="global",source="factory",vdom="root"} 0
        fortigate_certificate_cmdb_references{name="Fortinet_SSL_RSA4096",scope="global",source="factory",vdom="root"} 0
        fortigate_certificate_cmdb_references{name="Fortinet_Wifi",scope="global",source="factory",vdom="root"} 1
        # HELP fortigate_certificate_info Info metric containing meta information about the certificate
        # TYPE fortigate_certificate_info gauge
        fortigate_certificate_info{name="Fortinet_CA_SSL",scope="global",source="factory",status="valid",type="local-ca",vdom="root"} 1
        fortigate_certificate_info{name="Fortinet_CA_SSL",scope="vdom",source="factory",status="valid",type="local-ca",vdom="root"} 1
        fortigate_certificate_info{name="Fortinet_CA_Untrusted",scope="global",source="factory",status="valid",type="local-ca",vdom="root"} 1
        fortigate_certificate_info{name="Fortinet_Factory",scope="global",source="factory",status="valid",type="local-cer",vdom="root"} 1
        fortigate_certificate_info{name="Fortinet_SSL",scope="global",source="factory",status="valid",type="local-cer",vdom="root"} 1
        fortigate_certificate_info{name="Fortinet_SSL_DSA1024",scope="global",source="factory",status="valid",type="local-cer",vdom="root"} 1
        fortigate_certificate_info{name="Fortinet_SSL_DSA2048",scope="global",source="factory",status="valid",type="local-cer",vdom="root"} 1
        fortigate_certificate_info{name="Fortinet_SSL_ECDSA256",scope="global",source="factory",status="valid",type="local-cer",vdom="root"} 1
        fortigate_certificate_info{name="Fortinet_SSL_ECDSA384",scope="global",source="factory",status="valid",type="local-cer",vdom="root"} 1
        fortigate_certificate_info{name="Fortinet_SSL_ECDSA521",scope="global",source="factory",status="valid",type="local-cer",vdom="root"} 1
        fortigate_certificate_info{name="Fortinet_SSL_ED25519",scope="global",source="factory",status="valid",type="local-cer",vdom="root"} 1
        fortigate_certificate_info{name="Fortinet_SSL_ED448",scope="global",source="factory",status="valid",type="local-cer",vdom="root"} 1
        fortigate_certificate_info{name="Fortinet_SSL_RSA1024",scope="global",source="factory",status="valid",type="local-cer",vdom="root"} 1
        fortigate_certificate_info{name="Fortinet_SSL_RSA2048",scope="global",source="factory",status="valid",type="local-cer",vdom="root"} 1
        fortigate_certificate_info{name="Fortinet_SSL_RSA4096",scope="global",source="factory",status="valid",type="local-cer",vdom="root"} 1
        fortigate_certificate_info{name="Fortinet_Wifi",scope="global",source="factory",status="valid",type="local-cer",vdom="root"} 1
        # HELP fortigate_certificate_valid_from_seconds Unix timestamp from which this certificate is valid
        # TYPE fortigate_certificate_valid_from_seconds gauge
        fortigate_certificate_valid_from_seconds{name="Fortinet_CA_SSL",scope="global",source="factory",vdom="root"} 1.472285182e+09
        fortigate_certificate_valid_from_seconds{name="Fortinet_CA_SSL",scope="vdom",source="factory",vdom="root"} 1.472285182e+09
        fortigate_certificate_valid_from_seconds{name="Fortinet_CA_Untrusted",scope="global",source="factory",vdom="root"} 1.472285185e+09
        fortigate_certificate_valid_from_seconds{name="Fortinet_Factory",scope="global",source="factory",vdom="root"} 1.468370862e+09
        fortigate_certificate_valid_from_seconds{name="Fortinet_SSL",scope="global",source="factory",vdom="root"} 1.47228519e+09
        fortigate_certificate_valid_from_seconds{name="Fortinet_SSL_DSA1024",scope="global",source="factory",vdom="root"} 1.51007442e+09
        fortigate_certificate_valid_from_seconds{name="Fortinet_SSL_DSA2048",scope="global",source="factory",vdom="root"} 1.510074429e+09
        fortigate_certificate_valid_from_seconds{name="Fortinet_SSL_ECDSA256",scope="global",source="factory",vdom="root"} 1.510074429e+09
        fortigate_certificate_valid_from_seconds{name="Fortinet_SSL_ECDSA384",scope="global",source="factory",vdom="root"} 1.510074429e+09
        fortigate_certificate_valid_from_seconds{name="Fortinet_SSL_ECDSA521",scope="global",source="factory",vdom="root"} 1.582830187e+09
        fortigate_certificate_valid_from_seconds{name="Fortinet_SSL_ED25519",scope="global",source="factory",vdom="root"} 1.582830187e+09
        fortigate_certificate_valid_from_seconds{name="Fortinet_SSL_ED448",scope="global",source="factory",vdom="root"} 1.582830187e+09
        fortigate_certificate_valid_from_seconds{name="Fortinet_SSL_RSA1024",scope="global",source="factory",vdom="root"} 1.510074404e+09
        fortigate_certificate_valid_from_seconds{name="Fortinet_SSL_RSA2048",scope="global",source="factory",vdom="root"} 1.510074417e+09
        fortigate_certificate_valid_from_seconds{name="Fortinet_SSL_RSA4096",scope="global",source="factory",vdom="root"} 1.582830187e+09
        fortigate_certificate_valid_from_seconds{name="Fortinet_Wifi",scope="global",source="factory",vdom="root"} 1.606176e+09
        # HELP fortigate_certificate_valid_to_seconds Unix timestamp till which this certificate is valid
        # TYPE fortigate_certificate_valid_to_seconds gauge
        fortigate_certificate_valid_to_seconds{name="Fortinet_CA_SSL",scope="global",source="factory",vdom="root"} 1.787904382e+09
        fortigate_certificate_valid_to_seconds{name="Fortinet_CA_SSL",scope="vdom",source="factory",vdom="root"} 1.787904382e+09
        fortigate_certificate_valid_to_seconds{name="Fortinet_CA_Untrusted",scope="global",source="factory",vdom="root"} 1.787904385e+09
        fortigate_certificate_valid_to_seconds{name="Fortinet_Factory",scope="global",source="factory",vdom="root"} 2.147483647e+09
        fortigate_certificate_valid_to_seconds{name="Fortinet_SSL",scope="global",source="factory",vdom="root"} 1.78790439e+09
        fortigate_certificate_valid_to_seconds{name="Fortinet_SSL_DSA1024",scope="global",source="factory",vdom="root"} 1.82569362e+09
        fortigate_certificate_valid_to_seconds{name="Fortinet_SSL_DSA2048",scope="global",source="factory",vdom="root"} 1.825693629e+09
        fortigate_certificate_valid_to_seconds{name="Fortinet_SSL_ECDSA256",scope="global",source="factory",vdom="root"} 1.825693629e+09
        fortigate_certificate_valid_to_seconds{name="Fortinet_SSL_ECDSA384",scope="global",source="factory",vdom="root"} 1.825693629e+09
        fortigate_certificate_valid_to_seconds{name="Fortinet_SSL_ECDSA521",scope="global",source="factory",vdom="root"} 1.898449387e+09
        fortigate_certificate_valid_to_seconds{name="Fortinet_SSL_ED25519",scope="global",source="factory",vdom="root"} 1.898449387e+09
        fortigate_certificate_valid_to_seconds{name="Fortinet_SSL_ED448",scope="global",source="factory",vdom="root"} 1.898449387e+09
        fortigate_certificate_valid_to_seconds{name="Fortinet_SSL_RSA1024",scope="global",source="factory",vdom="root"} 1.825693604e+09
        fortigate_certificate_valid_to_seconds{name="Fortinet_SSL_RSA2048",scope="global",source="factory",vdom="root"} 1.825693617e+09
        fortigate_certificate_valid_to_seconds{name="Fortinet_SSL_RSA4096",scope="global",source="factory",vdom="root"} 1.898449387e+09
        fortigate_certificate_valid_to_seconds{name="Fortinet_Wifi",scope="global",source="factory",vdom="root"} 1.640476799e+09
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
