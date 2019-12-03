package main

import (
	"fmt"
	"sort"
)

type Route struct {
	Objects []Object `json:"objects"`
	Points  []Point  `json:"points"`
	Id      int      `json:"id"`
	Length  float64  `json:"length"` //meters
	Time    int      `json:"time"`   //seconds
	Name    string   `json:"name"`
}

func NameByRoute(route *Route) string {
	// хотим генерить что-то умное и триггерное
	// а пока по номеру
	return fmt.Sprintf("Route %v", route.Id)
}

func ABRoute(a, b Point) (*Route, error) {
	objects := RandomObjectsInRange(a, b, 10)
	sort.Slice(objects, func(i, j int) bool {
		return a.distance(objects[i].Position) < a.distance(objects[j].Position)
	})
	result, err := getOSRM(objects).Route()
	if err != nil {
		//ищем маршрут между парами
	}
	result.Objects = objects
	// пушить в бд и брать нормальный id
	result.Id = 0
	result.Name = NameByRoute(result)
	return result, nil
}

func CircularRoute(start Point, radius int) (*Route, error) {
	a := Point{
		Lat: start.Lat - float64(radius),
		Lon: start.Lon - float64(radius),
	}
	b := Point{
		Lat: start.Lat + float64(radius),
		Lon: start.Lon + float64(radius),
	}
	// пока в маршрут выбираем случайные объекты
	objects := RandomObjectsInRange(a, b, 10)
	// сортируем по полярному углу относительно старта
	sort.Slice(objects, func(i, j int) bool {
		x1 := objects[i].Position.Lat - start.Lat
		y1 := objects[i].Position.Lon - start.Lon
		x2 := objects[j].Position.Lat - start.Lat
		y2 := objects[j].Position.Lon - start.Lon
		return (x1*y2 - x2*y1) < 0
	})
	return routeByObjects(objects)
}

func routeByObjects(objects []Object) (*Route, error) {
	result, err := getOSRM(objects).Route()
	if err != nil {
		//ищем маршрут между парами
	}
	result.Objects = objects
	result.Id = 0
	result.Name = NameByRoute(result)
	return result, nil
}
