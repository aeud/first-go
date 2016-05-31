package main

import (
    "bufio"
    "fmt"
    "os"
)

const (
    n = 6
)

func getHourglass(arr [][]int, x, y int) int {
    h := arr[x][y] + arr[x][y+1] + arr[x][y+2]
    h += arr[x+1][y+1]
    h += arr[x+2][y] + arr[x+2][y+1] + arr[x+2][y+2]
    var nX, nY int
    if y < 3 {
        nX = x
        nY = y + 1
    } else {
        nX = x + 1
        nY = 0
    }
    if nX > 3 {
        return h
    } else {
        if h2 := getHourglass(arr, nX, nY); h >= h2 {
            return h
        } else {
            return h2
        }
    }
}

func main() {
    io := bufio.NewReader(os.Stdin)
    arr := make([][]int, n)
    for i := 0; i < n; i++ {
        arr[i] = make([]int, n)
        for j := 0; j < n; j++ {
            fmt.Fscan(io, &arr[i][j])
        }
    }
    fmt.Println(getHourglass(arr, 0, 0))
}
