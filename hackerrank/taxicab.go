package main

import (
    "bufio"
    "fmt"
    "math"
    "os"
)

var N, H, V int

var js [][]int

var sum int

func main() {
    io := bufio.NewReader(os.Stdin)

    fmt.Fscan(io, &N)
    fmt.Fscan(io, &H)
    fmt.Fscan(io, &V)
    js = make([][]int, N)
    for i := 0; i < N; i++ {
        js[i] = make([]int, 2)
        fmt.Fscan(io, &js[i][0])
        fmt.Fscan(io, &js[i][1])
    }
    for i := 0; i < N-1; i++ {
        var o, t int
        fmt.Fscan(io, &o)
        fmt.Fscan(io, &t)
        if int(math.Abs(float64(js[o-1][0]-js[t-1][0]))) > V || int(math.Abs(float64(js[o-1][1]-js[t-1][1]))) > H {
            sum++
        }
    }

    fmt.Println(js)
    fmt.Println(sum)
}
