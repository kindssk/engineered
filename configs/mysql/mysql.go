package mysql

import (
	"database/sql"
	"fmt"
)

type DialOption func(*dialOptions)

type dialOptions struct {
	database string
	host     string
	port     string
	username string
	password string
	encoding string
	timeOut  int32
}

var MysqlClient *sql.DB

func DialEncoding(encoding string) DialOption {
	return func(option *dialOptions) {
		option.encoding = encoding
	}
}

func DialTimeout(timeOut int32) DialOption {
	return func(option *dialOptions) {
		option.timeOut = timeOut
	}
}

func Dial(database, host, port, username, password string, options ...DialOption) ( error) {
	do := &dialOptions{
		database: database,
		host:     host,
		port:     port,
		username: username,
		password: password,
		encoding: "UTF-8",
		timeOut:  1,
	}
	for _, option := range options {
		option(do)
	}

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", do.username, do.password, do.host, do.port, do.database, do.encoding)
	var err error
	MysqlClient, err = sql.Open("mysql", conn)
	if err != nil {
		return  err
	}
	//return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&timeout=%d", do.username, do.password, do.host, do.port, do.database, do.encoding, do.timeOut), nil
	return  nil
}
func MySqlClient() *sql.DB {
	return MysqlClient
}