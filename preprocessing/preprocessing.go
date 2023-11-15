package preprocessing

import (
	"go.temporal.io/sdk/workflow"

	"github.com/jraddaoui/preprocessing/spec"
)

var _ spec.Preprocessing = (*Preprocessing)(nil)

type Preprocessing struct{}

func (s *Preprocessing) Activities() []spec.Activity {
	return []spec.Activity{}
}

func (s *Preprocessing) Execute(ctx workflow.Context, params spec.Params) (spec.Result, error) {
	return spec.Result{Path: params.Path}, nil
}
