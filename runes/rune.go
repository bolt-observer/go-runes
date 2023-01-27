package runes

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
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

	if len(authbase) != 32 {
		return nil, fmt.Errorf("invalid authbase length %d", len(authbase))
	}

	midState := &MidState{}
	for i := 0; i < 8; i++ {
		authbase, midState.H[i] = consumeUint32(authbase)
	}
	err := ret.Sha256.SetMidState(midState)
	if err != nil {
		return nil, err
	}
	ret.Sha256.SetLen(64)

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

	runelength := ret.Sha256.GetLen()
	for _, r := range restrictions {
		runelength += uint64(len(r.String()))
		runelength += padLen(runelength)
	}
	ret.Sha256.SetLen(runelength)

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
	return r.ToBase64Internal(true)
}

// ToBase64Internal returns the base64 encoded representation of rune
func (r *Rune) ToBase64Internal(trim bool) string {
	rest := make([]string, 0, len(r.Restrictions))
	for _, one := range r.Restrictions {
		rest = append(rest, one.String())
	}

	s := strings.Join(rest, "&")
	b := new(bytes.Buffer)
	b.Write(r.GetAuthCode())
	b.WriteString(s)

	ret := base64.URLEncoding.EncodeToString(b.Bytes())
	if trim {
		ret = strings.TrimRight(ret, "=")
	}

	return ret
}

// MustGetFromString returns a new rune from string representation
func MustGetFromString(str string) Rune {
	ret, err := FromString(str)
	if err != nil {
		panic(err)
	}
	return *ret
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

// MustGetFromBase64 returns a new rune from base64 representation
func MustGetFromBase64(str string) Rune {
	ret, err := FromBase64(str)
	if err != nil {
		panic(err)
	}

	return *ret
}

// FromBase64 returns a new rune from base64 encoded string representation
func FromBase64(str string) (*Rune, error) {
	str = strings.TrimRight(str, "=")
	addendum := strings.Repeat("=", (4-(len(str)%4))%4)

	data, err := base64.URLEncoding.DecodeString(str + addendum)
	if err != nil {
		return nil, err
	}

	return FromString(hex.EncodeToString(data[:32]) + ":" + string(data[32:]))
}

// GetRestricted obtains a restricted rune
func (r *Rune) GetRestricted(restrictions ...Restriction) (*Rune, error) {
	//rune, err := MakeRune(r.GetAuthCode(), nil, nil, nil)
	rune, err := FromAuthCode(r.GetAuthCode(), r.Restrictions)
	if err != nil {
		return nil, err
	}

	for _, r := range restrictions {
		rune.AddRestriction(r)
	}

	return rune, nil
}

// MustGetRestrictedFromString obtains a restricted rune
func (r *Rune) MustGetRestrictedFromString(str string) Rune {
	ret, err := r.GetRestricted(MustMakeRestrictionsFromString(str)...)
	if err != nil {
		panic(err)
	}

	return *ret
}

// Check checks rune
func (r *Rune) Check(vals map[string]any) error {
	ok, msg := r.Evaluate(vals)
	if ok {
		return nil
	}

	return fmt.Errorf(msg)
}

func (r *Rune) getID() string {
	if len(r.Restrictions) < 1 {
		return ""
	}
	if len(r.Restrictions[0].Alternatives) < 1 {
		return ""
	}

	// uniqueID restriction is the first one by definition
	a := r.Restrictions[0].Alternatives[0]
	if a.Field == "" && a.Cond == "=" {
		return fmt.Sprintf("%v", a.Value)
	}

	return ""
}

func (r *Rune) getIDPart(num int) int {
	s := r.getID()
	split := strings.Split(s, "-")
	if num >= len(split) {
		return -1
	}

	data, err := strconv.Atoi(split[num])
	if err != nil {
		return -1
	}

	return data
}

// GetVersion gets the version of a rune or default (0)
func (r *Rune) GetVersion() int {
	result := r.getIDPart(1)
	if result == -1 {
		return 0
	}

	return result
}

// GetUniqueID gets the uniqueID of a rune or -1
func (r *Rune) GetUniqueID() int {
	return r.getIDPart(0)
}
