package handler

import "go.uber.org/fx"

// Module invokes mainserver
var Module = fx.Options(
	fx.Provide(
		newUserHandler,
	),
)
