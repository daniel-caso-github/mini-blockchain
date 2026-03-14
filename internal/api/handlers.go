package api

import (
	"net/http"
	"strconv"

	"github.com/danielcaso/mini-blockchain/internal/blockchain"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	bc  *blockchain.Blockchain
	hub *Hub
}

func NewHandler(bc *blockchain.Blockchain, hub *Hub) *Handler {
	return &Handler{bc: bc, hub: hub}
}

// GetChain godoc
// @ID           get-chain
// @Summary      Get the blockchain with pagination
// @Description  Returns blocks with pagination. Use ?page=1&limit=20 query params. Defaults to page=1, limit=20.
// @Tags         blockchain
// @Produce      json
// @Param        page   query     int  false  "Page number (1-based)"  default(1)
// @Param        limit  query     int  false  "Blocks per page"        default(20)
// @Success      200  {object}  ChainResponse
// @Router       /chain [get]
func (h *Handler) GetChain(c *gin.Context) {
	page := 1
	limit := 20

	if v := c.Query("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			page = p
		}
	}
	if v := c.Query("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 {
			if l > 100 {
				l = 100
			}
			limit = l
		}
	}

	blocks, total := h.bc.GetChainPaginated(page, limit)
	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, ChainResponse{
		Chain:      blocks,
		Length:     len(blocks),
		Difficulty: h.bc.Difficulty,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// MineBlock godoc
// @ID           mine-block
// @Summary      Mine a new block
// @Description  Creates and mines a new block with the provided data using Proof of Work. The block is appended to the chain and persisted to the database
// @Tags         blockchain
// @Accept       json
// @Produce      json
// @Param        request  body      MineRequest  true  "Block data"
// @Success      201      {object}  blockchain.Block
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /mine [post]
func (h *Handler) MineBlock(c *gin.Context) {
	var req MineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "field 'data' is required"})
		return
	}

	block, adjusted, err := h.bc.AddBlock(req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	if h.hub != nil {
		h.hub.BroadcastBlock(block)
		if adjusted {
			h.hub.BroadcastDifficulty(h.bc.Difficulty)
		}
	}
	c.JSON(http.StatusCreated, block)
}

func (h *Handler) HandleWS(c *gin.Context) {
	h.hub.HandleWS(c)
}

// ValidateChain godoc
// @ID           validate-chain
// @Summary      Validate chain integrity
// @Description  Verifies the integrity of the entire blockchain by checking hashes, links, and proof of work for every block
// @Tags         blockchain
// @Produce      json
// @Success      200  {object}  ValidateResponse
// @Router       /validate [get]
func (h *Handler) ValidateChain(c *gin.Context) {
	c.JSON(http.StatusOK, ValidateResponse{Valid: h.bc.IsValid()})
}

// GetBlockByID godoc
// @ID           get-block-by-id
// @Summary      Get a block by index
// @Description  Returns a single block from the blockchain by its numeric index (0 = genesis block)
// @Tags         blockchain
// @Produce      json
// @Param        id   path      int  true  "Block index"
// @Success      200  {object}  blockchain.Block
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /block/{id} [get]
func (h *Handler) GetBlockByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid block id"})
		return
	}

	block, err := h.bc.GetBlock(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, block)
}

// HealthCheck godoc
// @ID           health-check
// @Summary      Health check
// @Description  Returns the current health status of the API server
// @Tags         system
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{Status: "ok"})
}
