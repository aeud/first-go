package main

import (
    "fmt"
)

type List []int

func mergeSort(l List, ch chan List) {
    length := len(l)
    if length > 1 {
        l1, l2 := l[:length / 2], l[length /2 :]

        newch := make(chan List)
        go mergeSort(l1, newch)
        go mergeSort(l2, newch)
        l1 = <-newch
        l2 = <-newch

        i, j := 0, 0

        newList := make(List, length)
        for i + j < length {
            if i >= len(l1) {
                newList[i + j] = l2[j]
                j++
            } else if j >= len(l2) {
                newList[i + j] = l1[i]
                i++
            } else if l1[i] < l2[j] {
                newList[i + j] = l1[i]
                i++
            } else {
                newList[i + j] = l2[j]
                j++
            }
        }

        ch <- newList
    } else {
        ch <- l
    }
}

func main() {
    var l List = List{5, 4, 3, 2, 1, 5, 4, 3, 4, 5, 4, 3, 4, 2, 4, 5, 6, 5}

    ch := make(chan List)
    go mergeSort(l, ch)
    l = <-ch

    fmt.Println(l)
}