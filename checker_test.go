package zbaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeabur/action/environment"
)

func TestAccessMeta(t *testing.T) {
	softwareList := environment.SoftwareList{
		"go": "1.21",
	}

	assert.Equal(t, "1.21", accessMapByDot(softwareList, "go"))
	assert.Equal(t, nil, accessMapByDot(softwareList, "node"))
	assert.Equal(t, nil, accessMapByDot(softwareList, "go.version"))
}
