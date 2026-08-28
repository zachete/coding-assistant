package tools

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ReadFileTool struct{}

func (t ReadFileTool) Name() string {
	return "read_file"
}

func (t ReadFileTool) Description() string {
	return "Read a file content. Use when need to read a file content."
}

func (t ReadFileTool) Params() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{"type": "string", "description": "A file path"},
		},
		"required":             []string{"path"},
		"additionalProperties": false,
	}
}

func (t ReadFileTool) Execute(args map[string]any) (string, error) {
	path, ok := args["path"].(string)
	if !ok {
		return "", errors.New("empty path")
	}

	root, err := os.Getwd()
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(filepath.Join(root, path))
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(absPath, root) {
		return "", fmt.Errorf("attempt to access file outside of workspace")
	}

	res, err := os.ReadFile(absPath)
	if err != nil {
		return "", err
	}

	return string(res), nil
}
