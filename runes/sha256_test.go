package runes

import (
	"encoding/base64"
	"reflect"
	"testing"
)

func TestMidStatesEqual(t *testing.T) {
	x := NewSha256()

	mid1 := x.GetMidState()

	x.Write([]byte("random"))

	temp := x.GetMidState()

	if reflect.DeepEqual(temp, mid1) {
		t.Fatalf("Temporary state should not be equal as before")
	}

	x.SetMidState(mid1)

	mid2 := x.GetMidState()

	if !reflect.DeepEqual(mid1, mid2) {
		t.Fatalf("Midstates are not equal")
	}
}

func TestEmpty(t *testing.T) {
	CORRECT := "N0cI__dxndWXnsh11WzSKG9tPPfsMXo7JWMqqyjsN7s="

	x := NewSha256()
	empty := make([]byte, 16)
	x.Write(empty)
	x.AddPadding()

	sum := x.GetSum()
	result := base64.URLEncoding.EncodeToString(sum[:])
	if result != CORRECT {
		t.Fatalf("Correct %s vs. %s", CORRECT, result)
	}

	midsum := x.GetMidState().GetSum()
	result = base64.URLEncoding.EncodeToString(midsum[:])
	if result != CORRECT {
		t.Fatalf("Correct %s vs. %s", CORRECT, result)
	}
}
