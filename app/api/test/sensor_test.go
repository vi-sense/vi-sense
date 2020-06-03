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

	var m map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &m)

	assert.Equal(t, 1.0, m["id"])
	assert.Equal(t, 1.0, m["room_model_id"])
}

func TestQuerySensorIDNotFound(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1000", nil)
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

	var a []interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, 3, len(a))
}

func TestQuerySensorDataStartDate(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/data?start_date=2019-10-01 00:10:00", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "[{\"id\":3,\"sensor_id\":1,\"value\":11.910317,\"gradient\":" +
		"-0.00018,\"date\":\"2019-10-01T00:10:07Z\"}]"
	assert.Equal(t, expected, w.Body.String())
}

func TestQuerySensorDataStartAndEndDate(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/data?start_date=2019-10-01 00:05:00&"+
		"end_date=2019-10-01 00:10:00", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "[{\"id\":2,\"sensor_id\":1,\"value\":11.965417,\"gradient\":" +
		"-0.00018,\"date\":\"2019-10-01T00:05:02Z\"}]"
	assert.Equal(t, expected, w.Body.String())
}

func TestQuerySensorDataIDNotFound(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1000/data", nil)
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
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1000", strings.NewReader(AsJSON(i)))
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

	i := map[string]interface{}{"lower_bound": 11.95}
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"types\":[\"Below Lower Limit\"],\"date\":\"2019-10-01T00:10:07Z\",\"value\":" +
		"11.910317,\"gradient\":-0.00018}]"

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
	i := map[string]interface{}{"upper_bound": 12}

	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"types\":[\"Above Upper Limit\"],\"date\":\"2019-10-01T00:00:00Z\",\"value\":12.02,\"gradient\":0}]"

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
	i := map[string]interface{}{"gradient_bound": 0.0001}

	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"types\":[\"High Downward Gradient\"],\"date\":\"2019-10-01T00:05:02Z\",\"value\":" +
		"11.965417,\"gradient\":-0.00018},{\"types\":[\"High Downward Gradient\"]," +
		"\"date\":\"2019-10-01T00:10:07Z\",\"value\":11.910317,\"gradient\":-0.00018}]"

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
	i := map[string]interface{}{"gradient_bound": 0.0001, "upper_bound": 11.96}

	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"types\":[\"Above Upper Limit\"],\"date\":\"2019-10-01T00:00:00Z\",\"value\":12.02," +
		"\"gradient\":0},{\"types\":[\"Above Upper Limit\",\"High Downward Gradient\"]," +
		"\"date\":\"2019-10-01T00:05:02Z\",\"value\":11.965417,\"gradient\":-0.00018}," +
		"{\"types\":[\"High Downward Gradient\"],\"date\":\"2019-10-01T00:10:07Z\",\"value\":11.910317," +
		"\"gradient\":-0.00018}]"

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
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies?start_date=2019-10-01 00:05:00&end_date=2019-10-01 00:10:00", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"types\":[\"Above Upper Limit\"],\"date\":\"2019-10-01T00:05:02Z\"," +
		"\"value\":11.965417,\"gradient\":-0.00018}]"

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
