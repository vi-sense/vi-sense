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