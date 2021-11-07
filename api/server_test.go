package api

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mjevans93308/avoxi-demo-app/config"
)

var a App

func TestHomeHandler(t *testing.T) {
	a.Initialize()
	req, _ := http.NewRequest(http.MethodGet, config.HOME, nil)

	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	if http.StatusOK != rr.Code {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, rr.Code)
	}

	log.Print(rr.Body.String())
	expected := "hello, world!"

	if body := rr.Body.String(); body != expected {
		t.Errorf("Expected %s. Got %s", expected, body)
	}
}

func TestAliveHandler(t *testing.T) {
	a.Initialize()
	req, _ := http.NewRequest(http.MethodGet, config.API_GROUP+config.V1_GROUP+config.ALIVE, nil)

	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	if http.StatusOK != rr.Code {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, rr.Code)
	}

	log.Print(rr.Body.String())
	expected := "It's...ALIVE!!!"

	if body := rr.Body.String(); body != expected {
		t.Errorf("Expected %s. Got %s", expected, body)
	}
}