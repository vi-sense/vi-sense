package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/vi-sense/vi-sense/app/model"
	"math"
	"net/http"
	"strconv"
	"time"
)

type UpdateSensor struct {
	MeshID string			`json:"mesh_id"`
}

type ParamParseError struct {
	Param string
	Value string
}

func (e *ParamParseError) Error() string {
	return fmt.Sprintf("Error parsing parameter '%s' with value '%s'.", e.Param, e.Value)
}

type Anomaly struct {
	Type       AnomalyType	`json:"type"`
	Date       time.Time	`json:"date"`
	Value      float64 		`json:"value"`
	Gradient   float64		`json:"gradient"`
}

type AnomalyType string

const (
	HighGradient    AnomalyType = "High Gradient"
	AboveUpperLimit AnomalyType = "Above Upper Limit"
	BelowLowerLimit AnomalyType = "Below Lower Limit"
)

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
	var r []Sensor
	DB.Find(&r)
	return http.StatusOK, AsJSON(&r)
}

//QuerySensor godoc
//@Summary Query sensor
//@Description Query a single sensor by id
//@Tags sensors
//@Produce json
//@Param id path int true "Sensor ID"
//@Success 200 {object} model.Sensor
//@Failure 400 {string} string "bad request"
//@Failure 404 {string} string "not found"
//@Failure 500 {string} string "internal server error"
//@Router /sensors/{id} [get]
func QuerySensor(c *gin.Context) (int, string) {
	var r Sensor
	id := c.Param("id")
	DB.First(&r, id)
	if r.ID == 0 {
		return http.StatusNotFound, AsJSON(gin.H{"error": fmt.Sprintf("Sensor %s not found.", id)})
	}
	return http.StatusOK, AsJSON(&r)
}

//QuerySensorData godoc
//@Summary Query sensor data
//@Description Query data for a specific sensor
//@Tags sensors
//@Produce json
//@Param id path int true "Sensor ID"
//@Param start_date query string false "Start Date"
//@Param end_date query string false "End Date"
//@Success 200 {array} model.Data
//@Failure 400 {string} string "bad request"
//@Failure 404 {string} string "not found"
//@Failure 500 {string} string "internal server error"
//@Router /sensors/{id}/data [get]
func QuerySensorData(c *gin.Context) (int, string) {
	id := c.Param("id")

	queryParams := map[string]interface{}{
		"start_date":  "",
		"end_date":    "",
	}

	// check if sensor exists
	var s Sensor
	DB.First(&s, id)
	if s.ID == 0 {
		return http.StatusNotFound, AsJSON(gin.H{"error": fmt.Sprintf("Sensor '%s' not found.", id)})
	}

	err := fillQueryParams(c, &queryParams)
	if err != nil {
		return http.StatusBadRequest, AsJSON(gin.H{"error": err.Error()})
	}

	// query sensor data within defined time period
	r := make([]Data, 0)

	q := DB.Where("sensor_id = ?", id)

	if queryParams["start_date"] != "" {
		q = q.Where("date >= ?", queryParams["start_date"])
	}

	if queryParams["end_date"] != "" {
		q = q.Where("date <= ?", queryParams["end_date"])
	}

	q.Find(&r)

	return http.StatusOK, AsJSON(r)
}

//QuerySensor godoc
//@Summary Query anomalies
//@Description Query anomalies for a specific sensor
//@Tags sensors
//@Produce json
//@Param id path int true "Sensor ID"
//@Param start_date query string false "Start Date"
//@Param end_date query string false "End Date"
//@Param max_grad query number false "Maximum Gradient"
//@Param lower_limit query number false "Lower Value Limit"
//@Param upper_limit query number false "Upper Value Limit"
//@Success 200 {array} Anomaly
//@Failure 400 {string} string "bad request"
//@Failure 404 {string} string "not found"
//@Failure 500 {string} string "internal server error"
//@Router /sensors/{id}/anomalies [get]
func QueryAnomalies(c *gin.Context) (int, string) {
	id := c.Param("id")

	queryParams := map[string]interface{}{
		"start_date":  "",
		"end_date":    "",
		"max_grad":    math.MaxFloat64,
		"lower_limit": -math.MaxFloat64,
		"upper_limit": math.MaxFloat64,
	}

	// check if sensor exists
	var s Sensor
	DB.First(&s, id)
	if s.ID == 0 {
		return http.StatusNotFound, AsJSON(gin.H{"error": fmt.Sprintf("Sensor '%s' not found.", id)})
	}

	err := fillQueryParams(c, &queryParams)
	if err != nil {
		return http.StatusBadRequest, AsJSON(gin.H{"error": err.Error()})
	}

	// query sensor data within defined time period
	var r []Data
	q := DB.Where("sensor_id = ?", id)

	if queryParams["start_date"] != "" {
		q = q.Where("date >= ?", queryParams["start_date"])
	}

	if queryParams["end_date"] != "" {
		q = q.Where("date <= ?", queryParams["end_date"])
	}

	q.Find(&r)

	if len(r) == 0 {
		return http.StatusNotFound, AsJSON(gin.H{"error": fmt.Sprintf("Data for Sensor %s not found.", id)})
	}

	// create empty object
	a := make([]Anomaly, 0)

	//loop over pre-filtered data to find anomalies
	for i := 0; i < len(r); i++ {

		// search for data below the specified lower value limit
		if r[i].Value < queryParams["lower_limit"].(float64) {
			a = append(a, Anomaly{Value: r[i].Value, Gradient: 0.0,
				Date: r[i].Date, Type: BelowLowerLimit})
		}

		// search for data above the specified upper value limit
		if r[i].Value > queryParams["upper_limit"].(float64) {
			a = append(a, Anomaly{Value: r[i].Value, Gradient: 0.0,
				Date: r[i].Date, Type: AboveUpperLimit})
		}

		// when the user didn't query neither max_grad nor max_diff this calculation can be skipped
		if queryParams["max_grad"] != math.MaxFloat64 {

			// calculate gradient and add anomaly if they exceed the specified maximum values
			if i > 0 {
				timeDiff := r[i].Date.Unix() - r[i-1].Date.Unix()
				valDiff := r[i].Value - r[i-1].Value
				grad := valDiff / float64(timeDiff)

				d := time.Unix(r[i-1].Date.Unix()+(timeDiff/2), 0).UTC()

				if math.Abs(grad) > queryParams["max_grad"].(float64) {
					a = append(a, Anomaly{Value: r[i-1].Value + valDiff/2, Gradient: grad,
						Date: d, Type: HighGradient})
				}
			}
		}
	}

	return http.StatusOK, AsJSON(a)
}

func fillQueryParams(c *gin.Context, m *map[string]interface{}) (e error) {
	var err error

	// loop over every entry of the passed map
	for k, v := range *m {
		p := c.Query(k)

		// using Type assertions / Type switches to either validate the query param as date string...
		// https://yourbasic.org/golang/type-assertion-switch/ and https://golang.org/ref/spec#Type_assertions
		switch tv := v.(type) {
		case string:
			(*m)[k], err = validateDateParam(p)
			if err != nil {
				return &ParamParseError{Param: k, Value: p}
			}
		// or parse the query param to a float value
		case float64:
			(*m)[k], err = parseFloatParam(p, tv)
			if err != nil {
				return &ParamParseError{Param: k, Value: p}
			}
		}
	}

	return nil
}

func validateDateParam(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	_, err := time.Parse(Layout, s)
	if err == nil {
		return s, nil
	} else {
		return "", err
	}
}

func parseFloatParam(s string, def float64) (float64, error) {
	if s == "" {
		return def, nil
	}

	return strconv.ParseFloat(s, 64)

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
	var r Sensor
	id := c.Param("id")
	DB.First(&r, id)

	if r.ID == 0 {
		return http.StatusNotFound, AsJSON(gin.H{"error": fmt.Sprintf("Sensor %s not found.", id)})
	}

	var input UpdateSensor
	if err := c.ShouldBindJSON(&input); err != nil {
		return http.StatusBadRequest, AsJSON(gin.H{"error": err.Error()})
	}

	DB.Model(&r).Update(input)

	return http.StatusOK, AsJSON(&r)
}
