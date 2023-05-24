package errors

import "fmt"

type Error struct {
	Code    string `json:"code"`
	Message string `json:"error"`
	Err     error  `json:"-"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
}

func New(code, message string) *Error {
	return &Error{Code: code, Message: message}
}

func NewNotFound() *Error {
	return &Error{Code: NotFound, Message: notFoundMessage}
}

func Status(err error) int {
	if err == nil {
		return 500
	}
	e, ok := err.(*Error)
	if !ok {
		return 500
	}
	s, ok := codeToStatus[e.Code]
	if !ok {
		panic(fmt.Sprintf("No status defined for error code: %s", e.Code))
	}
	return s
}

func Ensure(err error) *Error {
	if err == nil {
		return &Error{Code: Unexpected, Message: "an unexpected error occurred"}
	}
	e, ok := err.(*Error)
	if !ok {
		return &Error{Code: Unexpected, Message: "an unexpected error occurred", Err: err}
	}
	if e.Code == Internal && e.Message == "" {
		e.Message = internalMessage
	}
	if e.Code == NotFound && e.Message == "" {
		e.Message = notFoundMessage
	}
	return e
}

func Code(err error) string {
	return Ensure(err).Code
}

func Details(err error) string {
	e := Ensure(err)
	if e.Err != nil {
		return fmt.Sprintf("%s: %s. Original: %s", e.Code, e.Message, e.Err.Error())
	}
	return e.Error()
}
