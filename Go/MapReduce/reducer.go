package main

import (
    "fmt"
    "bufio"
    "os"
)

func reducer() {
    prevWord := ""
    curWord := ""
    count := 0
    reader := bufio.NewReader(os.Stdin)
    for  {
        line, _, err := reader.ReadLine()
        if err != nil {
            break
        }
        curWord = string(line)
        if prevWord != "" && curWord != prevWord {
            fmt.Printf("%v\t%v\n", prevWord, count)
            count = 0
        }
        prevWord = curWord
        count += 1
    }
    if curWord == prevWord {
        fmt.Printf("%v\t%v\n", prevWord, count)
    }
}

func main() {
    reducer()
}

