package main

import (
	"fmt"
	"net/http"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func handlerAdminCount(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(db.HostURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	count, err := session.DB(db.Databasename).C(db.TrackCollectionName).Count()
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		panic(err)
	}
}

func handlerAdminDeleteTrack(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(db.HostURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	_, err = session.DB(db.Databasename).C(db.TrackCollectionName).RemoveAll(bson.M{})
	if err != nil {
		panic(err)
	}

}
