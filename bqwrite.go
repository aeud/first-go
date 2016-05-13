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
	tableId        = "test"
	BQModeNullable = "NULLABLE"
	BQTypeInteger  = "INTEGER"
	objectName     = "test-file.txt.gz"
	bucketName     = "lx-ga"
)

type Transaction struct {
	TransId   int
	AccountId int
}

func insertRows(transactions []*Transaction) {
	var dbConf string
	dbConfBytes, err := ioutil.ReadFile("./mssql.txt")
	if err != nil {
		log.Fatal(err)
	} else {
		dbConf = string(dbConfBytes)
	}

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

	tablesService := bigquery.NewTablesService(service)

	tableDeleteCall := tablesService.Delete(projectId, datasetId, tableId)
	if tableDeleteCall.Do() != nil {
		fmt.Println("Nothing to delete")
	}

	tableReference := new(bigquery.TableReference)
	tableReference.DatasetId = datasetId
	tableReference.ProjectId = projectId
	tableReference.TableId = tableId

	tableSchemaField := new(bigquery.TableFieldSchema)
	tableSchemaField.Mode = BQModeNullable
	tableSchemaField.Name = "testField"
	tableSchemaField.Type = BQTypeInteger

	tableSchema := new(bigquery.TableSchema)
	tableSchema.Fields = []*bigquery.TableFieldSchema{tableSchemaField}

	table := new(bigquery.Table)
	table.TableReference = tableReference
	table.Schema = tableSchema
	tableInsertCall := tablesService.Insert(projectId, datasetId, table)
	table, err = tableInsertCall.Do()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(table.Id)

	tableDataInsertService := bigquery.NewTabledataService(service)
	tableDataInsertAllRequest := new(bigquery.TableDataInsertAllRequest)

	tableDataInsertAllRequestRow := new(bigquery.TableDataInsertAllRequestRows)

	tableDataInsertAllRequestRows := make([]*bigquery.TableDataInsertAllRequestRows, len(transactions))

	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]
		v := make(map[string]bigquery.JsonValue)
		v["testField"] = transaction.TransId
		tableDataInsertAllRequestRow.Json = v
		tableDataInsertAllRequestRows = append(tableDataInsertAllRequestRows, tableDataInsertAllRequestRow)
	}

	tableDataInsertAllRequest.Rows = tableDataInsertAllRequestRows
	tableInsertAllCall := tableDataInsertService.InsertAll(projectId, datasetId, tableId, tableDataInsertAllRequest)

	_, err = tableInsertAllCall.Do()
	if err != nil {
		log.Fatal(err)
	}
}

func queryDb() []*Transaction {
	t0 := time.Now()
	//q := "select top 10 t.account_id as account_id from dbo.dimaccount t"
	q := `
select top 100 t.trans_id, t.account_id
from dbo.dimtrans t
    `
	// , t.account_id, t.store_id, t.total_qtys, t.total_sales, t.trans_time
	db, err := sql.Open("mssql", dbConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(q)
	rows, err := db.Query(q)
	defer rows.Close()
	defer db.Close()
	transactions := make([]*Transaction, 0)
	for rows.Next() {
		var transaction = new(Transaction)
		if err := rows.Scan(&transaction.TransId, &transaction.AccountId); err != nil {
			log.Fatal(err)
		}
		transactions = append(transactions, transaction)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Executed in %v\n", time.Now().Sub(t0))
	return transactions
}

func writeFile(transactions []*Transaction) {
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
	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]
		textJson, _ := json.Marshal(transaction)
		w.Write(textJson)
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
	writeFile(queryDb())
}
