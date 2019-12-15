package controllers

import (
	"encoding/json"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"routes"
	"strconv"
)

func getRoute(c echo.Context) error {
	req := RouteRequest{}
	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateError(0, "wrong request format"))
	}

	var route *routes.Route

	switch req.Type {
	case string(routes.Direct):
		route, err = routes.ABRoute(req.Points[0], req.Points[1], req.Filters)
	case string(routes.Round):
		route, err = routes.RoundRoute(req.Points[0], req.Radius, req.Filters)
	default:
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
	route, err := routes.RouteById(id)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, CreateError(0, "no route with given id"))
	}
	return c.JSON(http.StatusOK, *route)
}

func removePoint(c echo.Context) error {
	req := RemovePointRequest{}
	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateError(0, "wrong request format"))
	}
	route, err := routes.RemovePoint(1, 2)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateError(0, err.Error()))
	}
	return c.JSON(http.StatusOK, *route)
}
