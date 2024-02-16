package zbaction_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	zbaction "github.com/zeabur/action"
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

	_, ok = mc.GetRawVariable("test3")
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

	_, ok = mc.GetRawVariable("test5")
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

	_, ok = mc.GetVariable("test4")
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

	_, ok = mc.GetVariable("test4")
	assert.False(t, ok)
}

func TestListEnvironmentVariable(t *testing.T) {
	mc := zbaction.NewMapContainer(map[string]string{
		"test":     "test",
		"test2":    "${test}",
		"out.test": "1234",
	})

	lec := zbaction.ListEnvironmentVariables(mc)
	assert.Len(t, lec, 2)
	assert.Contains(t, lec, "test")
	assert.Contains(t, lec, "test2")
	assert.NotContains(t, lec, "out.test")
}

func TestListEnvironmentVariable_Nil(t *testing.T) {
	lec := zbaction.ListEnvironmentVariables(nil)
	assert.Len(t, lec, 0)
}

func TestListEnvironmentVariable_ParentNil(t *testing.T) {
	lec := zbaction.ListEnvironmentVariables(
		zbaction.NewVariableContainerWithExtraParameters(map[string]string{
			"test": "1",
		}, zbaction.NewMapContainer(nil)),
	)

	assert.Len(t, lec, 1)
	assert.Equal(t, "1", lec["test"])
}

func TestEnvironmentVariables_ToList(t *testing.T) {
	lec := zbaction.EnvironmentVariables{
		"test":  "1",
		"test2": "2",
	}

	assert.ElementsMatch(t, []string{"test=1", "test2=2"}, lec.ToList())
}
