package tools

import (
	"errors"
	"os"
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
		return "", errors.New("empty file path")
	}

	content, ok := args["content"].(string)
	if !ok {
		return "", errors.New("empty file content")
	}

	err := os.WriteFile(path, []byte(content), 0644)

	if err != nil {
		return "", err
	}

	return "", nil
}
