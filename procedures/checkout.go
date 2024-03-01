package procedures

import (
	"context"
	"strconv"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	zbaction "github.com/zeabur/action"
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
			AuthUsername: zbaction.NewArgumentStr(args["authUsername"]),
			AuthPassword: zbaction.NewArgumentStr(args["authPassword"]),
		}, nil
	})
}

type CheckoutAction struct {
	URL          zbaction.Argument[string]
	Branch       zbaction.Argument[string]
	Depth        zbaction.Argument[int]
	AuthUsername zbaction.Argument[string]
	AuthPassword zbaction.Argument[string]
}

func (i *CheckoutAction) Run(ctx context.Context, sc *zbaction.StepContext) (zbaction.CleanupFn, error) {
	url := i.URL.Value(sc.ExpandString)
	branch := i.Branch.Value(sc.ExpandString)
	depth := i.Depth.Value(sc.ExpandString)
	authUsername := i.AuthUsername.Value(sc.ExpandString)
	authPassword := i.AuthPassword.Value(sc.ExpandString)

	var auth *http.BasicAuth
	if authUsername != "" && authPassword != "" {
		auth = &http.BasicAuth{
			Username: authUsername,
			Password: authPassword,
		}
	}

	_, err := git.PlainCloneContext(
		ctx,
		sc.Root(),
		false,
		&git.CloneOptions{
			URL:               url,
			Progress:          sc.Stderr(),
			Depth:             depth,
			ReferenceName:     plumbing.ReferenceName("refs/heads/" + branch),
			SingleBranch:      true,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			Auth:              auth,
		},
	)

	return nil, err
}

var _ zbaction.ProcedureStep = (*CheckoutAction)(nil)
