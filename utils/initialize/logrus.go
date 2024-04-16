package initialize

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitLogrus(conf *viper.Viper) *logrus.Logger {
	// logglyToken := conf.GetString("LOGGLY_TOKEN")
	log := logrus.New()
	log.WithField("Source", "go")
	// if !conf.GetBool("sentry_enable") {
	// 	return log
	// }
	// hook := logrusly.NewLogglyHook(logglyToken, "http://logs-01.loggly.com/inputs", logrus.InfoLevel, "http")
	// log.Hooks.Add(hook)
	formatter := &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg: "message",
		},
	}
	log.SetFormatter(formatter)

	log.SetReportCaller(true)
	// log.AddHook(sentryhook.New([]logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}))
	// hook.Flush()
	return log
}
