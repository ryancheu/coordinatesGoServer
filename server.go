package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Coordinate struct {
	Lat  string
	Long string
	Date int64
}

//This should really use a list but I don't know how to do that yet in Golang
var coords map[int]Coordinate
var curIndex = 0

var curCommand string

func coordinateAddHandler(rw http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	lat := query.Get("lat")
	long := query.Get("long")

	coord := Coordinate{lat, long, time.Now().Unix()}
	fmt.Println("New coordinate: ", coord)
	coords[curIndex] = coord
	curIndex++
}

func coordinateListHandler(rw http.ResponseWriter, r *http.Request) {

	//This is some bad code that transfers the map into a fixed length array
	//This should really be refactored to use the go equivalent of vector/arraylist
	//so this step isn't needed
	var coordArray []Coordinate = make([]Coordinate, len(coords))
	for i := 0; i < len(coords); i++ {
		coordArray[len(coords)-1-i] = coords[i]
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

func commandSetHandler(rw http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	command := query.Get("command")
	curCommand = command
}

func commandViewHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(rw, curCommand)
}

func main() {
	coords = make(map[int]Coordinate)
	http.HandleFunc("/coordinate/new", coordinateAddHandler)
	http.HandleFunc("/coordinate/list", coordinateListHandler)
	http.HandleFunc("/coordinate/clear", coordinateClearHandler)
	http.HandleFunc("/command/set", commandSetHandler)
	http.HandleFunc("/command/view", commandViewHandler)
	http.ListenAndServe(":"+strconv.Itoa(8080), nil)
}
