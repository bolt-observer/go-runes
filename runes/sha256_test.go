package runes

import (
	"encoding/base64"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMidStatesEqual(t *testing.T) {
	x := NewSha256()

	mid1 := x.GetMidState()

	x.Write([]byte("random"))

	temp := x.GetMidState()

	if reflect.DeepEqual(temp, mid1) {
		t.Fatalf("Temporary state should not be equal as before")
	}

	err := x.SetMidState(mid1)
	assert.NoError(t, err)

	mid2 := x.GetMidState()

	if !reflect.DeepEqual(mid1, mid2) {
		t.Fatalf("Midstates are not equal")
	}
}

func TestHashChanges(t *testing.T) {
	zeros := "0000000000000000000000000000000000000000000000000000000000000000"
	x := NewSha256()

	s := x.GetMidState()
	for i := 0; i < 8; i++ {
		s.H[i] = 0
	}
	s.Len = 0

	err := x.SetMidState(s)
	assert.NoError(t, err)

	// Setting midstate to 0 should return zero hash
	b := x.GetSum()
	result1 := hex.EncodeToString(b[:])
	assert.Equal(t, zeros, result1)

	// Write does not change hash (until size < 64)
	x.Write([]byte("burek"))
	b = x.GetSum()
	result2 := hex.EncodeToString(b[:])
	assert.Equal(t, zeros, result2)

	// Adding padding finalizes chunk -> hash should change
	x.AddPadding()
	b = x.GetSum()
	result3 := hex.EncodeToString(b[:])
	assert.NotEqual(t, zeros, result3)

	// Another write (there was a bug with another chunk)
	x.Write([]byte("mesni"))
	b = x.GetSum()
	result4 := hex.EncodeToString(b[:])
	assert.Equal(t, result3, result4)

	// And another padding
	x.AddPadding()
	b = x.GetSum()
	result5 := hex.EncodeToString(b[:])
	assert.NotEqual(t, result4, result5)
}

func TestEmpty(t *testing.T) {
	CORRECT := "N0cI__dxndWXnsh11WzSKG9tPPfsMXo7JWMqqyjsN7s="

	x := NewSha256()
	empty := make([]byte, 16)
	x.Write(empty)
	x.AddPadding()

	sum := x.GetSum()
	result := base64.URLEncoding.EncodeToString(sum[:])
	assert.Equal(t, CORRECT, result)

	midsum := x.GetMidState().GetSum()
	result = base64.URLEncoding.EncodeToString(midsum[:])
	assert.Equal(t, CORRECT, result)
}
