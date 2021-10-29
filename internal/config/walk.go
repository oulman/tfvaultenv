/*
Copyright © 2021 James Oulman

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// formats a path string based on parent depth
func getPath(depth int, filename string) string {
	sep := string(os.PathSeparator)
	if depth == 0 {
		return fmt.Sprintf(".%s%s", sep, filename)
	}

	var pathPrefix string
	for i := 0; i < depth; i++ {
		pathPrefix += fmt.Sprintf("..%s", sep)
	}

	return filepath.Join(pathPrefix, filename)
}

// walks backwards depth n parent directories looking for filename
func FindConfigFile(depth int, filename string) (string, error) {
	for i := 0; i <= depth; i++ {
		path := getPath(i, filename)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("unable to locate config file %s", filename)
}
