package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"github.com/marni/goigc"
)

const (
	//VERSION of program
	VERSION = "v1"
	//DESC is a description of the program
	DESC = "Service for Paragliding tracks."
)

// TrackMongoDB stores the DB connection
type TrackMongoDB struct {
	HostURL             string
	Databasename        string
	TrackCollectionName string
}

//MetaInfo about the program
type MetaInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

//Track is glider track info
type Track struct {
	TrackID     int           `json:"TrackID" bson:"TrackID"`
	Hdate       time.Time     `json:"H_Date" bson:"H_Date"`
	Pilot       string        `json:"pilot" bson:"pilot"`
	Glider      string        `json:"glider" bson:"glider"`
	GliderID    string        `json:"glider_id" bson:"glider_id"`
	Tracklength float64       `json:"calculated total track length" bson:"calculated total track length"`
	TrackURL    string        `json:"url" bson:"url"`
	TimeStamp   bson.ObjectId `bson:"timestamp"`
}

//Ticker keeps info about the ticker response
type Ticker struct {
	Tlatest      bson.ObjectId `json:"t_latest"`
	Tstart       bson.ObjectId `json:"t_start"`
	Tstop        bson.ObjectId `json:"t_stop"`
	TracksIds    []int         `json:"tracks"`
	Responsetime time.Duration `json:"responsetime"`
}

//Temp is a struct to keep the timestamp casue weird reasons this was the only fix I found
type temp struct {
	Timestamp bson.ObjectId
}

//WebHook stores data about the webhook
type WebHook struct {
	WebhookID       int
	URL             string
	MinTriggerValue int
	ActualValue     int
}

var startTime time.Time
var db TrackMongoDB

//ID counter
var ID int

func init() {
	startTime = time.Now()
	ID = 1
	HookID = 1
	db = TrackMongoDB{
		"mongodb://user:test1234@ds143293.mlab.com:43293/igctracker",
		"igctracker",
		"Tracks",
	}
	//Makes sure that we can connect to database
	session, err := mgo.Dial(db.HostURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	count, err := session.DB(db.Databasename).C(db.TrackCollectionName).Count()
	if err != nil {
		panic(err)
	}
	if count != 0 {
		trackss := []Track{}
		ids := []int{}
		err = session.DB(db.Databasename).C(db.TrackCollectionName).Find(bson.M{}).Sort("TrackID").All(&trackss)
		if err != nil {
			panic(err)
		}
		for _, track := range trackss {
			ids = append(ids, track.TrackID)
		}
		if ids[len(ids)-1] != 0 {
			ID = ids[len(ids)-1] + 1
		}
	}

	//TODO put extra constraints on Track collection

	WebHookInit()
}

//Uptime calculates uptime of program
func Uptime() string {
	now := time.Now()
	now.Format(time.RFC3339)
	startTime.Format(time.RFC3339)

	return now.Sub(startTime).String()
}

func handlerRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/paraglider/api/", http.StatusSeeOther)
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
	_, err = w.Write(infoJSON)
	if err != nil {
		panic(err)
	}
}

//Displays track
func handlerGetTrack(w http.ResponseWriter, r *http.Request) {
	//Check if it's a GET request, gives error 405
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, err := mgo.Dial(db.HostURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	//Gets and parses the ID from URL
	vars := mux.Vars(r)
	varID := vars["id"]
	//Converts ID to string
	id, err := strconv.Atoi(varID)
	if err != nil {
		http.Error(w, "400 Bad request", http.StatusBadRequest)
	}

	var track Track
	err = session.DB(db.Databasename).C(db.TrackCollectionName).Find(bson.M{"TrackID": id}).One(&track)
	if err != nil {
		http.Error(w, "400 Bad request", http.StatusBadRequest)
	}

	TrackJSON, err := json.Marshal(track)
	if err != nil {
		http.Error(w, "400 Bad request", http.StatusBadRequest)
	}
	//Check if track ID exists
	if track.TrackID == 0 {
		http.Error(w, "404 Not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(TrackJSON)
	if err != nil {
		panic(err)
	}
}

//Display a field from a track
func handlerGetField(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "400 Bad request", http.StatusBadRequest)
	}

	session, err := mgo.Dial(db.HostURL)
	if err != nil {
		http.Error(w, "400 Bad request", http.StatusBadRequest)
	}
	defer session.Close()

	var track Track
	err = session.DB(db.Databasename).C(db.TrackCollectionName).Find(bson.M{"TrackID": id}).One(&track)
	if err != nil {
		http.Error(w, "400 Bad request", http.StatusBadRequest)
	}

	//Map to search for fields
	trackMap := map[string]string{
		"pilot":        track.Pilot,
		"glider":       track.Glider,
		"glider_id":    track.GliderID,
		"track_length": fmt.Sprintf("%f", track.Tracklength),
		"h_date":       track.Hdate.String(),
	}
	//Makes the fields to lowercase
	field := strings.ToLower(vars["field"])
	//Finds the given field, error 400 if field is invalid
	if val, ok := trackMap[field]; ok {
		//Writes field as plain/text
		_, err = fmt.Fprintf(w, val)
		if err != nil {
			panic(err)
		}
	} else {
		http.Error(w, "400 - Bad Request, field invalid", http.StatusBadRequest)
		return

	}

}

//Displays Ids or adds Track to database
func handlerIGC(w http.ResponseWriter, r *http.Request) {

	session, err := mgo.Dial(db.HostURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	//Check if request is GET or POST
	switch r.Method {
	//Displays Ids
	case ("GET"):
		//Slice of Ids
		trackss := []Track{}
		ids := []int{}

		err := session.DB(db.Databasename).C(db.TrackCollectionName).Find(bson.M{}).All(&trackss)
		if err != nil {
			panic(err)
		}
		for _, track := range trackss {
			ids = append(ids, track.TrackID)
		}
		IDJSON, err := json.Marshal(ids)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(IDJSON)
		if err != nil {
			panic(err)
		}

	//Creates a track
	case ("POST"):

		postFunction(w, r)

	default:
		//If request is not GET or POST error 405
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func postFunction(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(db.HostURL)
	if err != nil {
		panic(err)
	}
	//URL string
	var data map[string]string
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer r.Body.Close()
	//Parses IGC file from URL
	track, err := igc.ParseLocation(data["url"])
	if err != nil {
		http.Error(w, err.Error(), 406)
		return
	}
	//Fills the struct with data
	t := Track{ID, track.Date, track.Pilot, track.GliderType, track.GliderID, CalculateDistance(track), data["url"], bson.NewObjectIdWithTime(time.Now())}
	//Insert track into database
	err = session.DB(db.Databasename).C(db.TrackCollectionName).Insert(t)
	if err != nil {
		fmt.Printf("Error in insert(): %v", err.Error())
		return
	}

	infoJSON, err := json.Marshal(t.TrackID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(infoJSON)
	if err != nil {
		panic(err)
	}
	//Adds one to ID so every track that is created is uniqe
	ID++
	UpdateWebHooks()

}

func handlerLastTicker(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(db.HostURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	var track Track
	err = session.DB(db.Databasename).C(db.TrackCollectionName).Find(bson.M{"TrackID": (ID - 1)}).One(&track)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintf(w, track.TimeStamp.String())
	if err != nil {
		panic(err)
	}

}

func handlerTicker(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	session, err := mgo.Dial(db.HostURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	tracksperpage := 5
	pageStart := 0
	pageEnd := pageStart + tracksperpage
	tracks := []Track{}
	ids := []int{}
	//Gets all the tracks from database
	err = session.DB(db.Databasename).C(db.TrackCollectionName).Find(bson.M{}).All(&tracks)
	if err != nil {
		panic(err)
	}
	//Finds the end of the page if there is less than 5 tracks
	if len(tracks) < tracksperpage {
		pageEnd = (len(tracks) - 1) % tracksperpage
	}

	tstart := tracks[pageStart].TimeStamp
	tstop := tracks[pageEnd].TimeStamp
	tlatest := tracks[ID-2].TimeStamp
	//Appends Ids slice with max tracksperpage
	for i := pageStart; i <= pageEnd; i++ {
		ids = append(ids, tracks[i].TrackID)
	}

	responsetime := time.Since(start)

	ticker := Ticker{tlatest,
		tstart,
		tstop,
		ids,
		responsetime}

	TICKERJSON, err := json.Marshal(ticker)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(TICKERJSON)
	if err != nil {
		panic(err)
	}

}

func handlerTickerTimestamp(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	session, err := mgo.Dial(db.HostURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	vars := mux.Vars(r)
	timestampstring := vars["timestamp"]

	tracksperpage := 5

	tracks := []Track{}
	ids := []int{}
	timestamp := temp{}
	//Gets the specified timestamp
	err = session.DB(db.Databasename).C(db.TrackCollectionName).Find(bson.M{"timestamp": bson.ObjectIdHex(timestampstring)}).One(&timestamp)
	if err != nil {
		panic(err)
	}
	//Gets all the tracks from database
	err = session.DB(db.Databasename).C(db.TrackCollectionName).Find(bson.M{}).All(&tracks)
	if err != nil {
		panic(err)
	}

	pageStart := 0
	pageEnd := len(tracks) - 1

	//Finds the end of the page if there is less than 5 tracks
	if len(tracks) < tracksperpage {
		pageEnd = (len(tracks) - 1) % tracksperpage
	}

	tstart := tracks[pageStart].TimeStamp
	tstop := tracks[pageEnd].TimeStamp
	tlatest := tracks[ID-2].TimeStamp
	//Appends Ids slice with max tracksperpage
	for i := 0; i < tracksperpage; i++ {
		if timestamp.Timestamp < tracks[i].TimeStamp {
			ids = append(ids, tracks[i].TrackID)
		}
	}

	responsetime := time.Since(start)

	ticker := Ticker{tlatest,
		tstart,
		tstop,
		ids,
		responsetime}

	TICKERJSON, err := json.Marshal(ticker)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(TICKERJSON)
	if err != nil {
		panic(err)
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
	webhookmain()
	router := mux.NewRouter()
	router.HandleFunc("/paraglider/", handlerRedirect)
	router.HandleFunc("/paraglider/api/", handlerAPI)
	router.HandleFunc("/paraglider/api/igc/", handlerIGC)
	router.HandleFunc("/paraglider/api/igc/{id:[0-9]+}/", handlerGetTrack)
	router.HandleFunc("/paraglider/api/igc/{id:[0-9]+}/{field:[a-zA-Z_]+}/", handlerGetField)
	router.HandleFunc("/paraglider/api/ticker/", handlerTicker)
	router.HandleFunc("/paraglider/api/ticker/latest/", handlerLastTicker)
	router.HandleFunc("/paraglider/api/ticker/{timestamp:[a-zA-Z0-9_]+}/", handlerTickerTimestamp)
	router.HandleFunc("/paraglider/api/webhook/new_track/", handlerNewWebHook)
	router.HandleFunc("/paraglider/api/webhook/new_track/{webhook:[a-zA-Z0-9_]+}/", handlerGetWebHook).Methods("GET")
	router.HandleFunc("/paraglider/api/webhook/new_track/{webhook:[a-zA-Z0-9_]+}/", handlerDeleteWebHook).Methods("DELETE")
	router.HandleFunc("/paraglider/admin/api/tracks_count/", handlerAdminCount).Methods("GET")
	router.HandleFunc("/paraglider/admin/api/tracks/", handlerAdminDeleteTrack).Methods("DELETE")
	err := http.ListenAndServe(GetPort(), router)
	if err != nil {
		panic(err)
	}
}
