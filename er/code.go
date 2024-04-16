package er

// Code in `er code` type.
// All the project exceptions should be defined of this type
type Code int

// Generic exceptions
const (
	UncaughtException Code = iota // 0
	UserNotFound
	Unauthorized
)
