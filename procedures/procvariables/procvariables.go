// Package procvariables allows you finding the key of variables procedures use.
//
// Do not import other procedures in your library, or it will be "registered" (init) in our registry.
package procvariables

import (
	"os"

	zbaction "github.com/zeabur/action"
)

// WithProcVariables injects variables for procedures into the action.
func WithProcVariables(key string, value string) zbaction.ExecutorOptionsFn {
	return func(o *zbaction.ExecutorOptions) {
		if o.RuntimeVariables == nil {
			o.RuntimeVariables = make(map[string]string)
		}

		o.RuntimeVariables[key] = value
	}
}

// WithEnvBuildkitHost injects the buildkit address from environment variable into the action.
func WithEnvBuildkitHost() zbaction.ExecutorOptionsFn {
	host := os.Getenv("BUILDKIT_HOST")
	if host == "" {
		panic("buildkit host is not set")
	}

	return WithProcVariables(VarBuildkitHostKey, host)
}
