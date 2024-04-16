package user

import (
	"time"

	"go.uber.org/fx"
)

// Module provides all constructor and invocation methods to facilitate credits module
var Module = fx.Options(
	fx.Provide(
		NewDBRepository,
		NewService,
	),
)

type (
	// User represents the user entity
	User struct {
		tableName struct{}  `pg:"users,discard_unknown_columns"`
		ID        int       `json:"id" pg:"id,pk"`
		Username  string    `json:"username" pg:"username,unique"`
		IsActive  bool      `json:"is_active" pg:"is_active"`
		CreatedAt time.Time `json:"created_at" pg:"created_at"`
		UpdatedAt time.Time `json:"updated_at" pg:"updated_at"`
	}
	UserImages struct {
		tableName struct{}  `pg:"user_images,discard_unknown_columns"`
		ID        int       `json:"-" pg:"id,pk"`
		UserID    int       `json:"-" pg:"user_id"`
		Url       string    `json:"url" pg:"url"`
		IsActive  bool      `json:"-" pg:"is_active"`
		CreatedAt time.Time `json:"-" pg:"created_at"`
		UpdatedAt time.Time `json:"-" pg:"updated_at"`
	}
	HistoryLogs struct {
		ID        int       `json:"id" pg:"id"`
		Input     string    `json:"input" pg:"input"`
		Output    string    `json:"output" pg:"output"`
		IsActive  bool      `json:"is_active" pg:"is_active"`
		CreatedAt time.Time `json:"created_at" pg:"created_at"`
		UpdatedAt time.Time `json:"updated_at" pg:"updated_at"`
	}
)

var UserName string
