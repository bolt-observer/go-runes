package runes

import (
	"fmt"
	"strings"
)

// Restriction struct
type Restriction struct {
	Alternatives []Alternative
}

// MakeRestriction from alternatives
func MakeRestriction(alternatives []Alternative) (*Restriction, error) {
	if alternatives == nil || len(alternatives) < 1 {
		return nil, fmt.Errorf("restriction must have some alternative")
	}

	return &Restriction{
		Alternatives: alternatives,
	}, nil
}

// String returns a string representation
func (r *Restriction) String() string {
	str := make([]string, 0)
	for _, one := range r.Alternatives {
		str = append(str, one.String())
	}

	return strings.Join(str, "|")
}

// Evaluate evaluates the restriction
func (r *Restriction) Evaluate(vals map[string]any) (bool, string) {
	reasons := make([]string, 0)
	for _, one := range r.Alternatives {
		b, s := one.Evaluate(vals)
		if !b {
			reasons = append(reasons, s)
		} else {
			return true, ""
		}
	}

	return false, strings.Join(reasons, " AND ")
}

// MakeRestrictionFromString returns a new restriction from a string
func MakeRestrictionFromString(str string, allowIDField bool) (*Restriction, string, error) {

	alternatives := make([]Alternative, 0)

	s := strings.TrimSpace(str)
	allowID := allowIDField
	afterRestriction := ""

	for {
		if strings.HasPrefix(s, "&") {
			afterRestriction = s[1:]
			break
		}
		alt, rest, err := MakeAlternativeFromString(s, allowID)
		if err != nil {
			return nil, "", err
		}

		alternatives = append(alternatives, *alt)

		s = rest
		if len(s) < 1 {
			break
		}
		allowID = false
	}

	if len(alternatives) > 1 && alternatives[0].IsUniqueID() {
		return nil, "", fmt.Errorf("unique_id field cannot have alternatives")
	}

	ret, err := MakeRestriction(alternatives)
	return ret, afterRestriction, err
}

// UniqueID is a helper method to create an unique id restriction
func UniqueID(uniqueID any, version any) (*Restriction, error) {
	if uniqueID == nil {
		return nil, fmt.Errorf("nil unique_id")
	}

	id := fmt.Sprintf("%v", uniqueID)
	if strings.Contains(id, "-") {
		return nil, fmt.Errorf("hyphen not allowed in unique_id %s", id)
	}

	ver := ""
	if version != nil {
		ver = fmt.Sprintf("%v", version)
		id += fmt.Sprintf("-%s", ver)
	}

	alt, err := MakeAlternative("", "=", id, true)
	if err != nil {
		return nil, err
	}

	return &Restriction{
		Alternatives: []Alternative{
			*alt,
		}}, nil
}

// MakeRestrictionsFromString creates restrictionn from string representation
func MakeRestrictionsFromString(str string) ([]Restriction, error) {
	var err error
	rest := str
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

	return restrictions, nil
}

// MustMakeRestrictionsFromString creates restrictionn from string representation
func MustMakeRestrictionsFromString(str string) []Restriction {
	ret, err := MakeRestrictionsFromString(str)
	if err != nil {
		panic(err)
	}
	return ret
}
