package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/vi-sense/vi-sense/app/model"
	"net/http"
	"time"
)

type UpdateSensor struct {
	MeshID string
}

type Anomaly struct {
	Gradient    float64
	Difference    float64
	Date     time.Time
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

func QueryAnomalies(c *gin.Context) (int, string) {
	s := c.Query("start")
	e := c.Query("end")

	if s == "" {
		s = time.Time{}.Format(Layout)
	}
	if e == "" {
		e = time.Now().Format(Layout)
	}

	fmt.Println("start")
	fmt.Println(s)
	fmt.Println("end")
	fmt.Println(e)

	var q []Data
	id := c.Param("id")
	DB.Where("sensor_id = ? AND date >= ? AND date <= ?", id, s, e).Find(&q)

	if len(q) == 0 {
		return http.StatusNotFound, AsJSON(gin.H{"error": fmt.Sprintf("Data for Sensor %s not found.", id)})
	}

	var a []Anomaly

	for i := 1; i < len(q); i++ {
		timeDiff := q[i].Date.Unix() - q[i-1].Date.Unix()
		valDiff := q[i].Value - q[i-1].Value
		grad := valDiff / float64(timeDiff)

		d:= time.Unix(q[i-1].Date.Unix() + (timeDiff / 2), 0)
		a = append(a, Anomaly{Gradient: grad, Difference: valDiff, Date: d})
	}

	return http.StatusOK, AsJSON(a)
}