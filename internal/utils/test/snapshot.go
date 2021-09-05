package test

import (
	"path"

	"github.com/bradleyjkemp/cupaloy/v2"
)

var Snapshotting = cupaloy.New(
	cupaloy.SnapshotSubdirectory(path.Join("testdata", "snapshots")))
