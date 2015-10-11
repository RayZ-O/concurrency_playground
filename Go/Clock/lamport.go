package main

import (
"fmt"
"math/rand"
"time"
"os"
"strconv"
)

type Message struct {
    timestamps int
    sender *Clock
}

type Clock struct {
    id int
    timestamps int
    neighboor []Clock
    numClocks int
    msgQueue chan Message
}

func (c Clock) Start() {
    fmt.Printf("[START] Clock %d starts with timestamps: %d\n", c.id, c.timestamps)
    go c.Tick()
    go c.Send()
    for {
        m := <-c.msgQueue
        c.ProcessMessage(m)
    }
}

func (c *Clock) Tick() {
    for {
        fmt.Printf("[Tick] Clock %d timestamp: %d\n", c.id, c.timestamps)
        c.timestamps++
        time.Sleep(500 * time.Millisecond)
    }
}

func (c *Clock) Send() {
    for {
        time.Sleep(1000 * time.Millisecond)
        c.SendMessage()
    }
}

func (c Clock) SendMessage() {
    random := c.id
    for random == c.id {
         random = rand.Intn(c.numClocks)
     }
    c.neighboor[random].msgQueue <- Message{c.timestamps, &c}
}

func (c *Clock) ProcessMessage(msg Message) {
    // fmt.Printf("[LOG] Clock %d receive message from clock %d\n", c.id, msg.sender.id)
    if msg.timestamps > c.timestamps {
        c.timestamps = msg.timestamps
        fmt.Printf("[SYNCHRONIZE] Clock %d receives from clock %d, timestamps become %d\n", c.id, msg.sender.id, c.timestamps)
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
        clocks[i] = Clock{id:i, timestamps:i, numClocks:numOfClocks, msgQueue:make(chan Message)}
    }

    for i := 0; i < numOfClocks; i++ {
        clocks[i].neighboor = clocks[:];
        go clocks[i].Start()
    }
    time.Sleep(20 * time.Second)
}
