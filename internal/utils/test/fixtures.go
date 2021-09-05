package test

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bluecmd/fortigate_exporter/internal/utils/files"
)

func GetFixturePath(fixtureFileName string) (string, error) {
	callerDir, err := files.GetCallerDir(2)
	if err != nil {
		return "", err
	}
	fixturePath := filepath.Join(callerDir, "testdata", "fixtures", fixtureFileName)
	return fixturePath, nil
}

func GetFixturePathPanic(fixtureFileName string) string {
	callerDir, err := files.GetCallerDir(2)
	if err != nil {
		log.Panicf("failed to get caller stack %v", err)
	}
	fixturePath := filepath.Join(callerDir, "testdata", "fixtures", fixtureFileName)
	return fixturePath
}

func GetRelativeFixturePathPanic(fixtureFileName string) string {
	callerDir, err := files.GetCallerDir(2)
	if err != nil {
		log.Panicf("failed to get caller stack %v", err)
	}
	workDir, err := os.Getwd()
	if err != nil {
		log.Panicf("failed to get workdir %v", err)
	}
	relDir := strings.TrimPrefix(callerDir, workDir)
	fixturePath := filepath.Join(relDir, "testdata", "fixtures", fixtureFileName)
	return fixturePath
}
