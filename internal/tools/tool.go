package tools

type Tool interface {
	Name() string
	Description() string
	Params() map[string]any
	Execute(args map[string]any) (string, error)
}
