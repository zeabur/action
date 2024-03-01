package zbaction

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"golang.org/x/sync/errgroup"
)

type StepsOutputMap map[StepID]StepOutput
type StepOutput = map[string]any

func RunAction(ctx context.Context, action Action, options ...ExecutorOptionsFn) error {
	slog.Info("Running action", slog.String("action", action.String()))

	executorOptions := ExecutorOptions{
		/* defaults */
		RuntimeVariables: nil,
		Stdout:           os.Stdout,
		Stderr:           os.Stderr,
	}
	for _, fn := range options {
		fn(&executorOptions)
	}

	variables := NewMapContainer(action.Variables)
	if len(executorOptions.RuntimeVariables) > 0 {
		variables = NewVariableContainerWithExtraParameters(executorOptions.RuntimeVariables, variables)
	}

	ac := &ActionContext{
		variables: variables,
		action:    &action,
	}

	type CleanupFnContext struct {
		fn func() error
		id JobID
	}

	// we times len(action.Jobs) for 2 to prevent potential deadlock
	jobCleanupFn := make(chan CleanupFnContext, len(action.Jobs)*2)
	defer func() {
		for cleanup := range jobCleanupFn {
			if cleanup.fn == nil {
				continue
			}

			if err := cleanup.fn(); err != nil {
				slog.Error("Failed to cleanup job",
					slog.String("job", cleanup.id),
					slog.String("error", err.Error()))
			}
		}
	}()

	eg, ectx := errgroup.WithContext(ctx)

	for _, job := range action.Jobs {
		job := job

		eg.Go(func() error {
			jc := &JobContext{
				actionContext: ac,
				job:           &job,
				output:        make(StepsOutputMap),
				variables:     NewMapContainer(job.Variables),
			}
			defer func(jc *JobContext, job Job) {
				jobCleanupFn <- CleanupFnContext{
					fn: jc.Cleanup,
					id: job.ID,
				}
			}(jc, job)

			if err := jc.Run(ectx); err != nil {
				if errors.Is(err, context.Canceled) {
					return err // cancelled
				}

				slog.Error("Failed to run job",
					slog.String("job", job.String()),
					slog.String("error", err.Error()))
				return fmt.Errorf("run job %s: %w", job.String(), err)
			}

			return nil
		})
	}

	err := eg.Wait()
	close(jobCleanupFn)

	if err != nil {
		slog.Error("Failed to run action",
			slog.String("action", action.String()),
			slog.String("error", err.Error()))
		return fmt.Errorf("run action %s: %w", action.String(), err)
	}

	return nil
}

type ActionContext struct {
	variables VariableContainer
	action    *Action

	stdout io.Writer
	stderr io.Writer

	cachedID *ActionID `exhaustruct:"optional"`
}

func (ac *ActionContext) ID() ActionID {
	if ac.cachedID == nil {
		id := ac.action.String()
		ac.cachedID = &id
	}

	return *ac.cachedID
}

func (ac *ActionContext) VariableContainer() VariableContainer {
	return ac.variables
}

func (ac *ActionContext) Action() Action {
	return *ac.action
}

func (ac *ActionContext) GetVariable(key string) (string, bool) {
	return ac.VariableContainer().GetVariable(key)
}

func (ac *ActionContext) GetRawVariable(key string) (string, bool) {
	return ac.VariableContainer().GetRawVariable(key)
}

type JobContext struct {
	actionContext *ActionContext
	job           *Job

	output    map[StepID]StepOutput
	variables VariableContainer

	root *string `exhaustruct:"optional"`

	// cache

	cachedID *JobID `exhaustruct:"optional"`
}

func (jc *JobContext) ID() JobID {
	if jc.cachedID == nil {
		id := jc.job.String()
		jc.cachedID = &id
	}

	return *jc.cachedID
}

func (jc *JobContext) ActionContext() ActionContext {
	return *jc.actionContext
}

func (jc *JobContext) Job() Job {
	return *jc.job
}

func (jc *JobContext) VariableContainer() VariableContainer {
	root, err := jc.GetRoot()
	if err != nil {
		slog.Error("failed to inject context.root to job variables", slog.String("error", err.Error()))
		root = ""
	}

	return NewVariableContainerWithExtraParameters(
		map[string]string{
			// ${context.root}
			"context.root": root,
		},
		NewVariableContainerWithParent(jc.variables, jc.actionContext.VariableContainer()),
	)
}

func (jc *JobContext) GetVariable(key string) (string, bool) {
	return jc.VariableContainer().GetVariable(key)
}

func (jc *JobContext) GetRawVariable(key string) (string, bool) {
	return jc.VariableContainer().GetRawVariable(key)
}

func (jc *JobContext) GetRoot() (string, error) {
	if jc.root != nil {
		return *jc.root, nil
	}

	// create a temporary path
	tmpdir, err := os.MkdirTemp("", "zbaction-*")
	if err != nil {
		return "", fmt.Errorf("create temporary directory: %w", err)
	}

	jc.root = &tmpdir
	return tmpdir, nil
}

func (jc *JobContext) Cleanup() error {
	if jc.root == nil {
		return nil
	}

	return os.RemoveAll(*jc.root)
}

func (jc *JobContext) Run(ctx context.Context) error {
	job := jc.Job()
	slog.Info("Running job", slog.String("job", job.String()))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	root, err := jc.GetRoot()
	if err != nil {
		return fmt.Errorf("get job root: %w", err)
	}

	cleanupStack := CleanupStack{}
	defer cleanupStack.Run()

	for _, step := range job.Steps {
		step := step

		slog.Info("Running step", slog.String("step", step.HumanName()))

		sc := &StepContext{
			id:         step.String(),
			jobContext: jc,
			root:       root,
			variables:  NewMapContainer(step.Variables),
		}

		if err := ctx.Err(); err != nil {
			return err
		}

		cleanup, err := step.Run(ctx, sc)
		if cleanup != nil {
			cleanupStack.Push(cleanup)
		}
		if err != nil {
			slog.Error("Failed to run step",
				slog.String("step", step.String()),
				slog.String("error", err.Error()))
			return fmt.Errorf("failed to run step %s: %w", step.String(), err)
		}
	}

	return nil
}

type StepContext struct {
	id         StepID
	jobContext *JobContext

	root      string
	variables VariableContainer
}

func (sc *StepContext) Root() string {
	return sc.root
}

func (sc *StepContext) ID() StepID {
	return sc.id
}

func (sc *StepContext) JobContext() JobContext {
	return *sc.jobContext
}

func (sc *StepContext) VariableContainer() VariableContainer {
	return NewVariableContainerWithParent(sc.variables, sc.jobContext.VariableContainer())
}

func (sc *StepContext) SetThisOutput(key string, value any) {
	if sc.jobContext.output == nil {
		sc.jobContext.output = make(StepsOutputMap)
	}
	if sc.jobContext.output[sc.id] == nil {
		sc.jobContext.output[sc.id] = make(StepOutput)
	}

	sc.jobContext.output[sc.id][key] = value
}

func (sc *StepContext) GetThisOutput(key string) (any, bool) {
	return sc.GetOutput(sc.id, key)
}

func (sc *StepContext) GetOutput(id StepID, key string) (any, bool) {
	if sc.jobContext.output == nil {
		return nil, false
	}

	if sc.jobContext.output[id] == nil {
		return nil, false
	}

	if value, ok := sc.jobContext.output[id][key]; ok {
		return value, true
	}

	return nil, false
}

func (sc *StepContext) ExpandString(s string) string {
	return os.Expand(s, func(s string) string {
		if v, ok := sc.VariableContainer().GetVariable(s); ok {
			return v
		}

		// ${out.<step_id>.<key>}
		if after, found := strings.CutPrefix(s, "out."); found {
			if stepID, key, ok := strings.Cut(after, "."); ok {
				if v, ok := sc.GetOutput(stepID, key); ok {
					return fmt.Sprintf("%v", v)
				}
			}
		}

		return ""
	})
}

// GetWriter returns the general stdout and stderr writer for the action.
//
// The first return value is the stdout writer, and the second return value is the stderr writer.
func (sc *StepContext) GetWriter() (io.Writer, io.Writer) {
	return sc.jobContext.actionContext.stdout, sc.jobContext.actionContext.stderr
}
