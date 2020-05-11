package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

//@title vi-sense BIM API
//@version 0.1.2
//@description This API provides information about 3D room models with associated sensors and their data.

//@host visense.f4.htw-berlin.de:8080
//@BasePath /
//@schemes http

//SetupRouter initializes all available routes / endpoints and the access to static files
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	r.Static("/files", "/sample-data/models/")

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/models", func(c *gin.Context) {
		c.String(QueryRoomModels())
	})

	r.GET("/models/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(QueryRoomModel(id))
	})

	r.GET("/sensors", func(c *gin.Context) {
		c.String(QuerySensors())
	})

	r.GET("/sensors/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(QuerySensor(id))
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func asJSON(obj interface{}) string {
	b, err := json.Marshal(&obj)
	if err != nil {
		fmt.Println("[!]", err)
		return ""
	}

	return string(b)
}
