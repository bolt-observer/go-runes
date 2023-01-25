package runes

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

// Rune struct
type Rune struct {
	Sha256       *Sha256
	Restrictions []Restriction
}

// AddRestriction adds a new restriction
func (r *Rune) AddRestriction(restriction Restriction) error {
	r.Restrictions = append(r.Restrictions, restriction)

	// This way always a complete chunk is added
	_, err := r.Sha256.Write([]byte(restriction.String()))
	if err != nil {
		return err
	}
	err = r.Sha256.AddPadding()
	if err != nil {
		return err
	}

	return nil
}

// GetAuthCode gets the auth code
func (r *Rune) GetAuthCode() []byte {
	b := r.Sha256.GetSum()
	return b[:]
}

// MakeRune creates a new rune
func MakeRune(authbase []byte, uniqueid, version any, restrictions []Restriction) (*Rune, error) {
	rest := make([]Restriction, 0, len(restrictions))
	ret := &Rune{
		Sha256:       NewSha256(),
		Restrictions: rest,
	}

	/* Authcode is always 64 bytes */
	midState := &MidState{}
	for i := 0; i < 8; i++ {
		authbase, midState.H[i] = consumeUint32(authbase)
	}
	midState.Len = 64
	err := ret.Sha256.SetMidState(midState)
	if err != nil {
		return nil, err
	}

	if uniqueid != nil {
		u, err := UniqueID(uniqueid, version)
		if err != nil {
			return nil, err
		}
		err = ret.AddRestriction(*u)
		if err != nil {
			return nil, err
		}
	}

	for _, r := range restrictions {
		err = ret.AddRestriction(r)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func padLen(x uint64) uint64 {
	return (64 - (x % 64)) % 64
}

// FromAuthCode create a new rune from auth code
func FromAuthCode(authcode []byte, restrictions []Restriction) (*Rune, error) {
	ret, err := MakeRune(authcode, nil, nil, []Restriction{})
	if err != nil {
		return nil, err
	}

	runelength := ret.Sha256.GetMidState().Len

	for _, r := range restrictions {
		runelength += uint64(len(r.String()))
		runelength += padLen(runelength)
	}

	midState := ret.Sha256.GetMidState()
	midState.Len = runelength
	ret.Sha256.SetMidState(midState)

	ret.Restrictions = restrictions

	return ret, nil
}

// Evaluate evaluates the rune
func (r *Rune) Evaluate(vals map[string]any) (bool, string) {
	for _, one := range r.Restrictions {
		ok, msg := one.Evaluate(vals)
		if !ok {
			return false, msg
		}
	}

	return true, ""
}

// String returns a string representation of rune
func (r *Rune) String() string {
	rest := make([]string, 0, len(r.Restrictions))
	for _, one := range r.Restrictions {
		rest = append(rest, one.String())
	}

	return hex.EncodeToString(r.GetAuthCode()) + ":" + strings.Join(rest, "&")
}

// ToBase64 returns the base64 encoded representation of rune
func (r *Rune) ToBase64() string {
	rest := make([]string, 0, len(r.Restrictions))
	for _, one := range r.Restrictions {
		rest = append(rest, one.String())
	}

	s := strings.Join(rest, "&")
	b := new(bytes.Buffer)
	b.Write(r.GetAuthCode())
	b.WriteString(s)

	return base64.URLEncoding.EncodeToString(b.Bytes())
}

// FromString returns a new rune from string representation
func FromString(str string) (*Rune, error) {
	if len(str) < 64 || str[64] != ':' {
		return nil, fmt.Errorf("rune strings must start with 64 hex digits then '-'")
	}

	authcode, err := hex.DecodeString(str[0:64])
	if err != nil {
		return nil, err
	}

	rest := str[65:]
	restrictions := make([]Restriction, 0)

	var restriction *Restriction
	for len(rest) > 0 {
		allowIDField := len(restrictions) == 0

		restriction, rest, err = MakeRestrictionFromString(rest, allowIDField)
		if err != nil {
			return nil, err
		}

		restrictions = append(restrictions, *restriction)
	}

	return FromAuthCode(authcode, restrictions)
}

// FromBase64 returns a new rune from base64 encoded string representation
func FromBase64(str string) (*Rune, error) {
	data, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}

	return FromString(hex.EncodeToString(data[:32]) + ":" + string(data[32:]))
}

// GetRestricted obtains a restricted rune
func (r *Rune) GetRestricted(restrictions ...*Restriction) (*Rune, error) {
	rune, err := MakeRune(r.GetAuthCode(), nil, nil, nil)
	if err != nil {
		return nil, err
	}
	rune.Restrictions = r.Restrictions

	for _, r := range restrictions {
		rune.AddRestriction(*r)
	}

	return rune, nil
}

// Check checks rune
func (r *Rune) Check(rune *Rune, vals map[string]any) error {
	ok, msg := rune.Evaluate(vals)
	if ok {
		return nil
	}

	return fmt.Errorf(msg)
}
