package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.Use(CORS())

	r.Static("/files", "/sample-data/models/")

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// version check
	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, "0.1.1")
	})

	r.GET("/models", func(c *gin.Context) {
		c.String(http.StatusOK, queryRoomModels())
	})

	r.GET("/models/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(http.StatusOK, queryRoomModel(id))
	})

	r.GET("/sensors", func(c *gin.Context) {
		c.String(http.StatusOK, querySensors())
	})

	r.GET("/sensors/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(http.StatusOK, querySensor(id))
	})

	return r
}

// https://github.com/gin-contrib/cors/issues/29#issuecomment-397859488
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	setupDatabase(true)
	createMockData("/sample-data")

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
