package main

import (
"fmt"
"math/rand"
"time"
"os"
"strconv"
)

type Message struct {
    vector []int
    sender *Clock
}

type Clock struct {
    id int
    vector []int
    neighboor []Clock
    numClocks int
    msgQueue chan Message
}

func (c Clock) Start() {
    fmt.Printf("[START] Clock %d starts with vector: %v\n", c.id, c.vector)
    go c.Tick()
    go c.Send()
    for {
        m := <-c.msgQueue
        c.ProcessMessage(m)
    }
}

func (c *Clock) Tick() {
    for {
        fmt.Printf("[Tick] Clock %d vector: %v\n", c.id, c.vector)
        c.vector[c.id]++
        time.Sleep(1000 * time.Millisecond)
    }
}

func (c *Clock) Send() {
    for {
        time.Sleep(2000 * time.Millisecond)
        c.SendMessage()
    }
}

func (c Clock) SendMessage() {
    random := c.id
    for random == c.id {
         random = rand.Intn(c.numClocks)
     }
    c.neighboor[random].msgQueue <- Message{c.vector, &c}
}

func (c *Clock) ProcessMessage(msg Message) {
    // fmt.Printf("[LOG] Clock %d receive message from clock %d\n", c.id, msg.sender.id)
    c.vector[c.id]++
    updated := false
    for i, time := range msg.vector {
        if time > c.vector[i] {
            c.vector[i] = time
            updated = true
        }
    }
    if updated {
        fmt.Printf("[SYNCHRONIZE] Clock %d receives from clock %d, times vector become %v\n", c.id, msg.sender.id, c.vector)
    }
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run lamport.go [num of clocks]")
        return
    }
    // get number of clocks from input
    numOfClocks, err := strconv.Atoi(os.Args[1])
    if err != nil  {
        fmt.Println("Malform input argument")
        return
    }
    // inital clocks
    clocks := make([]Clock, numOfClocks)
    for i := 0; i < numOfClocks; i++ {
        clocks[i] = Clock{id:i, vector:make([]int, numOfClocks), numClocks:numOfClocks, msgQueue:make(chan Message)}
    }

    for i := 0; i < numOfClocks; i++ {
        clocks[i].neighboor = clocks[:];
        go clocks[i].Start()
    }
    terminate := make(chan bool)
    <-terminate
    // time.Sleep(20 * time.Second)
}
