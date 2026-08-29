package tools

type Tool interface {
	Name() string
	Description() string
	Params() map[string]any
	NeedConfirm() bool
	Execute() (string, error)
	SetArgs(args map[string]any) error
	GetNotice() string
}
