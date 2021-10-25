package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var Blockchain []Block

type Block struct {
	Index      int
	Timestamp  string
	BPM        int
	Hash       string
	PrevHash   string
	Difficulty int
	Nonce      string
}

func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, BPM int) (Block, error) {
	var newBlock Block

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = time.Now().String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	//newBlock.Hash = calculateHash(newBlock)
	newBlock.Difficulty = oldBlock.Difficulty

	for i := 0; ; i++ {
		ihex := fmt.Sprintf("%x", i)
		newBlock.Nonce = ihex
		newHash := calculateHash(newBlock)
		if !isHashValid(newHash, newBlock.Difficulty) {
			fmt.Println(newHash + " not valid hash.")
			//time.Sleep(time.Duration(500) * time.Millisecond) // to simulate the PoW duration
			continue
		} else {
			fmt.Printf("\n"+newHash+" work is done and it took:%d attempts", i)
			newBlock.Hash = newHash
			break
		}
	}

	return newBlock, nil
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func getLastBlock() Block {
	return Blockchain[len(Blockchain)-1]
}

func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}
