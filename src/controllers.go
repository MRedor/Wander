package main

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func getObjectById(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, Error{Code: 0, Message: "id should be an integer"})
	}
	o, err := ObjectById(id)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, Error{Code: 0, Message: "no object with given id"})
	}
	return c.JSON(http.StatusOK, *o)
}

func parseBoundingBox(box string) (*Point, *Point, error) {
	points := strings.Split(box, ";")
	if len(points) != 2 {
		return nil, nil, errors.New("two points should be passed")
	}

	a := strings.Split(points[0], ",")
	if len(a) != 2 {
		return nil, nil, errors.New("point should have two coordinates")
	}
	latA, err := strconv.ParseFloat(a[0], 10)
	if err != nil {
		return nil, nil, errors.New("point should have two float coordinates")
	}
	lonA, err := strconv.ParseFloat(a[1], 10)
	if err != nil {
		return nil, nil, errors.New("point should have two float coordinates")
	}

	b := strings.Split(points[1], ",")
	if len(b) != 2 {
		return nil, nil, errors.New("point should have two coordinates")
	}
	latB, err := strconv.ParseFloat(b[0], 10)
	if err != nil {
		return nil, nil, errors.New("point should have two float coordinates")
	}
	lonB, err := strconv.ParseFloat(b[1], 10)
	if err != nil {
		return nil, nil, errors.New("point should have two float coordinates")
	}

	return &Point{Lat: latA, Lon: lonA}, &Point{Lat: latB, Lon: lonB}, nil
}

func getRandomObjects(c echo.Context) error {
	boundingBox := c.Param("boundingBox")
	a, b, err := parseBoundingBox(boundingBox)
	if err != nil {
		log.Println("getRandomObjects: ", err)
		return c.JSON(http.StatusBadRequest, Error{Code: 0, Message: err.Error()})
	}

	if c.QueryParam("count") == "" {
		return c.JSON(http.StatusBadRequest, Error{Code: 0, Message: "count was not provided"})
	}

	count, err := strconv.ParseInt(c.QueryParam("count"), 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, Error{Code: 0, Message: "count should be an integer"})
	}
	return c.JSON(http.StatusOK, RandomObjectsInRange(*a, *b, count))
}

func getRoute(c echo.Context) error {
	// todo:
	req := RouteRequest{}
	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Code: 0, Message: "wrong request format"})
	}
	var route *Route
	if req.Type == "direct" {
		route, err = ABRoute(req.Points[0], req.Points[1])
	} else {
		route, err = CircularRoute(req.Points[0], req.Radius)
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Code: 0, Message: ""})
	}
	return c.JSON(http.StatusOK, *route)
}

func removePoint(c echo.Context) error {
	// todo:
	return c.JSON(http.StatusOK, "test response")
}

func getRouteById(c echo.Context) error {
	// todo:
	return c.JSON(http.StatusOK, "test response")
}

func getListById(c echo.Context) error {
	// todo:
	return c.JSON(http.StatusOK, c.Param("id"))
}

func getLists(c echo.Context) error {
	// todo:
	return c.JSON(http.StatusOK, "test response")
}

func saveFeedback(c echo.Context) error {
	// todo:
	return c.JSON(http.StatusOK, "test response")
}
