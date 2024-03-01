package zbaction

import (
	"bytes"
	"context"
	"io"
	"os/exec"

	"github.com/samber/lo"
)

type CommandStep struct {
	Command []string
}

func (c CommandStep) Run(ctx context.Context, sc *StepContext) (CleanupFn, error) {
	stdoutWriter, stderrWriter := sc.GetWriter()

	stdout := NewContextWriter(sc, "stdout", stdoutWriter)
	stderr := NewContextWriter(sc, "stderr", stderrWriter)

	// expand command
	expandedCommand := lo.Map(c.Command, func(s string, _ int) string {
		return sc.ExpandString(s)
	})

	cmd := exec.CommandContext(ctx, expandedCommand[0], expandedCommand[1:]...)
	cmd.Dir = sc.Root()
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = ListEnvironmentVariables(sc.VariableContainer()).ToList()

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// save to variable
	_ = stdout.Close()
	_ = stderr.Close()

	return nil, nil
}

type contextWriter struct {
	io.Writer

	variable    string
	sc          *StepContext
	bytesBuffer *bytes.Buffer
}

func NewContextWriter(sc *StepContext, variable string, target io.Writer) io.WriteCloser {
	buf := &bytes.Buffer{}
	mw := io.MultiWriter(buf, target)

	return &contextWriter{
		Writer:      mw,
		sc:          sc,
		variable:    variable,
		bytesBuffer: buf,
	}
}

func (cw *contextWriter) Close() error {
	cw.sc.SetThisOutput(cw.variable, cw.bytesBuffer.String())
	return nil
}

var _ RunnableStep = (*CommandStep)(nil)
