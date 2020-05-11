package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/vi-sense/vi-sense/app/docs"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.Use(cors.Default())

	r.Static("/files", "/sample-data/models/")

	//uses a middleware to serve static files to the root url.
	//this is needed by vue and can't be done with gin only, because gin complains if there is a wildcard route with conflicting child routes.

	r.Use(static.Serve("/", static.LocalFile("/static/", false)))

	//needed for vue-router, routes every route that wasn't found to index.html
	r.NoRoute(func(c *gin.Context) {
		c.File("/static/index.html")
	})
	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/models", func(c *gin.Context) {
		c.String(http.StatusOK, QueryRoomModels())
	})

	r.GET("/models/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(http.StatusOK, QueryRoomModel(id))
	})

	r.GET("/sensors", func(c *gin.Context) {
		c.String(http.StatusOK, QuerySensors())
	})

	r.GET("/sensors/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(http.StatusOK, QuerySensor(id))
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func main() {
	SetupDatabase(true)
	CreateMockData("/sample-data")

	//check if bind mount is working
	dat, err := ioutil.ReadFile("/sample-data/info.txt")
	if err != nil {
		panic(err)
	}

	fmt.Print(string(dat))

	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	err = r.Run(":8080")
	if err != nil {
		panic(r)
	}
}
