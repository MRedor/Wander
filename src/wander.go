package main

import (
	"controllers"
	"db"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	controllers.StartServer()
}

func init() {
	db.InitDB()
}
