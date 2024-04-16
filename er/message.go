package er

var messages = map[string]string{
	"1": "Oops! Something went wrong. Please try later",
	"2": "User not found",
	"3": "unauthorized",
}

var codes = map[Code]string{
	UncaughtException: "1",
	UserNotFound:      "2",
	Unauthorized:      "3",
}
