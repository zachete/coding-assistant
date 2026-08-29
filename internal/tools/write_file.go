package tools

import (
	"errors"
	"fmt"
	"os"
)

type WriteToolArgs struct {
	path    string
	content string
}

type WriteFileTool struct {
	rawArgs map[string]any
	args    WriteToolArgs
}

func (t *WriteFileTool) NeedConfirm() bool {
	return true
}

func (t *WriteFileTool) Name() string {
	return "write_file"
}

func (t *WriteFileTool) Description() string {
	return "Write a file content. Use when need to update the file."
}

func (t *WriteFileTool) Params() map[string]any {
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

func (t *WriteFileTool) Execute() (string, error) {
	err := os.WriteFile(t.args.path, []byte(t.args.content), 0644)
	if err != nil {
		return "", err
	}

	return "", nil
}

func (t *WriteFileTool) GetNotice() string {
	return fmt.Sprintf(`Write file(path:"%s")`, t.args.path)
}

func (t *WriteFileTool) SetArgs(rawArgs map[string]any) error {
	rawPath, ok := rawArgs["path"]
	if !ok {
		return errors.New("empty path")
	}

	path, err := sanitizePath(rawPath.(string))
	if err != nil {
		return err
	}

	content, ok := rawArgs["content"].(string)
	if !ok {
		return errors.New("empty file content")
	}

	t.args = WriteToolArgs{
		path:    path,
		content: content,
	}

	return nil
}
