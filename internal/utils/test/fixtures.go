package test

import (
	"path/filepath"
)

func GetRelativeFixturePath(fixtureFileName string) string {
	return filepath.Join("testdata", "fixtures", fixtureFileName)
}
