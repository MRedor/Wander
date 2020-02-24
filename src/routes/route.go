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
	Type    string           `json:"type"`
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
	if route != nil {
		db.UpdateRouteCounter(route.Id)
		return route, err
	}

	allObjects := objects.GetAllObjectInRange(a, b, filters)
	pathObjects := routes.ABRoute{Start: a, Finish: b}.Build(allObjects)
	routeMainPoints := []points.Point{a}
	routeMainPoints = append(routeMainPoints, objects.PointsByObjects(pathObjects)...)
	routeMainPoints = append(routeMainPoints, b)

	route, err = routeByPoints(routeMainPoints)
	if err != nil {
		return nil, err
	}
	route.Objects = pathObjects
	route.Type = string(Direct)
	route.Id, err = saveInDB(route, filters.Int())
	if err != nil {
		return nil, err
	}
	route.Name = nameByRoute(route)
	db.UpdateRouteName(route.Id, route.Name)
	return route, err
}

func RoundRoute(start points.Point, radius int, filters filters.StringFilter) (*Route, error) {
	route, err := getRoundRoute(start, radius, filters)
	if route != nil {
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
	route.Id, err = saveInDB(route, filters.Int())
	if err != nil {
		return nil, err
	}
	route.Name = nameByRoute(route)
	db.UpdateRouteName(route.Id, route.Name)
	return route, err
}

func saveInDB(route *Route, filters int) (int64, error) {
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
	var objectsIds []int64
	for i := range route.Objects {
		objectsIds = append(objectsIds, route.Objects[i].Id)
	}

	objectsJSON, _ := json.Marshal(objectsIds)
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

func getTimeByDistance(dist float64) int {
	speed := float64(25) / 18 // 5 km/h ~ 1.4 m/s
	return int(math.Round(dist / speed))
}

func routeByPoints(points_ []points.Point) (*Route, error) {
	result, err := routeByOSRMResponce(osrm.GetOSRMByPoints(points_))
	if err != nil {
		route := Route{
			Points: []points.Point{},
			Length: 0,
			Time:   0,
		}
		for i := 1; i < len(points_); i++ {
			result, err := routeByOSRMResponce(osrm.GetOSRMByPoints([]points.Point{points_[i-1], points_[i]}))
			if err != nil {
				route.Points = append(route.Points, points_[i-1], points_[i])
				dist := routes.GetMetersDistanceByPoints(points_[i-1], points_[i])
				route.Length += dist
				route.Time += getTimeByDistance(dist)
			} else {
				route.Length += result.Length
				route.Time += result.Time
				route.Points = append(route.Points, result.Points...)
			}
		}
		return &route, nil
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
		route.Points = append(route.Points, points.Point{g[1], g[0]})
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
	route.Id, err = saveInDB(route, routeFilters)
	if err != nil {
		return nil, err
	}
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
	route.Id, err = saveInDB(route, routeFilters)
	if err != nil {
		return nil, err
	}
	route.Name = nameByRoute(route)
	db.UpdateRouteName(route.Id, route.Name)

	return route, err
}
