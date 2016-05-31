package main

import (
    "bufio"
    "fmt"
    "os"
)

const (
    S = 3
)

var N, Q int
var lastAns int
var seq [][]int

func getIndex(x, y int) int {
    return (x ^ lastAns) % N
}

func queryType1(x, y int) {
    i := getIndex(x, y)
    seq[i] = append(seq[i], y)
}

func queryType2(x, y int) {
    i := getIndex(x, y)
    j := y % len(seq[i])
    lastAns = seq[i][j]
    fmt.Println(lastAns)
}

func main() {

    io := bufio.NewReader(os.Stdin)
    fmt.Fscan(io, &N)
    fmt.Fscan(io, &Q)

    arr := make([][]int, Q)
    for i := 0; i < Q; i++ {
        arr[i] = make([]int, S)
        for j := 0; j < S; j++ {
            fmt.Fscan(io, &arr[i][j])
        }
    }

    seq = make([][]int, N)
    for i := 0; i < N; i++ {
        seq[i] = make([]int, 0)
    }

    for i := 0; i < Q; i++ {
        if arr[i][0] == 1 {
            queryType1(arr[i][1], arr[i][2])
        } else {
            queryType2(arr[i][1], arr[i][2])
        }
    }
}
