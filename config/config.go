package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Module is config module
var Module = fx.Options(
	fx.Provide(
		New,
	),
)

type argvMeta struct {
	desc       string
	defaultVal string
}

// New returns a viper object.
// This object is used to read environment variables or command line arguments.
func New() (config *viper.Viper) {
	config = viper.New()

	confList := map[string]argvMeta{
		"env": {
			defaultVal: "development",
			desc:       "Environment",
		},
		"postgresql_db": {
			defaultVal: "userdb",
			desc:       "postgresql db name",
		},
		"postgresql_host": {
			defaultVal: "localhost",
			desc:       "postgresql host",
		},
		"postgresql_port": {
			defaultVal: "5432",
			desc:       "postgresql port",
		},
		"postgresql_user": {
			defaultVal: "postgres",
			desc:       "postgresql username",
		},
		"postgresql_password": {
			defaultVal: "",
			desc:       "postgresql password",
		},
		"port": {
			defaultVal: "8765",
			desc:       "Port number of user API server",
		},
		"mode": {
			defaultVal: "server",
			desc:       "App mode eg. consumer, server, worker",
		},
		"log_level": {
			defaultVal: "debug",
			desc:       "Log level to be printed. List of log level by Priority - debug, info, warn, error, dpanic, panic, fatal",
		},
		"open_ai_api_key": {
			defaultVal: "",
			desc:       "open ai api key",
		},
		"aws_access_key": {
			defaultVal: "",
			desc:       "aws access key",
		},
		"aws_secret_key": {
			defaultVal: "",
			desc:       "aws access token",
		},
		"aws_bucket": {
			defaultVal: "",
			desc:       "aws bucket name",
		},
		"aws_region": {
			defaultVal: "",
			desc:       "aws bucket region",
		},
		"redis_worker": {
			defaultVal: "localhost:6379",
			desc:       "redis server",
		},
	}

	for key, meta := range confList {
		// automatic conversion of environment var key to `UPPER_CASE` will happen.
		config.BindEnv(key)

		// read command-line arguments
		pflag.String(key, meta.defaultVal, meta.desc)
	}

	pflag.Parse()
	config.BindPFlags(pflag.CommandLine)
	return
}
