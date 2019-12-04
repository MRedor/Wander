package main

import (
	"fmt"
	"math"
	"math/rand"
)

type Object struct {
	Id       int    `json:"id"`
	Name     string `json:"title"`
	Position Point  `json:"position"`
	Image    string `json:"image"`
	Type     string `json:"type"`

	Address string `json:"address"`
	Url     string `json:"link"`
	Prices  string `json:"price"`
	//workingTime
	Description string `json:"description"`
}

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (p Point) String() string {
	return fmt.Sprintf("%f,%f", p.Lat, p.Lon)
}

func (p Point) distance(a Point) float64 {
	return math.Sqrt(math.Pow(p.Lon-a.Lon, 2) + math.Pow(p.Lat-a.Lat, 2))
}

func (p Object) positionString() string {
	return fmt.Sprintf("%f,%f", p.Position.Lat, p.Position.Lon)
}

func ObjectById(id int64) (*Object, error) {
	o, err := DBObjectById(id)
	if err != nil {
		return nil, err
	}
	return o.Object(), nil
}

// RandomObjectsInRange gets a slice of random objects from the given range.
func RandomObjectsInRange(a, b Point, count int64) []Object {
	objectsDB := getDBObjectsInRange(a, b)
	var objects []Object
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

func positionsToStrings(points []Object) []string {
	var result []string
	for _, p := range points {
		result = append(result, p.positionString())
	}
	return result
}

func pointsToStrings(points []Point) []string {
	var result []string
	for _, p := range points {
		result = append(result, p.String())
	}
	return result
}