package main

import (
	"log"

	_ "github.com/danielcaso/mini-blockchain/docs"
	"github.com/danielcaso/mini-blockchain/internal/api"
	"github.com/danielcaso/mini-blockchain/internal/blockchain"
	"github.com/danielcaso/mini-blockchain/internal/config"
)

// @title           Mini Blockchain API
// @version         1.0
// @description     A minimal blockchain implementation in Go with REST API.
// @host            localhost:8080
// @BasePath        /
func main() {
	cfg := config.Load()
	bc, err := blockchain.NewBlockchain(cfg.Difficulty, cfg.DBPath, blockchain.DifficultyConfig{
		AdjustInterval:  cfg.AdjustmentInterval,
		TargetBlockTime: cfg.TargetBlockTime,
		MinDifficulty:   cfg.MinDifficulty,
		MaxDifficulty:   cfg.MaxDifficulty,
	})
	if err != nil {
		log.Fatalf("Failed to initialize blockchain: %v", err)
	}
	defer bc.Close()
	r := api.SetupRouter(bc)
	log.Printf("Starting server on :%s (difficulty=%d)", cfg.Port, cfg.Difficulty)
	log.Fatal(r.Run(":" + cfg.Port))
}
