package runes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromString(t *testing.T) {
	for _, str := range []string{"374708fff7719dd5979ec875d56cd2286f6d3cf7ec317a3b25632aab28ec37bb:", "6035731a2cbb022cbeb67645aa0f8a26653d8cc454e0e087d4d19d282b8da4bd:=1",
		"4520773407c9658646326fdffe685ffbc3c8639a080dae4310b371830a205cf1:=2-1", "1edf4068e2b0b1e4e075e66751c2d3f5c9fc4515d114f875e6dc6e3e6704efa9:f1=1|f2=3&f3~v1"} {
		r, err := FromString(str)
		assert.NoError(t, err)
		assert.Equal(t, str, r.String())
	}
}

func TestFromBase64(t *testing.T) {
	for _, str := range []string{"N0cI__dxndWXnsh11WzSKG9tPPfsMXo7JWMqqyjsN7s", "YDVzGiy7Aiy-tnZFqg-KJmU9jMRU4OCH1NGdKCuNpL09MQ",
		"RSB3NAfJZYZGMm_f_mhf-8PIY5oIDa5DELNxgwogXPE9Mi0x", "Ht9AaOKwseTgdeZnUcLT9cn8RRXRFPh15txuPmcE76lmMT0xfGYyPTMmZjN-djE", "EMXekLFLz2z-I7bEOBkfQmR5bR_V78iaf-L-LeFu8Mc9MA",
		"hamQrEOb90nw5-BAZjrKCjIzniVfvi9nREJa_A6Tsx09MCZtZXRob2RebGlzdHxtZXRob2ReZ2V0fG1ldGhvZD1zdW1tYXJ5Jm1ldGhvZC9saXN0ZGF0YXN0b3JlJm1ldGhvZF5saXN0fG1ldGhvZF5nZXR8bWV0aG9kPXN1bW1hcnkmbWV0aG9kL2xpc3RkYXRhc3RvcmU"} {
		r, err := FromBase64(str)
		assert.NoError(t, err)
		assert.Equal(t, str, r.ToBase64Internal(true))
	}
}
