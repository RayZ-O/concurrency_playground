 package main

import (
    "fmt"
    "time"
    "math/rand"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Message struct {
    name string
    message string
}

type Publisher struct {
    name string
    pubCh chan Message
}

func (p Publisher) Start(t time.Duration) {
    ticker := time.NewTicker(t)
    go func() {
        for {
            <-ticker.C
            msg := RandString(10)
            p.Pub(msg)
        }
    }()
}

func (p Publisher) Pub(message string) {
    p.pubCh <- Message{p.name, message}
}

type PubSub struct {
    ch chan Message
    subs map[string]map[chan Message]bool
}

func (ps *PubSub) Sub(name string, subChs []chan Message) {
    if ps.subs[name] == nil {
        ps.subs[name] = make(map[chan Message]bool)
    }
    for  _, c := range subChs {
        ps.subs[name][c] = true
    }
}

func (ps *PubSub) UnSub(name string, subCh chan Message) {
    delete(ps.subs[name], subCh)
}

func (ps PubSub) Start() {
    for {
        msg := <- ps.ch
        if chs, ok := ps.subs[msg.name]; ok {
            for chs, _ := range chs {
                chs <- msg
            }
        }
    }
}

func RandString(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func main() {
    fmt.Println("start")
    pubCh := make(chan Message)
    pubsub := PubSub{pubCh, make(map[string]map[chan Message]bool)}

    p1 := Publisher{"sport", pubCh}
    p2 := Publisher{"hacker", pubCh}
    p3 := Publisher{"travell", pubCh}

    var subs [3]chan Message
    for i := range subs {
       subs[i] = make(chan Message)
    }

    pubsub.Sub("sport", subs[:])
    pubsub.Sub("hacker", subs[:2])
    pubsub.Sub("travell", subs[2:])

    go pubsub.Start()
    p1.Start(2 * time.Second)
    p2.Start(5 * time.Second)
    p3.Start(7 * time.Second)
    for {
        select {
        case msg0 := <- subs[0]:
            fmt.Printf("Subscriber 1 receive from %v\n", msg0.name)
        case msg1 := <- subs[1]:
            fmt.Printf("Subscriber 2 receive from %v\n", msg1.name)
        case msg2 := <- subs[2]:
            fmt.Printf("Subscriber 3 receive from %v\n", msg2.name)
        }
    }

}
