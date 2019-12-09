package routes

import (
	"data"
	"db"
	"fmt"
	"objects"
	"osrm"
	"sort"
)

type RouteType string

const (
	Direct RouteType = "direct"
	Round  RouteType = "round"
)

func NameByRoute(route *data.Route) string {
	// хотим генерить что-то умное и триггерное
	// а пока по номеру
	return fmt.Sprintf("Route %v", route.Id)
}

func ABRoute(a, b data.Point) (*data.Route, error) {
	if db.ExistDirectRoute(a, b) {
		return GetDirectRoute(a, b)
	}
	randomObjects := objects.RandomObjectsInRange(a, b, 10)
	sort.Slice(randomObjects, func(i, j int) bool {
		return a.Distance(randomObjects[i].Position) < a.Distance(randomObjects[j].Position)
	})
	route, err := RouteByObjects(randomObjects)
	route.Type = string(Direct)
	db.InsertRoute(*route)
	return route, err
}

func RoundRoute(start data.Point, radius int) (*data.Route, error) {
	if db.ExistRoundRoute(start, radius) {
		return GetRoundRoute(start, radius)
	}
	a := data.Point{
		Lat: start.Lat - float64(radius),
		Lon: start.Lon - float64(radius),
	}
	b := data.Point{
		Lat: start.Lat + float64(radius),
		Lon: start.Lon + float64(radius),
	}
	// пока в маршрут выбираем случайные объекты
	objects := objects.RandomObjectsInRange(a, b, 10)
	// сортируем по полярному углу относительно старта
	sort.Slice(objects, func(i, j int) bool {
		x1 := objects[i].Position.Lat - start.Lat
		y1 := objects[i].Position.Lon - start.Lon
		x2 := objects[j].Position.Lat - start.Lat
		y2 := objects[j].Position.Lon - start.Lon
		return (x1*y2 - x2*y1) < 0
	})
	route, err := RouteByObjects(objects)
	route.Type = string(Round)
	route.Radius = radius
	db.InsertRoute(*route)
	return route, err
}

func RouteByObjects(objects []data.Object) (*data.Route, error) {
	result, err := osrm.GetOSRMByObjects(objects).Route()
	if err != nil {
		//ищем маршрут между парами
	}
	result.Objects = objects
	result.Id, err = db.FreeIdInRoutes()
	if err != nil {
		return nil, err
	}
	result.Name = NameByRoute(result)
	return result, nil
}

func RouteById(id int64) (*data.Route, error) {
	dbroute, err := db.DBRouteById(id)
	if err != nil {
		return nil, err
	}
	return dbroute.Route(), nil
}

func GetDirectRoute(a, b data.Point) (*data.Route, error) {
	dbroute, err := db.GetDirectDBRoute(a, b)
	if err != nil {
		return nil, err
	}
	return dbroute.Route(), nil
}

func GetRoundRoute(start data.Point, radius int) (*data.Route, error) {
	dbroute, err := db.GetRoundDBRoute(start, radius)
	if err != nil {
		return nil, err
	}
	return dbroute.Route(), nil
}
