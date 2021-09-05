package test

import (
	"path/filepath"

	"github.com/bradleyjkemp/cupaloy/v2"
)

var Snapshotting = cupaloy.New(
	cupaloy.SnapshotSubdirectory(filepath.Join("testdata", "snapshots")))
