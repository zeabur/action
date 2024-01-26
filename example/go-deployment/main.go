package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/zeabur/builder/zbaction"

	_ "github.com/zeabur/builder/zbaction/procedures"
	_ "github.com/zeabur/builder/zbaction/procedures/artifact"
	_ "github.com/zeabur/builder/zbaction/procedures/golang"
)

func main() {
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(h))

	action := zbaction.Action{
		ID: "go-deployment",
		Jobs: []zbaction.Job{
			{
				Steps: []zbaction.Step{
					{
						Name: "Initialize module",
						RunnableStep: zbaction.CommandStep{
							Command: []string{"go", "mod", "init", "github.com/zeabur/zbaction-example/go-deployment"},
						},
					},
					{
						Name: "Write main.go",
						RunnableStep: zbaction.ProcStep{
							Uses: "action/write",
							With: zbaction.ProcStepArgs{
								"filename": "main.go",
								"content": `package main

											func main() {
												println("Hello world!")
											}`,
							},
						},
					},
					{
						Name: "Download dependencies",
						RunnableStep: zbaction.ProcStep{
							Uses: "action/golang/mod-download",
						},
					},
					{
						ID:   "go-binary-step",
						Name: "Build the binary",
						RunnableStep: zbaction.ProcStep{
							Uses: "action/golang/build",
						},
					},
					{
						Name: "DEBUG: see the artifact",
						RunnableStep: zbaction.CommandStep{
							Command: []string{"ls", "-l", "${out.go-binary-step.outDir}"},
						},
					},
					{
						ID:   "docker-image-step",
						Name: "Build the docker image",
						RunnableStep: zbaction.ProcStep{
							Uses: "action/artifact/docker",
							With: zbaction.ProcStepArgs{
								"context": "${out.go-binary-step.outDir}",
								"dockerfile": `
										FROM alpine
                              			COPY ./server /server
										CMD ["/server"]`,
							},
						},
					},
					{
						Name: "Print the artifact",
						RunnableStep: zbaction.ProcStep{
							Uses: "action/echo",
							With: zbaction.ProcStepArgs{
								"message": "Docker tar: ${out.docker-image-step.artifact}",
							},
						},
					},
				},
			},
		},
	}

	if err := zbaction.RunAction(context.TODO(), action); err != nil {
		panic(err)
	}
}
