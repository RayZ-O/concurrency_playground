package main

import (
"fmt"
"math/rand"
"time"
)

const numOfPeers = 4000 // greater than 1
const threshold = 10

type Message struct {
    msgType int
    content string
    sender *Peer
}

type Peer struct {
    id int
    neighboor []Peer
    count int
    state string
    msgQueue chan Message
}

func (p Peer) Start(convergent chan bool) {
    msg := <-p.msgQueue
    // fmt.Printf("[KNOWS] Peer %d knows the roumor\n", p.id)
    p.count++
    p.state = msg.content
    p.ProcessMessage(msg)
    p.Propagate(convergent, msg)
}

func (p Peer) Propagate(convergent chan bool, msg Message) {
    ticker := time.NewTicker(10 * time.Millisecond)
    cancelTick := make(chan bool)
    go func() {
        for {
            select {
            case <-ticker.C:
                if len(p.neighboor) > 0 {
                    p.SendMessage(p.state)
                } else {
                    cancelTick <- true
                }
            case <-cancelTick:
                ticker.Stop()
                return
            }
        }
    }()

    for {
        m := <-p.msgQueue
        if m.msgType == 0 {
            p.count++
            p.ProcessMessage(m)
            if p.count == threshold {
                fmt.Printf("[CONVERGE] Peer %d converged\n", p.id)
                convergent <- true
            }
            if p.count >= threshold {
                p.NotifyTerminate(m)
            }
        } else {
            for i, n := range p.neighboor {
                if m.sender.id == n.id {
                    p.neighboor = append(p.neighboor[:i], p.neighboor[i+1:]...)
                    break
                }
            }
            if len(p.neighboor) == 0 && p.count < threshold {
                p.count = threshold + 1
                fmt.Printf("[CONVERGE] Peer %d is isolated\n", p.id)
                convergent <- true
            }
        }
    }

}

func (p Peer) SendMessage(content string) {
    i := rand.Intn(len(p.neighboor))
    // fmt.Printf("Peer %d send to Peer %d\n" , p.id, p.neighboor[i].id)
    p.neighboor[i].msgQueue <- Message{0, content, &p}
}

func (p Peer) ProcessMessage(msg Message) {
    // fmt.Printf("Peer %d receive %s from Peer %d\n" , p.id, msg.content, msg.sender.id)
}

func (p Peer) NotifyTerminate(msg Message) {
    // for _, n := range p.neighboor {
    //     n.msgQueue <- Message{2, "", &p}
    // }
    msg.sender.msgQueue <- Message{2, "", &p}
    // fmt.Printf("Peer %d notify terminate to %d\n", p.id, msg.sender.id)
}

func BuildFullNetwork(peers []Peer) {
    for i :=1; i < numOfPeers + 1; i++ {
        for j :=1; j < numOfPeers + 1; j++ {
            if j != i {
                peers[i].neighboor = append(peers[i].neighboor, peers[j])
            }
        }
    }
}

func BuildLineNetwork(peers []Peer) {
    for i :=2; i < numOfPeers; i++ {
        peers[i].neighboor = append(peers[i].neighboor, peers[i - 1])
        peers[i].neighboor = append(peers[i].neighboor, peers[i + 1])
    }
    peers[1].neighboor = append(peers[1].neighboor, peers[2])
    peers[numOfPeers].neighboor = append(peers[numOfPeers].neighboor, peers[numOfPeers - 1])
}

func main() {
    // create peers
    var peers [numOfPeers + 1]Peer
    for i := 0; i < numOfPeers + 1; i++ {
        peers[i] = Peer{id:i, count:0, msgQueue:make(chan Message)}
    }
    // build topology
    BuildFullNetwork(peers[:])
    // BuildLineNetwork(peers[:])
    for i := 1; i < numOfPeers + 1; i++ {
        peers[0].neighboor = append(peers[0].neighboor, peers[i])
    }

    convergent := make(chan bool)
    rand.Seed(time.Now().UTC().UnixNano())
    for i := 1; i < numOfPeers + 1; i++ {
        go peers[i].Start(convergent)
    }
    peers[0].SendMessage("Hello, I am main process")
    for i := 0; i < numOfPeers; i++ {
        <-convergent
    }
}
