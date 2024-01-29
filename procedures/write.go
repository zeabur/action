package procedures

import (
	"context"
	"os"
	"path"

	zbaction "github.com/zeabur/action"
)

func init() {
	zbaction.RegisterProcedure("action/write", func(args zbaction.ProcStepArgs) (zbaction.ProcedureStep, error) {
		filename := args["filename"]
		if filename == "" {
			return nil, zbaction.NewErrRequiredArgument("filename")
		}

		return &WriteAction{
			Filename: zbaction.NewArgumentStr(filename),
			Content:  zbaction.NewArgumentStr(args["content"]),
		}, nil
	})
}

type WriteAction struct {
	Filename zbaction.Argument[string]
	Content  zbaction.Argument[string]
}

func (w *WriteAction) Run(_ context.Context, sc *zbaction.StepContext) (zbaction.CleanupFn, error) {
	outFilePath := path.Join(sc.Root(), w.Filename.Value(sc.ExpandString))
	err := os.WriteFile(outFilePath, []byte(w.Content.Value(sc.ExpandString)), 0644)
	if err != nil {
		return nil, err
	}

	sc.SetThisOutput("filepath", outFilePath)

	return func() {
		_ = os.Remove(outFilePath)
	}, nil
}
