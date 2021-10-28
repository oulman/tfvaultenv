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
	return "", fmt.Errorf("unable to locate %s at", filename)
}
