package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	v := new(interface{})
	j := `
{
	"name": "toto",
	"obj": {
		"key": "value"
	}
}
	`
	if err := json.Unmarshal([]byte(j), v); err != nil {
		fmt.Println(err)
	}
	fmt.Println(*v)
	fmt.Println((*v).(map[string]interface{})["name"])

	bs, err := json.Marshal(v)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bs))

	time.Sleep(1000 * time.Second)
}
