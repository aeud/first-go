package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func main() {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write([]byte("hello, world\n"))
	w.Close()
	fmt.Println(b.String())
}
