package errors

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

type Error struct {
	Code      string
	Message   string
	Status    int
	Err       error
	Stack     []string
	Context   map[string]interface{}
	Timestamp string
}

func New(code, message string) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Status:    500,
		Stack:     captureStack(),
		Context:   make(map[string]interface{}),
		Timestamp: getTimestamp(),
	}
}

func NewWithStatus(code, message string, status int) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Status:    status,
		Stack:     captureStack(),
		Context:   make(map[string]interface{}),
		Timestamp: getTimestamp(),
	}
}

func Wrap(err error, code, message string) *Error {
	if err == nil {
		return nil
	}
	
	return &Error{
		Code:      code,
		Message:   message,
		Status:    500,
		Err:       err,
		Stack:     captureStack(),
		Context:   make(map[string]interface{}),
		Timestamp: getTimestamp(),
	}
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) WithContext(key string, value interface{}) *Error {
	e.Context[key] = value
	return e
}

func (e *Error) WithStatus(status int) *Error {
	e.Status = status
	return e
}

func (e *Error) StackTrace() []string {
	return e.Stack
}

func (e *Error) GetContext() map[string]interface{} {
	return e.Context
}

func captureStack() []string {
	stack := make([]string, 0)
	
	for i := 2; i < 10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		
		funcName := fn.Name()
		if strings.Contains(funcName, "runtime.") {
			continue
		}
		
		stack = append(stack, fmt.Sprintf("%s:%d %s", file, line, funcName))
	}
	
	return stack
}

func getTimestamp() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

type ErrorHandler interface {
	Handle(err error) error
}

type ErrorReporter interface {
	Report(err *Error) error
}

type ErrorRegistry struct {
	handlers  []ErrorHandler
	reporters []ErrorReporter
}

func NewErrorRegistry() *ErrorRegistry {
	return &ErrorRegistry{
		handlers:  make([]ErrorHandler, 0),
		reporters: make([]ErrorReporter, 0),
	}
}

func (er *ErrorRegistry) RegisterHandler(handler ErrorHandler) {
	er.handlers = append(er.handlers, handler)
}

func (er *ErrorRegistry) RegisterReporter(reporter ErrorReporter) {
	er.reporters = append(er.reporters, reporter)
}

func (er *ErrorRegistry) Handle(err error) error {
	for _, handler := range er.handlers {
		if handled := handler.Handle(err); handled != nil {
			return handled
		}
	}
	return err
}

func (er *ErrorRegistry) Report(err *Error) {
	for _, reporter := range er.reporters {
		reporter.Report(err)
	}
}

var (
	ErrNotFound      = NewWithStatus("not_found", "Resource not found", 404)
	ErrBadRequest    = NewWithStatus("bad_request", "Bad request", 400)
	ErrUnauthorized  = NewWithStatus("unauthorized", "Unauthorized", 401)
	ErrForbidden     = NewWithStatus("forbidden", "Forbidden", 403)
	ErrInternalError = NewWithStatus("internal_error", "Internal server error", 500)
)

