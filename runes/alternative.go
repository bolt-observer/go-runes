package runes

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// KnownConditions are the currently known conditions
var KnownConditions = []string{"!", "=", "/", "^", "$", "~", "<", ">", "}", "{", "#"}

// Alternative struct
type Alternative struct {
	Field string
	Cond  string
	Value any
}

// ObtainValue is the signature of a function to get current value
type ObtainValue func() any

func containsPunctuation(s string) bool {
	for _, c := range s {
		if isPunct(c) {
			return true
		}
	}

	return false
}

func knownCondition(cond string) bool {
	for _, c := range KnownConditions {
		if cond == c {
			return true
		}
	}

	return false
}

// MakeAlternative returns a new Alternative
func MakeAlternative(field string, cond string, value any, allowIDField bool) (*Alternative, error) {
	if containsPunctuation(field) {
		return nil, fmt.Errorf("field not valid")
	}

	if field == "" {
		if !allowIDField {
			return nil, fmt.Errorf("uniqueId field not valid here")
		}
		if cond != "=" {
			return nil, fmt.Errorf("uniqueId condition must be '='")
		}
	}

	if !knownCondition(cond) {
		return nil, fmt.Errorf("cond not valid %s", cond)
	}

	return &Alternative{Field: field, Cond: cond, Value: value}, nil
}

// MakeAlternativeFromString returns a new alternativee from a string
func MakeAlternativeFromString(str string, allowIDField bool) (*Alternative, string, error) {

	offset := 0

	cond := ""
	for _, r := range str {
		if isPunct(r) {
			cond = string(r)
			break
		}
		offset++
	}

	if cond == "" {
		return nil, "", fmt.Errorf("%s does not contain any operator", str)
	}

	field := str[0:offset]
	offset++

	var sb strings.Builder

	offset2 := 0
	for _, r := range str[offset:] {
		if r == '|' {
			offset2++
			break
		}
		if r == '&' {
			break
		}
		if r == '\\' {
			offset2++
			continue
		}

		sb.WriteRune(r)
		offset2++
	}

	alt, err := MakeAlternative(field, cond, sb.String(), allowIDField)
	if err != nil {
		return nil, "", err
	}

	return alt, str[offset+offset2:], nil
}

func escape(s string) string {
	str := strings.ReplaceAll(s, "&", "\\&")
	str = strings.ReplaceAll(str, "|", "\\|")
	str = strings.ReplaceAll(str, "\\", "\\\\")
	return str
}

// IsUniqueID - is this alternative the unique id
func (a *Alternative) IsUniqueID() bool {
	return a.Field == ""
}

// String returns a string representation
func (a *Alternative) String() string {
	s, ok := a.Value.(string)
	if ok {
		return fmt.Sprintf("%s%s%s", a.Field, a.Cond, escape(s))
	}

	return fmt.Sprintf("%s%s%v", a.Field, a.Cond, a.Value)
}

// Evaluate evaluates the alternative
func (a *Alternative) Evaluate(vals map[string]any) (bool, string) {
	if a.Cond == "#" {
		return true, ""
	}

	if _, ok := vals[a.Field]; !ok {
		if a.IsUniqueID() {
			s, ok := a.Value.(string)
			if !ok {
				return false, "unique id should be string"
			}

			if strings.Contains(s, "-") {
				return false, fmt.Sprintf("unknown version %v", a.Value)
			}
		}
		if a.Cond != "!" {
			return false, fmt.Sprintf("%s is missing", a.Field)
		}
		return true, ""
	}

	actualValue := vals[a.Field]
	obtainer, ok := actualValue.(ObtainValue)
	if ok {
		actualValue = obtainer()
	}

	switch a.Cond {
	case "!":
		return false, fmt.Sprintf("%s is present", a.Field)
	case "=":
		ret, err := isEqual(actualValue, a.Value)
		if ret && err == nil {
			return true, ""
		}
		return false, fmt.Sprintf("!= %s", a.Value)
	case "/":
		ret, err := isEqual(actualValue, a.Value)
		if !ret && err == nil {
			return true, ""
		}
		return false, fmt.Sprintf("= %s", a.Value)
	case "^":
		// starts with
		val := fmt.Sprintf("%v", a.Value)
		entry := fmt.Sprintf("%v", actualValue)

		if strings.HasPrefix(entry, val) {
			return true, ""
		}
		return false, fmt.Sprintf("does not start with %s", val)
	case "$":
		// ends with
		val := fmt.Sprintf("%v", a.Value)
		entry := fmt.Sprintf("%v", actualValue)

		if strings.HasSuffix(entry, val) {
			return true, ""
		}
		return false, fmt.Sprintf("does not end with %s", val)
	case "~":
		// contains
		val := fmt.Sprintf("%v", a.Value)
		entry := fmt.Sprintf("%v", actualValue)

		if strings.Contains(entry, val) {
			return true, ""
		}
		return false, fmt.Sprintf("does not contain %s", val)
	case "<":
		ret, err := isLower(actualValue, a.Value)
		if ret && err == nil {
			return true, ""
		}
		return false, fmt.Sprintf(">= %v", a.Value)
	case ">":
		ret, err := isHigher(actualValue, a.Value)
		if ret && err == nil {
			return true, ""
		}
		return false, fmt.Sprintf("<= %v", a.Value)
	case "{":
		ret := lexoCmp(actualValue, a.Value)
		if ret < 0 {
			return true, ""
		}
		return false, fmt.Sprintf("is the same or ordered after %v", actualValue)
	case "}":
		ret := lexoCmp(vals[a.Field], a.Value)
		if ret > 0 {
			return true, ""
		}
		return false, fmt.Sprintf("is the same or ordered before %v", actualValue)
	default:
		return false, fmt.Sprintf("unhandled case: %v", a.Cond)
	}
}

func isPunct(r rune) bool {
	punc := `!"#$%&'()*+,-./:;<=>?@[\]^_{|}~"`
	punc += "`"

	// Because some chars like "+" are apparently not unicode punctuations
	return unicode.IsPunct(r) || strings.ContainsRune(punc, r)
}

// Wake me up when golang gets better generics, until then we do some ugly hacks with "any" (I'd rather use comparable and constraints.Ordered)

func toInt(a, b string) (int64, int64, error) {
	numA, err := strconv.ParseInt(a, 10, 64)
	if err != nil {
		return 0, 0, err
	}
	numB, err := strconv.ParseInt(b, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return numA, numB, nil
}

func toFloat(a, b string) (float64, float64, error) {
	numA, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return 0, 0, err
	}
	numB, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return 0, 0, err
	}

	return numA, numB, nil
}

func isLower(a, b any) (bool, error) {
	one := fmt.Sprintf("%v", a)
	two := fmt.Sprintf("%v", b)

	intA, intB, err := toInt(one, two)
	if err == nil {
		return intA < intB, nil
	}

	floatA, floatB, err := toFloat(one, two)
	if err == nil {
		return floatA < floatB, nil
	}

	return false, fmt.Errorf("could not compare")
}

func isHigher(a, b any) (bool, error) {
	one := fmt.Sprintf("%v", a)
	two := fmt.Sprintf("%v", b)

	intA, intB, err := toInt(one, two)
	if err == nil {
		return intA > intB, nil
	}

	floatA, floatB, err := toFloat(one, two)
	if err == nil {
		return floatA > floatB, nil
	}

	return false, fmt.Errorf("could not compare")
}

func isEqual(a, b any) (bool, error) {
	one := fmt.Sprintf("%v", a)
	two := fmt.Sprintf("%v", b)

	return one == two, nil
}

func lexoCmp(f, v any) int {
	field := fmt.Sprintf("%v", f)
	val := fmt.Sprintf("%v", v)

	length := len(field)
	if len(val) < length {
		length = len(val)
	}
	cmp := strings.Compare(field, val[0:length])

	/* If val is same but longer, field is < */
	if cmp == 0 && len(val) > len(field) {
		cmp = -1
	}

	return cmp
}
