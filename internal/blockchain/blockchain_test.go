package blockchain

import (
	"path/filepath"
	"testing"
	"time"
)

var defaultDiffCfg = DifficultyConfig{
	AdjustInterval:  5,
	TargetBlockTime: 10,
	MinDifficulty:   1,
	MaxDifficulty:   6,
}

func TestHashDeterminism(t *testing.T) {
	b := NewBlock(1, "test data", "prevhash")
	hash1 := b.CalculateHash()
	hash2 := b.CalculateHash()
	if hash1 != hash2 {
		t.Errorf("hash not deterministic: %s != %s", hash1, hash2)
	}
}

func TestHashChangesOnDataMutation(t *testing.T) {
	b1 := NewBlock(1, "data A", "prevhash")
	b2 := NewBlock(1, "data B", "prevhash")
	// Timestamps may differ, but even with same timestamp different data = different hash
	b2.Timestamp = b1.Timestamp
	b2.Hash = b2.CalculateHash()
	if b1.Hash == b2.Hash {
		t.Error("different data should produce different hashes")
	}
}

func TestProofOfWork(t *testing.T) {
	b := NewBlock(1, "pow test", "prevhash")
	difficulty := 2
	b.MineBlock(difficulty)
	prefix := "00"
	if b.Hash[:2] != prefix {
		t.Errorf("mined hash %s does not start with %s", b.Hash, prefix)
	}
}

func TestGenesisBlock(t *testing.T) {
	bc, err := NewBlockchain(2, "", defaultDiffCfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bc.Chain) != 1 {
		t.Fatalf("expected 1 block (genesis), got %d", len(bc.Chain))
	}
	genesis := bc.Chain[0]
	if genesis.Index != 0 {
		t.Errorf("genesis index should be 0, got %d", genesis.Index)
	}
	if genesis.PrevHash != "0" {
		t.Errorf("genesis prev_hash should be '0', got %s", genesis.PrevHash)
	}
}

func TestValidChain(t *testing.T) {
	bc, err := NewBlockchain(2, "", defaultDiffCfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	bc.AddBlock("block 1")
	bc.AddBlock("block 2")
	if !bc.IsValid() {
		t.Error("chain should be valid")
	}
}

func TestTamperedChainIsInvalid(t *testing.T) {
	bc, err := NewBlockchain(2, "", defaultDiffCfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	bc.AddBlock("block 1")
	bc.AddBlock("block 2")

	// Tamper with a block
	bc.Chain[1].Data = "tampered data"
	if bc.IsValid() {
		t.Error("tampered chain should be invalid")
	}
}

func TestPersistenceRoundtrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "blockchain.db")

	bc, err := NewBlockchain(2, path, defaultDiffCfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	bc.AddBlock("persist test")
	bc.Close()

	// Load from DB
	bc2, err := NewBlockchain(2, path, defaultDiffCfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer bc2.Close()

	if len(bc2.Chain) != len(bc.Chain) {
		t.Fatalf("expected %d blocks, got %d", len(bc.Chain), len(bc2.Chain))
	}
	if bc2.Chain[1].Data != "persist test" {
		t.Errorf("expected data 'persist test', got '%s'", bc2.Chain[1].Data)
	}
}

func TestInMemoryPersistence(t *testing.T) {
	bc, err := NewBlockchain(2, ":memory:", defaultDiffCfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer bc.Close()

	bc.AddBlock("memory test")
	if len(bc.Chain) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(bc.Chain))
	}
	if bc.Chain[1].Data != "memory test" {
		t.Errorf("expected data 'memory test', got '%s'", bc.Chain[1].Data)
	}
}

func TestDifficultyIncreasesWhenMiningTooFast(t *testing.T) {
	cfg := DifficultyConfig{
		AdjustInterval:  3,
		TargetBlockTime: 10,
		MinDifficulty:   1,
		MaxDifficulty:   6,
	}
	bc, err := NewBlockchain(2, ":memory:", cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer bc.Close()

	// Set timestamps very close together (fast mining)
	now := time.Now().UTC()
	bc.Chain[0].Timestamp = now.Format(time.RFC3339)

	for i := 0; i < 3; i++ {
		block := NewBlock(len(bc.Chain), "fast block", bc.Chain[len(bc.Chain)-1].Hash)
		block.Timestamp = now.Add(time.Duration(i+1) * time.Second).Format(time.RFC3339)
		block.MineBlock(bc.Difficulty)
		bc.mu.Lock()
		bc.Chain = append(bc.Chain, block)
		bc.store.SaveBlock(block)
		bc.adjustDifficulty()
		bc.mu.Unlock()
	}

	if bc.Difficulty <= 2 {
		t.Errorf("expected difficulty to increase above 2, got %d", bc.Difficulty)
	}
}

func TestDifficultyDecreasesWhenMiningTooSlow(t *testing.T) {
	cfg := DifficultyConfig{
		AdjustInterval:  3,
		TargetBlockTime: 10,
		MinDifficulty:   1,
		MaxDifficulty:   6,
	}
	bc, err := NewBlockchain(3, ":memory:", cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer bc.Close()

	// Set timestamps far apart (slow mining: >2x expected)
	now := time.Now().UTC()
	bc.Chain[0].Timestamp = now.Format(time.RFC3339)

	for i := 0; i < 3; i++ {
		block := NewBlock(len(bc.Chain), "slow block", bc.Chain[len(bc.Chain)-1].Hash)
		block.Timestamp = now.Add(time.Duration((i+1)*120) * time.Second).Format(time.RFC3339)
		block.MineBlock(bc.Difficulty)
		bc.mu.Lock()
		bc.Chain = append(bc.Chain, block)
		bc.store.SaveBlock(block)
		bc.adjustDifficulty()
		bc.mu.Unlock()
	}

	if bc.Difficulty >= 3 {
		t.Errorf("expected difficulty to decrease below 3, got %d", bc.Difficulty)
	}
}

func TestDifficultyStaysInBounds(t *testing.T) {
	cfg := DifficultyConfig{
		AdjustInterval:  3,
		TargetBlockTime: 10,
		MinDifficulty:   2,
		MaxDifficulty:   4,
	}

	// Test min bound
	bc, err := NewBlockchain(2, ":memory:", cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer bc.Close()

	now := time.Now().UTC()
	bc.Chain[0].Timestamp = now.Format(time.RFC3339)

	for i := 0; i < 3; i++ {
		block := NewBlock(len(bc.Chain), "slow", bc.Chain[len(bc.Chain)-1].Hash)
		block.Timestamp = now.Add(time.Duration((i+1)*200) * time.Second).Format(time.RFC3339)
		block.MineBlock(bc.Difficulty)
		bc.mu.Lock()
		bc.Chain = append(bc.Chain, block)
		bc.store.SaveBlock(block)
		bc.adjustDifficulty()
		bc.mu.Unlock()
	}

	if bc.Difficulty < cfg.MinDifficulty {
		t.Errorf("difficulty %d went below min %d", bc.Difficulty, cfg.MinDifficulty)
	}

	// Test max bound
	bc2, err := NewBlockchain(4, ":memory:", cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer bc2.Close()

	bc2.Chain[0].Timestamp = now.Format(time.RFC3339)

	for i := 0; i < 3; i++ {
		block := NewBlock(len(bc2.Chain), "fast", bc2.Chain[len(bc2.Chain)-1].Hash)
		block.Timestamp = now.Add(time.Duration(i+1) * time.Second).Format(time.RFC3339)
		block.MineBlock(bc2.Difficulty)
		bc2.mu.Lock()
		bc2.Chain = append(bc2.Chain, block)
		bc2.store.SaveBlock(block)
		bc2.adjustDifficulty()
		bc2.mu.Unlock()
	}

	if bc2.Difficulty > cfg.MaxDifficulty {
		t.Errorf("difficulty %d went above max %d", bc2.Difficulty, cfg.MaxDifficulty)
	}
}

func TestDifficultyUnchangedWithinInterval(t *testing.T) {
	cfg := DifficultyConfig{
		AdjustInterval:  5,
		TargetBlockTime: 10,
		MinDifficulty:   1,
		MaxDifficulty:   6,
	}
	bc, err := NewBlockchain(2, ":memory:", cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer bc.Close()

	// Add only 3 blocks (less than interval of 5), difficulty should not change
	for i := 0; i < 3; i++ {
		_, adjusted, err := bc.AddBlock("test block")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if adjusted {
			t.Error("difficulty should not adjust before reaching interval")
		}
	}

	if bc.Difficulty != 2 {
		t.Errorf("expected difficulty to remain 2, got %d", bc.Difficulty)
	}
}
