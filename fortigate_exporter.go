// Server executable of fortigate_exporter
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

package main

import (
	"log"
	"net/http"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/bluecmd/fortigate_exporter/pkg/probe"

	"github.com/bluecmd/fortigate_exporter/internal/config"
	fortiHTTP "github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Version = "(devel)"
	GitHash = "(no hash)"
)

type BuildInfo struct {
	version   string
	gitHash   string
	goVersion string
}

func setUpMetricsEndpoint(buildInfo BuildInfo) {
	fortigateExporterInfo := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "fortigate_exporter_build_info",
		Help: "This info metric contains build information for about the exporter",
	}, []string{"version", "revision", "goversion"})

	fortigateExporterInfo.With(prometheus.Labels{
		"version":   buildInfo.version,
		"revision":  buildInfo.gitHash,
		"goversion": buildInfo.goVersion,
	}).Set(1)
}

func getBuildInfo() BuildInfo {
	// don't overwrite the version if it was set by -ldflags=-X
	if info, ok := debug.ReadBuildInfo(); ok && Version == "(devel)" {
		mod := &info.Main
		if mod.Replace != nil {
			mod = mod.Replace
		}
		Version = mod.Version
	}
	// remove leading `v`
	massagedVersion := strings.TrimPrefix(Version, "v")
	buildInfo := BuildInfo{
		version:   massagedVersion,
		gitHash:   GitHash,
		goVersion: runtime.Version(),
	}
	return buildInfo
}

func main() {
	buildInfo := getBuildInfo()
	log.Printf("FortigateExporter %s ( %s )", buildInfo.version, buildInfo.gitHash)
	setUpMetricsEndpoint(buildInfo)

	if err := config.Init(); err != nil {
		log.Fatalf("Initialization error: %+v", err)
	}

	savedConfig := config.GetConfig()

	if err := fortiHTTP.Configure(savedConfig); err != nil {
		log.Fatalf("%+v", err)
	}

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/probe", probe.ProbeHandler)
	go func() {
		if err := http.ListenAndServe(savedConfig.Listen, nil); err != nil {
			log.Fatalf("Unable to serve: %v", err)
		}
	}()
	log.Printf("Fortigate exporter running, listening on %q", savedConfig.Listen)
	select {}
}
