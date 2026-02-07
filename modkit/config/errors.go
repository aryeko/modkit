package config

import (
	"fmt"

	"github.com/go-modkit/modkit/modkit/module"
)

// MissingRequiredError reports a required key that was unset.
type MissingRequiredError struct {
	Key       string
	Token     module.Token
	Sensitive bool
}

func (e *MissingRequiredError) Error() string {
	return fmt.Sprintf("missing required config: key=%q token=%q", e.Key, e.Token)
}

// ParseError reports a value parse failure.
type ParseError struct {
	Key       string
	Token     module.Token
	Type      string
	Sensitive bool
	Err       error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("invalid config value: key=%q token=%q type=%q: %v", e.Key, e.Token, e.Type, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// InvalidSpecError reports an invalid value specification.
type InvalidSpecError struct {
	Token  module.Token
	Reason string
}

func (e *InvalidSpecError) Error() string {
	return fmt.Sprintf("invalid config spec: token=%q reason=%s", e.Token, e.Reason)
}
