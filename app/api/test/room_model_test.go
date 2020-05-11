package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/vi-sense/vi-sense/app/api"
)

func TestQueryRoomModels(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/models", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestQueryRoomModel(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/models/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "{\"ID\":1,\"Sensors\":[{\"ID\":1,\"RoomModelID\":1,\"Data\":null,\"MeshID\":\"node358\",\"Name\":\"Flow Sensor\",\"Description\":\"A basic flow sensor.\",\"MeasurementUnit\":\"°C\"},{\"ID\":2,\"RoomModelID\":1,\"Data\":null,\"MeshID\":\"node422\",\"Name\":\"Return Flow Sensor\",\"Description\":\"A basic return flow sensor with a longer description. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua.\",\"MeasurementUnit\":\"°C\"},{\"ID\":3,\"RoomModelID\":1,\"Data\":null,\"MeshID\":\"node441\",\"Name\":\"Fuel Sensor\",\"Description\":\"A basic thermal sensor\",\"MeasurementUnit\":\"l\"},{\"ID\":4,\"RoomModelID\":1,\"Data\":null,\"MeshID\":\"node505\",\"Name\":\"Pressure Sensor\",\"Description\":\"A basic thermal sensor\",\"MeasurementUnit\":\"bar\"}],\"Name\":\"Facility Mechanical Room\",\"Description\":\"This model shows a facility mechanical room with lots of pipes and stuff.\",\"Url\":\"files/facility-mechanical-room/model.zip\",\"ImageUrl\":\"files/facility-mechanical-room/thumbnail.png\"}"
	assert.Equal(t, expected, w.Body.String())
}

func TestQueryRoomModelIDNotFound(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/models/4", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestQueryRoomModelIDMalformed(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/models/w", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}
