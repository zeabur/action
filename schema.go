package zbaction

import (
	"context"
	"fmt"

	"github.com/mitchellh/hashstructure/v2"
)

type MachineName = string
type ProcStepName = string
type ProcStepArgs = map[string]string
type ActionID = string
type JobID = string
type StepID = string

type Action struct {
	ID        ActionID
	Variables map[string]string
	Jobs      []Job
}

func (a Action) String() string {
	if a.ID != "" {
		return a.ID
	}

	uuid, err := hashstructure.Hash(a, hashstructure.FormatV2, nil)
	if err == nil {
		return fmt.Sprintf("%x", uuid)
	}

	return "<unknown action>"
}

type Job struct {
	ID        JobID
	RunOn     MachineName
	Variables map[string]string
	Steps     []Step
}

func (j Job) String() string {
	if j.ID != "" {
		return j.ID
	}

	uuid, err := hashstructure.Hash(j, hashstructure.FormatV2, nil)
	if err == nil {
		return fmt.Sprintf("%x", uuid)
	}

	return "<unknown job>"
}

type Step struct {
	ID        StepID
	Name      string
	Variables map[string]string
	RunnableStep
}

func (s Step) HumanName() string {
	if s.Name != "" {
		return s.Name
	}

	return s.String()
}

func (s Step) String() string {
	if s.ID != "" {
		return s.ID
	}

	uuid, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	if err == nil {
		return fmt.Sprintf("%x", uuid)
	}

	return "<unknown step>"
}

type RunnableStep interface {
	Run(ctx context.Context, sc *StepContext) (CleanupFn, error)
}

type ProcStep struct {
	Uses ProcStepName
	With ProcStepArgs
}

func (p ProcStep) Run(ctx context.Context, sc *StepContext) (CleanupFn, error) {
	step, err := ResolveProcedure(p.Uses, p.With)
	if err != nil {
		return nil, err
	}

	return step.Run(ctx, sc)
}

type CleanupFn func()
