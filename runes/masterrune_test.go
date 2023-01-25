package runes

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	var empty [16]byte

	correct := "374708fff7719dd5979ec875d56cd2286f6d3cf7ec317a3b25632aab28ec37bb:"

	rune, err := MakeMasterRune(empty[:], nil, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, correct, rune.String())
	hash := hex.EncodeToString(rune.GetAuthCode())
	assert.Equal(t, true, strings.HasPrefix(correct, hash))

	r, _, err := MakeRestrictionFromString("burek=3", false)
	assert.NoError(t, err)
	rune.Rune.AddRestriction(*r)
	assert.NotEqual(t, correct, rune.String())

	hash = hex.EncodeToString(rune.GetAuthCode())
	assert.Equal(t, false, strings.HasPrefix(correct, hash))
}

func TestIsRuneAuthorized(t *testing.T) {
	var empty [16]byte
	rune, err := MakeMasterRune(empty[:], nil, nil, nil)
	assert.NoError(t, err)

	r, _, err := MakeRestrictionFromString("burek=3", false)
	assert.NoError(t, err)
	restricted, err := rune.GetRestricted(r)
	assert.NoError(t, err)

	assert.Equal(t, true, rune.IsRuneAuthorized(restricted))

	s := restricted.String()
	s = strings.ReplaceAll(s, "burek=3", "burek=4")
	fake, err := FromString(s)
	assert.NoError(t, err)

	assert.Equal(t, false, rune.IsRuneAuthorized(fake))

}
