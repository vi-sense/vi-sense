package api

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	. "github.com/vi-sense/vi-sense/app/api"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestQuerySensors(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors", nil)
	r.ServeHTTP(w, req)

	fmt.Println(w.Body.String())

	assert.Equal(t, 200, w.Code)
}

func TestQuerySensor(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "{\"id\":1,\"room_model_id\":1,\"latest_data\":{\"id\":3,\"sensor_id\":1,\"value\":7.8," +
		"\"gradient\":-0.00093,\"date\":\"2020-01-01T00:02:00Z\"},\"mesh_id\":\"node358\",\"name\":\"Flow Sensor\"," +
		"\"description\":\"A basic flow sensor.\",\"measurement_unit\":\"°C\",\"upper_bound\":null," +
		"\"lower_bound\":null,\"gradient_bound\":null}"

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

func TestQuerySensorData(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/data", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "[{\"id\":1,\"sensor_id\":1,\"value\":7.836,\"gradient\":0,\"date\":\"2020-01-01T00:00:00Z\"}," +
		"{\"id\":2,\"sensor_id\":1,\"value\":7.856,\"gradient\":0.00033,\"date\":\"2020-01-01T00:01:00Z\"}," +
		"{\"id\":3,\"sensor_id\":1,\"value\":7.8,\"gradient\":-0.00093,\"date\":\"2020-01-01T00:02:00Z\"}]"

	assert.Equal(t, expected, w.Body.String())
}

func TestQuerySensorDataStartDate(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/data?start_date=2020-01-01 00:02:00", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "[{\"id\":3,\"sensor_id\":1,\"value\":7.8,\"gradient\":-0.00093,\"date\":\"2020-01-01T00:02:00Z\"}]"
	assert.Equal(t, expected, w.Body.String())
}

func TestQuerySensorDataStartAndEndDate(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/data?start_date=2020-01-01 00:01:00&"+
		"end_date=2020-01-01 00:02:00", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "[{\"id\":2,\"sensor_id\":1,\"value\":7.856,\"gradient\":0.00033,\"date\":\"2020-01-01T00:01:00Z\"}]"
	assert.Equal(t, expected, w.Body.String())
}

func TestQuerySensorDataIDNotFound(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/13/data", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestQuerySensorDataIDMalformed(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/malformed/data", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestPatchSensor(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := map[string]interface{}{"mesh_id": "node357"}
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var m map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &m)

	assert.Equal(t, i["mesh_id"], m["mesh_id"])

	// change data back
	i = map[string]interface{}{"mesh_id": "node357"}
	req, _ = http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestPatchSensorAllFields(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := map[string]interface{}{"mesh_id": "node357", "lower_bound": 7.6, "upper_bound": 7.6, "gradient_bound": 7.6}
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var m map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &m)

	assert.Equal(t, i["mesh_id"], m["mesh_id"])
	assert.Equal(t, i["lower_bound"], m["lower_bound"])
	assert.Equal(t, i["upper_bound"], m["upper_bound"])
	assert.Equal(t, i["gradient_bound"], m["gradient_bound"])

	w = httptest.NewRecorder()

	// change data back
	i = map[string]interface{}{"mesh_id": "358", "lower_bound": nil, "upper_bound": nil, "gradient_bound": nil}
	req, _ = http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	fmt.Println(w.Body.String())
	_ = json.Unmarshal(w.Body.Bytes(), &m)

	assert.Equal(t, i["mesh_id"], m["mesh_id"])
	assert.Equal(t, i["lower_bound"], m["lower_bound"])
	assert.Equal(t, i["upper_bound"], m["upper_bound"])
	assert.Equal(t, i["gradient_bound"], m["gradient_bound"])
}

func TestPatchSensorIDNotFound(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := map[string]interface{}{"mesh_id": "node357", "lower_bound": 7.6, "upper_bound": 7.6, "gradient_bound": 7.6}
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/13", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestPatchSensorIDMalformed(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := map[string]interface{}{"mesh_id": "node357", "lower_bound": 7.6, "upper_bound": 7.6, "gradient_bound": 7.6}
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/malformed", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestPatchSensorIgnoreInaccessibleFields(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := map[string]interface{}{"name": "name", "description": "description"}
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var m map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &m)

	assert.NotEqual(t, i["name"], m["name"])
	assert.NotEqual(t, i["description"], m["description"])
}

func TestPatchSensorIgnoreNewField(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := map[string]interface{}{"key": "value"}
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var m map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &m)

	assert.NotEqual(t, i["key"], m["key"])
	assert.Equal(t, nil, m["key"])
}

func TestPatchSensorWrongType(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := "{\"gradient_bound\":\"value\"}"
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestPatchSensorEmptyBody(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	resp := w.Body.String()

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader("{}"))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	assert.Equal(t, resp, w.Body.String())
}

func TestPatchSensorNilBody(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestQueryAnomaliesLowerBound(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()

	i := map[string]interface{}{"lower_bound": 7.81}
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"types\":[\"Below Lower Limit\"],\"date\":\"2020-01-01T00:02:00Z\"," +
		"\"value\":7.8,\"gradient\":-0.00093}]"

	assert.Equal(t, expected, w.Body.String())

	i = map[string]interface{}{"lower_bound": nil}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestQueryAnomaliesUpperBound(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := map[string]interface{}{"upper_bound": 7.85}

	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"types\":[\"Above Upper Limit\"],\"date\":\"2020-01-01T00:01:00Z\"," +
		"\"value\":7.856,\"gradient\":0.00033}]"

	assert.Equal(t, expected, w.Body.String())

	i = map[string]interface{}{"upper_bound": nil}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestQueryAnomaliesGradientBound(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := map[string]interface{}{"gradient_bound": 0.0005}

	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"types\":[\"High Downward Gradient\"],\"date\":\"2020-01-01T00:02:00Z\"," +
		"\"value\":7.8,\"gradient\":-0.00093}]"

	assert.Equal(t, expected, w.Body.String())

	i = map[string]interface{}{"gradient_bound": nil}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestQueryAnomaliesCombined(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := map[string]interface{}{"gradient_bound": 0.0001, "upper_bound": 7.8}

	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"types\":[\"Above Upper Limit\"],\"date\":\"2020-01-01T00:00:00Z\",\"value\":7.836," +
		"\"gradient\":0},{\"types\":[\"Above Upper Limit\",\"High Upward Gradient\"]," +
		"\"date\":\"2020-01-01T00:01:00Z\",\"value\":7.856,\"gradient\":0.00033}," +
		"{\"types\":[\"High Downward Gradient\"],\"date\":\"2020-01-01T00:02:00Z\"," +
		"\"value\":7.8,\"gradient\":-0.00093}]"

	assert.Equal(t, expected, w.Body.String())

	i = map[string]interface{}{"gradient_bound": nil, "upper_bound": nil}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestQueryAnomaliesTimePeriod(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := map[string]interface{}{"upper_bound": 7}

	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies?start_date=2020-01-01 00:01:00&end_date=2020-01-01 00:02:00", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"types\":[\"Above Upper Limit\"],\"date\":\"2020-01-01T00:01:00Z\",\"value\":7.856,\"gradient\":0.00033}]"

	assert.Equal(t, expected, w.Body.String())

	i = map[string]interface{}{"upper_bound": nil}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestQueryAnomaliesNoBoundaries(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()

	w = httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/anomalies", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[]"

	assert.Equal(t, expected, w.Body.String())
}
