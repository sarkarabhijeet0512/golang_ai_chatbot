package user

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type Repository interface {
	upsertUserRegistration(context.Context, *User) error
	fetchUserByUsername(context.Context, string) (*User, error)
	userUploadPhoto(context.Context, *UserImages) error
	retrievePhotos(context.Context, int) ([]UserImages, error)
}

// NewRepositoryIn is function param struct of func `NewRepository`
type NewRepositoryIn struct {
	fx.In

	Log *logrus.Logger
	DB  *pg.DB `name:"userdb"`
}

// PGRepo is postgres implementation
type PGRepo struct {
	log *logrus.Logger
	db  *pg.DB
}

// NewDBRepository returns a new persistence layer object which can be used for
// CRUD on db
func NewDBRepository(i NewRepositoryIn) (Repo Repository, err error) {

	Repo = &PGRepo{
		log: i.Log,
		db:  i.DB,
	}

	return
}

// IsActive checks if DB is connected
func (r *PGRepo) IsActive() (ok bool, err error) {

	ctx := context.Background()
	err = r.db.Ping(ctx)
	if err == nil {
		ok = true
	}
	return
}
func (r *PGRepo) GetDBConnection(dCtx context.Context) *pg.DB {
	return r.db
}

func (r *PGRepo) upsertUserRegistration(dCtx context.Context, req *User) error {
	_, err := r.db.ModelContext(dCtx, req).OnConflict("(username) DO UPDATE").Insert()
	return err
}

func (r *PGRepo) fetchUserByUsername(dCtx context.Context, username string) (res *User, err error) {
	res = &User{}
	err = r.db.ModelContext(dCtx, res).Where("username = ?", username).Select()
	return res, err
}

func (r *PGRepo) userUploadPhoto(ctx context.Context, userImages *UserImages) error {
	_, err := r.db.ModelContext(ctx, userImages).Insert()
	return err
}
func (r *PGRepo) retrievePhotos(ctx context.Context, userID int) ([]UserImages, error) {
	userImages := []UserImages{}
	err := r.db.ModelContext(ctx, &userImages).Where("user_id = ?", userID).Select()
	return userImages, err
}
