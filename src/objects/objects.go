package objects

import (
	"db"
	"filters"
	"fmt"
	"math/rand"
	"points"
)

type Object struct {
	Id       int64        `json:"id"`
	Name     string       `json:"title"`
	Position points.Point `json:"position"`
	Image    string       `json:"image"`
	Type     string       `json:"type"`

	Address     string `json:"address"`
	Url         string `json:"link"`
	Prices      string `json:"price"`
	WorkingTime string `json:"workingTime"`
	Description string `json:"description"`
}

func (o Object) PositionString() string {
	return fmt.Sprintf("%f,%f", o.Position.Lat, o.Position.Lon)
}

func (o Object) Point() points.Point {
	return o.Position
}

func PointsByObjects(objects []Object) []points.Point {
	result := []points.Point{}
	for _, o := range objects {
		result = append(result, o.Point())
	}
	return result
}

func ObjectById(id int64) (*Object, error) {
	o, err := db.DBObjectById(id)
	if err != nil {
		return nil, err
	}
	return ObjectByDBObject(o), nil
}

func GetAllObjectInRange(a, b points.Point, filters filters.StringFilter) []Object {
	objectsDB := db.GetDBObjectsInRange(a, b, filters)
	var objects []Object
	for _, o := range objectsDB {
		objects = append(objects, *(ObjectByDBObject(&o)))
	}
	return objects
}

// RandomObjectsInRange gets a slice of random objects from the given range.
func RandomObjectsInRange(a, b points.Point, count int64, filters filters.StringFilter) []Object {
	objectsDB := db.GetDBObjectsInRange(a, b, filters)
	var objects []Object
	for _, o := range objectsDB {
		objects = append(objects, *(ObjectByDBObject(&o)))
	}

	rand.Shuffle(
		len(objects),
		func(i, j int) {
			objects[i], objects[j] = objects[j], objects[i]
		})
	if int64(len(objects)) < count {
		return objects[:len(objects)]
	}
	return objects[:count]
}

func PositionsToStrings(points []Object) []string {
	var result []string
	for _, p := range points {
		result = append(result, p.PositionString())
	}
	return result
}

func PointsToStrings(points []points.Point) []string {
	var result []string
	for _, p := range points {
		result = append(result, p.String())
	}
	return result
}

func ObjectByDBObject(o *db.DBObject) *Object {
	image := ""
	if o.Image.Valid {
		image = "https://travelpath.ru" + o.Image.String
	}
	return &Object{
		Id:          o.Id,
		Name:        o.Name.String,
		Position:    points.Point{Lat: o.Lat, Lon: o.Lon},
		Image:       image,
		Type:        o.Type,
		Address:     o.Address.String,
		Url:         o.Url.String,
		Prices:      o.Prices.String,
		WorkingTime: o.Time.String,
		Description: o.Description.String,
	}
}
