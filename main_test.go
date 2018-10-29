package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marni/goigc"
)

/*
//Connects to the test database
func setupDB(t *testing.T) *TrackMongoDB {
	db := TrackMongoDB{
		"mongodb://user:test1234@ds217092.mlab.com:17092/testtracksdb",
		"testtracksdb",
		"Tracks",
	}
	_, err := mgo.Dial(db.HostURL)
	if err != nil {
		t.Error(err)
	}
	return &db
}

//Deletes contents and drops connection to test database
func tearDownDB(t *testing.T, db *TrackMongoDB) {
	session, err := mgo.Dial(db.HostURL)
	if err != nil {
		t.Error(err)
	}
	err = session.DB(db.Databasename).DropDatabase()
	if err != nil {
		t.Error(err)
	}
}*/

//Test if you get a 400 error when posting without payload
func TestHandlerIGCNoPayload(t *testing.T) {
	//Creates a mock server
	ts := httptest.NewServer(http.HandlerFunc(handlerIGC))
	defer ts.Close()
	//Mock requests to server
	client := &http.Client{}
	//Send POST request to server
	req, err := http.NewRequest("POST", ts.URL, nil)
	if err != nil {
		t.Errorf("Not able to make POST request, %s", err)
	}
	//Get the response from server
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Can't do POST request, %s", err)
	}
	//Make sure the response is correct
	if resp.StatusCode != 400 {
		t.Errorf("Expected status 400 got %d", resp.StatusCode)
	}

}

//Check if you get a error if sending wrong request
func TestWrongRequestIGC(t *testing.T) {

	//Creates a mock server
	ts := httptest.NewServer(http.HandlerFunc(handlerIGC))
	defer ts.Close()
	//Mock requests to server
	client := &http.Client{}
	//Send DELETE request to server
	req, err := http.NewRequest("DELETE", ts.URL, nil)
	if err != nil {
		t.Errorf("Not able to make POST request, %s", err)
	}
	//Get the response from server
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Can't do POST request, %s", err)
	}
	//Make sure the response is correct
	if resp.StatusCode != 405 {
		t.Errorf("Expected status 400 got %d", resp.StatusCode)
	}

}

func TestNonIGCURL(t *testing.T) {

	payload := []byte(`{"url": "http://skypolaris.org/wp-content/uploads/IGS%drid%20to%20Jerez.igc"}`)

	//Creates a mock server
	ts := httptest.NewServer(http.HandlerFunc(handlerIGC))
	defer ts.Close()
	//Mock requests to server
	client := &http.Client{}
	//Send DELETE request to server
	req, err := http.NewRequest("POST", ts.URL, bytes.NewReader(payload))
	if err != nil {
		t.Errorf("Not able to make POST request, %s", err)
	}
	//Get the response from server
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Can't do POST request, %s", err)
	}
	//Make sure the response is correct
	if resp.StatusCode != 406 {
		t.Errorf("Expected status 400 got %d", resp.StatusCode)
	}
}

func TestTracklenght(t *testing.T) {
	distance := CalculateDistance(igc.Track{})

	if distance != 0.0 {
		t.Errorf("Expected distance 0 but got %f", distance)
	}
}

func TestInserttoDB(t *testing.T) {

}

func TestHandlerAPI(t *testing.T) {
	req, err := http.NewRequest("GET", "/paragliding/api/", nil)
	if err != nil {
		t.Error(err)
	}

	// Create a ResponseRecorder to record the response
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerAPI)

	handler.ServeHTTP(resp, req)

	// Check the status code
	if resp.Code != http.StatusOK { // It should be 200 (OK)
		t.Errorf("Handler returned wrong status got %v want %v", resp.Code, http.StatusOK)
	}
}

//Check if correct error if using a nan when getting track
func TestHandlerGetIDNan(t *testing.T) {
	req, err := http.NewRequest("GET", "/paragliding/api/igc/abs/", nil)
	if err != nil {
		t.Error(err)
	}

	// Create a ResponseRecorder to record the response
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerGetTrack)

	handler.ServeHTTP(resp, req)

	// Check the status code
	if resp.Code != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status got %v want %v", resp.Code, http.StatusBadRequest)
	}
}

func TestHandlerGetFieldNan(t *testing.T) {
	req, err := http.NewRequest("GET", "/paragliding/api/igc/abs/pilot/", nil)
	if err != nil {
		t.Error(err)
	}

	// Create a ResponseRecorder to record the response
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerGetField)

	handler.ServeHTTP(resp, req)

	// Check the status code
	if resp.Code != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status got %v want %v", resp.Code, http.StatusBadRequest)
	}
}

func TestHandlerGetFieldNotValidFeild(t *testing.T) {
	req, err := http.NewRequest("GET", "/paragliding/api/igc/1/Mybrainhurts!/", nil)
	if err != nil {
		t.Error(err)
	}

	// Create a ResponseRecorder to record the response
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerGetField)

	handler.ServeHTTP(resp, req)

	// Check the status code
	if resp.Code != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status got %v want %v", resp.Code, http.StatusBadRequest)
	}
}

func TestRedirectCorrect(t *testing.T) {
	req, err := http.NewRequest("GET", "/paragliding/", nil)
	if err != nil {
		t.Error(err)
	}

	// Create a ResponseRecorder to record the response
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerRedirect)

	handler.ServeHTTP(resp, req)

	// Check the status code
	if resp.Code != http.StatusSeeOther {
		t.Errorf("Handler returned wrong status got %v want %v", resp.Code, http.StatusSeeOther)
	}
}
