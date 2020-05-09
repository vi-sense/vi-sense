package api

import (
	. "github.com/vi-sense/vi-sense/app/model"
	"net/http"
)

// QuerySensors godoc
// @Summary Query sensors
// @Description Query all available sensors.
// @Tags sensors
// @Produce json
// @Success 200 {array} model.Sensor
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} string "internal server error"
// @Router /sensors [get]
func QuerySensors() (int, string) {
	var q []Sensor
	DB.Find(&q)
	return http.StatusOK, asJSON(&q)
}

// QuerySensor godoc
// @Summary Query sensor
// @Description Query a single sensor by id with containing sensor data.
// @Tags sensors
// @Produce json
// @Param id path int true "SensorId"
// @Success 200 {object} model.Sensor
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /sensors/{id} [get]
func QuerySensor(id string) (int, string) {
	var q Sensor
	DB.Preload("Data").First(&q, id)
	if q.ID == 0 {
		return http.StatusNotFound, ""
	}
	return http.StatusOK, asJSON(&q)
}
