package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"

	sqlx "github.com/jmoiron/sqlx"
)

func initDB() *sqlx.DB {
	source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)
	db, err := sqlx.Open("mysql", source)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

type DBObject struct {
	Id                 int            `db:"id"`
	Lat                float64        `db:"lat"`
	Lon                float64        `db:"lon"`
	Type               string         `db:"type"`
	Name               sql.NullString `db:"name"`
	Description        sql.NullString `db:"description"`
	Time               sql.NullString `db:"time"`
	NameEnglish        sql.NullString `db:"name_en"`
	DescriptionEnglish sql.NullString `db:"description_en"`
	Image              sql.NullString `db:"img"`
	IdInGraph          int            `db:"id_point_in_graph"`
	Updated            string         `db:"updated"`
	Active             bool           `db:"active"`
	Address            sql.NullString `db:"address"`
	Prices             sql.NullString `db:"prices"`
	Url                sql.NullString `db:"url"`
	NightDescription   sql.NullString `db:"night_description"`
	NightImage         sql.NullString `db:"night_photo"`
	NightType          sql.NullString `db:"night_type"`
	ActiveOnlyAtNight  sql.NullInt64  `db:"active_only_at_night"`
}

func (o *DBObject) Object() *Object {
	return &Object{
		Id:          o.Id,
		Name:        o.Name.String,
		Position:    Point{o.Lat, o.Lon},
		Image:       "https://travelpath.ru" + o.Image.String,
		Type:        o.Type,
		Address:     o.Address.String,
		Url:         o.Url.String,
		Prices:      o.Prices.String,
		Description: o.Description.String,
	}
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

func getDBObjectsInRange(a, b Point) []DBObject {
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

func freeIdInRoutes() (int, error) {
	var id int

	err := db.Get(&id, "select count(*) from routes")
	if err != nil {
		return 0, err
	}

	return id + 1, nil
}

func insertRoute(route Route) {
	objectsJSON, _ := json.Marshal(route.Objects)
	pointsJSON, _ := json.Marshal(route.Points)

	db.MustExec(
		"insert into routes (id, type, start_lat, start_lon, finish_lat, finish_lon, radius, length, time, objects, points, name, count) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		route.Id,
		route.type_,
		route.Points[0].Lat, route.Points[0].Lon,
		route.Points[len(route.Points)].Lat, route.Points[len(route.Points)].Lon,
		route.radius,
		route.Length, route.Time,
		objectsJSON,
		pointsJSON,
		route.Name,
		1, // count
	)
}

func getRoundDBRoute(start Point, radius int) (*DBRoute, error) {
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

func existRoundRoute(start Point, radius int) bool {
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

func getDirectDBRoute(a, b Point) (*DBRoute, error) {
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

func existDirectRoute(a, b Point) bool {
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

func (r *DBRoute) Route() *Route {
	route := Route{
		Objects: []Object{},
		Points:  []Point{},
		Id:      r.Id,
		Length:  r.Length,
		Time:    r.Time,
		Name:    r.Name,
		type_:   r.Type,
		radius:  r.Radius,
	}

	json.Unmarshal([]byte(r.Objects), route.Objects)
	json.Unmarshal([]byte(r.Points), route.Points)

	return &route
}
