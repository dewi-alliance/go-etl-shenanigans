package main

import (
	"dewietl/database"
	"dewietl/handler"
	"dewietl/scheduler"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	// Start database
	database.Start()

	// Start scheduler
	scheduler.Start()

	e := echo.New()
	e.Pre(middleware.AddTrailingSlash())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from DeWi!")
	})

	e.GET("/location/hex/:hash/", handler.GetLocationAddress)
	e.GET("/location/hotspot/:hash/", handler.GetLocationHotspot)

	// Get parameter to know if running on dev or production
	serverPort := ":1323"
	if *scheduler.DEV {
		serverPort = ":8081"
	}

	e.Logger.Fatal(e.Start(serverPort))
}
