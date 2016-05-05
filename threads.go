package main

import (
    "fmt"
    "time"
)

func main() {
    go say("let's go", 3*time.Second)
    go say("let's go 2", 2*time.Second)
    go say("let's go 1", 1*time.Second)
    time.Sleep(4 * time.Second)
}

func say(text string, delay time.Duration) {
    time.Sleep(delay)
    fmt.Println(text)
}