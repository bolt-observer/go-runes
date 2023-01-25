package runes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLower(t *testing.T) {
	a := 123.45
	b := 12444

	ret, err := isLower(a, b)
	assert.NoError(t, err)
	assert.Equal(t, true, ret)
}

func TestIsEqual(t *testing.T) {
	ret, err := isEqual(uint16(12), 12.0)
	assert.NoError(t, err)
	assert.Equal(t, true, ret)

	ret, err = isEqual("true", true)
	assert.NoError(t, err)
	assert.Equal(t, true, ret)
}

func TestMakeAlternativeFromString(t *testing.T) {
	resp, rest, err := MakeAlternativeFromString("ab<cd|", false)
	assert.NoError(t, err)

	assert.Equal(t, "<", resp.Cond)
	assert.Equal(t, "ab", resp.Field)
	assert.Equal(t, "cd", resp.Value)
	assert.Equal(t, "", rest)
}

func TestSubstring(t *testing.T) {
	resp, err := MakeAlternative("burek", "~", "mesni", false)
	assert.NoError(t, err)

	eval, _ := resp.Evaluate(map[string]any{"burek": "mesni burek"})
	assert.Equal(t, true, eval)

	eval, _ = resp.Evaluate(map[string]any{"burek": "sirni burek"})
	assert.NotEqual(t, true, eval)
}
