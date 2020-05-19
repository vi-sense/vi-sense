package api

import (
	"github.com/stretchr/testify/assert"
	. "github.com/vi-sense/vi-sense/app/api"
	"github.com/vi-sense/vi-sense/app/model"
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

	expected := "{\"id\":1,\"room_model_id\":1,\"mesh_id\":\"node358\",\"name\":\"Flow Sensor\"" +
		",\"description\":\"A basic flow sensor.\",\"measurement_unit\":\"°C\"}"

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

	expected := "[{\"id\":1,\"sensor_id\":1,\"value\":7.836,\"date\":\"2020-01-01T00:00:00Z\"}," +
		"{\"id\":2,\"sensor_id\":1,\"value\":7.856,\"date\":\"2020-01-01T00:01:00Z\"}," +
		"{\"id\":3,\"sensor_id\":1,\"value\":7.8,\"date\":\"2020-01-01T00:02:00Z\"}]"
	assert.Equal(t, expected, w.Body.String())
}

func TestQuerySensorDataStartDate(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/data?start_date=2020-01-01 00:02:00", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "[{\"id\":3,\"sensor_id\":1,\"value\":7.8,\"date\":\"2020-01-01T00:02:00Z\"}]"
	assert.Equal(t, expected, w.Body.String())
}

func TestQuerySensorDataStartAndEndDate(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/data?start_date=2020-01-01 00:01:00&"+
		"end_date=2020-01-01 00:02:00", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "[{\"id\":2,\"sensor_id\":1,\"value\":7.856,\"date\":\"2020-01-01T00:01:00Z\"}]"
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

func TestQueryAnomaliesMaxGrad(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/anomalies?max_grad=0.0009", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "[{\"type\":\"High Gradient\",\"date\":\"2020-01-01T00:01:30Z\",\"value\":7.827999999999999,\"" +
		"gradient\":-0.0009333333333333342}]"

	assert.Equal(t, expected, w.Body.String())
}

func TestQueryAnomaliesLowerLimit(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/anomalies?lower_limit=7.81", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "[{\"type\":\"Below Lower Limit\",\"date\":\"2020-01-01T00:02:00Z\",\"value\":7.8,\"gradient\":0}]"

	assert.Equal(t, expected, w.Body.String())
}

func TestQueryAnomaliesUpperLimit(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/anomalies?upper_limit=7.85", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "[{\"type\":\"Above Upper Limit\",\"date\":\"2020-01-01T00:01:00Z\",\"value\"" +
		":7.856,\"gradient\":0}]"

	assert.Equal(t, expected, w.Body.String())
}

func TestQueryAnomaliesStartDate(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/anomalies?upper_limit=0.0&start_date=2020-01-01 00:02:00", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expected := "[{\"type\":\"Above Upper Limit\",\"date\":\"2020-01-01T00:02:00Z\",\"value\"" +
		":7.8,\"gradient\":0}]"

	assert.Equal(t, expected, w.Body.String())
}

func TestQueryAnomaliesStartEndDate(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sensors/1/anomalies?upper_limit=0.0"+
		"&start_date=2020-01-01 00:00:00&end_date=2020-01-01 00:01:00", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// attention: functionality between sqlite and postgres is different, therefore different results
	expected := "[{\"type\":\"Above Upper Limit\",\"date\":\"2020-01-01T00:00:00Z\",\"value\"" +
		":7.836,\"gradient\":0}]"

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
	expected := "{\"id\":1,\"room_model_id\":1,\"mesh_id\":\"node357\",\"name\":\"Flow Sensor\"" +
		",\"description\":\"A basic flow sensor.\",\"measurement_unit\":\"°C\"}"
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
	expected := "{\"id\":1,\"room_model_id\":1,\"mesh_id\":\"node358\",\"name\":\"Flow Sensor\"" +
		",\"description\":\"A basic flow sensor.\",\"measurement_unit\":\"°C\"}"
	assert.Equal(t, expected, string(body))
}

func TestPatchSensorEmptyBody(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", strings.NewReader("{}"))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body, _ := ioutil.ReadAll(w.Body)
	expected := "{\"id\":1,\"room_model_id\":1,\"mesh_id\":\"node358\",\"name\":\"Flow Sensor\"" +
		",\"description\":\"A basic flow sensor.\",\"measurement_unit\":\"°C\"}"
	assert.Equal(t, expected, string(body))
}

func TestPatchSensorNilBody(t *testing.T) {
	r := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/sensors/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}
