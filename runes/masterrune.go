package runes

import (
	"encoding/hex"
	"errors"
	"fmt"
)

var (
	// ErrTooLongSecret represents an error where secret was too long
	ErrTooLongSecret = errors.New("too long secret")
	// ErrTooShortSecret represents an error where secret was too short
	ErrTooShortSecret = errors.New("too short secret")
	// ErrUnauthorizedRune represents an error where rune was not authorized
	ErrUnauthorizedRune = errors.New("unauthorized rune")
)

// MasterRune struct
type MasterRune struct {
	SeedSecret []byte
	Rune
}

// MustMakeMasterRune is a helper constructor for creating a master rune
func MustMakeMasterRune(seedsecret []byte) MasterRune {
	rune, err := MakeMasterRune(seedsecret, nil, nil, nil)
	if err != nil {
		panic(err)
	}
	return *rune
}

// MakeMasterRune creates a new master rune
func MakeMasterRune(seedsecret []byte, uniqueid, version any, restrictions []Restriction) (*MasterRune, error) {
	if len(seedsecret)+1+8 > 64 {
		return nil, ErrTooLongSecret
	}
	if len(seedsecret) < 1 {
		return nil, ErrTooShortSecret
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
		u, err := UniqueID(uniqueid, version)
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

// IsRuneAuthorized check whether rune is authorized
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

// Check checks rune
func (r *MasterRune) Check(rune *Rune, vals map[string]any) error {
	if !r.IsRuneAuthorized(rune) {
		return ErrUnauthorizedRune
	}

	ok, msg := rune.Evaluate(vals)
	if ok {
		return nil
	}

	return fmt.Errorf(msg)
}
