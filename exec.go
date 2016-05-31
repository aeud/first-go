package main

import (
    "bytes"
    "fmt"
    "log"
    "os/exec"
    "regexp"
    "strings"
)

func runCommand(cmdName string, args ...string) string {
    cmd := exec.Command(cmdName, args...)
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        log.Fatal(err)
    }
    if err != nil {
        log.Fatal(err)
    }
    if err := cmd.Start(); err != nil {
        log.Fatal(err)
    }
    buf := new(bytes.Buffer)
    buf.ReadFrom(stdout)
    s := buf.String()
    return s
}

func main() {
    runCommand("rm", []string{"-rf", "selenium"}...)
    runCommand("git", []string{"clone", "git@github.com:luxola/selenium.git"}...)
    s := runCommand("go", []string{"test", "selenium/main_test.go"}...)
    runCommand("rm", []string{"-rf", "selenium"}...)
    ss := strings.Split(s, "\n")
    res := regexp.MustCompile("^\\w+").FindString(ss[len(ss)-2])
    fmt.Println(res)
    if res == "ok" {
        runCommand("go", []string{"get", "-u", "github.com/luxola/selenium"}...)
        runCommand("go", []string{"install", "-u", "github.com/luxola/selenium"}...)
    } else {
        fmt.Println("Error")
    }
}
