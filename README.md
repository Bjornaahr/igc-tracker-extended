[![Build Status](https://travis-ci.com/Bjornaahr/igc-tracker-extended.svg?branch=master)](https://travis-ci.com/Bjornaahr/igc-tracker-extended)

# go-igc-tracker

- ```GET /api/ ```API metainfo
- ```GET /api/igc/``` Shows all the track ids
- ```POST /api/igc/``` Takes ```{"url":"<url>"}``` and saves the IGC info in memory, returns assigned ID 
- ```GET /api/igc/<id>/``` Returns data about the track
- ```GET /api/igc/<id>/<field>/``` Returns data of type ```<field>```
- ```GET /api/ticker/latest/``` Latest timestamp added to DB
- ```GET /api/ticker/``` Ticker Info
- ```GET /api/ticker/<timestamp>/``` Ticker info by one timestamp
- ```POST /api/webhook/new_track/``` Takes ```{"webhookURL": "<url>", "minTriggerValue" "<number>"``` and adds the webhook to the database
- ```GET /api/webhook/new_track/<id>/``` Return the webhook with that id
- ```DELETE /api/webhook/new_track/<id>/``` Deletes webhook with that id
- ```GET /admin/api/tracks_count/``` Returns count of tracks in DB
- ```DELETE /admin/api/tracks/``` Deletes all tracks from DB


## Faults
- Can ```POST``` the same track several times
- ```GET /api/ ``` Uptime might not be 100% correct
- Weak testing 27% coverage

### Heroku link
https://igctrackerextended.herokuapp.com/

###### This was the second assignment for Cloud technologies NTNU Gjøvik 2018
