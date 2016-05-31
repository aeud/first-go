package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func reverse(arr []string) []string {
    n := len(arr)
    rarr := make([]string, n)
    for i := 0; i < n; i++ {
        rarr[n-i-1] = arr[i]
    }
    return rarr
}

func main() {
    var n int

    io := bufio.NewReader(os.Stdin)
    fmt.Fscan(io, &n)

    a := make([]string, n)

    for i := 0; i < n; i++ {
        fmt.Fscan(io, &a[i])
    }

    r := reverse(a)

    fmt.Println(strings.Join(r, " "))

}
