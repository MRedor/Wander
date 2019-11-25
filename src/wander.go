package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
	//osrm "github.com/maddevsio/osrm"
	//osrm "github.com/gojuno/go.osrm"
)

var (
	db = initDB()
)

func initDB() *sqlx.DB {
	db, err := sqlx.Open("mysql", "wander:wanderpassword@/wander")
	if err != nil {
		// наверно надо делать что-то более аккуратное, чем фатал
		log.Fatal(err)
	}
	//defer db.Close()

	return db
}

type Point struct {
	Id                 int            `db:"id"`
	Lat                float64        `db:"lat"`
	Lon                float64        `db:"lon"`
	ObjectType         string         `db:"type"`
	Name               sql.NullString `db:"name"`
	Description        sql.NullString `db:"description"`
	Time               sql.NullString `db:"time"`
	NameEnglish        sql.NullString `db:"name_en"`
	DescriptionEnglish sql.NullString `db:"description_en"`
	Img                sql.NullString `db:"img"`
	IdInGraph          int            `db:"id_point_in_graph"`
	// updated -- чет не очень понятно что это такое
	Updated           string         `db:"updated"`
	Active            bool           `db:"active"`
	Address           sql.NullString `db:"address"`
	Prices            sql.NullString `db:"prices"`
	Url               sql.NullString `db:"url"`
	NightDescription  sql.NullString `db:"night_description"`
	NightImage        sql.NullString `db:"night_photo"`
	NightType         sql.NullString `db:"night_type"`
	ActiveOnlyAtNight sql.NullInt64  `db:"active_only_at_night"`
}

type Path struct {
	name            string
	distanceMeters  int
	durationMinutes int
	points          []Point
	// что тут еще нужно?
}

func PointById(id int) Point {
	var point = Point{}

	err := db.Get(&point, "select * from points where id=?", id)
	if err != nil {
		log.Fatal(err)
	}

	return point
}

func PointsInRange( /* диапазон в каком-то виде. В каком? */ ) []Point {
	// Выбирать несколько рандомных
	// Сколько?
	return nil
}

func NameByPath(path Path) string {
	// хотим генерить что-то умное и триггерное
	return ""
}

type OSRMGeometry struct {
	Coordinates [][2]float64
	Type        string
}

type Route struct {
	Geometry OSRMGeometry
	Legs     []struct {
		Summary  string
		Duration float64
		Steps    []string
		Distance float64
	}
	Duration float64
	Distance float64
}

type OSRMResponse struct {
	Code   string
	Routes []Route
	// не используем
	WayPoints []struct {
		Hint     string
		Name     string
		location [2]float64
	}
}

func getOSRM() {
	resp, err := http.Get("http://travelpath.ru:5000/route/v1/foot/30.368802,59.937580;30.331551,59.927845;30.304600,59.932929?alternatives=false&steps=false&geometries=geojson")
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	var data OSRMResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(data)

	/*
		var result OSRMResponse
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(result)
		fmt.Println(result)
	*/
}

func ABPath(a, b Point) Path {
	return Path{}
}

func CircularPath(start Point) Path {
	return Path{}
}

func main() {
	e := echo.New()
	e.GET("/points/getRandom", pointsGetRandom)
	e.GET("/path/get", pathGet)
	e.GET("/path/getRound", pathGetRound)
	//...

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
	//fmt.Println("kek")
	//fmt.Println(PointById(10))
	//getOSRM()
}
