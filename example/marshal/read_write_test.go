package marshal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeabur/action"
	"github.com/zeabur/action/example/marshal"
)

func TestReadWrite(t *testing.T) {
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
						Name: "Print hello world with vars",
						RunnableStep: zbaction.CommandStep{
							Command: []string{"echo", "$HELLO_STEP | $HELLO_JOB | $HELLO_ACTION"},
						},
						Variables: map[string]string{
							"HELLO_STEP": "step world",
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

	if err := marshal.Write(action, "action.json"); err != nil {
		t.Fatal(err)
	}

	readAction, err := marshal.Read("action.json")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, action, readAction)
}
