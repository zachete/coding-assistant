package tools

import (
	"errors"
	"fmt"
	"os"
)

type ReadFileArgs struct {
	path string
}

type ReadFileTool struct {
	rawArgs map[string]any
	args    ReadFileArgs
}

func (t *ReadFileTool) NeedConfirm() bool {
	return true
}

func (t *ReadFileTool) Name() string {
	return "read_file"
}

func (t *ReadFileTool) Description() string {
	return "Read a file content. Use when need to read a file content."
}

func (t *ReadFileTool) Params() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{"type": "string", "description": "A file path"},
		},
		"required":             []string{"path"},
		"additionalProperties": false,
	}
}

func (t *ReadFileTool) Execute() (string, error) {
	res, err := os.ReadFile(t.args.path)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (t *ReadFileTool) SetArgs(rawArgs map[string]any) error {
	rawPath, ok := rawArgs["path"]
	if !ok {
		return errors.New("empty path")
	}

	path, err := sanitizePath(rawPath.(string))
	if err != nil {
		return err
	}

	t.args = ReadFileArgs{
		path: path,
	}

	return nil
}

func (t *ReadFileTool) GetNotice() string {
	return fmt.Sprintf(`Read file(path:"%s")`, t.args.path)
}
