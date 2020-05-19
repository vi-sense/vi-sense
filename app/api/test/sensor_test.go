package api

import (
	"github.com/vi-sense/vi-sense/app/model"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/vi-sense/vi-sense/app/api"
)

func TestQuerySensors(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestQuerySensor(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "{\"ID\":1,\"RoomModelID\":1,\"Data\":[{\"ID\":1,\"SensorID\":1,\"Value\":7.836,\"Date\":\"" +
		"2020-01-01T00:00:00Z\"},{\"ID\":2,\"SensorID\":1,\"Value\":7.856,\"Date\":\"2020-01-01T00:01:00Z\"}," +
		"{\"ID\":3,\"SensorID\":1,\"Value\":7.8,\"Date\":\"2020-01-01T00:02:00Z\"}],\"MeshID\":\"node358\",\"" +
		"Name\":\"Flow Sensor\",\"Description\":\"A basic flow sensor.\",\"MeasurementUnit\":\"°C\"}"

	assert.Equal(t, expected, w.Body.String())
}

func TestQuerySensorIDNotFound(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/13", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestQuerySensorIDMalformed(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/malformed", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestQueryAnomaliesMaxGrad(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/anomalies?max_grad=0.0009", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "[{\"Value\":7.827999999999999,\"Gradient\":-0.0009333333333333342,\"Difference\"" +
		":-0.05600000000000005,\"Type\":\"High Gradient\",\"Date\":\"2020-01-01T01:01:30+01:00\"}]"

	assert.Equal(t, expected, w.Body.String())
}

func TestQueryAnomaliesMaxDiff(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/anomalies?max_diff=0.05", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "[{\"Value\":7.827999999999999,\"Gradient\":-0.0009333333333333342,\"Difference\"" +
		":-0.05600000000000005,\"Type\":\"High Difference\",\"Date\":\"2020-01-01T01:01:30+01:00\"}]"

	assert.Equal(t, expected, w.Body.String())
}

func TestQueryAnomaliesLowerLimit(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/anomalies?lower_limit=7.81", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "[{\"Value\":7.8,\"Gradient\":0,\"Difference\":0,\"Type\":\"Below Lower Limit\",\"Date\"" +
		":\"2020-01-01T00:02:00Z\"}]"

	assert.Equal(t, expected, w.Body.String())
}

func TestQueryAnomaliesUpperLimit(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/anomalies?upper_limit=7.85", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "[{\"Value\":7.856,\"Gradient\":0,\"Difference\":0,\"Type\":\"Above Upper Limit\"" +
		",\"Date\":\"2020-01-01T00:01:00Z\"}]"

	assert.Equal(t, expected, w.Body.String())
}

func TestQueryAnomaliesStartDate(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/anomalies?upper_limit=0.0&start_date=2020-01-01 00:02:00", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "[{\"Value\":7.8,\"Gradient\":0,\"Difference\":0,\"Type\":\"Above Upper Limit\"" +
		",\"Date\":\"2020-01-01T00:02:00Z\"}]"

	assert.Equal(t, expected, w.Body.String())
}

func TestQueryAnomaliesStartEndDate(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/anomalies?upper_limit=0.0" +
		"&start_date=2020-01-01 00:00:00&end_date=2020-01-01 00:01:00", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// attention: functionality between sqlite and postgres is different, therefore different results
	expected := "[{\"Value\":7.836,\"Gradient\":0,\"Difference\":0,\"Type\":\"Above Upper Limit\"" +
		",\"Date\":\"2020-01-01T00:00:00Z\"}]"

	assert.Equal(t, expected, w.Body.String())
}

func TestQueryAnomaliesWithoutParameters(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/anomalies", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// attention: functionality between sqlite and postgres is different, therefore different results
	expected := "[]"

	assert.Equal(t, expected, w.Body.String())
}

func TestQueryAnomaliesNoResults(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/anomalies?upper_limit=10.0", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// attention: functionality between sqlite and postgres is different, therefore different results
	expected := "[]"

	assert.Equal(t, expected, w.Body.String())
}

func TestPatchSensor(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := AsJSON(&UpdateSensor{MeshID: "node357"})
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	expected := "{\"ID\":1,\"RoomModelID\":1,\"Data\":null,\"MeshID\":\"node357\",\"Name\":\"Flow " +
		"Sensor\",\"Description\":\"A basic flow sensor.\",\"MeasurementUnit\":\"°C\"}"
	assert.Equal(t, expected, string(body))

	// change data back
	i = AsJSON(&UpdateSensor{MeshID: "node358"})
	req, _ = http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestPatchSensorIDNotFound(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := AsJSON(&UpdateSensor{MeshID: "node357"})
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/13", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestPatchSensorIDMalformed(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := AsJSON(&UpdateSensor{MeshID: "node357"})
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/malformed", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestPatchSensorIgnoreInaccessibleFields(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := AsJSON(&model.Sensor{Name: "Changed Name", Description: "Changed Description"})
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	expected := "{\"ID\":1,\"RoomModelID\":1,\"Data\":null,\"MeshID\":\"node358\",\"Name\":\"Flow " +
		"Sensor\",\"Description\":\"A basic flow sensor.\",\"MeasurementUnit\":\"°C\"}"
	assert.Equal(t, expected, string(body))
}

func TestPatchSensorEmptyBody(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader("{}"))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	expected := "{\"ID\":1,\"RoomModelID\":1,\"Data\":null,\"MeshID\":\"node358\",\"Name\":\"Flow " +
		"Sensor\",\"Description\":\"A basic flow sensor.\",\"MeasurementUnit\":\"°C\"}"
	assert.Equal(t, expected, string(body))
}

func TestPatchSensorNilBody(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}