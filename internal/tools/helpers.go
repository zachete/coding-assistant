package tools

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var ignoredDirs = map[string]struct{}{
	".git":         {},
	".hg":          {},
	".svn":         {},
	".cache":       {},
	"node_modules": {},
	"vendor":       {},
	"dist":         {},
	"build":        {},
	"tmp":          {},
	"bin":          {},
}

func shouldSkip(path string, info os.FileInfo) bool {
	if info.IsDir() {
		_, skip := ignoredDirs[info.Name()]
		return skip
	}

	return isBinaryFile(path)
}

func isBinaryFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return true
	}
	defer file.Close()

	buf := make([]byte, 512)

	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return true
	}

	return bytes.IndexByte(buf[:n], 0) >= 0
}

func sanitizePath(path string) (string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(filepath.Join(root, path))
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(absPath, root) {
		return "", fmt.Errorf("invalid path")
	}

	return absPath, nil
}
