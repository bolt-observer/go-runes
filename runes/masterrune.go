package runes

import (
	"encoding/hex"
	"errors"
)

var (
	ErrTooLongSecret = errors.New("too long secret")
)

// MasterRune struct
type MasterRune struct {
	SeedSecret []byte
	Rune
}

// MakeMasterRune creates a new master rune
func MakeMasterRune(seedsecret []byte, uniqueid, version any, restrictions []Restriction) (*MasterRune, error) {
	if len(seedsecret)+1+8 > 64 {
		return nil, ErrTooLongSecret
	}

	if restrictions == nil {
		restrictions = make([]Restriction, 0)
	}

	ret := &MasterRune{}
	ret.SeedSecret = seedsecret

	ret.Sha256 = NewSha256()
	ret.Sha256.Write([]byte(seedsecret))
	ret.Sha256.AddPadding()

	if uniqueid != nil {
		u, err := UniqueId(uniqueid, version)
		if err != nil {
			return nil, err
		}
		ret.AddRestriction(*u)
	}

	for _, r := range restrictions {
		ret.AddRestriction(r)
	}

	return ret, nil
}

// GetRestricted obtains a restricted rune
func (r *MasterRune) GetRestricted(restrictions ...*Restriction) (*Rune, error) {
	rune, err := MakeRune(r.GetAuthCode(), nil, nil, nil)
	if err != nil {
		return nil, err
	}

	for _, r := range restrictions {
		rune.AddRestriction(*r)
	}

	return rune, nil
}

func (r *MasterRune) IsRuneAuthorized(other *Rune) bool {
	hasher := NewSha256()
	hasher.Write([]byte(r.SeedSecret))
	hasher.AddPadding()

	for _, restriction := range other.Restrictions {
		hasher.Write([]byte(restriction.String()))
		hasher.AddPadding()
	}

	sum := hasher.GetSum()

	return hex.EncodeToString(sum[:]) == hex.EncodeToString(other.GetAuthCode())
}

func (r *MasterRune) Evaluate(vals map[string]any) (bool, string) {
	return r.Rune.Evaluate(vals)
}
