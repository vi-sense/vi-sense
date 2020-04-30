package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.Use(cors.Default()) // when no longer in production change this line

	r.Static("/files", "/sample-data/models/")

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// version check
	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, "0.1")
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
