package controllers

import (
	"github.com/labstack/echo"
	"net/http"
)

func saveFeedback(c echo.Context) error {
	// todo:
	return c.JSON(http.StatusOK, "test response")
}
