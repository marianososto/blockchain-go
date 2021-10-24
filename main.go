package main

import (
	"bufio"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

// bcServer handles incoming concurrent blocks.
var bcServer chan []Block

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	bcServer = make(chan []Block)

	t := time.Now()
	genesisBlock := Block{0, t.String(), 0, "", "",1,""}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		fmt.Println("Received new connection")
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer func() {
		fmt.Println("\nClosing connection")
		conn.Close()
	}()
	io.WriteString(conn, "Enter a new BPM:")
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		bpm, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Printf("%v not a number: %v", scanner.Text(), err)
			continue
		}
		newBlock, err := generateBlock(getLastBlock(), bpm)
		if err != nil {
			log.Printf("Failed generating new block. error:%v", err)
			continue
		}
		if isBlockValid(newBlock, getLastBlock()) {
			newBlockchain := append(Blockchain, newBlock) //TODO:danger concurrently modifying a shared variable
			replaceChain(newBlockchain)
		}
		fmt.Printf("\nAdded new blockchain with index:%d and BPM:%d", newBlock.Index, bpm)
		//bcServer <- Blockchain
		io.WriteString(conn, "\nEnter a new BPM:")
	}

	/*
		go func() {
			every30Secs := time.NewTicker(time.Duration(30) * time.Second)
			defer every30Secs.Stop()

			for {
				select {
				case <-every30Secs.C:
					output, err := json.Marshal(Blockchain)
					if err != nil {
						log.Fatal(err)
						return // end goroutine
					}
					io.WriteString(conn, string(output))
					break
				}
			}
		}()
	*/
}
