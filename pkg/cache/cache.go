package cache

import (
	"time"
	"uber_fx_init_folder_structure/pkg/cache/persistence"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		persistence.NewRedisRepository,
		NewService,
	),
)

type Service struct {
	conf *viper.Viper
	log  *logrus.Logger
	Repo persistence.CacheStore
}

func NewService(conf *viper.Viper, log *logrus.Logger, Repo persistence.CacheStore) *Service {
	return &Service{
		conf: conf,
		log:  log,
		Repo: Repo,
	}
}

func (s *Service) Set(key, value string, expiry time.Duration) error {
	return s.Repo.Set(key, value, expiry)
}

func (s *Service) Get(key string, ptrValue interface{}) error {
	return s.Repo.Get(key, ptrValue)
}

func (s *Service) Delete(key string) error {
	return s.Repo.Delete(key)
}

const (
	CACHE_MIDDLEWARE_KEY = "gincontrib.cache"
)

var (
	PageCachePrefix = "gincontrib.page.cache"
)
