package main

import (
	"fmt"
	"reflect"
	"time"
)

type Thread chan int

type Threads []Thread

func initThreads(n int) Threads {
	t := make(Threads, n)
	for i := 0; i < len(t); i++ {
		t[i] = make(Thread, 0)
	}
	return t
}

var threads = initThreads(20)

func longFunction() {
	fmt.Println("Long function")
	time.Sleep(time.Second)
}

func push() {
	cases := make([]reflect.SelectCase, len(threads))
	for i := 0; i < len(threads); i++ {
		cases[i] = reflect.SelectCase{reflect.SelectSend, reflect.ValueOf(threads[i]), reflect.ValueOf(i)}
	}
	reflect.Select(cases)
}

func createWorkers() {
	for i := 0; i < len(threads); i++ {
		index := i
		ch := threads[i]
		go func() {
			for {
				select {
				case <-ch:
					fmt.Printf("Consuming %v\n", index)
					longFunction()
					fmt.Printf("Consumed %v\n", index)
				}
			}
		}()

	}
}

func main() {
	createWorkers()

	push()
	push()
	push()
	push()
	push()

	var input string
	fmt.Scanln(&input)
}
