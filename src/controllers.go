package main

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func pointsGetRandom(c echo.Context) error {
	// todo: parse params

	// fetch random points from database

	// return
	return c.JSON(http.StatusOK, "test response")
}

func pathGetRound(c echo.Context) error {
	// todo:
	return c.JSON(http.StatusOK, "test response")
}

func pathGet(c echo.Context) error {
	// parse params (here it's GET params)
	xLat, _ := strconv.ParseFloat(c.QueryParam("start_lat"), 10)
	xLon, _ := strconv.ParseFloat(c.QueryParam("start_lon"), 10)
	yLat, _ := strconv.ParseFloat(c.QueryParam("end_lat"), 10)
	yLon, _ := strconv.ParseFloat(c.QueryParam("end_lon"), 10)

	route := routesBuild(xLat, xLon, yLat, yLon)

	// return result
	return c.JSON(http.StatusOK, route)
}
