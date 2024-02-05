package golang

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	zbaction "github.com/zeabur/action"
	"golang.org/x/exp/slog"
)

func init() {
	zbaction.RegisterProcedure("action/golang/mod-download", func(args zbaction.ProcStepArgs) (zbaction.ProcedureStep, error) {
		optional := args["optional"]
		if optional == "" {
			optional = "false"
		}

		return &DownloadAction{
			Optional: zbaction.NewArgumentBool(optional),
		}, nil
	})
}

type DownloadAction struct {
	// Optional is a flag to download optional dependencies.
	//
	// Default: false
	Optional zbaction.Argument[bool]
}

func (d DownloadAction) Run(ctx context.Context, sc *zbaction.StepContext) (zbaction.CleanupFn, error) {
	optional := d.Optional.Value(sc.ExpandString)

	// Download dependencies. FIXME: cachable
	cmd := exec.CommandContext(ctx, "go", "mod", "download")
	cmd.Dir = sc.Root()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		if optional {
			slog.Error("failed to download mod (optional)", slog.String("error", err.Error()))
			return nil, nil
		}

		return nil, fmt.Errorf("download mod: %w", err)
	}

	return nil, nil
}
