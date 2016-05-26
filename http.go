package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type responeStruct struct {
	Test string
}

// First Test Handle for the Http Web Server
func TestHandle(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t responeStruct
	err := decoder.Decode(&t)
	if err != nil {
		io.WriteString(w, fmt.Sprintf("%v", "Error"))
	} else {
		io.WriteString(w, fmt.Sprintf("%v", t.Test))
	}
}

func main() {
	http.HandleFunc("/", TestHandle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
