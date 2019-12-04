package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
)

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

func getRoundRoute(start Point, radius int) *Route {
	return &Route{}
}

func existRoundRoute(start Point, radius int) bool {
	return false
}

func getDirectRoute(a, b Point) *Route {
	return &Route{}
}

func existDirectRoute(a, b Point) bool {
	return false
}

func routeById(id int64) (*Route, error) {
	return nil, nil
}