package zbaction

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

var ErrRequirementNotMet = errors.New("requirement not met")

type ExprString = string
type CompiledRequirement struct {
	Action Action

	compiledExpr map[ExprString]*vm.Program
}

type Environment struct {
	MatchVersion func(dependency string, constraintString string) bool `expr:"matchVersion"`
}

func CompileActionRequirement(action Action) (*CompiledRequirement, error) {
	cr := &CompiledRequirement{
		Action:       action,
		compiledExpr: make(map[ExprString]*vm.Program, len(action.Requirements)),
	}

	for _, req := range action.Requirements {
		if _, ok := cr.compiledExpr[req.Expr]; ok {
			continue
		}

		cExpr, err := expr.Compile(req.Expr, expr.Env(Environment{})) //nolint:exhaustruct
		if err != nil {
			slog.Error("failed to compile requirement",
				slog.String("requirement", req.String()),
				slog.String("error", err.Error()))
			return nil, fmt.Errorf("compile requirement %s: %w", req.String(), err)
		}

		cr.compiledExpr[req.Expr] = cExpr
	}

	return cr, nil
}

func (cr *CompiledRequirement) CheckRequirement(metadata map[string]any) error {
	env := Environment{
		MatchVersion: func(dependency string, constraintString string) bool {
			return matchSemver(metadata, dependency, constraintString)
		},
	}

	for _, req := range cr.Action.Requirements {
		slog.Debug("checking requirement", slog.String("requirement", req.String()))

		result, err := expr.Run(cr.compiledExpr[req.Expr], env)
		if err != nil {
			slog.Error("failed to evaluate requirement",
				slog.String("requirement", req.String()),
				slog.String("error", err.Error()))
			return err
		}

		slog.Info("requirement result", slog.String("requirement", req.String()), slog.Any("result", result))

		if result, ok := result.(bool); ok && !result {
			slog.Error("requirement not met", slog.String("requirement", req.String()))
			return ErrRequirementNotMet
		}
	}

	return nil
}

func matchSemver(meta map[string]any, dependency string, constraintString string) bool {
	dependencyVersionRaw := accessMeta(meta, dependency)
	if dependencyVersionRaw == nil {
		slog.Error("dependency not found", slog.String("dependency", dependency))
		return false
	}
	dependencyVersion, ok := dependencyVersionRaw.(string)
	if !ok {
		slog.Error("dependency version is invalid",
			slog.String("dependency", dependency),
			slog.Any("dependencyVersion", dependencyVersionRaw))
		return false
	}

	constraint, err := semver.NewConstraint(constraintString)
	if err != nil {
		slog.Error("invalid constraint", slog.String("constraint", constraintString))
		return false
	}

	version, err := semver.NewVersion(dependencyVersion)
	if err != nil {
		slog.Error("unparsable version",
			slog.String("dependency", dependency),
			slog.String("version", dependencyVersion))
		return false
	}

	return constraint.Check(version)
}

func accessMeta(meta map[string]any, key string) any {
	keys := strings.Split(key, ".")
	value := reflect.ValueOf(meta)

	for i, key := range keys {
		if meta == nil {
			return nil
		}

		value = value.MapIndex(reflect.ValueOf(key))
		if !value.IsValid() {
			return nil
		}

		if value.Kind() == reflect.Interface {
			value = value.Elem()
		}
		if value.Kind() != reflect.Map && i != len(keys)-1 {
			return nil // match a not string key
		}
	}

	return value.Interface()
}
