package controllers

import "github.com/labstack/echo"

func StartServer() {
	e := echo.New()

	e.GET("/api/objects/:id", getObjectById)
	e.GET("/api/objects/getFeatured/:boundingBox", getRandomObjects)
	e.POST("/api/routes/get", getRoute)
	e.GET("/api/routes/:id", getRouteById)

	//TODO:
	e.POST("/api/routes/removePoint", removePoint)
	e.GET("/api/list/:id", getListById)
	e.GET("/api/lists", getLists)
	e.POST("/api/feedback", saveFeedback)

	e.Logger.Fatal(e.Start(":1323"))
}
