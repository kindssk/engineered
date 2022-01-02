package data

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type UserRepo struct {
	db *sql.DB
}

//type UserRepo struct {
//	id   int32
//	name string
//	age  int32
//}

func (u *UserRepo) InsertUser(name string, age int32) {
	u.db.Query("INSERT INTO users (name, age) VALUES (?,?)", "ahh", 10)
}

func NewDateUser(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}
