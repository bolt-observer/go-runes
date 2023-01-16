package runes

import (
	"fmt"
	"strings"
)

type Rune struct {
	UniqueId string
	Version  string
	Sha256   *Sha256

	Restrictions []Restriction
}

type Restriction struct {
	FieldName string
	Value     string
	Condition RuneCondition
}

func escape(s string) string {
	str := strings.ReplaceAll(s, "&", "\\&")
	str = strings.ReplaceAll(str, "|", "\\|")
	str = strings.ReplaceAll(str, "\\", "\\\\")
	return str
}

func (restr *Restriction) Encode() string {
	return fmt.Sprintf("%s%s%s", restr.FieldName, string(restr.Condition), escape(restr.Value))
}

type RuneCondition string

const (
	RUNE_COND_IF_MISSING  RuneCondition = "!"
	RUNE_COND_EQUAL       RuneCondition = "="
	RUNE_COND_NOT_EQUAL   RuneCondition = "/"
	RUNE_COND_BEGINS      RuneCondition = "^"
	RUNE_COND_ENDS        RuneCondition = "$"
	RUNE_COND_CONTAINS    RuneCondition = "~"
	RUNE_COND_INT_LESS    RuneCondition = "<"
	RUNE_COND_INT_GREATER RuneCondition = ">"
	RUNE_COND_LEXO_BEFORE RuneCondition = "{"
	RUNE_COND_LEXO_AFTER  RuneCondition = "}"
	RUNE_COND_COMMENT     RuneCondition = "#"
)
