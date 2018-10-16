package app

import (
	"github.com/alexrios/challenge-api/db"
	"github.com/alexrios/challenge-api/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	app   App
	mongo db.MockTransport
)

//Testa o retorno de uma lista vazia
func TestEmptyListPlanets(t *testing.T) {
	mongo := &db.MockTransport{}
	app := &App{
		DB: mongo,
	}
	app.MakeRoutes()
	var planets []*models.Planet

	mongo.On("FindAll", &planets).Return(nil)

	req, _ := http.NewRequest("GET", "/planets", nil)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	if !strings.Contains(strings.TrimSpace(rr.Body.String()), "[]") {
		t.Fatal("Empty planets on DB should return empty slice as response")
	}
}
