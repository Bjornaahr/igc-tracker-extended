package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
)

func setupDB(t *testing.T) *TrackMongoDB {
	db := TrackMongoDB{
		"localhost",
		"testTracksDB",
		"Tracks",
	}
	_, err := mgo.Dial(db.HostURL)
	t.Error(err)

	return &db
}

func tearDownDB(t *testing.T, db *TrackMongoDB) {
	session, err := mgo.Dial(db.HostURL)
	t.Error(err)

	err = session.DB(db.Databasename).DropDatabase()
	if err != nil {
		t.Error(err)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	Touter := mux.NewRouter()
	Touter.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestHandlerIGC(t *testing.T) {
	db := setupDB(t)

	payload := []byte(`{"url": "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"}`)
	req, _ := http.NewRequest("POST", "localhost:5000/igcinfo/api/igc/", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)
	var m string
	json.Unmarshal(response.Body.Bytes(), &m)

	if m != "1" {
		t.Errorf("Expected id to be '1'. Got '%v'", m)
	}
	defer tearDownDB(t, db)

}
