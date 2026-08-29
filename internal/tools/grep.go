package tools

import (
	"bufio"
	"coding-assistant/internal/utils"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"errors"
	"io/fs"
	"path/filepath"
)

type GrepMatch struct {
	Path       string `json:"path"`
	LineNumber int    `json:"lineNumber"`
	Line       string `json:"line"`
}

type GrepResult struct {
	Matches   []GrepMatch `json:"matches"`
	Truncated bool        `json:"truncated"`
}

const MatchesLimit = 50

type GrepTool struct {
	rawArgs map[string]any
	args    GrepToolArgs
}

type GrepToolArgs struct {
	path    string
	pattern string
}

func (t *GrepTool) NeedConfirm() bool {
	return false
}

func (t *GrepTool) Name() string {
	return "grep"
}

func (t *GrepTool) Description() string {
	return `Search for a text pattern in files inside the workspace.
		Use this tool when you need to find:
		- where a function/type/variable is defined;
		- where a symbol is used;
		- files containing specific text.
	`
}

func (t *GrepTool) Params() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"pattern": map[string]any{"type": "string", "description": "The text to search for."},
			"path":    map[string]any{"type": "string", "description": "Directory or file to search in."},
		},
		"required":             []string{"path"},
		"additionalProperties": false,
	}
}

func (t *GrepTool) GetNotice() string {
	return fmt.Sprintf(`grep(pattern:"%s",path:"%s"`, t.args.pattern, t.args.path)
}

func (t *GrepTool) Execute() (string, error) {
	var matches []GrepMatch
	result := GrepResult{}

	err := filepath.WalkDir(t.args.path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if shouldSkip(path, info) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNumber := 0
		for scanner.Scan() {
			line := scanner.Text()
			lineNumber = lineNumber + 1

			if strings.Contains(line, t.args.pattern) {
				matches = append(matches, GrepMatch{
					Path:       path,
					LineNumber: lineNumber,
					Line:       line,
				})
			}
		}
		if err := scanner.Err(); err != nil {
			utils.LogToFile(fmt.Sprintf("error during scan: %s", err))
		}

		matchesLen := len(matches)
		if matchesLen > MatchesLimit {
			result.Truncated = true
			return filepath.SkipAll
		} else if matchesLen == MatchesLimit {
			return filepath.SkipAll
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	result.Matches = matches
	res, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (t *GrepTool) SetArgs(rawArgs map[string]any) error {
	pattern, ok := rawArgs["pattern"]
	if !ok {
		return errors.New("empty pattern")
	}

	rawPath, ok := rawArgs["path"]
	if !ok {
		return errors.New("empty path")
	}

	path, err := sanitizePath(rawPath.(string))
	if err != nil {
		return err
	}

	t.args = GrepToolArgs{
		path:    path,
		pattern: pattern.(string),
	}

	return nil
}
