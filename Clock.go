package main

/*
import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/globalsign/mgo/bson"

	"github.com/globalsign/mgo"
)

type TrackMongoDB struct {
	HostURL             string
	Databasename        string
	TrackCollectionName string
}

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

var timestamp bson.ObjectId
var db TrackMongoDB

func init() {
	db = TrackMongoDB{
		"mongodb://user:test1234@ds143293.mlab.com:43293/igctracker",
		"igctracker",
		"Tracks",
	}
}

func main() {
	var track Track
	session, err := mgo.Dial(db.HostURL)
	dbSize, err := session.DB(db.Databasename).C(db.TrackCollectionName).Count()
	if err != nil {
		panic(err)
	}

	err = session.DB(db.Databasename).C(db.TrackCollectionName).Find(nil).Skip(dbSize - 1).One(&track)
	if err != nil {
		panic(err)
	}
	timestamp = track.TimeStamp
	if err != nil {
		panic(err)
	}
	for {
		time.Sleep(1 * time.Minute)
		checkNewTimestamp()
	}
}

func checkNewTimestamp() {
	var track Track
	session, err := mgo.Dial(db.HostURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	dbSize, err := session.DB(db.Databasename).C(db.TrackCollectionName).Count()
	if err != nil {
		panic(err)
	}

	err = session.DB(db.Databasename).C(db.TrackCollectionName).Find(nil).Skip(dbSize - 1).One(&track)
	if err != nil {
		panic(err)
	}
	if track.TimeStamp > timestamp {
		timestamp = track.TimeStamp
		sendMessageClock()
	}
}

func sendMessageClock() {
	message := map[string]interface{}{
		"content": "New tracks added",
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = http.Post("https://discordapp.com/api/webhooks/506082624695959558/9TWZgzeIfM87Fh-cX5c3qbrsqDnnuXeSTDnIuQ0WbhWyUuaigdfJtnVTGIFDx4Q56gWr", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		println(err.Error())
	}
}

*/
