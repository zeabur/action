package main

import (
	"context"

	zbaction "github.com/zeabur/action"
	_ "github.com/zeabur/action/procedures"
	_ "github.com/zeabur/action/procedures/artifact"
)

func main() {
	err := zbaction.RunAction(context.TODO(), zbaction.Action{
		Jobs: []zbaction.Job{
			{
				Steps: []zbaction.Step{
					{
						RunnableStep: zbaction.ProcStep{
							Uses: "action/not-exist",
							With: zbaction.ProcStepArgs{},
						},
					},
				},
			},
		},
	})

	if err != nil {
		panic(err)
	}
}
