package preprocessingspec

import "go.temporal.io/sdk/workflow"

type Preprocessing interface {
	Activities() []Activity
	Execute(ctx workflow.Context, params Params) (Result, error)
}

type Activity struct {
	Name    string
	Execute interface{}
}

type Params struct {
	Path string
}

type Result struct {
	Path string
}
