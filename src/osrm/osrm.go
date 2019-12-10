package osrm

import (
	"points"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func GetOSRMByObjects(points []objects.Object) Response {
	pointParameters := strings.Join(objects.PositionsToStrings(points), ";")
	return getOSRM(pointParameters)
}

func GetOSRMByPoints(points []points.Point) Response {
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
