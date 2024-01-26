package zbaction

import (
	"fmt"
	"log/slog"
	"sync"
)

type ProcedureStep interface {
	RunnableStep
}

type ProcedureStepBuilder func(args ProcStepArgs) (ProcedureStep, error)

var resolver = NewProcedureStepResolver()

type ProcedureStepResolver interface {
	Register(name ProcStepName, builder ProcedureStepBuilder)
	Resolve(uses ProcStepName, with ProcStepArgs) (ProcedureStep, error)
}

type procedureStepResolver struct {
	registry map[ProcStepName]ProcedureStepBuilder
	mutex    sync.RWMutex
}

func NewProcedureStepResolver() ProcedureStepResolver {
	return &procedureStepResolver{
		registry: make(map[ProcStepName]ProcedureStepBuilder),
	}
}

func (p *procedureStepResolver) Register(name ProcStepName, builder ProcedureStepBuilder) {
	slog.Debug("Registering procedure step", slog.String("name", name))

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, ok := p.registry[name]; ok {
		slog.Error("Namespace conflict. Remove the conflicted module first.", slog.String("name", name))
		panic("namespace conflict")
	}

	p.registry[name] = builder
}

func (p *procedureStepResolver) Resolve(uses ProcStepName, with ProcStepArgs) (ProcedureStep, error) {
	slog.Debug("Resolving procedure step", slog.String("uses", uses), slog.Any("with", with))

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if builder, ok := p.registry[uses]; ok {
		step, err := builder(with)
		if err != nil {
			return nil, fmt.Errorf("build step %s: %w", uses, err)
		}

		return step, nil
	}

	return nil, fmt.Errorf("no procedure step builder found for %s", uses)
}

func RegisterProcedure(name ProcStepName, builder ProcedureStepBuilder) {
	resolver.Register(name, builder)
}

func ResolveProcedure(uses ProcStepName, with ProcStepArgs) (ProcedureStep, error) {
	return resolver.Resolve(uses, with)
}
