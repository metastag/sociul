package main

import (
	"log"
	"net/http"
	"sociul-api-gateway/internal/config"
	"sociul-api-gateway/internal/proxy"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Config Error-", err)
	}

	// Initialize HTTP client to forward requests to internal services
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Initialize proxy
	proxy := proxy.NewProxy(client, cfg)

	// Initialize Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// Map routes
	router.Any("/auth/*path", proxy.Auth) // unprotected routes

	// Define and start server
	srv := &http.Server{
		Addr:         cfg.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Println("Starting server on Port", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
