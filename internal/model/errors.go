package model

import (
	"errors"
	"fmt"
)

type DomainError struct {
	Code    string
	Message string
	cause   error
}

func (e *DomainError) Error() string { return e.Code + ": " + e.Message }
func (e *DomainError) Unwrap() error { return e.cause }

func Err(code, message string) error    { return &DomainError{Code: code, Message: message} }
func Wrap(code string, err error) error {
	if err == nil {
		return nil
	}
	return &DomainError{Code: code, Message: err.Error(), cause: err}
}
func IsCode(err error, code string) bool {
	var e *DomainError
	if errors.As(err, &e) {
		return e.Code == code
	}
	return false
}
func Invalid(message string) error  { return Err("invalid_request", message) }
func Conflict(message string) error { return Err("conflict", message) }
func NotFound(kind, id string) error {
	return Err("not_found", fmt.Sprintf("%s %s 不存在", kind, id))
}
