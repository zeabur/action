package golang

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/zeabur/action"
)

func init() {
	zbaction.RegisterProcedure("action/golang/mod-download", func(args zbaction.ProcStepArgs) (zbaction.ProcedureStep, error) {
		return &DownloadAction{}, nil
	})
}

type DownloadAction struct{}

func (d DownloadAction) Run(ctx context.Context, sc *zbaction.StepContext) (zbaction.CleanupFn, error) {
	// Download dependencies. FIXME: cachable
	cmd := exec.CommandContext(ctx, "go", "mod", "download")
	cmd.Dir = sc.Root()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("download mod: %w", err)
	}

	return nil, nil
}
