package tools

import (
	"errors"
	"os"
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

	res, err := os.ReadFile(path)
	if err != nil {
		panic("can not read the file")
	}

	return string(res), nil
}
