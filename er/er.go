// Package er is project's custom error package.
// It lets you add more information to a general error object.
package er

import (
	"errors"
	"fmt"
	"net/http"
)

//go:generate stringer -type Code

// var appName = strings.ReplaceAll(
// 	config.Config.GetString("new_relic_app_name"), " ", "",
// )

var appName = "ecommerce_backend"

// E is custom error type to pass more information along with the native error.
// E implements error interface and should be used in replacement of native error type.
type E struct {
	// Err is the native error object returned from failed function
	// this field will be excluded in JSON serialization
	Err error `json:"-"`

	// Code is er code number that is listed in `code.go` file in this directory.
	// A Code should be unique and used only once throughout the project
	// so that it can lead to the exact line of failed function when searched in the entire project.
	Code Code `json:"code"`

	// Exception is technical expection title. This should be mostly 2-3 words title of the
	// error/er to be understood by developers.
	// format: APP_NAME.DOMAIN.EXCEPTION_TITLE
	// example: [ADDSALE]SILVERBOLT-SERVER-MAIN.ITEM.INVALID_SLUG, [ADDSALE]SILVERBOLT-SERVER-MAIN.USER.AUTHORIZATION
	// note: APP_NAME will be prepended automatically.
	Exception string `json:"exception"`

	// TraceID received from API request
	TraceID string `json:"-"`

	// Status is HTTP status code that is set in API response
	Status int `json:"-"`

	// Message is user-friendly message that can be displayed on front-end for end-users
	Message string `json:"message"`

	// Info is link/URL to the logging system(Sentry/Kibana) to get more insight and trace error
	Info string `json:"info"`

	// ErrorMsg is actual technical error occurred
	ErrorMsg string `json:"-"`

	// NOP (no-operation) if set will not send error to sentry
	NOP bool `json:"-"`
}

// New constructs and returns new E object
func New(err error, code Code) *E {
	if err == nil {
		err = errors.New("uncaught exception")
	}

	return &E{
		Err:       err,
		Code:      code,
		Exception: fmt.Sprintf("%s.%s", appName, code),
		Status:    http.StatusInternalServerError,
		Message:   messageByCode(code),
		ErrorMsg:  err.Error(),
	}
}

// From takes a general error type and returns type-casted `*E` type.
// Checks if its actual value is our custom error `E` object.
// If not, returns a `*E` with exception `UncaughtException`
func From(err error) *E {
	e, ok := err.(*E)
	if !ok {
		e = New(err, UncaughtException)
	}
	return e
}

// Error returns technical error message with error code and er title
func (e *E) Error() string {
	return fmt.Sprintf("(#%d:%s) %s", e.Code, e.Exception, e.Err.Error())
}

func (e *E) String() string {
	return e.Error()
}

// SetStatus sets the HTTP Status in the error object
func (e *E) SetStatus(httpStatus int) *E {
	if httpStatus < 200 || httpStatus > 600 {
		return e
	}
	e.Status = httpStatus
	return e
}

// SetTraceID sets the TraceID in the error object
func (e *E) SetTraceID(traceID string) *E {
	e.TraceID = traceID
	return e
}

// Ignore sets `E.NOP` flag to avoid sending log to sentry
func (e *E) Ignore() *E {
	e.NOP = true
	return e
}

// messageByCode looks up for user-friendly message from messages map in this package
// and returns it with appending error code to it.
func messageByCode(code Code) string {
	n, ok := codes[code]
	if !ok {
		n = "1"
	}

	m, ok := messages[n]
	if !ok {
		m = messages["1"]
	}
	return fmt.Sprintf("(#%d) %s", code, m)
}

// IsCodeEq checks if error and and er.Code are equal
func IsCodeEq(err error, code Code) bool {
	return From(err).Code == code
}
