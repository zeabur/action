package procedures

import (
	"context"
	"os"
	"strconv"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/zeabur/builder/zbaction"
)

func init() {
	zbaction.RegisterProcedure("action/checkout", func(args zbaction.ProcStepArgs) (zbaction.ProcedureStep, error) {
		depth, ok := args["depth"]
		if !ok {
			depth = "1"
		}

		return &CheckoutAction{
			URL:    zbaction.NewArgumentStr(args["url"]),
			Branch: zbaction.NewArgumentStr(args["branch"]),
			Depth: zbaction.NewArgument(depth, func(s string) int {
				if i, err := strconv.Atoi(s); err == nil {
					return i
				}
				return 1
			}),
		}, nil
	})
}

type CheckoutAction struct {
	URL    zbaction.Argument[string]
	Branch zbaction.Argument[string]
	Depth  zbaction.Argument[int]
}

func (i *CheckoutAction) Run(ctx context.Context, sc *zbaction.StepContext) (zbaction.CleanupFn, error) {
	url := i.URL.Value(sc.ExpandString)
	branch := i.Branch.Value(sc.ExpandString)
	depth := i.Depth.Value(sc.ExpandString)

	_, err := git.PlainCloneContext(
		ctx,
		sc.Root(),
		false,
		&git.CloneOptions{
			URL:               url,
			Progress:          os.Stdout,
			Depth:             depth,
			ReferenceName:     plumbing.ReferenceName("refs/heads/" + branch),
			SingleBranch:      true,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		},
	)

	return nil, err
}

var _ zbaction.ProcedureStep = (*CheckoutAction)(nil)
