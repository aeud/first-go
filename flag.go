package main

import (
    "flag"
    "fmt"
)

var (
    port = flag.Int("port", 8000, "port to listen")
)

func main() {
    fmt.Printf("%v\n", *port)
}
