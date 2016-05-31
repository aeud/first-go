package main

import (
    "bufio"
    "fmt"
    "os"
)

var N, Q int

var strings, queries []string

func main() {

    io := bufio.NewReader(os.Stdin)

    fmt.Fscan(io, &N)
    strings = make([]string, N)
    for i := 0; i < N; i++ {
        fmt.Fscan(io, &strings[i])
    }

    fmt.Fscan(io, &Q)
    queries = make([]string, Q)
    for i := 0; i < Q; i++ {
        fmt.Fscan(io, &queries[i])
    }

    for i := 0; i < len(queries); i++ {
        var s int
        for j := 0; j < len(strings); j++ {
            if queries[i] == strings[j] {
                s += 1
            }
        }
        fmt.Println(s)
    }
}
