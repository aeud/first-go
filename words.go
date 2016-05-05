package main

import (
    "golang.org/x/tour/wc"
    "strings"
)

func WordCount(s string) map[string]int {
    words := strings.Fields(s)
    m := make(map[string]int)
    for i := 0; i < len(words); i++ {
        m[words[i]]++
    }
    return m
}

func main() {
    wc.Test(WordCount)
}
