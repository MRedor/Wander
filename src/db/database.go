package db

import (
	"data"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"math"
)

var (
	db *sqlx.DB
)

func InitDB() {
	cfg := readConfig()
	source := fmt.Sprintf("%s:%s@/%s",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Table)
	database, err := sqlx.Open("mysql", source)
	if err != nil {
		log.Fatal(err)
	}
	db = database
}

// DBObjectById gets object from database.
func DBObjectById(id int64) (*DBObject, error) {
	var point = DBObject{}

	err := db.Get(&point, "select * from points where (active=true and id=?)", id)
	if err != nil {
		return nil, err
	}

	return &point, nil
}

func GetDBObjectsInRange(a, b data.Point) []DBObject {
	maxLat := math.Max(a.Lat, b.Lat)
	minLat := math.Min(a.Lat, b.Lat)
	maxLon := math.Max(a.Lon, b.Lon)
	minLon := math.Min(a.Lon, b.Lon)

	result := []DBObject{}
	query := fmt.Sprintf(
		"select * from points where (%f <= lat and lat <= %f and %f <= lon and lon <= %f)",
		minLat, maxLat, minLon, maxLon)

	err := db.Select(&result, query)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return result
}

func FreeIdInRoutes() (int, error) {
	var id int

	err := db.Get(&id, "select count(*) from routes")
	if err != nil {
		return 0, err
	}

	return id + 1, nil
}

func InsertRoute(route data.Route) {
	objectsJSON, _ := json.Marshal(route.Objects)
	pointsJSON, _ := json.Marshal(route.Points)

	db.MustExec(
		"insert into routes (type, start_lat, start_lon, finish_lat, finish_lon, radius, length, time, objects, points, name, count) VALUES ((?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?))",
		route.Type,
		route.Points[0].Lat,
		route.Points[0].Lon,
		route.Points[len(route.Points)-1].Lat,
		route.Points[len(route.Points)-1].Lon,
		route.Radius,
		route.Length,
		route.Time,
		objectsJSON,
		pointsJSON,
		route.Name,
		1, // count
	)
}

func GetRoundDBRoute(start data.Point, radius int) (*DBRoute, error) {
	var route = DBRoute{}
	query := fmt.Sprintf(
		"select count(*) from routes where (start_lat=%f and start_lon=%f and radius=%f)",
		start.Lat, start.Lon, radius)

	err := db.Get(&route, query)
	if err != nil {
		return nil, err
	}

	return &route, nil
}

func ExistRoundRoute(start data.Point, radius int) bool {
	var cnt int
	query := fmt.Sprintf(
		"select count(*) from routes where (start_lat=%f and start_lon=%f and radius=%f)",
		start.Lat, start.Lon, radius)

	err := db.Get(&cnt, query)
	if err != nil {
		return false
	}
	if cnt == 0 {
		return false
	}

	return true
}

func GetDirectDBRoute(a, b data.Point) (*DBRoute, error) {
	var route = DBRoute{}
	query := fmt.Sprintf(
		"select * from routes where (start_lat=%f and start_lon=%f and finish_lat=%f and finish_lon=%f)",
		a.Lat, a.Lon, b.Lat, b.Lon)

	err := db.Get(&route, query)
	if err != nil {
		return nil, err
	}
	return &route, nil
}

func ExistDirectRoute(a, b data.Point) bool {
	var cnt int
	query := fmt.Sprintf(
		"select count(*) from routes where (start_lat=%f and start_lon=%f and finish_lat=%f and finish_lon=%f)",
		a.Lat, a.Lon, b.Lat, b.Lon)

	err := db.Get(&cnt, query)
	if err != nil {
		return false
	}
	if cnt == 0 {
		return false
	}

	return true
}

func DBRouteById(id int64) (*DBRoute, error) {
	var route = DBRoute{}

	err := db.Get(&route, "select * from routes where id=?", id)
	if err != nil {
		return nil, err
	}

	return &route, nil
}

type DBRoute struct {
	Id         int     `db:"id"`
	Type       string  `db:"type"`
	Start_lat  float64 `db:"start_lat"`
	Start_lon  float64 `db:"start_lon"`
	Finish_lat float64 `db:"finish_lat"`
	Finish_lon float64 `db:"finish_lon"`
	Radius     int     `db:"radius"`
	Length     float64 `db:"length"`
	Time       int     `db:"time"`
	Objects    string  `db:"objects"`
	Points     string  `db:"points"`
	Name       string  `db:"name"`
	Count      int     `db:"count"`
}

func (r *DBRoute) Route() *data.Route {
	route := data.Route{
		Objects: []data.Object{},
		Points:  []data.Point{},
		Id:      r.Id,
		Length:  r.Length,
		Time:    r.Time,
		Name:    r.Name,
		Type:    r.Type,
		Radius:  r.Radius,
	}

	json.Unmarshal([]byte(r.Objects), route.Objects)
	json.Unmarshal([]byte(r.Points), route.Points)

	return &route
}
