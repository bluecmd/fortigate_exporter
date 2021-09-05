package files

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// GetCallerDir stacktraceStep is the number calls we travel up to identify the source Go file. This will usually be 2.
func GetCallerDir(stacktraceStep int) (string, error) {
	_, filename, _, ok := runtime.Caller(stacktraceStep)
	if ok {
		return filepath.Dir(filename), nil
	}
	return "", fmt.Errorf("could not retrieve current dir")
}

// ReadRelativeFile read file relative to the calling go file
func ReadRelativeFile(relativePath string) ([]byte, error) {
	sourceDir, callErr := GetCallerDir(2)
	if callErr != nil {
		return nil, callErr
	}

	filePath := filepath.Join(sourceDir, relativePath)

	return os.ReadFile(filePath)
}
