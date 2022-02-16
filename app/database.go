package app

import (
	"database/sql"
	"golangrestfulapi/helper"
	"time"
)

// bikin function
// DB itu struct, makanya kita set sebagai pointer
func NewDB() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/golang_restful_api")
	helper.PanicIfError(err)

	// set connection pooling
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
