package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var testDB *gorm.DB

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestModelsRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/models", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

//This is a hack way to add test database for each case, as whole test will just share one database.
//You can read TestWithoutAuth's comment to know how to not share database each case.
func TestMain(m *testing.M) {
	SetupTestDatabase()
	CreateMockData("../sample-data")
	exitVal := m.Run()
	DeleteTestDatabase()
	os.Exit(exitVal)
}
