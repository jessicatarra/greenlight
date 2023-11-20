package database

import (
	"database/sql"
	"errors"
)

var (
	ErrEditConflict   = errors.New("edit conflict")
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Movies      MovieModel
	Tokens      TokenModel
	Permissions PermissionModel
	Users       UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies:      MovieModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Users:       UserModel{DB: db},
	}
}
