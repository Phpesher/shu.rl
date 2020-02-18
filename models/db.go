package models

import (
	"database/sql"
	"fmt"
)

type Database struct {
	MYSQL_USER     string
	MYSQL_PASSWORD string
	MYSQL_HOST     string
	DATABASE       *sql.DB
}

func NewDatabase(user, password, host string) *Database {
	db, err := sql.Open("mysql", user+":"+password+"@tcp("+host+":3306)/urls")
	if err != nil {
		fmt.Println(err.Error())
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}

	Database := new(Database)
	Database.DATABASE = db
	return Database
}