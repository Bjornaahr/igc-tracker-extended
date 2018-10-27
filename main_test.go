package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
)

func setupDB(t *testing.T) *TrackMongoDB {
	db := TrackMongoDB{
		"mongodb://user:test1234@ds217092.mlab.com:17092/testtracksdb",
		"testtracksdb",
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
	req, _ := http.NewRequest("POST", "/api", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	router := mux.NewRouter()

	router.ServeHTTP(w, req)

	defer tearDownDB(t, db)

}

func TestHandlerAPI(t *testing.T) {

	//req := httptest.NewRequest("GET", "igcinfo/api/", nil)
	res := httptest.NewRecorder()

	assert.Equal(t, http.StatusOK, res.Code)

}
