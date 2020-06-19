package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/s12i/gin-throttle"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vi-sense/vi-sense/app/docs"
	"net/http"
	"os"
)

func GetEnv(key string, defVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	fmt.Printf("Could not read %s returning default value %s.\n", key, defVal)
	return defVal
}

//@title vi-sense BIM API
//@version 0.1.7
//@description This API provides information about 3D room models with associated sensors and their data.

//@BasePath /

//SetupRouter initializes all available routes / endpoints and the access to static files
func SetupRouter() *gin.Engine {

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", GetEnv("HOST", "localhost"), GetEnv("PORT", "8080"))
	docs.SwaggerInfo.Schemes = []string{GetEnv("SCHEME", "http")}
	r := gin.Default()

	if GetEnv("PRODUCTION", "false") == "false" {
		r.Use(gzip.Gzip(gzip.BestSpeed))
		fmt.Println("[i] Using gzip.")
	}

	//to limit the number of requests per second
	r.Use(middleware.Throttle(100, 100))

	// cors settings
	/*r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081", "https://visense.f4.htw-berlin.de"},
		AllowMethods:     []string{"GET", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))*/
	r.Use(cors.Default())

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
