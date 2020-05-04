package types

import (
	"fmt"
)

// SystemError captures all errors returned from the Rust code as SystemError.
// Exactly one of the fields should be set.
type SystemError struct {
	InvalidRequest     *InvalidRequest     `json:"invalid_request,omitempty"`
	NoSuchContract     *NoSuchContract     `json:"no_such_contract,omitempty"`
	Unknown            *Unknown            `json:"unknown,omitempty"`
	UnsupportedRequest *UnsupportedRequest `json:"unsupported_request,omitempty"`
}

var (
	_ error = SystemError{}
	_ error = InvalidRequest{}
	_ error = NoSuchContract{}
	_ error = Unknown{}
	_ error = UnsupportedRequest{}
)

func (a SystemError) Error() string {
	switch {
	case a.InvalidRequest != nil:
		return a.InvalidRequest.Error()
	case a.NoSuchContract != nil:
		return a.NoSuchContract.Error()
	case a.Unknown != nil:
		return a.Unknown.Error()
	case a.UnsupportedRequest != nil:
		return a.UnsupportedRequest.Error()
	default:
		panic("unknown error variant")
	}
}

type InvalidRequest struct {
	Err string `json:"error,omitempty"`
}

func (e InvalidRequest) Error() string {
	return fmt.Sprintf("invalid request: %s", e.Err)
}

type NoSuchContract struct {
	Addr string `json:"addr,omitempty"`
}

func (e NoSuchContract) Error() string {
	return fmt.Sprintf("no such contract: %s", e.Addr)
}

type Unknown struct{}

func (e Unknown) Error() string {
	return "unknown system error"
}

type UnsupportedRequest struct {
	Kind string `json:"kind,omitempty"`
}

func (e UnsupportedRequest) Error() string {
	return fmt.Sprintf("unsupported request: %s", e.Kind)
}

// ToSystemError will try to convert the given error to an SystemError.
// This is important to returning any Go error back to Rust.
//
// If it is already StdError, return self.
// If it is an error, which could be a sub-field of StdError, embed it.
// If it is anything else, **return nil**
//
// This may return nil on an unknown error, whereas ToStdError will always create
// a valid error type.
func ToSystemError(err error) *SystemError {
	if isNil(err) {
		return nil
	}
	switch t := err.(type) {
	case SystemError:
		return &t
	case *SystemError:
		return t
	case InvalidRequest:
		return &SystemError{InvalidRequest: &t}
	case *InvalidRequest:
		return &SystemError{InvalidRequest: t}
	case NoSuchContract:
		return &SystemError{NoSuchContract: &t}
	case *NoSuchContract:
		return &SystemError{NoSuchContract: t}
	case Unknown:
		return &SystemError{Unknown: &t}
	case *Unknown:
		return &SystemError{Unknown: t}
	case UnsupportedRequest:
		return &SystemError{UnsupportedRequest: &t}
	case *UnsupportedRequest:
		return &SystemError{UnsupportedRequest: t}
	default:
		return nil
	}
}