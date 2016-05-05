package main

import (
    "database/sql"
    "log"
    "fmt"
    _ "github.com/lib/pq"
)

func main() {
    db, err := sql.Open("postgres", "postgres://postgres:root@localhost:5432/mike?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    rows, err := db.Query("SELECT count(*) AS \"count\" FROM omega.users")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        var count string
        if err := rows.Scan(&count); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%T", count)
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }

}