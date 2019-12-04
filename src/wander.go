package main

import (
	"fmt"
	"github.com/labstack/echo"
	"log"

	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
)

var (
	db  *sqlx.DB
	cfg Config
)

func startServer() {
	e := echo.New()

	e.GET("/api/objects/:id", getObjectById)
	e.GET("/api/objects/getFeatured/:boundingBox", getRandomObjects)
	e.POST("/api/routes/get", getRoute)
	e.GET("/api/routes/:id", getRouteById)

	//TODO:
	e.POST("/api/routes/removePoint", removePoint)
	e.GET("/api/list/:id", getListById)
	e.GET("/api/lists", getLists)
	e.POST("/api/feedback", saveFeedback)

	e.Logger.Fatal(e.Start(":1323"))
}

func main() {
	startServer()
}

func initDB() *sqlx.DB {
	source := fmt.Sprintf("%s:%s@/%s",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Table)
	db, err := sqlx.Open("mysql", source)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func init() {
	readConfig()
	db = initDB()
}
