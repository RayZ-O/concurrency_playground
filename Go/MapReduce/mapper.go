package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
)

func mapper(filename string) {
    f, err := os.Open(filename)
    if err != nil {
        fmt.Println("Failed to open input file")
        return
    }
    defer f.Close()
    // trim leading and trailing non alphanumeric
    rmPunc := func(r rune) bool {
        if r < 'A' || (r > 'Z' && r < 'a') || r > 'z'{
            return true;
        } else {
            return false;
        }
    }
    // read input file and print to stdout
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        words := strings.Fields(line)
        for _, w := range(words) {
            stdWord := strings.ToLower(strings.TrimFunc(w, rmPunc))
            fmt.Printf("%v\n", stdWord)
        }
    }
    if err = scanner.Err(); err != nil {
        fmt.Println("Failed to read input file")
        return
    }
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run mapper.go [input file name]")
        return
    }
    mapper(os.Args[1])
}
