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
	http.ListenAndServe(":"+strconv.Itoa(8080), nil)
}
