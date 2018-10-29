package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

//HookID global variable for id of the webhook
var HookID int
var dbh TrackMongoDB

//WebHookInit connects to a database and set's correct id
func WebHookInit() {
	dbh = TrackMongoDB{
		"mongodb://user:test1234@ds143293.mlab.com:43293/igctracker",
		"igctracker",
		"WebHooks",
	}
	WebHookSetID()
}

//WebHookSetID returns the correct id from the database
func WebHookSetID() {
	session, err := mgo.Dial(dbh.HostURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	count, err := session.DB(dbh.Databasename).C(dbh.TrackCollectionName).Count()
	if err != nil {
		panic(err)
	}
	HookID = 1

	if count != 0 {
		webhooks := []WebHook{}
		ids := []int{}
		err = session.DB(dbh.Databasename).C(dbh.TrackCollectionName).Find(bson.M{}).Sort("webhookid").All(&webhooks)
		if err != nil {
			panic(err)
		}
		for _, hook := range webhooks {
			ids = append(ids, hook.WebhookID)
		}
		if ids[len(ids)-1] != 0 {
			HookID = ids[len(ids)-1] + 1
		}
	}
}

func handlerNewWebHook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	session, err := mgo.Dial(dbh.HostURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var data map[string]string
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	defer r.Body.Close()
	triggerValue, err := strconv.Atoi(data["minTriggerValue"])
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	//Creates and insert webhook to db
	web := WebHook{HookID, data["webhookURL"], triggerValue, 0}
	err = session.DB(dbh.Databasename).C(dbh.TrackCollectionName).Insert(&web)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	infoJSON, err := json.Marshal(HookID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	WebHookSetID()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(infoJSON)
	if err != nil {
		panic(err)
	}
}

//Gets specified webhook
func handlerGetWebHook(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(dbh.HostURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	vars := mux.Vars(r)
	webhookID := vars["webhook"]

	id, err := strconv.Atoi(webhookID)
	if err != nil {
		panic(err)
	}

	webhook := WebHook{}
	//Gets webhook from db and puts data in webhook
	err = session.DB(dbh.Databasename).C(dbh.TrackCollectionName).Find(bson.M{"webhookid": id}).One(&webhook)
	if err != nil {
		panic(err)
	}
	//Maps the data we want to display
	webhookmap := map[string]string{
		"webhookURL":      webhook.URL,
		"minTriggerValue": strconv.Itoa(webhook.MinTriggerValue),
	}

	WEBJSON, err := json.Marshal(webhookmap)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(WEBJSON)
	if err != nil {
		panic(err)
	}

}

//Deletes specified webhook
func handlerDeleteWebHook(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(dbh.HostURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	vars := mux.Vars(r)
	webhookID := vars["webhook"]

	id, err := strconv.Atoi(webhookID)

	if err != nil {
		panic(err)
	}
	handlerGetWebHook(w, r)
	//Deletes webhook from db
	err = session.DB(dbh.Databasename).C(dbh.TrackCollectionName).Remove(bson.M{"webhookid": id})
	if err != nil {
		panic(err)
	}
}

//sendMessage sends message to webhook to dispaly
func sendMessage(token string, amountadded int, start time.Time) {
	session, err := mgo.Dial(dbh.HostURL)
	if err != nil {
		panic(err)
	}

	trackss := []Track{}
	ids := []string{}
	//Gets all tracks from db
	err = session.DB(db.Databasename).C(db.TrackCollectionName).Find(bson.M{}).Sort("TrackID").All(&trackss)
	if err != nil {
		panic(err)
	}
	//Gets the ids added since the last time
	for i := len(trackss) - 1; i >= len(trackss)-amountadded; i-- {
		ids = append(ids, strconv.Itoa(trackss[i].TrackID))
	}
	//Creates message to send
	message := map[string]interface{}{
		"content": fmt.Sprintf("Latest timestamp: %s, Newest %d added are %s [Processing time %s]",
			trackss[len(trackss)-1].TimeStamp,
			amountadded,
			ids,
			time.Since(start)),
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}
	//Sends message to webhook
	_, err = http.Post(token, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		println(err.Error())
	}

}

//UpdateWebHooks increases the actual value in all webhooks
func UpdateWebHooks() {
	session, err := mgo.Dial(dbh.HostURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	//Increses the value of actualvalue in the db
	_, err = session.DB(dbh.Databasename).C(dbh.TrackCollectionName).UpdateAll(bson.M{}, bson.M{"$inc": bson.M{"actualvalue": 1}})
	if err != nil {
		panic(err)
	}
	webhookmain()
}

//webhookmain check if it is time to send the message
func webhookmain() {
	start := time.Now()
	session, err := mgo.Dial(dbh.HostURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	webhooks := []WebHook{}
	//Gets all the webhooks
	err = session.DB(dbh.Databasename).C(dbh.TrackCollectionName).Find(bson.M{}).All(&webhooks)
	if err != nil {
		panic(err)
	}
	//Loops through all webhooks and send message if actualvalue > mintrigger
	for _, hook := range webhooks {
		if hook.ActualValue >= hook.MinTriggerValue {
			sendMessage(hook.URL, hook.ActualValue, start)
			//Sets actual value to 0
			_, err = session.DB(dbh.Databasename).C(dbh.TrackCollectionName).UpdateAll(bson.M{}, bson.M{"$set": bson.M{"actualvalue": 0}})
			if err != nil {
				panic(err)
			}
		}
	}
}
