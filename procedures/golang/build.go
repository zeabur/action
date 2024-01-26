package golang

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/zeabur/builder/zbaction"
)

func init() {
	zbaction.RegisterProcedure("action/golang/build", func(args zbaction.ProcStepArgs) (zbaction.ProcedureStep, error) {
		entry, ok := args["entry"]
		if !ok {
			entry = "."
		}

		return &BuildAction{
			Entry: zbaction.NewArgumentStr(entry),
		}, nil
	})
}

type BuildAction struct {
	Entry zbaction.Argument[string]
}

func (b BuildAction) Run(ctx context.Context, sc *zbaction.StepContext) (zbaction.CleanupFn, error) {
	entry := b.Entry.Value(sc.ExpandString)

	// Make a directory for storing the binaries.
	outDir, err := os.MkdirTemp("", "zbpack-go-out-*")
	if err != nil {
		return nil, fmt.Errorf("make temp dir: %w", err)
	}

	outFile := path.Join(outDir, "server")
	// Build the binary.
	{
		cmd := exec.CommandContext(ctx, "go", "build", "-o", outFile, entry)
		cmd.Env = []string{
			"CGO_ENABLED=0",
		}
		cmd.Dir = sc.Root()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("build: %w", err)
		}
	}

	// Set the output
	sc.SetThisOutput("outDir", outDir)
	sc.SetThisOutput("outFile", outFile)

	// Clean up
	return func() {
		_ = os.RemoveAll(outDir)
	}, nil
}
