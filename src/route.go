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
	type_   string
	radius  int
}

func NameByRoute(route *Route) string {
	// хотим генерить что-то умное и триггерное
	// а пока по номеру
	return fmt.Sprintf("Route %v", route.Id)
}

func ABRoute(a, b Point) (*Route, error) {
	if existDirectRoute(a, b) {
		return getDirectRoute(a, b)
	}
	objects := RandomObjectsInRange(a, b, 10)
	sort.Slice(objects, func(i, j int) bool {
		return a.distance(objects[i].Position) < a.distance(objects[j].Position)
	})
	route, err := routeByObjects(objects)
	route.type_ = "direct"
	route.radius = 0
	insertRoute(*route)
	return route, err
}

func RoundRoute(start Point, radius int) (*Route, error) {
	if existRoundRoute(start, radius) {
		return getRoundRoute(start, radius)
	}
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
	route, err := routeByObjects(objects)
	route.type_ = "round"
	route.radius = radius
	insertRoute(*route)
	return route, err
}

func routeByObjects(objects []Object) (*Route, error) {
	result, err := getOSRMByObjects(objects).Route()
	if err != nil {
		//ищем маршрут между парами
	}
	result.Objects = objects
	result.Id, err = freeIdInRoutes()
	if err != nil {
		return nil, err
	}
	result.Name = NameByRoute(result)
	return result, nil
}

func routeById(id int64) (*Route, error) {
	dbroute, err := DBRouteById(id)
	if err != nil {
		return nil, err
	}
	return dbroute.Route(), nil
}

func getDirectRoute(a, b Point) (*Route, error) {
	dbroute, err := getDirectDBRoute(a, b)
	if err != nil {
		return nil, err
	}
	return dbroute.Route(), nil
}

func getRoundRoute(start Point, radius int) (*Route, error) {
	dbroute, err := getRoundDBRoute(start, radius)
	if err != nil {
		return nil, err
	}
	return dbroute.Route(), nil
}
