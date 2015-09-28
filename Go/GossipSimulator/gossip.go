package main

import (
"fmt"
"math"
"math/rand"
"time"
"os"
"strconv"
)

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
                // fmt.Printf("[CONVERGE] Peer %d converged\n", p.id)
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
                // fmt.Printf("[CONVERGE] Peer %d is isolated\n", p.id)
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

func BuildFullNetwork(peers []Peer, numOfPeers int) {
    for i := 0; i < numOfPeers; i++ {
        for j := 0; j < numOfPeers; j++ {
            if j != i {
                peers[i].neighboor = append(peers[i].neighboor, peers[j])
            }
        }
    }
}

func BuildLineNetwork(peers []Peer, numOfPeers int) {
    for i :=1; i < numOfPeers - 1; i++ {
        peers[i].neighboor = append(peers[i].neighboor, peers[i - 1])
        peers[i].neighboor = append(peers[i].neighboor, peers[i + 1])
    }
    peers[0].neighboor = append(peers[0].neighboor, peers[1])
    peers[numOfPeers - 1].neighboor = append(peers[numOfPeers - 1].neighboor, peers[numOfPeers - 2])
}

func Map2DTo1D(x int, y int, size int) int {
    return x + y * size
}

func Build2DGrid(peers []Peer, numOfPeers int) {
    size := int(math.Sqrt(float64(numOfPeers)))
    for i := 0; i < size; i++ {
        for j := 0; j < size; j++ {
            pos := Map2DTo1D(j, i, size)
            if j > 0 {
                peers[pos].neighboor = append(peers[pos].neighboor, peers[Map2DTo1D(j - 1, i, size)])
            }
            if j < size - 1 {
                peers[pos].neighboor = append(peers[pos].neighboor, peers[Map2DTo1D(j + 1, i, size)])
            }
            if i > 0 {
                peers[pos].neighboor = append(peers[pos].neighboor, peers[Map2DTo1D(j, i - 1, size)])
            }
            if i < size - 1 {
                peers[pos].neighboor = append(peers[pos].neighboor, peers[Map2DTo1D(j, i + 1, size)])
            }
        }
    }
}

func IdInNeighbour(id int, neighboor []Peer) bool {
    for _, n := range neighboor {
        if id == n.id {
            return true
        }
    }
    return false
}

func Map3DTo1D(x int, y int, z int, size int) int {
    return x + y * size + z * size * size
}

func Build3DGrid(peers []Peer, numOfPeers int) {
    size := int(math.Pow(float64(numOfPeers), 1.0/3))
    for i := 0; i < size; i++ {
        for j := 0; j < size; j++ {
            for k := 0; k < size; k++ {
                pos := Map3DTo1D(k, j, i, size)
                if k > 0 {
                    peers[pos].neighboor = append(peers[pos].neighboor, peers[Map3DTo1D(k - 1, j, i, size)])
                }
                if k < size - 1 {
                    peers[pos].neighboor =  append(peers[pos].neighboor, peers[Map3DTo1D(k + 1, j, i, size)])
                }
                if j > 0 {
                    peers[pos].neighboor =  append(peers[pos].neighboor, peers[Map3DTo1D(k, j - 1, i, size)])
                }
                if j < size - 1 {
                     peers[pos].neighboor =  append(peers[pos].neighboor, peers[Map3DTo1D(k, j + 1, i, size)])
                }
                if i > 0 {
                    peers[pos].neighboor =  append(peers[pos].neighboor, peers[Map3DTo1D(k, j, i - 1, size)])
                }
                if i < size - 1 {
                    peers[pos].neighboor =  append(peers[pos].neighboor, peers[Map3DTo1D(k, j, i + 1, size)])
                }
            }
        }
    }
}

func BuildImpGrid(peers []Peer, numOfPeers int, buildGrid func(peers []Peer, numOfPeers int)) {
    buildGrid(peers, numOfPeers)
    for i := 0; i < numOfPeers; i++ {
        for {
            random := rand.Intn(numOfPeers)
            if (IdInNeighbour(random, peers[i].neighboor)) {
                continue
            } else {
                peers[i].neighboor = append(peers[i].neighboor, peers[random])
                break
            }
        }
    }
}

func BuildNetwork(peers []Peer, topology string, numOfPeers int) {
    switch topology {
    case "full":
        BuildFullNetwork(peers, numOfPeers)
    case "line":
        BuildLineNetwork(peers, numOfPeers)
    case "2D":
        Build2DGrid(peers, numOfPeers)
    case "imp2D":
        BuildImpGrid(peers, numOfPeers, Build2DGrid)
    case "3D":
        Build3DGrid(peers, numOfPeers)
    case "imp3D":
        BuildImpGrid(peers, numOfPeers, Build3DGrid)
    default:
        panic("Topology is not implemented")
    }
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: go run gossip.go [num of peers][topology]")
        return
    }
    // get number of peers and topology from input
    numOfPeers, err := strconv.Atoi(os.Args[1])
    if err != nil  { //|| num0s < 1 || num0s > 32
        fmt.Println("Usage: go run gossip.go [num of peers][topology]")
        return
    }
    topology := os.Args[2]
    // create peers
    peers := make([]Peer, numOfPeers)
    for i := 0; i < numOfPeers; i++ {
        peers[i] = Peer{id:i, count:0, msgQueue:make(chan Message)}
    }
    // build topology
    BuildNetwork(peers[:], topology, numOfPeers)

    convergent := make(chan bool)
    rand.Seed(time.Now().UTC().UnixNano())
    for i := 0; i < numOfPeers; i++ {
        go peers[i].Start(convergent)
    }
    // start gossip
    start := time.Now()
    peers[0].msgQueue <- Message{0, "Hello, I am main process", &peers[0]}
    for i := 0; i < numOfPeers; i++ {
        <-convergent
    }
    elapsed := time.Since(start)
    fmt.Printf("[END] %v peers, %v network spent %v\n", numOfPeers, topology, elapsed)
}
