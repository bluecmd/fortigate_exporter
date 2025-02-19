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

package files

import (
	"fmt"
	"os"
	"path"
	"runtime"
)

// GetCallerDir stacktraceStep is the number calls we travel up to identify the source Go file. This will usually be 2.
func GetCallerDir(stacktraceStep int) (string, error) {
	_, filename, _, ok := runtime.Caller(stacktraceStep)
	if ok {
		return path.Dir(filename), nil
	}
	return "", fmt.Errorf("could not retrieve current dir")
}

// ReadRelativeFile read file relative to the calling go file
func ReadRelativeFile(relativePath string) ([]byte, error) {
	sourceDir, callErr := GetCallerDir(2)
	if callErr != nil {
		return nil, callErr
	}

	filepath := path.Join(sourceDir, relativePath)

	return os.ReadFile(filepath)
}
