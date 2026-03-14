package api

import "github.com/danielcaso/mini-blockchain/internal/blockchain"

type MineRequest struct {
	Data string `json:"data" binding:"required" example:"Hello Blockchain"`
}

type ChainResponse struct {
	Chain      []blockchain.Block `json:"chain"`
	Length     int                `json:"length" example:"5"`
	Difficulty int                `json:"difficulty" example:"2"`
	Total      int                `json:"total" example:"42"`
	Page       int                `json:"page" example:"1"`
	Limit      int                `json:"limit" example:"20"`
	TotalPages int                `json:"total_pages" example:"3"`
}

type ValidateResponse struct {
	Valid bool `json:"valid" example:"true"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}

type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}

type WSMessage struct {
	Type  string           `json:"type"`
	Block blockchain.Block `json:"block"`
}

type WSDifficultyMessage struct {
	Type       string `json:"type"`
	Difficulty int    `json:"difficulty"`
}
