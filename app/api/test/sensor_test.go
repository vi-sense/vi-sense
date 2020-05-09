package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/vi-sense/vi-sense/app/api"
)

func TestQuerySensors(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/sensors", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestQuerySensor(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/sensors/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "{\"ID\":1,\"RoomModelID\":1,\"Data\":[{\"ID\":1,\"SensorID\":1,\"Value\":7.836,\"Date\":\"2020-01-01T00:00:00Z\"},{\"ID\":2,\"SensorID\":1,\"Value\":7.856,\"Date\":\"2020-01-01T00:01:00Z\"},{\"ID\":3,\"SensorID\":1,\"Value\":7.8,\"Date\":\"2020-01-01T00:02:00Z\"}],\"MeshID\":\"node358\",\"Name\":\"Flow Sensor\",\"Description\":\"A basic flow sensor.\",\"MeasurementUnit\":\"°C\"}"
	assert.Equal(t, expected, w.Body.String())
}

func TestQuerySensorIDNotFound(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/sensors/13", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestQuerySensorIDMalformed(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/sensors/malformed", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}