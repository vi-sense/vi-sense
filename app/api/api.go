package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
)

//@title vi-sense BIM API
//@version 0.1.5
//@description This API provides information about 3D room models with associated sensors and their data.

//@host visense.f4.htw-berlin.de:44344
//@BasePath /
//@schemes https

//SetupRouter initializes all available routes / endpoints and the access to static files
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081", "https://visense.f4.htw-berlin.de"},
		AllowMethods:     []string{"GET", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

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

		sensors.GET(":id/data", func(c *gin.Context) {
			c.String(QuerySensorData(c))
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
