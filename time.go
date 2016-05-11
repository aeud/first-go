package main

import (
	"fmt"
	"time"
)

func main() {
	t0 := time.Now()
	time.Sleep(time.Second)
	t1 := time.Now()
	fmt.Printf("%v", t1.Sub(t0))
}
