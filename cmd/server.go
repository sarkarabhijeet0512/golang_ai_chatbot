package main

import (
	"uber_fx_init_folder_structure/config"
	server "uber_fx_init_folder_structure/internal"
	"uber_fx_init_folder_structure/internal/handler"
	"uber_fx_init_folder_structure/pkg/cache"
	"uber_fx_init_folder_structure/pkg/user"
	"uber_fx_init_folder_structure/utils/initialize"

	"go.uber.org/fx"
)

func serverRun() {
	app := fx.New(
		fx.Provide(
			// postgres server
			initialize.NewDB,
			initialize.NewRedisWorker,
		),
		config.Module,
		initialize.Module,
		server.Module,
		handler.Module,
		user.Module,
		cache.Module,
	)

	// Run app forever
	app.Run()
}
