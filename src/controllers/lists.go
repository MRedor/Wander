package controllers

import (
	"github.com/labstack/echo"
	"net/http"
)

func getListById(c echo.Context) error {
	// todo:
	return c.JSON(http.StatusOK, c.Param("id"))
}

func getLists(c echo.Context) error {
	// todo:
	return c.JSON(http.StatusOK, "test response")
}
