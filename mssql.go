package main

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	var dbConf string
	dbConfBytes, err := ioutil.ReadFile("./mssql.txt")
	if err != nil {
		log.Fatal(err)
	} else {
		dbConf = string(dbConfBytes)
	}
	t0 := time.Now()
	//q := "select top 10 t.account_id as account_id from dbo.dimaccount t"
	q := `
select top 100 t.trans_id, t.account_id, t.store_id, t.total_qtys, t.total_sales, t.trans_time
from dbo.dimtrans t
    `
	db, err := sql.Open("mssql", dbConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(q)
	rows, err := db.Query(q)
	defer rows.Close()
	defer db.Close()
	for rows.Next() {
		var trans_id int
		var account_id int
		var store_id int
		var total_qtys int
		var total_sales float64
		var trans_time string
		if err := rows.Scan(&trans_id, &account_id, &store_id, &total_qtys, &total_sales, &trans_time); err != nil {
			log.Fatal(err)
		}
		fmt.Println(trans_time)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Executed in %v\n", time.Now().Sub(t0))
}
