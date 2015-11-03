package main

import (
"fmt"
"math/rand"
"time"
"os"
"strconv"
)

type Message struct {
    msgType string
    content int
}

type Coordinator  struct {
    cohorts []Cohort
    msgQueue chan Message
}

type Cohort struct {
    id int
    state int
    undoLog int
    coordinator Coordinator
    msgQueue chan Message
}

func (c Coordinator) Start() {
    repCount := 0
    ackCount := 0
    success := true
    c.SendToAll(Message{"query to commit", 5})
    go func() {
        for {
            m := <-c.msgQueue
            if m.msgType == "agreement" {
                fmt.Printf("[Receive] Coordinator receives %v agreements\n", repCount)
                repCount += 1
                c.CheckReply(repCount, success)
            } else if m.msgType == "acknowledgment" {
                ackCount += 1
                fmt.Printf("[Receive] Coordinator receives %v acknowledgment\n", ackCount)
                if ackCount >= len(c.cohorts) {
                    if success {
                        fmt.Println("[Info] Commit success!\n")
                    } else {
                        fmt.Println("[Info] Commit failed!\n")
                    }
                    ackCount, repCount, success = 0, 0, true
                    time.Sleep(2000 * time.Millisecond)
                    random := rand.Intn(19) + 1
                    c.SendToAll(Message{"query to commit", random})
                }
            } else if m.msgType == "abort" {
                fmt.Printf("[Receive] Coordinator receives abort\n")
                success = false
                repCount += 1
                c.CheckReply(repCount, success)
            }
        }
    } ()
}

func (c Coordinator) CheckReply(repCount int, success bool) {
    if repCount >= len(c.cohorts) {
        if success {
            fmt.Printf("[Send] Coordinator send commit to all cohorts\n")
            c.SendToAll(Message{"commit", 0})
        } else {
            fmt.Printf("[Send] Coordinator send rollback to all cohorts\n")
            c.SendToAll(Message{"rollback", 0})
        }
    }
}

func (c Coordinator) SendToAll(msg Message) {
    for _, cohort := range c.cohorts {
        cohort.msgQueue <- msg
    }
}

func (c *Cohort) Start() {
    fmt.Printf("[START] Cohort %d starts with state: %d\n", c.id, c.state)
    go func() {
        for {
            m := <-c.msgQueue
            if m.msgType == "query to commit" {
                success := rand.Intn(5) >= 1
                fmt.Printf("[Receive] Cohort %v receives commit request: +%v, current state: %v\n", c.id, m.content, c.state)
                if success {
                    c.undoLog = m.content
                    c.state += m.content
                    fmt.Printf("[Info] Cohort %v succeed, new state: %v\n", c.id, c.state)
                    fmt.Printf("[Send] Cohort %v sends agreement to coordinator\n", c.id)
                    c.coordinator.msgQueue <- Message{"agreement", 0}
                } else {
                    c.undoLog = -1
                    fmt.Printf("[Send] Cohort %v failed, sends abort to coordinator\n", c.id)
                    c.coordinator.msgQueue <- Message{"abort", 0}
                }
            } else if m.msgType == "rollback" {
                if c.undoLog > 0 {
                    c.state -= c.undoLog
                }
                fmt.Printf("[Receive] Cohort %v receives rollback, state becomes: %v\n", c.id, c.state)
                c.coordinator.msgQueue <- Message{"acknowledgment", 0}
            } else if m.msgType == "commit" {
                fmt.Printf("[Receive] Cohort %v receives commit\n", c.id)
                c.coordinator.msgQueue <- Message{"acknowledgment", 0}
            }
        }
    } ()
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run commit.go [num of cohort]")
        return
    }
    // get number of cohort from input
    numOfCohorts, err := strconv.Atoi(os.Args[1])
    if err != nil  {
        fmt.Println("Malform input argument")
        return
    }
    coordinator := Coordinator{msgQueue:make(chan Message)}
    cohorts := make([]Cohort, numOfCohorts)
    for i := 0; i < numOfCohorts; i++ {
        cohorts[i] = Cohort{id:i, state:i, undoLog:0, coordinator:coordinator, msgQueue:make(chan Message)}
        cohorts[i].Start()
    }
    coordinator.cohorts = cohorts
    coordinator.Start()
    terminate := make(chan bool)
    <-terminate
}

