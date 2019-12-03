package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
)

type OSRMGeometry struct {
	Coordinates [][2]float64
	Type        string
}

type OSRMRoute struct {
	Geometry OSRMGeometry
	Legs     []struct {
		Summary  string
		Duration float64
		Steps    []string
		Distance float64
	}
	Duration float64
	Distance float64
}

type OSRMResponse struct {
	Code   string
	Routes []OSRMRoute
	// не используем
	WayPoints []struct {
		Hint     string
		Name     string
		location [2]float64
	}
}

func (osrm OSRMResponse) Route() (*Route, error) {
	if osrm.Code != "Ok" {
		return nil, errors.New("bad OSRM")
	}
	osrmRoute := osrm.Routes[0]
	route := Route{
		Points: []Point{},
		Length: osrmRoute.Distance,
		Time:   int(math.Round(osrmRoute.Duration)),
	}
	for _, g := range osrmRoute.Geometry.Coordinates {
		route.Points = append(route.Points, Point{g[0], g[1]})
	}
	return &route, nil
}

func getOSRM(points []Object) OSRMResponse {
	pointParameters := strings.Join(positionsToStrings(points), ";")

	url := fmt.Sprintf(
		"http://travelpath.ru:5000/route/v1/foot/%s?alternatives=false&steps=false&geometries=geojson",
		pointParameters,
	)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	var data OSRMResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}

	//fmt.Println(data)
	return data
}
