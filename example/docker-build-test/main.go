package main

import (
	"context"
	"log/slog"
	"os"

	zbaction "github.com/zeabur/action"
	_ "github.com/zeabur/action/procedures"
	_ "github.com/zeabur/action/procedures/artifact"
	"github.com/zeabur/action/procedures/procvariables"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))

	err := zbaction.RunAction(context.TODO(), zbaction.Action{
		Jobs: []zbaction.Job{
			{
				Steps: []zbaction.Step{
					{
						RunnableStep: zbaction.ProcStep{
							Uses: "action/artifact/docker",
							With: zbaction.ProcStepArgs{
								"context": "${context.root}",
								"dockerfile": `
										FROM alpine
										RUN echo "Hello, World!"
								`,
							},
						},
					},
				},
			},
		},
	}, procvariables.WithEnvBuildkitHost())

	if err != nil {
		panic(err)
	}
}
