package tools

import (
	"errors"
	"maps"
	"slices"
)

type Registry struct {
	items map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{
		items: map[string]Tool{
			"read_file":  &ReadFileTool{},
			"write_file": &WriteFileTool{},
			"grep":       &GrepTool{},
		},
	}
}

func (r *Registry) ResolveTool(name string) (Tool, error) {
	if r.items[name] == nil {
		return nil, errors.New("Tool not found")
	}
	return r.items[name], nil
}

func (r *Registry) AsSlice() []Tool {
	return slices.Collect(maps.Values(r.items))
}
