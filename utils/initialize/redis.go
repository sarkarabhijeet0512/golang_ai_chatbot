package initialize

import (
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type RedisWorkerOut struct {
	fx.Out
	Client *redis.Pool `name:"redisWorker"`
}

// NewRedisWorker returns a new orbis-redis connection
func NewRedisWorker(conf *viper.Viper, log *logrus.Logger) (o RedisWorkerOut, err error) {
	redisConn := conf.GetString("REDIS_WORKER")

	pool := &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisConn)
		},
	}
	o = RedisWorkerOut{
		Client: pool,
	}
	log.WithFields(logrus.Fields{
		"url": redisConn,
	}).Info("redis connected successfully")
	return
}
