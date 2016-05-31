package main

import (
    "bufio"
    "fmt"
    "os"
)

var N int

type answer struct {
    l int
    r int
}

var numbers []int
var answers []answer

var max int

func Left(i int) int {
    if i == 0 || i == N {
        return 0
    } else if numbers[i] < numbers[i-1] {
        return i - 1 + 1
    } else {
        for k := answers[i-1].l - 1; k >= 0; k-- {
            if numbers[k] > numbers[i] {
                return k + 1
            }
        }
    }
    return 0
}
func Right(i int) int {
    if i == 0 || i == N {
        return 0
    } else if numbers[i] > numbers[i-1] {
        if answers[i-1].r == 0 {
            return 0
        }
        for k := answers[i-1].r - 1; k < N; k++ {
            if numbers[k] > numbers[i] {
                return k + 1
            }
        }
    } else {
        for k := i + 1; k < N; k++ {
            if numbers[k] > numbers[i] {
                return k + 1
            }
        }
    }
    return 0
}

func IndexProduct(i int) int {
    var firstOp, secondOp (func(int) int)
    if i > N/2 {
        firstOp = Right
        secondOp = Left
    } else {
        firstOp = Left
        secondOp = Right
    }
    if l := firstOp(i); l > 0 {
        if r := secondOp(i); r > 0 {
            if i > N/2 {
                answers[i] = answer{r, l}
            } else {
                answers[i] = answer{l, r}
            }
            return l * r
        }
    }
    return 0
}

//func IndexProduct(i int) int {
//    l := Left(i)
//    r := Right(i)
//    answers[i] = answer{l, r}
//    return l * r
//}

func MaxProduct(i int) int {
    return N * (i + 1)
}

func MaxIndexProduct() int {
    for i := 0; i < N; i++ {
        if MaxProduct(i) > max {
            if p := IndexProduct(i); p > max {
                max = p
            }
        }
    }
    return max
}

func main() {
    io := bufio.NewReader(os.Stdin)

    fmt.Fscan(io, &N)
    numbers = make([]int, N)
    answers = make([]answer, N)
    for i := 0; i < N; i++ {
        fmt.Fscan(io, &numbers[i])
    }

    fmt.Println(MaxIndexProduct())
}
