package initialize

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides common application dependencies
var Module = fx.Options(
	fx.Provide(
		InitLogrus,
	),
	fx.Invoke(
		LivenessProbe,
	),
)

// LivenessProbe writes an empty at /tmp/_healthz at every 30 seconds
// K8s checks if the last modified of the file is <= 30 seconds, if not, pod is restarted.
func LivenessProbe(logger *logrus.Logger) {
	go func() {
		for {
			time.Sleep(30 * time.Second)
			err := os.WriteFile("/tmp/_healthz", []byte{}, os.ModePerm)
			if err != nil {
				logger.Info("ERROR: livenessprobe /tmp/_healthz file creation failed",
					zap.Error(err),
				)
			}
		}
	}()
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
