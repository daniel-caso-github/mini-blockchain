package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/danielcaso/mini-blockchain/internal/blockchain"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(bc *blockchain.Blockchain) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	hub := NewHub()
	go hub.Run()

	h := NewHandler(bc, hub)

	r.GET("/health", h.HealthCheck)
	r.GET("/chain", h.GetChain)
	r.POST("/mine", h.MineBlock)
	r.GET("/validate", h.ValidateChain)
	r.GET("/block/:id", h.GetBlockByID)
	r.GET("/ws", h.HandleWS)
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve frontend static files in production
	distPath := filepath.Join("frontend", "dist")
	if _, err := os.Stat(distPath); err == nil {
		r.Static("/assets", filepath.Join(distPath, "assets"))
		r.StaticFile("/vite.svg", filepath.Join(distPath, "vite.svg"))
		r.NoRoute(func(c *gin.Context) {
			c.File(filepath.Join(distPath, "index.html"))
		})
	} else {
		r.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "not found"})
		})
	}

	return r
}
