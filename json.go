package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Column struct {
	Name string
	Type string
}

type Configuration []Column

func stringType(s string) interface{} {
	switch {
	case s == "int":
		return new(int)
	case s == "string":
		return new(string)
	case s == "float":
		return new(float32)
	case s == "bool":
		return new(bool)
	}
	return new(string)
}

func (c *Configuration) Scanable() []interface{} {
	l := len(*c)
	s := make([]interface{}, l)
	for i := 0; i < l; i++ {
		t := (*c)[i].Type
		fmt.Println(t)
		s[i] = stringType(t)
	}
	return s
}

func main() {
	m := map[string]int{
		"rsc": 3711,
		"r":   2138,
		"gri": 1908,
		"adg": 912,
	}
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
	a := Configuration{
		Column{"test", "int"},
		Column{"test2", "string"},
		Column{"test3", "float"},
	}
	fmt.Println(a.Scanable())
}
