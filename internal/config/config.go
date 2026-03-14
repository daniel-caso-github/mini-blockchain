package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port               string
	Difficulty         int
	DBPath             string
	AdjustmentInterval int
	TargetBlockTime    int
	MinDifficulty      int
	MaxDifficulty      int
}

func Load() Config {
	cfg := Config{
		Port:               "8080",
		Difficulty:         2,
		DBPath:             "data/blockchain.db",
		AdjustmentInterval: 5,
		TargetBlockTime:    10,
		MinDifficulty:      1,
		MaxDifficulty:      6,
	}

	if v := os.Getenv("PORT"); v != "" {
		cfg.Port = v
	}
	if v := os.Getenv("DIFFICULTY"); v != "" {
		if d, err := strconv.Atoi(v); err == nil && d > 0 {
			cfg.Difficulty = d
		}
	}
	if v := os.Getenv("ADJUSTMENT_INTERVAL"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.AdjustmentInterval = n
		}
	}
	if v := os.Getenv("TARGET_BLOCK_TIME"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.TargetBlockTime = n
		}
	}
	if v := os.Getenv("MIN_DIFFICULTY"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.MinDifficulty = n
		}
	}
	if v := os.Getenv("MAX_DIFFICULTY"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.MaxDifficulty = n
		}
	}
	if v := os.Getenv("DB_PATH"); v != "" {
		cfg.DBPath = v
	} else if v := os.Getenv("CHAIN_FILE"); v != "" {
		cfg.DBPath = v
	}

	return cfg
}
