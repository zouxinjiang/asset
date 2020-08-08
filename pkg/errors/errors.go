package errors

import (
	"fmt"
	"runtime"
)

type (
	CError interface {
		Code() string
		Message() string
		Unwrap() error
	}
	Error struct {
		error
		code       string
		message    string
		fileInfo   string
		line       int
		callerFunc string
	}
)

func (e Error) Code() string {
	return e.code
}

func (e Error) Message() string {
	return e.message
}

func (e Error) Unwrap() error {
	return e.error
}

func (e Error) Error() string {
	str := fmt.Sprintf(`[%s] %s`, e.code, e.message)
	return str
}

func Wrap(err error, code string, message ...string) error {
	e := Error{}

	pc, file, line, _ := runtime.Caller(-1)
	fn := runtime.FuncForPC(pc)
	if fn != nil {
		e.callerFunc = fn.Name()
	}
	e.fileInfo = file
	e.line = line
	e.error = err
	e.code = code
	if len(message) >= 1 {
		e.message = message[0]
	}
	return e
}

func WithCode(code string, message ...string) error {
	e := Error{}

	pc, file, line, _ := runtime.Caller(-1)
	fn := runtime.FuncForPC(pc)
	if fn != nil {
		e.callerFunc = fn.Name()
	}
	e.fileInfo = file
	e.line = line
	e.error = nil
	e.code = code
	if len(message) >= 1 {
		e.message = message[0]
	}
	return e
}
