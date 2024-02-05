package procedures

import (
	"context"
	"path"

	cp "github.com/otiai10/copy"
	zbaction "github.com/zeabur/action"
)

func init() {
	zbaction.RegisterProcedure("action/copy-local-dir", func(args zbaction.ProcStepArgs) (zbaction.ProcedureStep, error) {
		src := args["src"]
		if src == "" {
			return nil, zbaction.NewErrRequiredArgument("src")
		}

		dest := args["dest"]
		if dest == "" {
			return nil, zbaction.NewErrRequiredArgument("dest")
		}

		return &CopyLocalDirAction{
			Src:  zbaction.NewArgumentStr(src),
			Dest: zbaction.NewArgumentStr(dest),
		}, nil
	})
}

type CopyLocalDirAction struct {
	Src  zbaction.Argument[string]
	Dest zbaction.Argument[string]
}

func (c CopyLocalDirAction) Run(_ context.Context, sc *zbaction.StepContext) (zbaction.CleanupFn, error) {
	err := cp.Copy(c.Src.Value(sc.ExpandString), path.Join(sc.Root(), c.Dest.Value(sc.ExpandString)))
	return nil, err
}
