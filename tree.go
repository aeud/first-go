package main

import (
    "golang.org/x/tour/tree"
    "fmt"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
    RecWalk(t, ch)
    close(ch)
}

func RecWalk(t *tree.Tree, ch chan int) {
    if t != nil {
        RecWalk(t.Left, ch)
        ch <- t.Value
        RecWalk(t.Right, ch)
    }
}

// Same determines whether the trees
// t1 and t2 contain the same values.
//func Same(t1, t2 *tree.Tree) bool

func main() {
    ch := make(chan int)
    go Walk(tree.New(1), ch)
    for v := range ch {
        fmt.Println(v)
    }
}
