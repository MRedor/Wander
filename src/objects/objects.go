package objects

import (
	"points"
	"db"
	"fmt"
	"math/rand"
)

type Object struct {
	Id       int          `json:"id"`
	Name     string       `json:"title"`
	Position points.Point `json:"position"`
	Image    string       `json:"image"`
	Type     string       `json:"type"`

	Address string `json:"address"`
	Url     string `json:"link"`
	Prices  string `json:"price"`
	//workingTime
	Description string `json:"description"`
}

func (p Object) PositionString() string {
	return fmt.Sprintf("%f,%f", p.Position.Lat, p.Position.Lon)
}

func ObjectById(id int64) (*Object, error) {
	o, err := db.DBObjectById(id)
	if err != nil {
		return nil, err
	}
	return ObjectByDBObject(o), nil
}

// RandomObjectsInRange gets a slice of random objects from the given range.
func RandomObjectsInRange(a, b points.Point, count int64) []Object {
	objectsDB := db.GetDBObjectsInRange(a, b)
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
	return &Object{
		Id:          o.Id,
		Name:        o.Name.String,
		Position:    points.Point{Lat: o.Lat, Lon: o.Lon},
		Image:       "https://travelpath.ru" + o.Image.String,
		Type:        o.Type,
		Address:     o.Address.String,
		Url:         o.Url.String,
		Prices:      o.Prices.String,
		Description: o.Description.String,
	}
}
