package main

import (
"fmt"
"math"
"math/rand"
"time"
"os"
"strconv"
)

const threshold = 3

type Message struct {
    msgType int
    s float64
    w float64
    sender *Peer
}

type Peer struct {
    id int
    neighboor []Peer
    count int
    s float64
    w float64
    msgQueue chan Message
}

func (p Peer) Start(convergent chan float64) {
    msg := <-p.msgQueue
    // fmt.Printf("[KNOWS] Peer %d knows the roumor\n", p.id)
    p.ProcessMessage(msg)
    p.Propagate(convergent, msg)
}

func (p Peer) Propagate(convergent chan float64, msg Message) {
    ticker := time.NewTicker(10 * time.Millisecond)
    cancelTick := make(chan bool)
    go func() {
        for {
            select {
            case <-ticker.C:
                p.SendMessage()
            case <-cancelTick:
                ticker.Stop()
                return
            }
        }
    }()

    for {
        m := <-p.msgQueue
        p.ProcessMessage(m)
        if p.count >= threshold {
            // fmt.Printf("Peer %d converge, sum is %v\n", p.id, p.s / p.w)
            cancelTick <- true
            convergent <- p.s / p.w
        }
    }
}

func (p *Peer) SendMessage() {
    p.s /= 2.0
    p.w /= 2.0
    i := rand.Intn(len(p.neighboor))
    // fmt.Printf("Peer %d send to Peer %d, s: %f, w: %f\n" , p.id, p.neighboor[i].id, p.s, p.w)
    p.neighboor[i].msgQueue <- Message{0, p.s, p.w, p}
}

func (p *Peer) ProcessMessage(msg Message) {
    ratio := p.s / p.w
    p.s += msg.s
    p.w += msg.w
    newRatio := p.s / p.w
    if math.Abs(ratio - newRatio) < 1e-10 {
        p.count++
    } else {
        p.count = 0
    }
    // fmt.Printf("Peer %d receive from Peer %d, sum: %f\n" , p.id, msg.sender.id, p.s / p.w)
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

func Map2DTo1D(x int, y int, size int) int {
    return x + y * size
}

func Build2DGrid(peers []Peer, numOfPeers int) {
    size := int(math.Sqrt(float64(numOfPeers)) + 0.5)
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
    size := int(math.Pow(float64(numOfPeers), 1.0/3) + 0.5)
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
            if IdInNeighbour(random, peers[i].neighboor) || random == i {
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
        fmt.Println("Usage: go run pushsum.go [num of peers][topology]")
        return
    }
    // get number of peers and topology from input
    numOfPeers, err := strconv.Atoi(os.Args[1])
    if err != nil  { //|| num0s < 1 || num0s > 32
        fmt.Println("Usage: go run pushsum.go [num of peers][topology]")
        return
    }
    topology := os.Args[2]
    // if topology == "2D" || topology == "imp2D" {
    //     numOfPeers = int(math.Pow(float64(int(math.Sqrt(float64(numOfPeers)) + 0.5)), 2.0))
    // } else if topology == "3D" || topology == "imp3D"{
    //     numOfPeers = int(math.Pow(float64(int(math.Pow(float64(numOfPeers), 1.0/3) + 0.5)), 3.0))
    // }
    // create peers
    peers := make([]Peer, numOfPeers)
    for i := 0; i < numOfPeers; i++ {
        peers[i] = Peer{id:i, count:0, s:float64(i) + 1.0, w:0.0, msgQueue:make(chan Message)}
    }
    // build topology
    BuildNetwork(peers[:], topology, numOfPeers)

    convergent := make(chan float64)
    rand.Seed(time.Now().UTC().UnixNano())
    for i := 0; i < numOfPeers; i++ {
        go peers[i].Start(convergent)
    }
    // start push sum
    start := time.Now()
    peers[0].msgQueue <- Message{0, 0.0, 1.0, &peers[0]}
    fmt.Printf("Sum from 1 to %d: %f\n", numOfPeers, <-convergent)
    elapsed := time.Since(start)
    fmt.Printf("[END] %v peers, %v network spent %v\n", numOfPeers, topology, elapsed)
}
