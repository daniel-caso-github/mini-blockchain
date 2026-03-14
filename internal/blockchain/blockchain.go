package blockchain

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type DifficultyConfig struct {
	AdjustInterval  int
	TargetBlockTime int
	MinDifficulty   int
	MaxDifficulty   int
}

type Blockchain struct {
	Chain      []Block `json:"chain"`
	Difficulty int     `json:"difficulty"`
	store      *Store
	mu         sync.RWMutex
	diffCfg    DifficultyConfig
}

func NewBlockchain(difficulty int, dsn string, diffCfg DifficultyConfig) (*Blockchain, error) {
	store, err := OpenStore(dsn)
	if err != nil {
		return nil, fmt.Errorf("opening store: %w", err)
	}

	bc := &Blockchain{
		Difficulty: difficulty,
		store:      store,
		diffCfg:    diffCfg,
	}

	if store != nil {
		blocks, savedDiff, err := store.LoadChain()
		if err == nil && len(blocks) > 0 {
			bc.Chain = blocks
			if savedDiff > 0 {
				bc.Difficulty = savedDiff
			}
			return bc, nil
		}
	}

	genesis := NewBlock(0, "Genesis Block", "0")
	genesis.MineBlock(difficulty)
	bc.Chain = []Block{genesis}

	if err := bc.store.SaveBlock(genesis); err != nil {
		return nil, fmt.Errorf("saving genesis block: %w", err)
	}
	if err := bc.store.SaveDifficulty(difficulty); err != nil {
		return nil, fmt.Errorf("saving difficulty: %w", err)
	}

	return bc, nil
}

func (bc *Blockchain) Close() error {
	return bc.store.Close()
}

func (bc *Blockchain) AddBlock(data string) (Block, bool, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	prev := bc.Chain[len(bc.Chain)-1]
	block := NewBlock(prev.Index+1, data, prev.Hash)
	block.MineBlock(bc.Difficulty)
	bc.Chain = append(bc.Chain, block)

	if err := bc.store.SaveBlock(block); err != nil {
		return block, false, fmt.Errorf("block mined but persistence failed: %w", err)
	}

	adjusted := bc.adjustDifficulty()
	return block, adjusted, nil
}

func (bc *Blockchain) adjustDifficulty() bool {
	interval := bc.diffCfg.AdjustInterval
	if interval <= 0 {
		return false
	}
	// Only adjust at interval boundaries (skip genesis)
	if len(bc.Chain) < interval+1 || (len(bc.Chain)-1)%interval != 0 {
		return false
	}

	chainLen := len(bc.Chain)
	startBlock := bc.Chain[chainLen-interval]
	endBlock := bc.Chain[chainLen-1]

	startTime, err := time.Parse(time.RFC3339, startBlock.Timestamp)
	if err != nil {
		log.Printf("adjustDifficulty: parsing start timestamp: %v", err)
		return false
	}
	endTime, err := time.Parse(time.RFC3339, endBlock.Timestamp)
	if err != nil {
		log.Printf("adjustDifficulty: parsing end timestamp: %v", err)
		return false
	}

	actual := endTime.Sub(startTime).Seconds()
	expected := float64(interval * bc.diffCfg.TargetBlockTime)

	oldDifficulty := bc.Difficulty
	if actual < expected/2 {
		bc.Difficulty++
	} else if actual > expected*2 {
		bc.Difficulty--
	}

	if bc.diffCfg.MinDifficulty > 0 && bc.Difficulty < bc.diffCfg.MinDifficulty {
		bc.Difficulty = bc.diffCfg.MinDifficulty
	}
	if bc.diffCfg.MaxDifficulty > 0 && bc.Difficulty > bc.diffCfg.MaxDifficulty {
		bc.Difficulty = bc.diffCfg.MaxDifficulty
	}

	if bc.Difficulty != oldDifficulty {
		if err := bc.store.SaveDifficulty(bc.Difficulty); err != nil {
			log.Printf("adjustDifficulty: saving difficulty: %v", err)
		}
		log.Printf("Difficulty adjusted: %d -> %d (actual=%.1fs, expected=%.1fs)",
			oldDifficulty, bc.Difficulty, actual, expected)
		return true
	}
	return false
}

func (bc *Blockchain) IsValid() bool {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	for i := 1; i < len(bc.Chain); i++ {
		curr := bc.Chain[i]
		prev := bc.Chain[i-1]

		if curr.Hash != curr.CalculateHash() {
			return false
		}
		if curr.PrevHash != prev.Hash {
			return false
		}
		if !strings.HasPrefix(curr.Hash, strings.Repeat("0", bc.Difficulty)) {
			return false
		}
	}
	return true
}

func (bc *Blockchain) GetBlock(index int) (Block, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if index < 0 || index >= len(bc.Chain) {
		return Block{}, fmt.Errorf("block %d not found", index)
	}
	return bc.Chain[index], nil
}

func (bc *Blockchain) GetLastBlock() Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) GetChain() []Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	result := make([]Block, len(bc.Chain))
	copy(result, bc.Chain)
	return result
}

func (bc *Blockchain) GetChainPaginated(page, limit int) ([]Block, int) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	total := len(bc.Chain)
	start := (page - 1) * limit
	if start >= total {
		return []Block{}, total
	}
	end := start + limit
	if end > total {
		end = total
	}

	result := make([]Block, end-start)
	copy(result, bc.Chain[start:end])
	return result, total
}

func (bc *Blockchain) Len() int {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return len(bc.Chain)
}
