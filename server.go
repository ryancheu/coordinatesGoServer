package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type QualityType struct {
	Play    string
	Token   string
	Connect string `"json:"connect"`
	Bitrate float64
}

type Coordinate struct {
	Lat  string
	Long string
}

//Snippet for json responses taken from
//http://nesv.blogspot.com/2012/09/super-easy-json-http-responses-in-go.html
type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

const startPort = 6000
const maxPorts = 1000
const queryPort = 8080
const bitrate = 128
const refreshPeriod = 600 //In minutes

//This should really use a list but I don't know how to do that yet in Golang
var coords map[int]Coordinate
var curIndex = startPort

func coordinateAddHandler(rw http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	lat := query.Get("lat")
	long := query.Get("long")

	coord := Coordinate{lat, long}
	fmt.Println("New coordinate: ", coord)
	coords[curIndex] = coord
	curIndex++
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
}

func coordinateListHandler(rw http.ResponseWriter, r *http.Request) {

	//This is some bad code that transfers the map into a fixed length array
	//This should really be refactored to use the go equivalent of vector/arraylist
	//so this step isn't needed
	var coordArray []Coordinate = make([]Coordinate, len(coords)) 
	i := 0
	for coordinate := range coords { 
		coordArray[i] = coords[coordinate]
		i++
	}

	jsonList, err := json.Marshal(coordArray)
	if err != nil {
		fmt.Println(err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(rw, string(jsonList))
}

func coordinateClearHandler(rw http.ResponseWriter, r *http.Request) {
	coords = make(map[int]Coordinate)
}

func main() {
	coords = make(map[int]Coordinate)
	http.HandleFunc("/coordinate/new", coordinateAddHandler)
	http.HandleFunc("/coordinate/list", coordinateListHandler)
	http.HandleFunc("/coordinate/clear", coordinateClearHandler)
	http.ListenAndServe(":"+strconv.Itoa(queryPort), nil)
}
