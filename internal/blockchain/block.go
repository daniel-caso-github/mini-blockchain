package blockchain

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	Index     int    `json:"index" example:"1"`
	Timestamp string `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Data      string `json:"data" example:"Hello Blockchain"`
	PrevHash  string `json:"prev_hash" example:"0000abcdef1234567890"`
	Hash      string `json:"hash" example:"0000fedcba0987654321"`
	Nonce     int    `json:"nonce" example:"12345"`
}

func NewBlock(index int, data, prevHash string) Block {
	b := Block{
		Index:     index,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
		PrevHash:  prevHash,
	}
	b.Hash = b.CalculateHash()
	return b
}

func (b *Block) CalculateHash() string {
	record := fmt.Sprintf("%d%s%s%s%d", b.Index, b.Timestamp, b.Data, b.PrevHash, b.Nonce)
	h := sha256.Sum256([]byte(record))
	return fmt.Sprintf("%x", h)
}

func (b *Block) MineBlock(difficulty int) {
	prefix := strings.Repeat("0", difficulty)
	for !strings.HasPrefix(b.Hash, prefix) {
		b.Nonce++
		b.Hash = b.CalculateHash()
	}
}
