package main

import (
	"fmt"
	"github.com/labstack/echo"
	"log"
	"os"
	"src/gopkg.in/yaml.v2"

	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
)

var (
	db  *sqlx.DB
	cfg Config
)

type Config struct {
	Database struct {
		Username string `yaml:"user"`
		Password string `yaml:"password"`
		Table    string `yaml:"table"`
	} `yaml:"database"`
}

func startServer() {
	e := echo.New()

	e.GET("/api/objects/:id", getObjectById)
	e.GET("/api/objects/getFeatured/:boundingBox", getRandomObjects)
	e.POST("/api/routes/get", getRoute)

	//TODO:
	e.POST("/api/routes/removePoint", removePoint)
	e.GET("/api/routes/:id", getRouteById)
	e.GET("/api/list/:id", getListById)
	e.GET("/api/lists", getLists)
	e.POST("/api/feedback", saveFeedback)

	e.Logger.Fatal(e.Start(":1323"))
}

func main() {
	startServer()

}

func readConfig() {
	f, err := os.Open("config.yml")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}
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
