package data

import "database/sql"

type Models struct {
	UserModel UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		UserModel: UserModel{DB: db},
	}
}
