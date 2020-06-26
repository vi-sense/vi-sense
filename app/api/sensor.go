package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/vi-sense/vi-sense/app/model"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type UpdateSensor struct {
	MeshId        string  `json:"mesh_id"`
	LowerBound    float64 `json:"lower_bound"`
	UpperBound    float64 `json:"upper_bound"`
	GradientBound float64 `json:"gradient_bound"`
}

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
	Type      AnomalyType `json:"type"`
	StartData *Data       `json:"start_data"`
	EndData   *Data       `json:"end_data"`
	PeakData  *Data       `json:"peak_data"`
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

	for i := range r {
		r[i].LatestData = findLatestData(&r[i])
	}

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
	r.LatestData = findLatestData(&r)
	return http.StatusOK, AsJSON(&r)
}

//QuerySensorData godoc
//@Summary Query sensor data
//@Description Query data for a specific sensor
//@Tags sensors
//@Produce json
//@Param id path int true "Sensor ID"
//@Param limit query int false "Data Limit"
//@Param density query int false "Include only every nth element [1-16]"
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
		"limit":      int64(1000),
		"density":    int64(1),
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

	if queryParams["start_date"] != "" && queryParams["end_date"] == "" {
		q = q.Where("date >= ?", queryParams["start_date"]).Limit(queryParams["limit"]).Find(&r)
	} else if queryParams["start_date"] == "" && queryParams["end_date"] != "" {
		q = q.Where("date <= ?", queryParams["end_date"]).
			Order("date desc").Limit(queryParams["limit"]).Find(&r)
	} else if queryParams["start_date"] != "" && queryParams["end_date"] != "" {
		q = q.Where("date >= ?", queryParams["start_date"]).
			Where("date <= ?", queryParams["end_date"]).Find(&r)
	} else {
		q = q.Order("date desc").Limit(queryParams["limit"]).Find(&r)
	}

	sort.Slice(r, func(i, j int) bool {
		return r[i].Date.Time.Before(r[j].Date.Time)
	})

	density := int(queryParams["density"].(int64))
	if density < 1 || density > 16 {
		return http.StatusBadRequest, AsJSON(gin.H{"error":
		fmt.Sprintf("Data density is out of its range [1-16] value=%d.", density)})
	}

	var result []Data
	if density != 1 {
		for i := 0; i < len(r); i += density {
			result = append(result, r[i])
		}
	} else {
		result = r
	}

	return http.StatusOK, AsJSON(result)
}

//QuerySensor godoc
//@Summary Query anomalies
//@Description Query anomalies for a specific sensor
//@Tags sensors
//@Produce json
//@Param id path int true "Sensor ID"
//@Param start_date query string false "Start Date"
//@Param end_date query string false "End Date"
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

	anomalies := make([]Anomaly, 0)
	currAnomalies := make([]*Anomaly, 4)

	for i := 0; i < len(r); i++ {
		if s.LowerBound != nil && r[i].Value < *s.LowerBound {
			// new anomaly occurred
			if currAnomalies[0] == nil {
				currAnomalies[0] = &Anomaly{Type: BelowLowerLimit, StartData: &r[i], EndData: nil}
				currAnomalies[0].PeakData = &r[i]
			// anomaly goes on
			} else {
				// new peak value
				if r[i].Value < currAnomalies[0].PeakData.Value {
					currAnomalies[0].PeakData = &r[i]
				}
				currAnomalies[0].EndData = &r[i]
			}

		// anomaly ended
		} else if currAnomalies[0] != nil {
			anomalies = append(anomalies, *currAnomalies[0])
			currAnomalies[0] = nil
		}

		if s.UpperBound != nil && r[i].Value > *s.UpperBound {
			// new anomaly occurred
			if currAnomalies[1] == nil {
				currAnomalies[1] = &Anomaly{Type: AboveUpperLimit, StartData: &r[i], EndData: nil}
				currAnomalies[1].PeakData = &r[i]
			// anomaly goes on
			} else {
				// new peak value
				if r[i].Value > currAnomalies[1].PeakData.Value {
					currAnomalies[1].PeakData = &r[i]
				}
				currAnomalies[1].EndData = &r[i]
			}
		// anomaly ended
		} else if currAnomalies[1] != nil {
			anomalies = append(anomalies, *currAnomalies[1])
			currAnomalies[1] = nil
		}

		if s.GradientBound != nil {
			if r[i].Gradient >= 0 && r[i].Gradient > *s.GradientBound {
				// new anomaly occurred
				if currAnomalies[2] == nil {
					currAnomalies[2] = &Anomaly{Type: UpwardGradient, StartData: &r[i], EndData: nil}
					currAnomalies[2].PeakData = &r[i]
				// anomaly goes on
				} else {
					// new peak gradient
					if r[i].Gradient > currAnomalies[2].PeakData.Gradient {
						currAnomalies[2].PeakData = &r[i]
					}
					currAnomalies[2].EndData = &r[i]
				}
			// anomaly ended
			} else if currAnomalies[2] != nil {
				anomalies = append(anomalies, *currAnomalies[2])
				currAnomalies[2] = nil
			}

			if r[i].Gradient < 0 && r[i].Gradient < -*s.GradientBound {
				// new anomaly occurred
				if currAnomalies[3] == nil {
					currAnomalies[3] = &Anomaly{Type: DownwardGradient, StartData: &r[i], EndData: nil}
					currAnomalies[3].PeakData = &r[i]
				// anomaly goes on
				} else {
					// new peak gradient
					if r[i].Gradient < currAnomalies[3].PeakData.Gradient {
						currAnomalies[3].PeakData = &r[i]
					}
					currAnomalies[3].EndData = &r[i]
				}
			// anomaly ended
			} else if currAnomalies[3] != nil {
				anomalies = append(anomalies, *currAnomalies[3])
				currAnomalies[3] = nil
			}
		}
	}

	for _, a := range currAnomalies {
		if a != nil {
			anomalies = append(anomalies, *a)
		}
	}

	return http.StatusOK, AsJSON(anomalies)
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
		case int64:
			(*m)[k], err = parseIntParam(p, tv)
			if err != nil {
				return &ParamParseError{Param: k, Value: p}
			}
		}
	}

	return nil
}

//Patch	Sensor godoc
//@Summary Update sensor preferences
//@Description Updates the mesh id and anomaly preferences of a single sensor.
//@Tags sensors
//@Accept json
//@Produce json
//@Param id path int true "SensorId"
//@Param update_sensor body UpdateSensor true "UpdateSensor"
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

	r.LatestData = findLatestData(&r)

	return http.StatusOK, AsJSON(&r)
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

func parseIntParam(s string, def int64) (int64, error) {
	if s == "" {
		return def, nil
	}

	return strconv.ParseInt(s, 10, 64)
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

func findLatestData(s *Sensor) Data {
	var d Data
	DB.Where("sensor_id = ?", s.ID).Order("date desc").First(&d)
	return d
}
