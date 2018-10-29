package main

import (
	"fmt"
	"net/http"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

//handlerAdminCount dispaly count of all tracks
func handlerAdminCount(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(db.HostURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	//Get the count of all the tracks in db
	count, err := session.DB(db.Databasename).C(db.TrackCollectionName).Count()
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		panic(err)
	}
}

//handlerAdminDeleteTrack deletes all the tracks in the db
func handlerAdminDeleteTrack(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(db.HostURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	//Deletes all tracks in db
	_, err = session.DB(db.Databasename).C(db.TrackCollectionName).RemoveAll(bson.M{})
	if err != nil {
		panic(err)
	}

}
