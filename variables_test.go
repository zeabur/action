package zbaction_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeabur/builder/zbaction"
)

func TestMapContainer_GetRawVariable(t *testing.T) {
	varExample := map[string]string{
		"test":  "test",
		"test2": "${test}",
	}
	mc := zbaction.NewMapContainer(varExample)

	v, ok := mc.GetRawVariable("test")
	assert.True(t, ok)
	assert.Equal(t, "test", v)

	v, ok = mc.GetRawVariable("test2")
	assert.True(t, ok)
	assert.Equal(t, "${test}", v)

	v, ok = mc.GetRawVariable("test3")
	assert.False(t, ok)
}

func TestVariableContainerWithParent_GetRawVariable(t *testing.T) {
	parentExample := map[string]string{
		"test":  "test",
		"test2": "${test}",
	}
	pmc := zbaction.NewMapContainer(parentExample)

	thisExample := map[string]string{
		"test3": "${test2}",
		"test4": "hi",
	}
	tmc := zbaction.NewMapContainer(thisExample)

	mc := zbaction.NewVariableContainerWithParent(tmc, pmc)

	v, ok := mc.GetRawVariable("test")
	assert.True(t, ok)
	assert.Equal(t, "test", v)

	v, ok = mc.GetRawVariable("test2")
	assert.True(t, ok)
	assert.Equal(t, "${test}", v)

	v, ok = mc.GetRawVariable("test3")
	assert.True(t, ok)
	assert.Equal(t, "${test2}", v)

	v, ok = mc.GetRawVariable("test4")
	assert.True(t, ok)
	assert.Equal(t, "hi", v)

	v, ok = mc.GetRawVariable("test5")
	assert.False(t, ok)
}

func TestVariableContainerWithParent_GetVariable(t *testing.T) {
	varExample := map[string]string{
		"test":  "test",
		"test2": "${test}1",
		"test3": "${test2}2",
	}
	mc := zbaction.NewMapContainer(varExample)

	v, ok := mc.GetVariable("test")
	assert.True(t, ok)
	assert.Equal(t, "test", v)

	v, ok = mc.GetVariable("test2")
	assert.True(t, ok)
	assert.Equal(t, "test1", v)

	v, ok = mc.GetVariable("test3")
	assert.True(t, ok)
	assert.Equal(t, "test12", v)

	v, ok = mc.GetVariable("test4")
	assert.False(t, ok)
}

func TestVariableContainerWithParent_GetVariableWithParent(t *testing.T) {
	parentExample := map[string]string{
		"test":  "test",
		"test2": "${test}1",
	}
	pmc := zbaction.NewMapContainer(parentExample)

	thisExample := map[string]string{
		"test2": "${test}2",
		"test3": "${test2}3",
	}
	tmc := zbaction.NewMapContainer(thisExample)

	mc := zbaction.NewVariableContainerWithParent(tmc, pmc)

	v, ok := mc.GetVariable("test")
	assert.True(t, ok)
	assert.Equal(t, "test", v)

	v, ok = mc.GetVariable("test2")
	assert.True(t, ok)
	assert.Equal(t, "test2", v)

	v, ok = mc.GetVariable("test3")
	assert.True(t, ok)
	assert.Equal(t, "test23", v)

	v, ok = mc.GetVariable("test4")
	assert.False(t, ok)
}
