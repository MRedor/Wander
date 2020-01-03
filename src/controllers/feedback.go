package controllers

import (
	"encoding/json"
	"feedback"
	"github.com/labstack/echo"
	"net/http"
)

func saveFeedback(c echo.Context) error {
	req := FeedbackRequest{}
	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateError(0, "wrong request format"))
	}
	err = feedback.Send(feedback.CreateMessage(req))
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateError(0, "wrong request format"))
	}
	return c.JSON(http.StatusOK, "test response")
}
