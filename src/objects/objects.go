package objects

import (
	"data"
	"db"
	"math/rand"
)

func ObjectById(id int64) (*data.Object, error) {
	o, err := db.DBObjectById(id)
	if err != nil {
		return nil, err
	}
	return o.Object(), nil
}

// RandomObjectsInRange gets a slice of random objects from the given range.
func RandomObjectsInRange(a, b data.Point, count int64) []data.Object {
	objectsDB := db.GetDBObjectsInRange(a, b)
	var objects []data.Object
	for _, o := range objectsDB {
		objects = append(objects, *((&o).Object()))
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

func PositionsToStrings(points []data.Object) []string {
	var result []string
	for _, p := range points {
		result = append(result, p.PositionString())
	}
	return result
}

func PointsToStrings(points []data.Point) []string {
	var result []string
	for _, p := range points {
		result = append(result, p.String())
	}
	return result
}
