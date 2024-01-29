package main

import (
	"context"
	"log/slog"
	"os"

	zbaction "github.com/zeabur/action"
	_ "github.com/zeabur/action/procedures"
)

func main() {
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(h))

	action := zbaction.Action{
		ID: "echo-hello-world",
		Jobs: []zbaction.Job{
			{
				ID: "print-hello-world",
				Steps: []zbaction.Step{
					{
						Name: "Print hello world",
						RunnableStep: zbaction.CommandStep{
							Command: []string{"echo", "hello world from echo!"},
						},
					},
				},
			},
			{
				ID: "print-hello-world-with-procedure",
				Steps: []zbaction.Step{
					{
						Name: "Print hello world with procedure",
						RunnableStep: zbaction.ProcStep{
							Uses: "builtin/echo",
							With: zbaction.ProcStepArgs{
								"message": "hello world from procedure!",
							},
						},
					},
				},
			},
			{
				ID: "print-hello-world-with-vars",
				Steps: []zbaction.Step{
					{
						ID:   "print-hello-world-with-vars",
						Name: "Print hello world with vars",
						RunnableStep: zbaction.CommandStep{
							Command: []string{"echo", "$HELLO_STEP | $HELLO_JOB | $HELLO_ACTION"},
						},
						Variables: map[string]string{
							"HELLO_STEP": "step world",
						},
					},
					{
						Name: "Print the output again",
						RunnableStep: zbaction.CommandStep{
							Command: []string{"echo", "${out.print-hello-world-with-vars.stdout}"},
						},
					},
				},
				Variables: map[string]string{
					"HELLO_JOB": "job world",
				},
			},
		},
		Variables: map[string]string{
			"HELLO_ACTION": "action world",
		},
	}

	if err := zbaction.RunAction(context.TODO(), action); err != nil {
		panic(err)
	}
}
