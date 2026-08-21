package model

import "fmt"

type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string    { return e.Code + ": " + e.Message }
func Err(code, message string) error    { return &DomainError{Code: code, Message: message} }
func Wrap(code string, err error) error { return &DomainError{Code: code, Message: err.Error()} }
func IsCode(err error, code string) bool {
	if e, ok := err.(*DomainError); ok {
		return e.Code == code
	}
	return false
}
func Invalid(message string) error  { return Err("invalid_request", message) }
func Conflict(message string) error { return Err("conflict", message) }
func NotFound(kind, id string) error {
	return Err("not_found", fmt.Sprintf("%s %s 不存在", kind, id))
}
