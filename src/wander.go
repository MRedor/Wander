package main

import (
	"config"
	"controllers"
	"db"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	controllers.StartServer()
}

func init() {
	config.InitConfig()
	db.InitDB()
}
