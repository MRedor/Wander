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
		return c.JSON(http.StatusBadRequest, CreateError(0, "id should be an integer"))
	}
	o, err := ObjectById(id)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, CreateError(0, "no object with given id"))
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
		return c.JSON(http.StatusBadRequest, CreateError(0, err.Error()))
	}

	if c.QueryParam("count") == "" {
		return c.JSON(http.StatusBadRequest, CreateError(0, "count was not provided"))
	}

	count, err := strconv.ParseInt(c.QueryParam("count"), 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, CreateError(0, "count should be an integer"))
	}
	return c.JSON(http.StatusOK, RandomObjectsInRange(*a, *b, count))
}

func getRoute(c echo.Context) error {
	req := RouteRequest{}
	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateError(0, "wrong request format"))
	}

	var route *Route
	if req.Type == string(Direct) {
		route, err = ABRoute(req.Points[0], req.Points[1])
	}
	if req.Type == string(Round) {
		route, err = RoundRoute(req.Points[0], req.Radius)
	}
	if route == nil {
		return c.JSON(http.StatusBadRequest, CreateError(0, "unsupported route type"))
	}


	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateError(0, ""))
	}
	return c.JSON(http.StatusOK, *route)
}

func getRouteById(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, CreateError(0, "id should be an integer"))
	}
	route, err := routeById(id)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, CreateError(0, "no route with given id"))
	}
	return c.JSON(http.StatusOK, *route)
}

func removePoint(c echo.Context) error {
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
