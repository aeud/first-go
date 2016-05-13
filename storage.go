package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
	"io/ioutil"
	"log"
	"strings"
)

type Test struct {
	Id int
}

const (
	bucketName = "lx-ga"
	projectId  = "luxola.com:luxola-analytics"
	objectName = "test-file.txt.gz"
)

func writeFile([]int) {
	data, err := ioutil.ReadFile("/Users/adrien/.ssh/google.json")
	if err != nil {
		log.Fatal(err)
	}

	conf, err := google.JWTConfigFromJSON(data, storage.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(oauth2.NoContext)
	service, err := storage.New(client)
	if err != nil {
		log.Fatalf("Unable to create storage service: %v", err)
	}

	object := &storage.Object{Name: objectName}

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	textJson, err := json.Marshal(Test{1})
	w.Write(textJson)
	w.Write([]byte("\n"))
	w.Close()

	file := strings.NewReader(b.String())

	if res, err := service.Objects.Insert(bucketName, object).Media(file).Do(); err == nil {
		fmt.Printf("Created object %v at location %v\n\n", res.Name, res.SelfLink)
	} else {
		log.Fatalf("Objects.Insert failed: %v", err)
	}
}

func main() {

}
