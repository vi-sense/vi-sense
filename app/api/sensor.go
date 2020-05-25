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

type ParamParseError struct {
	Param string
	Value string
}

func (e *ParamParseError) Error() string {
	if e.Value != "" {
		return fmt.Sprintf("Error parsing parameter '%s' with value '%s'.", e.Param, e.Value)
	}
	return fmt.Sprintf("Error parsing parameter '%s'.", e.Param)
}

type Anomaly struct {
	Types    []AnomalyType `json:"types"`
	Date     time.Time     `json:"date"`
	Value    float64       `json:"value"`
	Gradient float64       `json:"gradient"`
}

type AnomalyType string

const (
	UpwardGradient   AnomalyType = "High Upward Gradient"
	DownwardGradient AnomalyType = "High Downward Gradient"
	AboveUpperLimit  AnomalyType = "Above Upper Limit"
	BelowLowerLimit  AnomalyType = "Below Lower Limit"
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
		"start_date": "",
		"end_date":   "",
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
		"start_date": "",
		"end_date":   "",
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

	q := DB.Where("sensor_id = ?", id)

	if queryParams["start_date"] != "" {
		q = q.Where("date >= ?", queryParams["start_date"])
	}

	if queryParams["end_date"] != "" {
		q = q.Where("date <= ?", queryParams["end_date"])
	}

	var r []Data
	q.Find(&r)

	if len(r) == 0 {
		return http.StatusNotFound, AsJSON(gin.H{"error": fmt.Sprintf("Data for Sensor %s not found.", id)})
	}

	a := make([]Anomaly, 0)

	for _, d := range r {
		var types []AnomalyType

		// search for data below the specified lower value limit
		if s.LowerBound != nil && d.Value < *s.LowerBound {
			types = append(types, BelowLowerLimit)
		}

		// search for data below the specified lower value limit
		if s.UpperBound != nil && d.Value > *s.UpperBound {
			types = append(types, AboveUpperLimit)
		}

		if s.GradientBound != nil && (math.Abs(d.Gradient) > *s.GradientBound) {
			if d.Gradient >= 0 {
				types = append(types, UpwardGradient)
			} else {
				types = append(types, DownwardGradient)
			}
		}

		if len(types) > 0 {
			a = append(a, Anomaly{Gradient: d.Gradient, Types: types, Value: d.Value, Date: d.Date})
		}
	}

	return http.StatusOK, AsJSON(a)
}

func fillQueryParams(c *gin.Context, m *map[string]interface{}) error {
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

	var i map[string]interface{}
	if err := c.ShouldBindJSON(&i); err != nil {
		return http.StatusBadRequest, AsJSON(gin.H{"error": err.Error()})
	}

	if err := validateUpdateValues(&i); err != nil {
		return http.StatusBadRequest, AsJSON(gin.H{"error": err.Error()})
	}

	DB.Model(&r).Update(i)

	return http.StatusOK, AsJSON(&r)
}

func validateUpdateValues(m *map[string]interface{}) error {
	var unknown []string

	for k, v := range *m {
		switch k {
		case "mesh_id":
			if _, ok := v.(string); !ok {
				return &ParamParseError{
					Param: k,
				}
			}

		case "lower_bound", "upper_bound", "gradient_bound":
			if _, ok := v.(float64); !ok && v != nil {
				return &ParamParseError{
					Param: k,
				}
			}
		default:
			unknown = append(unknown, k)
		}
	}

	for _, u := range unknown {
		delete(*m, u)
	}

	return nil
}
