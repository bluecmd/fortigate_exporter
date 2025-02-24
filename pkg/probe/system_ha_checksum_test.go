// Copyright 2025 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	# HELP fortigate_ha_member_has_role Master/Slave information
	# TYPE fortigate_ha_member_has_role gauge
	fortigate_ha_member_has_role{role="manage_master", serial="SERIAL111111111"} 1
	fortigate_ha_member_has_role{role="manage_master", serial="SERIAL222222222"} 0
	fortigate_ha_member_has_role{role="root_master", serial="SERIAL111111111"} 1
	fortigate_ha_member_has_role{role="root_master", serial="SERIAL222222222"} 0
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
