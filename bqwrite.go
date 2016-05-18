package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	bigquery "google.golang.org/api/bigquery/v2"
	storage "google.golang.org/api/storage/v1"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

const (
	projectId      = "luxola.com:luxola-analytics"
	datasetId      = "go"
	tableId        = "toto"
	BQModeNullable = "NULLABLE"
	BQTypeInteger  = "INTEGER"
	objectName     = "transactions-test.json.gz"
	bucketName     = "lx-ga"
	//googleCredentialsPath = "/Users/adrien/.ssh/google.json"
	googleCredentialsPath = "/home/ubuntu/.ssh/google.json"
)

type Column struct {
	Name string
	Type string
}

type Schema []Column

func (schema *Schema) GetRow() Row {
	l := len(*schema)
	s := make([]interface{}, l)
	m := make(map[string]interface{})
	for i := 0; i < l; i++ {
		t := (*schema)[i].Type
		n := (*schema)[i].Name
		v := stringType(t)
		s[i] = v
		m[n] = v
	}
	return Row{m, s}
}

type Row struct {
	mapable  map[string]interface{}
	scanable []interface{}
}

func (r *Row) ToJson() []byte {
	s, err := json.Marshal(r.mapable)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

type Export struct {
	Schema Schema
	Query  string
}

var dbConf string

var export Export

func stringType(s string) interface{} {
	switch {
	case s == "int":
		return new(int)
	case s == "string":
		return new(string)
	case s == "float":
		return new(float32)
	case s == "bool":
		return new(bool)
	}
	return new(string)
}

func stringBQType(s string) string {
	switch {
	case s == "int":
		return "INTEGER"
	case s == "string":
		return "STRING"
	case s == "float":
		return "FLOAT"
	case s == "bool":
		return "BOOLEAN"
	}
	return "STRING"
}

func createTable() {

	data, err := ioutil.ReadFile(googleCredentialsPath)
	if err != nil {
		log.Fatal(err)
	}

	conf, err := google.JWTConfigFromJSON(data, bigquery.BigqueryScope)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(oauth2.NoContext)
	service, err := bigquery.New(client)
	if err != nil {
		log.Fatal(err)
	}

	jobsService := bigquery.NewJobsService(service)

	tableReference := bigquery.TableReference{
		DatasetId: datasetId,
		ProjectId: projectId,
		TableId:   tableId,
	}

	tableFields := make([]*bigquery.TableFieldSchema, len(export.Schema))

	for i := 0; i < len(export.Schema); i++ {
		c := export.Schema[i]
		tableSchemaField := new(bigquery.TableFieldSchema)
		tableSchemaField.Mode = BQModeNullable
		tableSchemaField.Name = c.Name
		tableSchemaField.Type = stringBQType(c.Type)
		tableFields[i] = tableSchemaField
	}

	tableSchema := new(bigquery.TableSchema)
	tableSchema.Fields = tableFields

	bs, _ := tableSchema.MarshalJSON()
	fmt.Println(string(bs))

	jobConfigurationLoad := bigquery.JobConfigurationLoad{
		DestinationTable: &tableReference,
		Schema:           tableSchema,
		SourceFormat:     "NEWLINE_DELIMITED_JSON",
		SourceUris:       []string{fmt.Sprintf("gs://%v/%v", bucketName, objectName)},
		WriteDisposition: "WRITE_TRUNCATE",
	}

	jogConfiguration := bigquery.JobConfiguration{
		Load: &jobConfigurationLoad,
	}

	job := bigquery.Job{
		Configuration: &jogConfiguration,
	}

	jobsInsertCall := jobsService.Insert(projectId, &job)
	insertJob, err := jobsInsertCall.Do()
	if err != nil {
		log.Fatal(err)
	}

	jobsGetCall := jobsService.Get(projectId, insertJob.JobReference.JobId)
	gotJob, err := jobsGetCall.Do()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gotJob.Status)

	for gotJob.Status.State != "DONE" {
		jobsGetCall = jobsService.Get(projectId, insertJob.JobReference.JobId)
		gotJob, err = jobsGetCall.Do()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(gotJob.Status)
		if gotJob.Status.ErrorResult != nil {
			fmt.Println(gotJob.Status.ErrorResult)
			for i := 0; i < len(gotJob.Status.Errors); i++ {
				fmt.Println(gotJob.Status.Errors[i])
			}
		}
		time.Sleep(time.Second)
	}
}

func queryDb() []*Row {
	t0 := time.Now()
	schema := Schema{
		Column{"AccountId", "int"},
		Column{"TransactionId", "int"},
		Column{"TransactionTime", "string"},
		Column{"SalesLCY", "float"},
		Column{"SalesSGD", "float"},
		Column{"ProductId", "int"},
		Column{"Country", "int"},
		Column{"StoreId", "int"},
	}
	q := `
select
    account_id as AccountId,
    trans_id as TransactionId,
    trans_time as TransactionTime,
    sales as SalesLCY,
    sales_basic_currency as SalesSGD,
    product_id as ProductId,
    subscription_country as SubscriptionCountry,
    store_id as StoreId
from
    dbo.FactTrans
	`
	export = Export{schema, q}
	db, err := sql.Open("mssql", dbConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(q)
	rows, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	defer db.Close()
	rs := make([]*Row, 0)
	for rows.Next() {
		r := export.Schema.GetRow()
		s := r.scanable
		if err := rows.Scan(s...); err != nil {
			log.Fatal(err)
		}
		rs = append(rs, &r)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Executed in %v\n", time.Now().Sub(t0))
	return rs
}

func writeFile(rs []*Row) {
	data, err := ioutil.ReadFile(googleCredentialsPath)
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
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		w.Write(r.ToJson())
		w.Write([]byte("\n"))
	}
	w.Close()

	file := strings.NewReader(b.String())

	if res, err := service.Objects.Insert(bucketName, object).Media(file).Do(); err == nil {
		fmt.Printf("Created object %v at location %v\n\n", res.Name, res.SelfLink)
	} else {
		log.Fatalf("Objects.Insert failed: %v", err)
	}
}

func main() {
	dbConfBytes, err := ioutil.ReadFile("./mssql.txt")
	if err != nil {
		log.Fatal(err)
	} else {
		dbConf = string(dbConfBytes)
	}
	fmt.Println(dbConf)
	writeFile(queryDb())
	createTable()
}
