package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

//@title vi-sense BIM API
//@version 0.1.3
//@description This API provides information about 3D room models with associated sensors and their data.

//@host visense.f4.htw-berlin.de:8080
//@BasePath /
//@schemes http

//SetupRouter initializes all available routes / endpoints and the access to static files
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	//uses a middleware to serve static files to the root url.
	//this is needed by vue and can't be done with gin only, because gin complains if there is
	//a wildcard route with conflicting child routes.
	r.Use(static.Serve("/", static.LocalFile("/static/", false)))

	//needed for vue-router, routes every route that wasn't found to index.html
	r.NoRoute(func(c *gin.Context) {
		c.File("/static/index.html")
	})

	r.Static("/files", "/sample-data/models/")

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	models := r.Group("/models")
	{
		models.GET("", func(c *gin.Context) {
			c.String(QueryRoomModels())
		})
		models.GET(":id", func(c *gin.Context) {
			c.String(QueryRoomModel(c))
		})
	}

	sensors := r.Group("/sensors")
	{
		sensors.GET("", func(c *gin.Context) {
			c.String(QuerySensors())
		})

		sensors.GET(":id", func(c *gin.Context) {
			c.String(QuerySensor(c))
		})

		sensors.GET(":id/anomalies", func(c *gin.Context) {
			c.String(QueryAnomalies(c))
		})

		sensors.PATCH(":id", func(c *gin.Context) {
			c.String(PatchSensor(c))
		})
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func AsJSON(obj interface{}) string {
	b, err := json.Marshal(&obj)
	if err != nil {
		fmt.Println("[!]", err)
		return ""
	}

	return string(b)
}
