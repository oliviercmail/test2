package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/automation/corredor_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
)

var _ wfexec.ExecResponse

type (
	corredorHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h corredorHandler) register() {
	h.reg.AddFunctions(
		h.Exec(),
	)
}

type (
	corredorExecArgs struct {
		hasScript bool
		Script    string

		hasArgs bool
		Args    interface{}
	}

	corredorExecResults struct {
		Results interface{}
	}
)

// Exec function Executes script in Corredor Automation server
//
// expects implementation of exec function:
// func (h corredorHandler) exec(ctx context.Context, args *corredorExecArgs) (results *corredorExecResults, err error) {
//    return
// }
func (h corredorHandler) Exec() *atypes.Function {
	return &atypes.Function{
		Ref:    "corredorExec",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short: "Executes script in Corredor Automation server",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "script",
				Types: []string{"String"}, Required: true,
			},
			{
				Name:  "args",
				Types: []string{"Any"},
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "results",
				Types: []string{"Any"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &corredorExecArgs{
					hasScript: in.Has("script"),
					hasArgs:   in.Has("args"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *corredorExecResults
			if results, err = h.exec(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Results (interface{}) to Any
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Any").Cast(results.Results); err != nil {
					return
				} else if err = expr.Assign(out, "results", tval); err != nil {
					return
				}
			}

			return
		},
	}
}
