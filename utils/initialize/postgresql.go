package initialize

import (
	"context"
	"fmt"
	"os"
	"uber_fx_init_folder_structure/pkg/user"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	_ "github.com/lib/pq"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/viper"
)

const (
	envPgDB       = "postgresql_db"
	envPgUser     = "postgresql_user"
	envPgPassword = "postgresql_password"
	envPgHost     = "postgresql_host"
	envPgPort     = "postgresql_port"
)

type DBOut struct {
	fx.Out

	DB *pg.DB `name:"userdb"`
}
type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	qry, e := q.FormattedQuery()
	logrus.Info(string(qry), "\nerror: ", e)
	return nil
}

// NewDeliveryRiderDB creates a connection to DeliveryRider MongoDB
func NewDB(conf *viper.Viper, log *logrus.Logger) (out DBOut, err error) {
	pgDB := conf.GetString(envPgDB)
	pgUser := conf.GetString(envPgUser)
	pgPassword := conf.GetString(envPgPassword)
	pgHost := conf.GetString(envPgHost)
	pgPort := conf.GetString(envPgPort)
	db, err := postgresqlInit(pgDB, pgUser, pgPassword, pgHost, pgPort, log)
	if err != nil {
		log.Error(err)
		return
	}
	out = DBOut{
		DB: db,
	}
	return
}

func postgresqlInit(dbName, dbUser, dbPassword, dbHost, dbPort string, log *logrus.Logger) (
	DB *pg.DB, err error) {

	//the DB variable below is a connection pool.
	DB = pg.Connect(
		&pg.Options{
			Addr:     fmt.Sprintf(dbHost + ":" + dbPort),
			User:     dbUser,
			Password: dbPassword,
			Database: dbName,
			OnConnect: func(ctx context.Context, db *pg.Conn) error {
				_, err := db.Exec("SET timezone = 'Asia/Calcutta'")
				return err
			},
		},
	)

	// if os.Getenv("MODE") == "development" {
	DB.AddQueryHook(dbLogger{})
	// }
	err = DB.Ping(DB.Context())

	if err != nil {
		sentry.CaptureException(err)
		log.WithFields(logrus.Fields{
			"error": err.Error(),
			"uri":   fmt.Sprint(dbName, " ", dbHost),
		}).Fatal("postgresql connection failed")
		return
	}

	if os.Getenv("MODE") == "server" {
		createSchema(DB)
	}
	log.Info("Successfully connected!")
	log.WithFields(logrus.Fields{
		"database": dbName,
	}).Info("postgresql connected")
	return
}

func createSchema(db *pg.DB) error {
	models := []interface{}{

		(*user.User)(nil),
		(*user.UserImages)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
