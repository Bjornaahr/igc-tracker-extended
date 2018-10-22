# go-igc-tracker

- ```GET /api ```API metainfo
- ```GET /api/igc``` Shows all the track ids
- ```POST /api/igc``` Takes ```{"url":"<url>"}``` and saves the IGC info in memory, returns assigned ID 
- ```GET /api/igc/<id>``` Returns data about the track
- ```GET /api/igc/<id>/<field>``` Returns data of type ```<field>```


## Faults
- Can ```POST``` the same track several times
- ```GET /api/igc``` can appear in non ascending order
- ```GET /api ``` Uptime might not be 100% correct
- No tests

###### This was the first assignment for Cloud technologies NTNU Gjøvik 2018
