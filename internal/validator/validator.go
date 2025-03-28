package validator

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"
)

type Validator interface {
	Valid(context.Context) Evaluator
}

type Evaluator map[string]string

func (e *Evaluator) AddFieldError(key, message string) {
	if *e == nil {
		*e = make(map[string]string)
	}

	if _, exists := (*e)[key]; !exists {
		(*e)[key] = message
	}
}

func (e *Evaluator) CheckField(ok bool, key, message string) {
	if !ok {
		e.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func InEnum(value string, options []any) bool {
	for _, option := range options {
		if value == fmt.Sprintf("%v", option) {
			return true
		}
	}
	return false
}
