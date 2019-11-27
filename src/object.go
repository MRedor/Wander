package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
)

type Object struct {
	Id       string `json:"id"`
	Name     string `json:"title"`
	Position Point  `json:"position"`
	Image    string `json:"image"`
	Type     string `json:"type"`

	Address string `json:"adress"`
	Url     string `json:"link"`
	Prices  string `json:"price"`
	//workingTime
	Description string `json:"description"`
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

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (p DBObject) positionString() string {
	return fmt.Sprintf("%f,%f", p.Lat, p.Lon)
}

// DBObjectById gets object from database.
func DBObjectById(id int64) DBObject {
	var point = DBObject{}

	err := db.Get(&point, "select * from points where id=?", id)
	if err != nil {
		log.Fatal(err)
	}

	return point
}

func ObjectByDBObject(o DBObject) Object {
	return Object{
		Id:          fmt.Sprintf("%v", o.Id),
		Name:        o.Name.String,
		Position:    Point{o.Lat, o.Lon},
		Image:       o.Image.String,
		Type:        o.Type,
		Address:     o.Address.String,
		Url:         o.Url.String,
		Prices:      o.Prices.String,
		Description: o.Description.String,
	}
}

func ObjectById(id int64) Object {
	return ObjectByDBObject(DBObjectById(id))
}

func getDBObjectsInRange(a, b Point) []DBObject {
	maxLat, maxLon := a.Lat, a.Lon
	minLat, minLon := b.Lat, b.Lon
	if minLat > maxLat {
		minLat, maxLat = maxLat, minLat
	}
	if minLon > maxLon {
		minLon, maxLon = maxLon, minLon
	}

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

// RandomObjectsInRange gets a slice of random objects from the given range.
func RandomObjectsInRange(a, b Point, count int64) []Object {
	objectsDB := getDBObjectsInRange(a, b)
	var objects []Object
	for _, o := range objectsDB {
		objects = append(objects, ObjectByDBObject(o))
	}

	rand.Shuffle(
		len(objects),
		func(i, j int) {
			objects[i], objects[j] = objects[j], objects[i]
		})
	if int64(len(objects)) < count {
		count = int64(len(objects))
	}
	return objects[:count]
}

func pointsToStrings(points []DBObject) []string {
	var result []string
	for _, p := range points {
		result = append(result, p.positionString())
	}
	return result
}
