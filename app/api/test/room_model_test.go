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
	req, _ := http.NewRequest(http.MethodGet, "/models", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestQueryRoomModel(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/models/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "{\"id\":1,\"sensors\":[{\"id\":1,\"room_model_id\":1,\"latest_data\":{\"id\":0,\"sensor_id\":0," +
		"\"value\":0,\"gradient\":0,\"date\":\"0001-01-01T00:00:00Z\"},\"mesh_id\":\"node358\"," +
		"\"name\":\"Flow Sensor\",\"description\":\"A basic flow sensor.\",\"measurement_unit\":\"°C\"," +
		"\"upper_bound\":null,\"lower_bound\":null,\"gradient_bound\":null},{\"id\":2,\"room_model_id\":1," +
		"\"latest_data\":{\"id\":0,\"sensor_id\":0,\"value\":0,\"gradient\":0,\"date\":\"0001-01-01T00:00:00Z\"}," +
		"\"mesh_id\":\"node422\",\"name\":\"Return Flow Sensor\"," +
		"\"description\":\"A basic return flow sensor with a longer description. Lorem ipsum dolor sit amet," +
		" consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat," +
		" sed diam voluptua.\",\"measurement_unit\":\"°C\",\"upper_bound\":null," +
		"\"lower_bound\":null,\"gradient_bound\":null},{\"id\":3,\"room_model_id\":1," +
		"\"latest_data\":{\"id\":0,\"sensor_id\":0,\"value\":0,\"gradient\":0,\"date\":\"0001-01-01T00:00:00Z\"}," +
		"\"mesh_id\":\"node441\",\"name\":\"Fuel Sensor\",\"description\":\"A basic thermal sensor\"," +
		"\"measurement_unit\":\"l\",\"upper_bound\":null,\"lower_bound\":null,\"gradient_bound\":null}," +
		"{\"id\":4,\"room_model_id\":1,\"latest_data\":{\"id\":0,\"sensor_id\":0,\"value\":0,\"gradient\":0," +
		"\"date\":\"0001-01-01T00:00:00Z\"},\"mesh_id\":\"node505\",\"name\":\"Pressure Sensor\"," +
		"\"description\":\"A basic thermal sensor\",\"measurement_unit\":\"bar\",\"upper_bound\":null," +
		"\"lower_bound\":null,\"gradient_bound\":null}],\"name\":\"Facility Mechanical Room\"," +
		"\"description\":\"This model shows a facility mechanical room with lots of pipes and stuff.\"," +
		"\"url\":\"files/facility-mechanical-room/model.zip\"," +
		"\"image_url\":\"files/facility-mechanical-room/thumbnail.png\"}"

	assert.Equal(t, expected, w.Body.String())
}

func TestQueryRoomModelIDNotFound(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/models/4", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestQueryRoomModelIDMalformed(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/models/w", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}
