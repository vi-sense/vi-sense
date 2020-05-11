package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/vi-sense/vi-sense/app/model"
	"net/http"
)

type UpdateSensor struct {
	MeshID string
}

//QuerySensors godoc
//@Summary Query sensors
//@Description Query all available sensors.
//@Tags sensors
//@Produce json
//@Success 200 {array} model.Sensor
//@Failure 400 {string} string "bad request"
//@Failure 500 {string} string "internal server error"
//@Router /sensors [get]
func QuerySensors() (int, string) {
	var q []Sensor
	DB.Find(&q)
	return http.StatusOK, AsJSON(&q)
}

//QuerySensor godoc
//@Summary Query sensor
//@Description Query a single sensor by id with containing sensor data.
//@Tags sensors
//@Produce json
//@Param id path int true "SensorId"
//@Success 200 {object} model.Sensor
//@Failure 400 {string} string "bad request"
//@Failure 404 {string} string "not found"
//@Failure 500 {string} string "internal server error"
//@Router /sensors/{id} [get]
func QuerySensor(c *gin.Context) (int, string) {
	var q Sensor
	id := c.Param("id")
	DB.Preload("Data").First(&q, id)
	if q.ID == 0 {
		return http.StatusNotFound, AsJSON(gin.H{"error": fmt.Sprintf("Sensor %s not found.", id)})
	}
	return http.StatusOK, AsJSON(&q)
}

//Patch	Sensor godoc
//@Summary Update sensor location
//@Description Updates the mesh id of a single sensor.
//@Tags sensors
//@Accept json
//@Produce json
//@Param id path int true "SensorId"
//@Param sensor body UpdateSensor true "Update sensor"
//@Success 200 {object} model.Sensor
//@Failure 400 {string} string "bad request"
//@Failure 500 {string} string "internal server error"
//@Router /sensors/{id} [patch]
func PatchSensor(c *gin.Context) (int, string) {
	var q Sensor
	id := c.Param("id")
	DB.First(&q, id)

	if q.ID == 0 {
		return http.StatusNotFound, AsJSON(gin.H{"error": fmt.Sprintf("Sensor %s not found.", id)})
	}

	var input UpdateSensor
	if err := c.ShouldBindJSON(&input); err != nil {
		return http.StatusBadRequest, AsJSON(gin.H{"error": err.Error()})
	}

	DB.Model(&q).Update(input)

	return http.StatusOK, AsJSON(&q)
}
