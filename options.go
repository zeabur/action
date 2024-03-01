package zbaction

import (
	"io"
	"maps"
	"os"
	"strings"
)

type ExecutorOptionsFn func(*ExecutorOptions)
type ExecutorOptions struct {
	RuntimeVariables map[string]string
	Stdout           io.Writer
	Stderr           io.Writer
}

// WithRuntimeVariables injects custom runtime variables into the action.
//
// Its behavior is similar to adding variables to the Variables of the action yourself,
// but it will not affect the original definition of the action.
func WithRuntimeVariables(vars map[string]string) ExecutorOptionsFn {
	return func(o *ExecutorOptions) {
		if o.RuntimeVariables != nil {
			for k, v := range vars {
				o.RuntimeVariables[k] = v
			}
		} else {
			o.RuntimeVariables = maps.Clone(vars)
		}
	}
}

// WithCurrentEnvironmentVariable injects the current environment variables into the action.
//
// Useful if the application needs to access the environment variables.
func WithCurrentEnvironmentVariable() ExecutorOptionsFn {
	envs := os.Environ()
	envMap := make(map[string]string, len(envs))

	for _, env := range envs {
		k, v, ok := strings.Cut(env, "=")
		if !ok {
			continue
		}

		envMap[k] = v
	}

	return WithRuntimeVariables(envMap)
}

// WithCustomStdout injects a custom stdout writer into the action.
func WithCustomStdout(w io.Writer) ExecutorOptionsFn {
	return func(o *ExecutorOptions) {
		o.Stdout = w
	}
}

// WithCustomStderr injects a custom stderr writer into the action.
func WithCustomStderr(w io.Writer) ExecutorOptionsFn {
	return func(o *ExecutorOptions) {
		o.Stderr = w
	}
}
