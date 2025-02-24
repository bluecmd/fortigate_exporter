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

func TestUserFsso(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/user/fsso", "testdata/user-fsso.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeUserFsso, c, r) {
		t.Errorf("probeUserFsso() returned non-success")
	}

	em := `
	# HELP fortigate_user_fsso_info Info on Fsso defined connectors
	# TYPE fortigate_user_fsso_info gauge
	fortigate_user_fsso_info{id="",name="FSSO-VDOM2",status="connected",type="fsso",vdom="vdom2"} 1
	fortigate_user_fsso_info{id="",name="FSSO-VDOM4_2",status="disconnected",type="fsso",vdom="vdom4"} 1
	fortigate_user_fsso_info{id="",name="FSSO_VDOM1",status="disconnected",type="fsso",vdom="vdom1"} 1
	fortigate_user_fsso_info{id="",name="FSSO_VDOM3_1",status="connected",type="fsso",vdom="vdom3"} 1
	fortigate_user_fsso_info{id="",name="FSSO_VDOM3_2",status="connected",type="fsso",vdom="vdom3"} 1
	fortigate_user_fsso_info{id="",name="FSSO_VDOM4_1",status="disconnected",type="fsso",vdom="vdom4"} 1
	fortigate_user_fsso_info{id="1",name="",status="disconnected",type="fsso-polling",vdom="vdom5"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
