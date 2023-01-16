package runes

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrTooLongSecret = errors.New("too long secret")
)

func MakeNewRune(secret []byte, version string) (*Rune, error) {
	if len(secret)+1+8 > 64 {
		return nil, ErrTooLongSecret
	}

	ret := &Rune{
		UniqueId: "",
		Version:  version,
		Sha256:   NewSha256(),
	}

	_, err := ret.Sha256.Write(secret)
	if err != nil {
		return nil, err
	}

	err = ret.Sha256.AddPadding()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *Rune) UniqueIdRestriction(uniqueid, version string) *Restriction {
	if strings.Contains(uniqueid, "-") {
		return nil
	}

	id := uniqueid
	if version != "" {
		id = fmt.Sprintf("%s-%s", uniqueid, version)
	}

	// TODO: fix me
	fmt.Printf("%s", id)
	return nil
}

func (r *Rune) getRestrictions() []string {
	ret := make([]string, 0)
	for _, one := range r.Restrictions {
		ret = append(ret, one.Encode())
	}
	return ret
}

func (r *Rune) ToBase64() string {
	restrictions := ""
	if len(r.Restrictions) > 0 {
		restrictions = "&" + strings.Join(r.getRestrictions(), "|")
	}

	var buf bytes.Buffer

	hash := r.Sha256.GetSum()
	buf.Write(hash[:])
	buf.Write([]byte(restrictions))

	return base64.URLEncoding.EncodeToString(buf.Bytes())
}

/*
func FromBase64(s string) *Rune {
	data, err := base64.URLEncoding.DecodeString(s)

	ret := &Rune{
		UniqueId: "",
		Version:  version,
		Sha256:   NewSha256(),
	}

	data[:64]

	if err != nil {
		return nil
	}

}
*/

func (r *Rune) AddRestriction(restriction Restriction) {
	r.Restrictions = append(r.Restrictions, restriction)

	r.Sha256.Write([]byte(restriction.Encode()))
	r.Sha256.AddPadding()
}
