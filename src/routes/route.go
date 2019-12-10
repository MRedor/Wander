package routes

import (
	"points"
	"db"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"objects"
	"osrm"
	"sort"
)

type RouteType string

const (
	Direct RouteType = "direct"
	Round  RouteType = "round"
)

type Route struct {
	Objects []objects.Object `json:"objects"`
	Points  []points.Point   `json:"points"`
	Id      int              `json:"id"`
	Length  float64          `json:"length"` //meters
	Time    int              `json:"time"`   //seconds
	Name    string           `json:"name"`
	Type    string
	Radius  int
}

func RouteByDBRoute(r *db.DBRoute) *Route {
	route := Route{
		Objects: []objects.Object{},
		Points:  []points.Point{},
		Id:      r.Id,
		Length:  r.Length,
		Time:    r.Time,
		Name:    r.Name,
		Type:    r.Type,
		Radius:  r.Radius,
	}

	json.Unmarshal([]byte(r.Objects), route.Objects)
	json.Unmarshal([]byte(r.Points), route.Points)

	return &route
}

func NameByRoute(route *Route) string {
	// хотим генерить что-то умное и триггерное
	// а пока по номеру
	return fmt.Sprintf("Route %v", route.Id)
}

func ABRoute(a, b points.Point) (*Route, error) {
	route, err := getDirectRoute(a, b)
	if route != nil || err != nil {
		return route, err
	}
	randomObjects := objects.RandomObjectsInRange(a, b, 10)
	sort.Slice(randomObjects, func(i, j int) bool {
		return a.Distance(randomObjects[i].Position) < a.Distance(randomObjects[j].Position)
	})
	route, err = RouteByObjects(randomObjects)
	route.Type = string(Direct)
	route.Id = InsertRoute(route)
	return route, err
}

func RoundRoute(start points.Point, radius int) (*Route, error) {
	route, err := getRoundRoute(start, radius)
	if route != nil || err != nil {
		return route, err
	}
	a := points.Point{
		Lat: start.Lat - float64(radius),
		Lon: start.Lon - float64(radius),
	}
	b := points.Point{
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
	route, err = RouteByObjects(objects)
	route.Type = string(Round)
	route.Radius = radius
	route.Id = InsertRoute(route)
	return route, err
}

func InsertRoute(route *Route) int {
	return 0
}

func RouteByObjects(objects []objects.Object) (*Route, error) {
	result, err := RouteByOSRMResponce(osrm.GetOSRMByObjects(objects))
	if err != nil {
		//ищем маршрут между парами
	}
	result.Objects = objects
	//result.Id, err = db.FreeIdInRoutes()
	if err != nil {
		return nil, err
	}
	result.Name = NameByRoute(result)
	return result, nil
}

func RouteById(id int64) (*Route, error) {
	dbroute, err := db.DBRouteById(id)
	if err != nil {
		return nil, err
	}
	return RouteByDBRoute(dbroute), nil
}

func getDirectRoute(a, b points.Point) (*Route, error) {
	dbroute, err := db.GetDirectDBRoute(a, b)
	if err != nil {
		return nil, err
	}
	return RouteByDBRoute(dbroute), nil
}

func getRoundRoute(start points.Point, radius int) (*Route, error) {
	dbroute, err := db.GetRoundDBRoute(start, radius)
	if err != nil {
		return nil, err
	}
	return RouteByDBRoute(dbroute), nil
}

func RouteByOSRMResponce(resp osrm.Response) (*Route, error) {
	if resp.Code != "Ok" {
		return nil, errors.New("bad OSRM")
	}
	osrmRoute := resp.Routes[0]
	route := Route{
		Points: []points.Point{},
		Length: osrmRoute.Distance,
		Time:   int(math.Round(osrmRoute.Duration)),
	}
	for _, g := range osrmRoute.Geometry.Coordinates {
		route.Points = append(route.Points, points.Point{g[0], g[1]})
	}
	return &route, nil
}
