package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "regexp"
    "strconv"
)

func main() {
    var n int
    _, err := fmt.Scanf("%d", &n)
    if err != nil {
        log.Fatal("Cannot read n")
    }
    reader := bufio.NewReader(os.Stdin)
    rows := make([][]int, n)
    for i := 0; i < n; i++ {
        text, _ := reader.ReadString('\n')
        r := regexp.MustCompile("\\-?\\d+")
        numbers := r.FindAllString(text, n)
        row := make([]int, n)
        for h := 0; h < n; h++ {
            j, err := strconv.ParseInt(numbers[h], 10, 64)
            if err != nil {
                log.Fatal("Cannot read a number")
            }
            row[h] = int(j)
        }
        rows[i] = row
    }
    var d1, d2 int
    for i := 0; i < n; i++ {
        d1 += rows[i][i]
        d2 += rows[i][n-i-1]
    }
    if d1 > d2 {
        fmt.Println(d1 - d2)
    } else {
        fmt.Println(d2 - d1)
    }
}
