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

	s := w.Body.String()
	fmt.Println(s)
	assert.Equal(t, 5, len(a))
}

func TestQuerySensorDataStartDate(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/data?start_date=2019-10-01 00:10:00", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "[{\"id\":3,\"sensor_id\":1,\"value\":58.599918,\"gradient\":-0.00291," +
		"\"date\":\"2019-10-01T00:10:31Z\"},{\"id\":4,\"sensor_id\":1,\"value\":58.553765," +
		"\"gradient\":-0.00015,\"date\":\"2019-10-01T00:15:32Z\"},{\"id\":5,\"sensor_id\":1," +
		"\"value\":58.572021,\"gradient\":0.00006,\"date\":\"2019-10-01T00:20:33Z\"}]"

	assert.Equal(t, expected, w.Body.String())
}

func TestQuerySensorDataStartAndEndDate(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/data?start_date=2019-10-01 00:05:00&"+
		"end_date=2019-10-01 00:10:00", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "[{\"id\":2,\"sensor_id\":1,\"value\":59.50921,\"gradient\":0.00207,\"date\":\"2019-10-01T00:05:18Z\"}]"
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

	i := map[string]interface{}{"lower_bound": 59.0}
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"type\":\"Below Lower Limit\",\"start_data\":{\"id\":1,\"sensor_id\":1,\"value\":58.85," +
		"\"gradient\":0,\"date\":\"2019-10-01T00:00:00Z\"},\"end_data\":null},{\"type\":\"Below Lower Limit\"," +
		"\"start_data\":{\"id\":3,\"sensor_id\":1,\"value\":58.599918,\"gradient\":-0.00291," +
		"\"date\":\"2019-10-01T00:10:31Z\"},\"end_data\":{\"id\":5,\"sensor_id\":1,\"value\":58.572021," +
		"\"gradient\":0.00006,\"date\":\"2019-10-01T00:20:33Z\"}}]"

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
	i := map[string]interface{}{"upper_bound": 58.8}

	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"type\":\"Above Upper Limit\",\"start_data\":{\"id\":1,\"sensor_id\":1,\"value\":58.85," +
		"\"gradient\":0,\"date\":\"2019-10-01T00:00:00Z\"},\"end_data\":{\"id\":2,\"sensor_id\":1," +
		"\"value\":59.50921,\"gradient\":0.00207,\"date\":\"2019-10-01T00:05:18Z\"}}]"

	assert.Equal(t, expected, w.Body.String())

	i = map[string]interface{}{"upper_bound": nil}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestQueryAnomaliesUpwardGradient(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := map[string]interface{}{"gradient_bound": 0.002}

	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"type\":\"High Upward Gradient\",\"start_data\":{\"id\":2,\"sensor_id\":1,\"value\":59.50921," +
		"\"gradient\":0.00207,\"date\":\"2019-10-01T00:05:18Z\"},\"end_data\":null}," +
		"{\"type\":\"High Downward Gradient\",\"start_data\":{\"id\":3,\"sensor_id\":1,\"value\":58.599918," +
		"\"gradient\":-0.00291,\"date\":\"2019-10-01T00:10:31Z\"},\"end_data\":null}]"
	assert.Equal(t, expected, w.Body.String())

	i = map[string]interface{}{"gradient_bound": nil}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestQueryAnomaliesDownwardGradient(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := map[string]interface{}{"gradient_bound": 0.0029}

	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"type\":\"High Downward Gradient\",\"start_data\":{\"id\":3,\"sensor_id\":1,\"value\":58.599918," +
		"\"gradient\":-0.00291,\"date\":\"2019-10-01T00:10:31Z\"},\"end_data\":null}]"

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
	i := map[string]interface{}{"gradient_bound": 0.002, "upper_bound": 58.59}

	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"type\":\"High Upward Gradient\",\"start_data\":{\"id\":2,\"sensor_id\":1,\"value\":59.50921," +
		"\"gradient\":0.00207,\"date\":\"2019-10-01T00:05:18Z\"},\"end_data\":null},{\"type\":\"Above Upper Limit\"," +
		"\"start_data\":{\"id\":1,\"sensor_id\":1,\"value\":58.85,\"gradient\":0,\"date\":\"2019-10-01T00:00:00Z\"}," +
		"\"end_data\":{\"id\":3,\"sensor_id\":1,\"value\":58.599918,\"gradient\":-0.00291," +
		"\"date\":\"2019-10-01T00:10:31Z\"}},{\"type\":\"High Downward Gradient\",\"start_data\":{\"id\":3," +
		"\"sensor_id\":1,\"value\":58.599918,\"gradient\":-0.00291,\"date\":\"2019-10-01T00:10:31Z\"}," +
		"\"end_data\":null}]"

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
	i := map[string]interface{}{"upper_bound": 59}

	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies?start_date=2019-10-01 00:05:00&end_date=2019-10-01 00:15:00", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"type\":\"Above Upper Limit\",\"start_data\":{\"id\":2,\"sensor_id\":1,\"value\":59.50921," +
		"\"gradient\":0.00207,\"date\":\"2019-10-01T00:05:18Z\"},\"end_data\":null}]"
	assert.Equal(t, expected, w.Body.String())

	i = map[string]interface{}{"upper_bound": nil}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestQueryAnomaliesOutsideTimePeriod(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := map[string]interface{}{"upper_bound": 59}

	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(AsJSON(i)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies?start_date=2019-10-01 00:05:00&end_date=2019-10-01 00:10:00", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"type\":\"Above Upper Limit\",\"start_data\":{\"id\":2,\"sensor_id\":1,\"value\":59.50921," +
		"\"gradient\":0.00207,\"date\":\"2019-10-01T00:05:18Z\"},\"end_data\":null}]"

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
