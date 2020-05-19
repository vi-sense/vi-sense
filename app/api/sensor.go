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
	MeshID string
}

type ParamParseError struct {
	Param string
	Value string
}

func (e *ParamParseError) Error() string {
	return fmt.Sprintf("Error parsing parameter %s with value %s.", e.Param, e.Value)
}

type AnomalyQueryParams struct {
	StartDate  string
	EndDate    string
	MaxDiff    float64
	MaxGrad    float64
	LowerLimit float64
	UpperLimit float64
}

type Anomaly struct {
	Value      float64
	Gradient   float64
	Difference float64
	Type       AnomalyType
	Date       time.Time
}

type AnomalyType string

const (
	HighDifference  AnomalyType = "High Difference"
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
	var r Sensor
	id := c.Param("id")
	DB.Preload("Data").First(&r, id)
	if r.ID == 0 {
		return http.StatusNotFound, AsJSON(gin.H{"error": fmt.Sprintf("Sensor %s not found.", id)})
	}
	return http.StatusOK, AsJSON(&r)
}

//QuerySensor godoc
//@Summary Query anomalies
//@Description Query anomalies for a specific sensor
//@Tags sensors
//@Produce json
//@Param start_date query string false "Start Date"
//@Param end_date query string false "End Date"
//@Param max_diff query number false "Maximum Difference"
//@Param max_grad query number false "Maximum Gradient"
//@Param lower_limit query number false "Lower Value Limit"
//@Param upper_limit query number false "Upper Value Limit"
//@Success 200 {array} Anomaly
//@Failure 400 {string} string "bad request"
//@Failure 404 {string} string "not found"
//@Failure 500 {string} string "internal server error"
//@Router /sensors/{id}/anomalies [get]
func QueryAnomalies(c *gin.Context) (int, string) {
	params, err := getQueryParams(c)
	if err != nil {
		return http.StatusBadRequest, AsJSON(gin.H{"error": err.Error()})
	}

	id := c.Param("id")

	// check if sensor exists
	var s Sensor
	DB.First(&s, id)
	if s.ID == 0 {
		return http.StatusNotFound, AsJSON(gin.H{"error": fmt.Sprintf("Sensor %s not found.", id)})
	}

	// query sensor data within defined time period
	var r []Data
	q := DB.Where("sensor_id = ?", id)

	if params.StartDate != "" {
		q = q.Where("date >= ?", params.StartDate)
	}

	if params.EndDate != "" {
		q = q.Where("date <= ?", params.EndDate)
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
		if r[i].Value < params.LowerLimit {
			a = append(a, Anomaly{Value: r[i].Value, Gradient: 0.0, Difference: 0.0,
				Date: r[i].Date, Type: BelowLowerLimit})
		}

		// search for data above the specified upper value limit
		if r[i].Value > params.UpperLimit {
			a = append(a, Anomaly{Value: r[i].Value, Gradient: 0.0, Difference: 0.0,
				Date: r[i].Date, Type: AboveUpperLimit})
		}

		// when the user didn't query neither max_grad nor max_diff this calculation can be skipped
		if params.MaxGrad != math.MaxFloat64 || params.MaxDiff != math.MaxFloat64 {

			// calculate gradient and value difference and add anomaly if they exceed the specified maximum values
			if i > 0 {
				timeDiff := r[i].Date.Unix() - r[i-1].Date.Unix()
				valDiff := r[i].Value - r[i-1].Value
				grad := valDiff / float64(timeDiff)

				d := time.Unix(r[i-1].Date.Unix()+(timeDiff/2), 0).UTC()

				if math.Abs(valDiff) > params.MaxDiff {
					a = append(a, Anomaly{Value: r[i-1].Value + valDiff/2, Gradient: grad, Difference: valDiff,
						Date: d, Type: HighDifference})
				}
				if math.Abs(grad) > params.MaxGrad {
					a = append(a, Anomaly{Value: r[i-1].Value + valDiff/2, Gradient: grad, Difference: valDiff,
						Date: d, Type: HighGradient})
				}
			}
		}
	}

	return http.StatusOK, AsJSON(a)
}

func getQueryParams(c *gin.Context) (p *AnomalyQueryParams, e error) {
	var params AnomalyQueryParams
	var err error

	v := c.Query("start_date")
	params.StartDate, err = validateDateParam(v)
	if err != nil {
		return nil, &ParamParseError{Param: "start_date", Value: v}
	}

	v = c.Query("end_date")
	params.EndDate, err = validateDateParam(v)
	if err != nil {
		return nil, &ParamParseError{Param: "end_date", Value: v}
	}

	v = c.Query("max_diff")
	params.MaxDiff, err = parseFloatParam(v, math.MaxFloat64)
	if err != nil {
		return nil, &ParamParseError{Param: "max_diff"}
	}

	v = c.Query("max_grad")
	params.MaxGrad, err = parseFloatParam(v, math.MaxFloat64)
	if err != nil {
		return nil, &ParamParseError{Param: "max_grad"}
	}

	v = c.Query("upper_limit")
	params.UpperLimit, err = parseFloatParam(v, math.MaxFloat64)
	if err != nil {
		return nil, &ParamParseError{Param: "upper_limit"}
	}

	v = c.Query("lower_limit")
	params.LowerLimit, err = parseFloatParam(v, -math.MaxFloat64)
	if err != nil {
		return nil, &ParamParseError{Param: "lower_limit"}
	}

	return &params, nil
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
