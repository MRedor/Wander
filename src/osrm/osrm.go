package osrm

import (
	"data"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"objects"
	"strings"
)

type Geometry struct {
	Coordinates [][2]float64
	Type        string
}

type Route struct {
	Geometry Geometry
	Legs     []struct {
		Summary  string
		Duration float64
		Steps    []string
		Distance float64
	}
	Duration float64
	Distance float64
}

type Response struct {
	Code   string
	Routes []Route
	// не используем
	WayPoints []struct {
		Hint     string
		Name     string
		location [2]float64
	}
}

func (osrm Response) Route() (*data.Route, error) {
	if osrm.Code != "Ok" {
		return nil, errors.New("bad OSRM")
	}
	osrmRoute := osrm.Routes[0]
	route := data.Route{
		Points: []data.Point{},
		Length: osrmRoute.Distance,
		Time:   int(math.Round(osrmRoute.Duration)),
	}
	for _, g := range osrmRoute.Geometry.Coordinates {
		route.Points = append(route.Points, data.Point{g[0], g[1]})
	}
	return &route, nil
}

func GetOSRMByObjects(points []data.Object) Response {
	pointParameters := strings.Join(objects.PositionsToStrings(points), ";")
	return getOSRM(pointParameters)
}

func GetOSRMByPoints(points []data.Point) Response {
	pointParameters := strings.Join(objects.PointsToStrings(points), ";")
	return getOSRM(pointParameters)
}

func getOSRM(parameters string) Response {
	url := fmt.Sprintf(
		"http://travelpath.ru:5000/route/v1/foot/%s?alternatives=false&steps=false&geometries=geojson",
		parameters,
	)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	var res Response
	err = json.Unmarshal(body, &res)
	if err != nil {
		panic(err.Error())
	}

	return res
}
