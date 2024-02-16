package zbaction

import (
	"os"
	"strings"
)

type VariableContainer interface {
	GetRawVariable(key string) (string, bool)
	GetVariable(key string) (string, bool)
	ListRawVariables() map[string]string
}

type mapContainer struct {
	values map[string]string
}

type variableContainerWithParent struct {
	parent VariableContainer
	this   VariableContainer
}

type variableContainerWithExtraParameters struct {
	extra  map[string]string
	parent VariableContainer // prioritized
}

func NewMapContainer(values map[string]string) VariableContainer {
	return &mapContainer{
		values: values,
	}
}

func NewVariableContainerWithParent(this VariableContainer, parent VariableContainer) VariableContainer {
	return variableContainerWithParent{
		parent: parent,
		this:   this,
	}
}

func NewVariableContainerWithExtraParameters(extra map[string]string, parent VariableContainer) VariableContainer {
	return variableContainerWithExtraParameters{
		extra:  extra,
		parent: parent,
	}
}

func (m mapContainer) GetRawVariable(key string) (string, bool) {
	v, ok := m.values[key]
	return v, ok
}

func (m mapContainer) GetVariable(key string) (string, bool) {
	v, ok := m.GetRawVariable(key)
	if !ok {
		return "", false
	}

	return expandValue(v, m.GetVariable), true
}

func (m mapContainer) ListRawVariables() map[string]string {
	if m.values == nil {
		return map[string]string{}
	}

	return m.values
}

func (m variableContainerWithParent) GetRawVariable(key string) (string, bool) {
	v, ok := m.this.GetRawVariable(key)
	if ok {
		return v, true
	}
	return m.parent.GetRawVariable(key)
}

func (m variableContainerWithParent) GetVariable(key string) (string, bool) {
	v, ok := m.GetRawVariable(key)
	if !ok {
		return "", false
	}

	return expandValue(v, m.GetVariable), true
}

func (m variableContainerWithParent) ListRawVariables() map[string]string {
	parent := m.parent.ListRawVariables()
	this := m.this.ListRawVariables()

	merged := make(map[string]string, len(parent)+len(this))
	for k, v := range parent {
		merged[k] = v
	}
	for k, v := range this {
		merged[k] = v
	}

	return merged
}

func (v variableContainerWithExtraParameters) GetRawVariable(key string) (string, bool) {
	if value, ok := v.parent.GetRawVariable(key); ok {
		return value, true
	}

	if v.extra != nil {
		if v, ok := v.extra[key]; ok {
			return v, true
		}
	}

	return "", false
}

func (v variableContainerWithExtraParameters) GetVariable(key string) (string, bool) {
	value, ok := v.GetRawVariable(key)
	if !ok {
		return "", false
	}

	return expandValue(value, v.GetVariable), true
}

func (v variableContainerWithExtraParameters) ListRawVariables() map[string]string {
	parent := v.parent.ListRawVariables()
	merged := make(map[string]string, len(v.extra)+len(parent))

	for k, v := range parent {
		merged[k] = v
	}
	for k, v := range v.extra {
		merged[k] = v
	}

	return merged
}

func expandValue(currentValue string, getNextExpandedVariableFn func(referencedKey string) (string, bool)) string {
	return os.Expand(currentValue, func(keyReference string) string {
		if keyReference == currentValue {
			return "" // cycle detected
		}

		v, ok := getNextExpandedVariableFn(keyReference)
		if ok {
			return v
		}
		return ""
	})
}

type EnvironmentVariables map[string]string

func ListEnvironmentVariables(vc VariableContainer) EnvironmentVariables {
	if vc == nil {
		return EnvironmentVariables{}
	}

	allVariables := vc.ListRawVariables()
	filteredVariables := make(EnvironmentVariables)

	for key, value := range allVariables {
		// if key contains `.`, we consider it as an internal variable and skip it
		if strings.ContainsRune(key, '.') {
			continue
		}

		filteredVariables[strings.ToUpper(key)] = expandValue(value, vc.GetVariable)
	}

	return filteredVariables
}

func (v EnvironmentVariables) ToList() []string {
	var list []string
	for k, v := range v {
		list = append(list, k+"="+v)
	}
	return list
}
