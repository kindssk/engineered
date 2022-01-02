package data

import (
	"database/sql"
	"engineered/configs/mysql"
	"fmt"
)

func NewMysqlDB() *sql.DB {
	fmt.Println(mysql.MySqlClient())
	return mysql.MySqlClient()
}
