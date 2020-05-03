package main

import (
	"encoding/json"
	"fmt"
)

// @title vi-sense BIM API
// @version 0.1.1
// @description This API provides information about 3D room models with associated sensors and their data.

// @host visense.f4.htw-berlin.de:8080
// @BasePath /
// @schemes http

// QueryRoomModels godoc
// @Summary Query models
// @Description Query all avaiable room models.
// @Tags models
// @Produce  json
// @Success 200 {array} main.RoomModel
// @Failure 400 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /models [get]
func QueryRoomModels() string {
	var q []RoomModel
	db.Find(&q)
	return asJson(&q)
}

// QueryRoomModel godoc
// @Summary Query room model
// @Description Query a single room model by id with containing sensors.
// @Tags models
// @Produce json
// @Param id path int true "RoomModelID"
// @Success 200 {object} main.RoomModel
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /models/{id} [get]
func QueryRoomModel(id string) string {
	var q RoomModel
	db.Preload("Sensors").First(&q, id)
	return asJson(&q)
}

// QuerySensors godoc
// @Summary Query sensors
// @Description Query all avaiable sensors.
// @Tags sensors
// @Produce json
// @Success 200 {array} main.Sensor
// @Failure 400 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /sensors [get]
func QuerySensors() string {
	var q []Sensor
	db.Find(&q)
	return asJson(&q)
}

// QuerySensor godoc
// @Summary Query sensor
// @Description Query a single sensor by id with containing sensor data.
// @Tags sensors
// @Produce json
// @Param id path int true "SensorId"
// @Success 200 {object} main.Sensor
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /sensors/{id} [get]
func QuerySensor(id string) string {
	var q Sensor
	db.Preload("Data").First(&q, id)
	return asJson(&q)
}

func asJson(obj interface{}) string {
	b, err := json.Marshal(&obj)
	if err != nil {
		fmt.Println("[!]", err)
		return ""
	}

	return string(b)
}
