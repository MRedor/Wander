package main

import (
	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
)

var (
	db  *sqlx.DB
	cfg Config
)

func main() {
	startServer()
}

func init() {
	readConfig()
	db = initDB()
}
