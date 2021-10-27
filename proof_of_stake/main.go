package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
	"time"
)

var done chan bool  // this is channel is for signalling cancellation

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	// create genesis block
	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{0, t.String(), 0, calculateBlockHash(genesisBlock), "", ""}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	// start TCP and serve TCP server
	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer server.Close()

	done = make(chan bool)
	defer close(done)

	go func() {
		for candidate := range candidateBlocks {
			mutex.Lock()
			tempBlocks = append(tempBlocks, candidate)
			mutex.Unlock()
		}
	}()

	go func() {
		every30Secs := time.Tick(time.Duration(30) * time.Second)
		for {
			select {
			case <-every30Secs:
				pickWinner()
			case <-done:
				return
			}

		}
	}()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}
