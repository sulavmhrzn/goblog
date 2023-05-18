package data

import "database/sql"

type Models struct {
	UserModel  UserModel
	TokenModel TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		UserModel:  UserModel{DB: db},
		TokenModel: TokenModel{DB: db},
	}
}
