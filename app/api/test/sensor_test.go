package api

import (
	"github.com/stretchr/testify/assert"
	. "github.com/vi-sense/vi-sense/app/api"
	"io/ioutil"
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

	assert.Equal(t, 200, w.Code)
}

func TestQuerySensor(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "{\"id\":1,\"room_model_id\":1,\"latest_value\":0,\"mesh_id\":\"node358\",\"name\":\"Flow Sensor\"" +
		",\"description\":\"A basic flow sensor.\",\"measurement_unit\":\"°C\",\"upper_bound\":null,\"lower_bound\"" +
		":null,\"gradient_bound\":null}"

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
	i := "{\"mesh_id\":\"node357\"}"
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	expected := "{\"id\":1,\"room_model_id\":1,\"latest_value\":0,\"mesh_id\":\"node357\",\"name\":\"Flow Sensor\"" +
		",\"description\":\"A basic flow sensor.\",\"measurement_unit\":\"°C\",\"upper_bound\":null,\"lower_bound\"" +
		":null,\"gradient_bound\":null}"

	assert.Equal(t, expected, string(body))

	// change data back
	i = "{\"mesh_id\":\"node358\"}"
	req, _ = http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestPatchSensorAllFields(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := "{\"mesh_id\":\"node357\",\"lower_bound\":7.6,\"upper_bound\":7.6,\"gradient_bound\":7.6}"
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	expected := "{\"id\":1,\"room_model_id\":1,\"latest_value\":0,\"mesh_id\":\"node357\",\"name\":\"Flow Sensor\"" +
		",\"description\":\"A basic flow sensor.\",\"measurement_unit\":\"°C\",\"upper_bound\":7.6,\"lower_bound\"" +
		":7.6,\"gradient_bound\":7.6}"

	assert.Equal(t, expected, string(body))

	// change data back
	i = "{\"mesh_id\":\"node358\",\"lower_bound\":null,\"upper_bound\":null,\"gradient_bound\":null}"
	req, _ = http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body, _ = ioutil.ReadAll(w.Body)
	expected = "{\"id\":1,\"room_model_id\":1,\"latest_value\":0,\"mesh_id\":\"node358\",\"name\":\"Flow Sensor\"" +
		",\"description\":\"A basic flow sensor.\",\"measurement_unit\":\"°C\",\"upper_bound\":null,\"lower_bound\"" +
		":null,\"gradient_bound\":null}"

	assert.Equal(t, expected, string(body))
}

func TestPatchSensorIDNotFound(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := "{\"mesh_id\":\"node357\",\"lower_bound\":7.6,\"upper_bound\":7.6,\"gradient_bound\":7.6}"
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/13", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestPatchSensorIDMalformed(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := "{\"mesh_id\":\"node357\",\"lower_bound\":7.6,\"upper_bound\":7.6,\"gradient_bound\":7.6}"
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/malformed", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestPatchSensorIgnoreInaccessibleFields(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := "{\"name\":\"name\", \"description\":\"description\"}"
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	expected := "{\"id\":1,\"room_model_id\":1,\"latest_value\":0,\"mesh_id\":\"node358\",\"name\":\"Flow Sensor\"," +
		"\"description\":\"A basic flow sensor.\",\"measurement_unit\":\"°C\",\"upper_bound\":null," +
		"\"lower_bound\":null,\"gradient_bound\":null}"

	assert.Equal(t, expected, string(body))
}

func TestPatchSensorIgnoreNewField(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := "{\"key\":\"value\"}"
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	expected := "{\"id\":1,\"room_model_id\":1,\"latest_value\":0,\"mesh_id\":\"node358\",\"name\":\"Flow Sensor\"," +
		"\"description\":\"A basic flow sensor.\",\"measurement_unit\":\"°C\",\"upper_bound\":null," +
		"\"lower_bound\":null,\"gradient_bound\":null}"

	assert.Equal(t, expected, string(body))
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
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader("{}"))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	expected := "{\"id\":1,\"room_model_id\":1,\"latest_value\":0,\"mesh_id\":\"node358\",\"name\":\"Flow Sensor\"," +
		"\"description\":\"A basic flow sensor.\",\"measurement_unit\":\"°C\"," +
		"\"upper_bound\":null,\"lower_bound\":null,\"gradient_bound\":null}"

	assert.Equal(t, expected, string(body))
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
	i := "{\"lower_bound\":7.81}"

	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"types\":[\"Below Lower Limit\"],\"date\":\"2020-01-01T00:02:00Z\"," +
		"\"value\":7.8,\"gradient\":-0.00093}]"

	assert.Equal(t, expected, w.Body.String())

	i = "{\"lower_bound\":null}"

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestQueryAnomaliesUpperBound(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := "{\"upper_bound\":7.85}"

	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"types\":[\"Above Upper Limit\"],\"date\":\"2020-01-01T00:01:00Z\"," +
		"\"value\":7.856,\"gradient\":0.00033}]"

	assert.Equal(t, expected, w.Body.String())

	i = "{\"upper_bound\":null}"

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestQueryAnomaliesGradientBound(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := "{\"gradient_bound\":0.0005}"

	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"types\":[\"High Downward Gradient\"],\"date\":\"2020-01-01T00:02:00Z\"," +
		"\"value\":7.8,\"gradient\":-0.00093}]"

	assert.Equal(t, expected, w.Body.String())

	i = "{\"gradient_bound\":null}"

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestQueryAnomaliesCombined(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := "{\"gradient_bound\":0.0001,\"upper_bound\":7.8}"

	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
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

	i = "{\"gradient_bound\":null, \"upper_bound\":null}"

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestQueryAnomaliesTimePeriod(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	i := "{\"upper_bound\":7}"

	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/sensors/1/anomalies?start_date=2020-01-01 00:01:00&end_date=2020-01-01 00:02:00", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	expected := "[{\"types\":[\"Above Upper Limit\"],\"date\":\"2020-01-01T00:01:00Z\",\"value\":7.856,\"gradient\":0.00033}]"

	assert.Equal(t, expected, w.Body.String())

	i = "{\"upper_bound\":null}"

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader(i))
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
