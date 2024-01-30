package zbaction

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestAccessMeta(t *testing.T) {
	meta := map[string]any{
		"sdk": map[string]any{
			"go": "1.21",
		},
		"base": "alpine",
	}

	assert.Equal(t, "1.21", accessMeta(meta, "sdk.go"))
	assert.DeepEqual(t, map[string]any{
		"go": "1.21",
	}, accessMeta(meta, "sdk"))
	assert.Equal(t, "alpine", accessMeta(meta, "base"))
	assert.Equal(t, nil, accessMeta(meta, "sdk.node"))
	assert.Equal(t, nil, accessMeta(meta, "sdk.go.version"))
}
