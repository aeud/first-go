package main

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	bigquery "google.golang.org/api/bigquery/v2"
	"io/ioutil"
	"log"
	"strconv"
	"sync"
	"time"
)

const (
	projectId = "luxola.com:luxola-analytics"
)

type Order struct {
	account        string
	orderdId       int64
	campaign       string
	source         string
	medium         string
	channel        string
	deviceCategory string
}

func (o Order) String() string {
	return fmt.Sprintf("Order %v (%v): %v", o.orderdId, o.account, o.channel)
}

type Orders []*Order

var ch = make(chan *Orders)

func (o *Order) fill(f []*bigquery.TableCell) {
	oid, _ := strconv.ParseInt(f[1].V.(string), 10, 64)
	o.orderdId = oid
	o.account = f[0].V.(string)
	if f[2].V != nil {
		o.campaign = f[2].V.(string)
	}
	if f[3].V != nil {
		o.source = f[3].V.(string)
	}
	if f[4].V != nil {
		o.medium = f[4].V.(string)
	}
	if f[5].V != nil {
		o.channel = f[5].V.(string)
	}
	if f[6].V != nil {
		o.deviceCategory = f[6].V.(string)
	}
}

func (os *Orders) fill(rows []*bigquery.TableRow) {
	for i := 0; i < len(rows); i++ {
		el := new(Order)
		el.fill(rows[i].F)
		(*os)[i] = el
	}
}

func checkJob(c *bigquery.JobsGetQueryResultsCall, j *bigquery.Job, wg *sync.WaitGroup) {
	q, err := c.Do()
	if err != nil {
		log.Fatal(err)
	}
	if q.JobComplete {
		rows := q.Rows
		os := make(Orders, len(rows))
		os.fill(rows)
		ch <- &os
		if q.PageToken != "" {
			wg.Add(1)
			checkJob(c.PageToken(q.PageToken), j, wg)
		}
	} else {
		fmt.Println("Try again")
		time.Sleep(time.Second)
		checkJob(c, j, wg)
	}
}

func main() {
	wg := new(sync.WaitGroup)
	go func() {
		for {
			os := <-ch
			// Do something with the orders
			fmt.Println(len(*os))
			wg.Done()
		}
	}()

	data, err := ioutil.ReadFile("/Users/adrien/.ssh/google.json")
	if err != nil {
		log.Fatal(err)
	}

	conf, err := google.JWTConfigFromJSON(data, bigquery.BigqueryScope)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(oauth2.NoContext)
	service, _ := bigquery.New(client)

	//  datasetsService := bigquery.NewDatasetsService(service)
	//  datasetsServiceList := datasetsService.List(projectId)
	//  list, err := datasetsServiceList.Do()
	//  if err != nil {
	//      log.Fatal(err)
	//  }
	//  datasets := list.Datasets
	//  for i := 0; i < len(datasets); i++ {
	//      fmt.Println(datasets[i].Id)
	//  }

	//  queryRequest := bigquery.QueryRequest{Query: "select account, integer(orderId), campaign, source, medium, channel from colors.ga_orders limit 10"}
	jobsService := bigquery.NewJobsService(service)
	//  resp, err := queryService.Query(projectId, &queryRequest).Do()
	//  if err != nil {
	//      log.Fatal(err)
	//  }
	//  rows := resp.Rows
	//  for i := 0; i < len(rows); i++ {
	//      f := rows[i].F
	//      o := new(Order)
	//      o.fill(f)
	//      fmt.Println(o)
	//  }
	job := new(bigquery.Job)
	configuration := new(bigquery.JobConfiguration)
	query := new(bigquery.JobConfigurationQuery)
	query.Query = "select account, integer(orderId), campaign, source, medium, channel, deviceCategory from colors.ga_orders limit 60000"
	configuration.Query = query
	job.Configuration = configuration
	j, err := jobsService.Insert(projectId, job).Do()
	if err != nil {
		log.Fatal(err)
	}
	wg.Add(1)
	checkJob(jobsService.GetQueryResults(projectId, j.JobReference.JobId), j, wg)
	wg.Wait()
}
