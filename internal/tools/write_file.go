package tools

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type WriteFileTool struct{}

func (t WriteFileTool) Name() string {
	return "write_file"
}

func (t WriteFileTool) Description() string {
	return "Write a file content. Use when need to update the file."
}

func (t WriteFileTool) Params() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{"type": "string", "description": "A file path"},
			"content": map[string]any{
				"type":        "string",
				"description": "A file content",
			},
		},
		"required":             []string{"path", "content"},
		"additionalProperties": false,
	}
}

func (t WriteFileTool) Execute(args map[string]any) (string, error) {
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

	content, ok := args["content"].(string)
	if !ok {
		return "", errors.New("empty file content")
	}

	err = os.WriteFile(absPath, []byte(content), 0644)
	if err != nil {
		return "", err
	}

	return "", nil
}
