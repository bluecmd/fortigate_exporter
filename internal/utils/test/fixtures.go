package test

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
)

func GetCallerFilePath(stackLevel int) (string, error) {
	_, filename, _, ok := runtime.Caller(stackLevel)
	if !ok {
		return "", fmt.Errorf("no caller information")
	}
	return filename, nil
}

func GetFixturePath(fixtureFileName string) (string, error) {
	filename, err := GetCallerFilePath(2)
	if err != nil {
		return "", err
	}
	fixturePath := path.Join(path.Dir(filename), "testdata", "fixtures", fixtureFileName)
	return fixturePath, nil
}

func GetFixturePathPanic(fixtureFileName string) string {
	filename, err := GetCallerFilePath(2)
	if err != nil {
		log.Panicf("failed to get caller stack %v", err)
	}
	fixturePath := path.Join(path.Dir(filename), "testdata", "fixtures", fixtureFileName)
	return fixturePath
}

func GetRelativeFixturePathPanic(fixtureFileName string) string {
	filename, err := GetCallerFilePath(2)
	if err != nil {
		log.Panicf("failed to get caller stack %v", err)
	}
	workDir, err := os.Getwd()
	if err != nil {
		log.Panicf("failed to get workdir %v", err)
	}
	relDir := strings.TrimPrefix(path.Dir(filename), workDir)
	fixturePath := path.Join(relDir, "testdata", "fixtures", fixtureFileName)
	return fixturePath
}
