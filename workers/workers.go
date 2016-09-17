package workers

import (
    "fmt"
)

func Test() {
    fmt.Println("Done.")
}

type Workable struct {
    Action     func()
    MaxWorkers int
}
