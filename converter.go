package zbaction

import (
	"fmt"

	"github.com/zeabur/action/proto"
)

func ActionToProto(action Action) (*proto.Action, error) {
	p := &proto.Action{
		Id:           action.ID,
		Jobs:         make([]*proto.Job, len(action.Jobs)),
		Variables:    action.Variables,
		Requirements: make([]*proto.Requirement, len(action.Requirements)),
	}

	for requirementIndex, requirement := range action.Requirements {
		p.Requirements[requirementIndex] = &proto.Requirement{
			Expr:        requirement.Expr,
			Description: requirement.Description,
			Required:    requirement.Required,
		}
	}

	for jobIndex, job := range action.Jobs {
		pj := &proto.Job{
			Id:        job.ID,
			RunOn:     job.RunOn,
			Steps:     make([]*proto.Step, len(job.Steps)),
			Variables: job.Variables,
		}

		for stepIndex, step := range job.Steps {
			ps := &proto.Step{
				Id:        step.ID,
				Name:      step.Name,
				Step:      nil,
				Variables: step.Variables,
			}

			if err := exactStepToProto(step, ps); err != nil {
				return nil, fmt.Errorf("failed to convert step %s to proto: %w", step.ID, err)
			}

			pj.Steps[stepIndex] = ps
		}

		p.Jobs[jobIndex] = pj
	}

	return p, nil
}

func ActionFromProto(p *proto.Action) (Action, error) {
	action := Action{
		ID:           p.Id,
		Jobs:         make([]Job, len(p.Jobs)),
		Variables:    p.Variables,
		Requirements: make([]Requirement, len(p.Requirements)),
	}

	for requirementIndex, requirement := range p.Requirements {
		action.Requirements[requirementIndex] = Requirement{
			Expr:        requirement.Expr,
			Description: requirement.Description,
			Required:    requirement.Required,
		}
	}

	for jobIndex, job := range p.Jobs {
		j := Job{
			ID:        job.Id,
			RunOn:     job.RunOn,
			Steps:     make([]Step, len(job.Steps)),
			Variables: job.Variables,
		}

		for stepIndex, step := range job.Steps {
			s, err := exactStepFromProto(step)
			if err != nil {
				return Action{}, fmt.Errorf("failed to convert step %s from proto: %w", step.Id, err)
			}

			j.Steps[stepIndex] = Step{
				ID:           step.Id,
				Name:         step.Name,
				RunnableStep: s,
				Variables:    step.Variables,
			}
		}

		action.Jobs[jobIndex] = j
	}

	return action, nil
}

func exactStepToProto(step Step, out *proto.Step) error {
	runnableStep := step.RunnableStep

	switch runnableStep := runnableStep.(type) {
	case CommandStep:
		out.Step = &proto.Step_Command{
			Command: &proto.CommandStep{
				Command: runnableStep.Command,
			},
		}
	case ProcStep:
		out.Step = &proto.Step_Proc{
			Proc: &proto.ProcStep{
				Uses: runnableStep.Uses,
				With: runnableStep.With,
			},
		}
	default:
		return fmt.Errorf("unknown step type received: %T (%+v)", runnableStep, runnableStep)
	}

	return nil
}

func exactStepFromProto(p *proto.Step) (RunnableStep, error) {
	var step RunnableStep

	switch p := p.Step.(type) {
	case *proto.Step_Command:
		step = CommandStep{
			Command: p.Command.Command,
		}
	case *proto.Step_Proc:
		step = ProcStep{
			Uses: p.Proc.Uses,
			With: p.Proc.With,
		}
	default:
		return Step{}, fmt.Errorf("unknown step type received: %T (%+v)", p, p)
	}

	return step, nil
}
