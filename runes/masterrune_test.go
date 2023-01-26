package runes

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	var empty [16]byte

	// 16 * 0 byte is known test vector
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

func TestSecretLength(t *testing.T) {
	long := make([]byte, 56)
	_, err := rand.Read(long)
	assert.NoError(t, err)
	_, err = MakeMasterRune(long, nil, nil, nil)
	assert.ErrorIs(t, err, ErrTooLongSecret)

	_, err = MakeMasterRune(long[:55], nil, nil, nil)
	assert.NoError(t, err)

	_, err = MakeMasterRune(long[:0], nil, nil, nil)
	assert.ErrorIs(t, err, ErrTooShortSecret)
}

func TestIsRuneAuthorized(t *testing.T) {
	secret := make([]byte, 55)
	_, err := rand.Read(secret)
	assert.NoError(t, err)

	rune := MustMakeMasterRune(secret)

	// Simple case
	r, _, err := MakeRestrictionFromString("burek=3", false)
	assert.NoError(t, err)
	restricted, err := rune.GetRestricted(*r)
	assert.NoError(t, err)
	assert.Equal(t, true, rune.IsRuneAuthorized(restricted))

	// Tampering with rune fails
	s := restricted.String()
	s = strings.ReplaceAll(s, "burek=3", "burek=4")
	fake, err := FromString(s)
	assert.NoError(t, err)
	assert.Equal(t, false, rune.IsRuneAuthorized(fake))

	// Further restricting restricted rune works
	r, _, err = MakeRestrictionFromString("mesni{2", false)
	assert.NoError(t, err)
	restricted2, err := restricted.GetRestricted(*r)
	assert.NoError(t, err)
	s = restricted2.String()
	assert.Equal(t, true, rune.IsRuneAuthorized(restricted2))

	// Tampering with restricted rune also fails
	s = strings.ReplaceAll(s, "{", "}")
	fake, err = FromString(s)
	assert.NoError(t, err)
	assert.Equal(t, false, rune.IsRuneAuthorized(fake))

	// Check returns unauthorized for tampered rune
	err = rune.Check(fake, map[string]any{"burek": 3, "mesni": "1"})
	assert.ErrorIs(t, err, ErrUnauthorizedRune)

	// Check returns error for wrong conditions
	err = rune.Check(restricted2, map[string]any{"burek": 4, "mesni": "4"})
	assert.Error(t, err)

	// But still works for correct ones
	err = rune.Check(restricted2, map[string]any{"burek": 3.0, "mesni": 1})
	assert.NoError(t, err)
}

func TestRealRune(t *testing.T) {
	rune, err := FromBase64("EMXekLFLz2z-I7bEOBkfQmR5bR_V78iaf-L-LeFu8Mc9MA")
	assert.NoError(t, err)

	restricted, err := FromBase64("uxzKjgrPj6rwr0ySqNP--p2ggNmnb7URM0Awj-Zr56E9MCZtZXRob2RebGlzdHxtZXRob2ReZ2V0fG1ldGhvZD1zdW1tYXJ5Jm1ldGhvZC9saXN0ZGF0YXN0b3Jl")
	assert.NoError(t, err)

	fresh := rune.MustGetRestrictedFromString("method^list|method^get|method=summary&method/listdatastore")

	assert.Equal(t, restricted.String(), fresh.String())
}
