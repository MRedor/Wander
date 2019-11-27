package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func getOSRM(points []DBObject) OSRMResponse {
	pointParameters := strings.Join(pointsToStrings(points), ";")

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
