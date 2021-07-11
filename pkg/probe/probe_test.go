// Tests of fortigate_exporter
//
// Copyright (C) 2020  Christian Svensson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package probe

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/google/go-jsonnet"
	"github.com/prometheus/client_golang/prometheus"
)

type preparedResp struct {
	d []byte
	q url.Values
}

type fakeClient struct {
	data map[string][]preparedResp
}

func (c *fakeClient) prepare(path string, jfile string) {
	u, err := url.Parse(path)
	if err != nil {
		panic(err)
	}
	vm := jsonnet.MakeVM()
	output, err := vm.EvaluateFile(jfile)
	if err != nil {
		log.Fatalf("Failed to evaluate jsonnet %q: %v", jfile, err)
	}
	c.data[u.Path] = append(c.data[u.Path], preparedResp{
		d: []byte(output),
		q: u.Query(),
	})
}

func (c *fakeClient) Get(path string, query string, obj interface{}) error {
	rs, ok := c.data[path]
	if !ok {
		log.Fatalf("Tried to get unprepared URL %q", path)
	}
	q, err := url.ParseQuery(query)
	if err != nil {
		log.Fatalf("Unable to parse DUT query: %v", err)
	}
alt:
	for _, r := range rs {
		for k, v := range r.q {
			if len(q[k]) == 0 || q[k][0] != v[0] {
				continue alt
			}
		}
		return json.Unmarshal(r.d, obj)
	}
	log.Fatalf("No prepared response matched URL %q, query %q", path, query)
	return nil
}

type Registry interface {
	MustRegister(...prometheus.Collector)
}

type testProbeCollector struct {
	metrics []prometheus.Metric
}

func (p *testProbeCollector) Collect(c chan<- prometheus.Metric) {
	for _, m := range p.metrics {
		c <- m
	}
}

func (p *testProbeCollector) Describe(c chan<- *prometheus.Desc) {
}

func testProbe(pf probeFunc, c http.FortiHTTP, r Registry) bool {
	meta := &TargetMetadata{
		VersionMajor: 7,
		VersionMinor: 0,
	}
	return testProbeWithMetadata(pf, c, meta, r)
}

func testProbeWithMetadata(pf probeFunc, c http.FortiHTTP, meta *TargetMetadata, r Registry) bool {
	m, ok := pf(c, meta)
	if !ok {
		return false
	}
	p := &testProbeCollector{metrics: m}
	r.MustRegister(p)
	return true
}

func newFakeClient() *fakeClient {
	return &fakeClient{data: map[string][]preparedResp{}}
}
