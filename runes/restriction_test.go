package runes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestrictionSubstring(t *testing.T) {
	a, err := MakeAlternative("burek", "~", "mesni", false)
	assert.NoError(t, err)

	resp, err := MakeRestriction([]Alternative{*a})
	assert.NoError(t, err)

	eval, _ := resp.Evaluate(map[string]any{"burek": "mesni burek"})
	assert.Equal(t, true, eval)

	eval, _ = resp.Evaluate(map[string]any{"burek": "sirni burek"})
	assert.NotEqual(t, true, eval)
}
