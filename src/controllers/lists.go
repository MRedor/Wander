package controllers

import (
	"encoding/json"
	"github.com/labstack/echo"
	"lists"
	"log"
	"net/http"
	"strconv"
)

func getListById(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, CreateError(0, "id should be an integer"))
	}
	list, err := lists.ListById(id)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, CreateError(0, err.Error()))
	}
	return c.JSON(http.StatusOK, *list)
}

func getLists(c echo.Context) error {
	req := GetListsRequest{}
	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateError(0, "wrong request format"))
	}
	lists, err := lists.GetSliceOfLists(req.Count, req.Offset)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateError(0, err.Error()))
	}
	return c.JSON(http.StatusOK, lists)
}
