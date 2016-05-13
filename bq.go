package main

import (
	db "database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	bigquery "google.golang.org/api/bigquery/v2"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	projectId = "luxola.com:luxola-analytics"
	dbConf    = "postgres://postgres:root@localhost:5432/mike?sslmode=disable"
)

type Order struct {
	account        string
	orderId        int64
	campaign       string
	source         string
	medium         string
	channel        string
	deviceCategory string
}

func (o Order) String() string {
	return fmt.Sprintf("Order %v (%v): %v", o.orderId, o.account, o.channel)
}

type OrderPrep struct {
	stmt *db.Stmt
	args []interface{}
}

func (op *OrderPrep) exec() {
	stmt := *op.stmt
	_, err := stmt.Exec(op.args...)
	if err != nil {
		log.Fatal(err)
	}
}

type Orders []*Order

var ch = make(chan *Orders)

func (o *Order) fill(f []*bigquery.TableCell) {
	oid, _ := strconv.ParseInt(f[1].V.(string), 10, 64)
	o.orderId = oid
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

func (os *Orders) toSQLValues() string {
	osv := *os
	a := make([]string, len(osv))
	for i := 0; i < len(osv); i++ {
		ov := *osv[i]
		a[i] = fmt.Sprintf("('%v', %v, '%v', '%v', '%v', '%v', '%v')", ov.account, ov.orderId, ov.campaign, ov.medium, ov.source, ov.channel, ov.deviceCategory)
	}
	return strings.Join(a, ",\n")
}

func (os *Orders) insertAll() {
	osv := *os
	db, _ := db.Open("postgres", dbConf)
	for i := 0; i < len(osv); i++ {
		o := *osv[i]
		stmt, _ := db.Prepare("insert into bq.ga_orders (account, local_order_id, campaign, medium, source, channel, device_category) values ($1, $2, $3, $4, $5, $6, $7)")
		args := []interface{}{o.account, o.orderId, o.campaign, o.medium, o.source, o.channel, o.deviceCategory}
		op := OrderPrep{stmt, args}
		op.exec()
		defer stmt.Close()
	}
	defer db.Close()
}

func executeQuery(sql string) {
	t0 := time.Now()
	db, err := db.Open("postgres", dbConf)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(sql)
	db.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Executed in %v\n", time.Now().Sub(t0))
	return
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
	sql := `
create schema if not exists bq;
drop table if exists bq.ga_orders cascade;
create table bq.ga_orders (
	account varchar(10),
	local_order_id integer,
	campaign text,
	medium text,
	source text,
	channel text,
	device_category text,
	constraint pk_bq_ga_orders primary key(account, local_order_id)
)
`
	executeQuery(sql)
	wg := new(sync.WaitGroup)
	go func() {
		for {
			os := <-ch
			fmt.Println("Working...")
			os.insertAll()
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
	query.Query = "select account, integer(orderId), campaign, source, medium, channel, deviceCategory from colors.ga_orders"
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
