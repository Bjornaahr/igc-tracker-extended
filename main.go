package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/marni/goigc"
)

const (
	//VERSION of program
	VERSION = "1.0"
	//DESC is a description of the program
	DESC = "Service for IGC tracks."
)

//MetaInfo about the program
type MetaInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

//Track is glider track info
type Track struct {
	ID          int       `json:"ID"`
	Hdate       time.Time `json:"H_Date"`
	Pilot       string    `json:"pilot"`
	Glider      string    `json:"glider"`
	GliderID    string    `json:"glider_id"`
	Tracklength float64   `json:"calculated total track length"`
}

var startTime time.Time
var tracks map[int]Track

//ID counter
var ID int

func init() {
	startTime = time.Now()
	tracks = make(map[int]Track)
	ID = 1
}

//Uptime calculates uptime of program
func Uptime() string {
	now := time.Now()
	now.Format(time.RFC3339)
	startTime.Format(time.RFC3339)

	return now.Sub(startTime).String()
}

//Displays metadata
func handlerAPI(w http.ResponseWriter, r *http.Request) {
	//Gets uptime without decimal points
	time := strings.Split(Uptime(), ".")
	info := MetaInfo{time[0],
		DESC,
		VERSION}
	infoJSON, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}
	//Set headertype, status and write the metadata
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(infoJSON, ErrBodyNotAllowed)

}

//Displays track
func handlerGetTrack(w http.ResponseWriter, r *http.Request) {
	//Check if it's a GET request, gives error 405
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//Gets and parses the ID from URL
	vars := mux.Vars(r)
	varID := vars["id"]
	//Converts ID to string
	id, err := strconv.Atoi(varID)
	if err != nil {
		panic(err)
	}

	TrackJSON, err := json.Marshal(tracks[id])
	if err != nil {
		panic(err)
	}
	//Check if track ID exists
	if tracks[id].ID == 0 {
		http.Error(w, "404 Not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(TrackJSON)
}

//Display a field from a track
func handlerGetField(w http.ResponseWriter, r *http.Request) {

	//Gets and parses the ID from URL
	vars := mux.Vars(r)
	varID := vars["id"]
	//Converts ID to string
	id, err := strconv.Atoi(varID)
	if err != nil {
		panic(err)
	}

	//Map to search for fields
	trackMap := map[string]string{
		"pilot":        tracks[id].Pilot,
		"glider":       tracks[id].Glider,
		"glider_id":    tracks[id].GliderID,
		"track_length": fmt.Sprintf("%f", tracks[id].Tracklength),
		"h_date":       tracks[id].Hdate.String(),
	}
	//Makes the fields to lowercase
	field := strings.ToLower(vars["field"])
	//Finds the given field, error 400 if field is invalid
	if val, ok := trackMap[field]; ok {
		//Writes field as plain/text
		fmt.Fprintf(w, val)
	} else {
		http.Error(w, "400 - Bad Request, field invalid", http.StatusBadRequest)
		return

	}

}

//Displays Ids or adds Track to memory
func handlerIGC(w http.ResponseWriter, r *http.Request) {

	//Check if request is GET or POST
	switch r.Method {
	//Displays Ids
	case ("GET"):
		//Slice of Ids
		ids := []int{}
		//Add id of track to slice
		for index := range tracks {
			ids = append(ids, tracks[index].ID)
		}
		IDJSON, err := json.Marshal(ids)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(IDJSON)

	//Creates a track
	case ("POST"):
		//URL string
		var url string
		err := json.NewDecoder(r.Body).Decode(&url)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		//Parses IGC file from URL
		track, err := igc.ParseLocation(url)
		//Fills in values in track
		tracks[ID] = Track{ID, track.Date, track.Pilot, track.GliderType, track.GliderID, CalculateDistance(track)}

		infoJSON, err := json.Marshal(tracks[ID].ID)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(infoJSON)
		//Adds one to ID so every track that is created is uniqe
		ID++

	default:
		//If request is not GET or POST error 405
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

//CalculateDistance calculates the track distance
func CalculateDistance(track igc.Track) float64 {

	trackdistance := 0.0
	//Loops through all the points and find the distance between them
	for i := 0; i < len(track.Points)-1; i++ {
		trackdistance += track.Points[i].Distance(track.Points[i+1])
	}

	return trackdistance
}

//GetPort retrives the port from the enviorment
func GetPort() string {
	//Gets the port
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "5000"
		fmt.Println("Could not find port in enviorment, setting port to: " + port)
	}
	return ":" + port
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/igcinfo/api/", handlerAPI)
	router.HandleFunc("/igcinfo/api/igc/", handlerIGC)
	router.HandleFunc("/igcinfo/api/igc/{id:[0-9]+}/", handlerGetTrack)
	router.HandleFunc("/igcinfo/api/igc/{id:[0-9]+}/{field:[a-zA-Z_]+}/", handlerGetField)
	http.ListenAndServe(GetPort(), router)
}
