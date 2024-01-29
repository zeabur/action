package procedures

import (
	"context"
	"fmt"

	zbaction "github.com/zeabur/action"
)

func init() {
	zbaction.RegisterProcedure("action/echo", func(args zbaction.ProcStepArgs) (zbaction.ProcedureStep, error) {
		return &EchoAction{
			Message: zbaction.NewArgumentStr(args["message"]),
		}, nil
	})
}

type EchoAction struct {
	Message zbaction.Argument[string]
}

func (i *EchoAction) Run(_ context.Context, sc *zbaction.StepContext) (zbaction.CleanupFn, error) {
	fmt.Println(i.Message.Value(sc.ExpandString))
	return nil, nil
}

var _ zbaction.ProcedureStep = (*EchoAction)(nil)
