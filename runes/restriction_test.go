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

func TestTwoRestrictions(t *testing.T) {
	a, err := MakeAlternative("burek", "=", "3", false)
	assert.NoError(t, err)
	b, err := MakeAlternative("mesni", "<", 2, false)
	assert.NoError(t, err)

	resp, err := MakeRestriction([]Alternative{*a, *b})
	assert.NoError(t, err)

	eval, _ := resp.Evaluate(map[string]any{"burek": 3, "mesni": "1"})
	assert.Equal(t, true, eval)

	eval, _ = resp.Evaluate(map[string]any{"burek": "3", "mesni": 2})
	assert.Equal(t, false, eval)
}

func TestTwoRestrictionsMixedFail(t *testing.T) {
	a, err := MakeAlternative("burek", "=", "3", false)
	assert.NoError(t, err)
	b, err := MakeAlternative("mesni", "<", 2, false)
	assert.NoError(t, err)

	resp, err := MakeRestriction([]Alternative{*a, *b})
	assert.NoError(t, err)

	eval, _ := resp.Evaluate(map[string]any{"burek": "4", "mesni": 1})
	assert.Equal(t, false, eval)

	eval, _ = resp.Evaluate(map[string]any{"burek": "3", "mesni": 2})
	assert.Equal(t, false, eval)

	eval, _ = resp.Evaluate(map[string]any{"burek": "3", "mesni": "1"})
	assert.Equal(t, true, eval)
}
