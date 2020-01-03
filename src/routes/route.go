package routes

import (
	"db"
	"encoding/json"
	"errors"
	"filters"
	"fmt"
	"math"
	"objects"
	"osrm"
	"points"
	"routes/types"
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
	Id      int64            `json:"id"`
	Length  float64          `json:"length"` //meters
	Time    int              `json:"time"`   //seconds
	Name    string           `json:"name"`
	Type    string
	radius  int
	filters int
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
		radius:  r.Radius,
		filters: r.Filters,
	}

	json.Unmarshal([]byte(r.Objects), route.Objects)
	json.Unmarshal([]byte(r.Points), route.Points)

	return &route
}

func nameByRoute(route *Route) string {
	// хотим генерить что-то умное и триггерное
	// а пока по номеру
	return fmt.Sprintf("Route %v", route.Id)
}

func ABRoute(a, b points.Point, filters filters.StringFilter) (*Route, error) {
	route, err := getDirectRoute(a, b, filters)
	if route != nil || err != nil {
		db.UpdateRouteCounter(route.Id)
		return route, err
	}
	routeObjects := objects.RandomObjectsInRange(a, b, 10, filters)

	var routeMainPoints []points.Point
	if false {
		// todo: implement
		//routeObjects = objects.getAllObjectInRange(a, b, filters)
		pathObjects := routes.ABRoute{Start: a, Finish: b}.Build(routeObjects)
		routeMainPoints = []points.Point{a}
		routeMainPoints = append(routeMainPoints, objects.PointsByObjects(pathObjects)...)
		routeMainPoints = append(routeMainPoints, b)
	} else {
		// old solution
		sort.Slice(routeObjects, func(i, j int) bool {
			return a.Distance(routeObjects[i].Position) < a.Distance(routeObjects[j].Position)
		})
		routeMainPoints = append([]points.Point{a}, objects.PointsByObjects(routeObjects)...)
		routeMainPoints = append(routeMainPoints, b)
	}
	route, err = routeByPoints(routeMainPoints)
	if err != nil {
		return nil, err
	}
	route.Objects = routeObjects
	route.Type = string(Direct)
	route.Id = saveInDB(route, filters.Int())
	route.Name = nameByRoute(route)
	db.UpdateRouteName(route.Id, route.Name)
	return route, err
}

func RoundRoute(start points.Point, radius int, filters filters.StringFilter) (*Route, error) {
	route, err := getRoundRoute(start, radius, filters)
	if route != nil || err != nil {
		db.UpdateRouteCounter(route.Id)
		return route, err
	}
	a := points.Point{
		Lat: start.Lat - routes.MetersToLat(float64(radius)),
		Lon: start.Lon - routes.MetersToLon(start, float64(radius)),
	}
	b := points.Point{
		Lat: start.Lat + routes.MetersToLat(float64(radius)),
		Lon: start.Lon + routes.MetersToLon(start, float64(radius)),
	}
	allObjects := objects.RandomObjectsInRange(a, b, 100, filters)
	pathObjects := routes.RoundRoute{Center: start, Radius: radius}.Build(allObjects)
	routeMainPoints := append([]points.Point{start}, objects.PointsByObjects(pathObjects)...)
	routeMainPoints = append(routeMainPoints, start)

	route, err = routeByPoints(routeMainPoints)
	if err != nil {
		return nil, err
	}
	route.Objects = pathObjects
	route.Type = string(Round)
	route.radius = radius
	route.Id = saveInDB(route, filters.Int())
	route.Name = nameByRoute(route)
	db.UpdateRouteName(route.Id, route.Name)
	return route, err
}

func saveInDB(route *Route, filters int) int64 {
	dbroute := db.DBRoute{
		Start_lat:  route.Points[0].Lat,
		Start_lon:  route.Points[0].Lon,
		Finish_lat: route.Points[len(route.Points)-1].Lat,
		Finish_lon: route.Points[len(route.Points)-1].Lon,
		Length:     route.Length,
		Time:       route.Time,
		Name:       route.Name,
		Filters:    filters,
	}
	objectsJSON, _ := json.Marshal(route.Objects)
	dbroute.Objects = string(objectsJSON)
	pointsJSON, _ := json.Marshal(route.Points)
	dbroute.Points = string(pointsJSON)

	if route.Type == string(Direct) {
		dbroute.Type = string(Direct)
		return db.InsertDirectRoute(dbroute)
	} else {
		dbroute.Type = string(Round)
		dbroute.Radius = route.radius
		return db.InsertRoundRoute(dbroute)
	}
}

func routeByObjects(objects []objects.Object) (*Route, error) {
	result, err := routeByOSRMResponce(osrm.GetOSRMByObjects(objects))
	if err != nil {
		//ищем маршрут между парами
	}
	result.Objects = objects
	result.Name = nameByRoute(result)
	return result, nil
}

func routeByPoints(points []points.Point) (*Route, error) {
	result, err := routeByOSRMResponce(osrm.GetOSRMByPoints(points))
	if err != nil {
		//ищем маршрут между парами
	}
	return result, nil
}

func RouteById(id int64) (*Route, error) {
	dbroute, err := db.DBRouteById(id)
	if err != nil {
		return nil, err
	}
	return RouteByDBRoute(dbroute), nil
}

func getDirectRoute(a, b points.Point, filters filters.StringFilter) (*Route, error) {
	dbroute, err := db.GetDirectDBRoute(a, b, filters.Int())
	if err != nil {
		return nil, err
	}
	return RouteByDBRoute(dbroute), nil
}

func getRoundRoute(start points.Point, radius int, filters filters.StringFilter) (*Route, error) {
	dbroute, err := db.GetRoundDBRoute(start, radius, filters.Int())
	if err != nil {
		return nil, err
	}
	return RouteByDBRoute(dbroute), nil
}

func routeByOSRMResponce(resp osrm.Response) (*Route, error) {
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

func RemovePoint(routeId, objectId int64) (*Route, error) {
	route, err := RouteById(routeId)
	if err != nil {
		return nil, err
	}

	if route.Type == string(Round) {
		return removePointFromRoundRoute(route, objectId)
	} else {
		return removePointFromDirectRoute(route, objectId)
	}
}

func searchInObjectSlice(objectId int64, slice []objects.Object) int {
	idInSlice := -1
	for i := 0; i < len(slice); i++ {
		if slice[i].Id == objectId {
			idInSlice = i
			break
		}
	}
	return idInSlice
}

func removePointFromRoundRoute(route *Route, objectId int64) (*Route, error) {
	idInSlice := searchInObjectSlice(objectId, route.Objects)
	if idInSlice == -1 {
		return nil, errors.New("no object with given id in the route")
	}
	routeObjects := append(route.Objects[:idInSlice], route.Objects[idInSlice+1:]...) //удаляем
	routeMainPoints := append([]points.Point{route.Points[0]}, objects.PointsByObjects(routeObjects)...)
	routeMainPoints = append(routeMainPoints, route.Points[0])

	radius := route.radius
	routeFilters := route.filters

	route, err := routeByPoints(routeMainPoints)
	if err != nil {
		return nil, err
	}

	route.Objects = routeObjects
	route.Type = string(Round)
	route.radius = radius
	route.Id = saveInDB(route, routeFilters)
	route.Name = nameByRoute(route)
	db.UpdateRouteName(route.Id, route.Name)

	return route, err
}

func removePointFromDirectRoute(route *Route, objectId int64) (*Route, error) {
	idInSlice := searchInObjectSlice(objectId, route.Objects)
	if idInSlice == -1 {
		return nil, errors.New("no object with given id in the route")
	}
	routeObjects := append(route.Objects[:idInSlice], route.Objects[idInSlice+1:]...) //удаляем
	routeMainPoints := append([]points.Point{route.Points[0]}, objects.PointsByObjects(routeObjects)...)
	routeMainPoints = append(routeMainPoints, route.Points[len(route.Points)-1])

	routeFilters := route.filters

	route, err := routeByPoints(routeMainPoints)
	if err != nil {
		return nil, err
	}

	route.Objects = routeObjects
	route.Type = string(Direct)
	route.Id = saveInDB(route, routeFilters)
	route.Name = nameByRoute(route)
	db.UpdateRouteName(route.Id, route.Name)

	return route, err
}
