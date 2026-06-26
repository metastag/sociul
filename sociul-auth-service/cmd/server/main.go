package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sociul-auth-service/internal/config"
	"sociul-auth-service/internal/handlers"
	"sociul-auth-service/internal/repository"
	"sociul-auth-service/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mailgun/mailgun-go/v5"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load config from environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error initializing config - ", err)
	}

	// Initialize Postgres connection pool
	ctx := context.Background()
	pool, err := repository.NewPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Error initializing database - ", err)
	}
	defer pool.Close()

	// Initialize Redis client
	redisCache := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: "",
		DB:       0,
	})
	defer redisCache.Close()

	// Initialize Mailgun email client
	mg := mailgun.NewMailgun(cfg.MailgunAPIKey)

	// Initialize repositories
	authRepo := repository.NewAuthRepository(pool.Pool)
	cache := repository.NewCache(redisCache)

	// Initialize services
	authService := service.NewAuthService(authRepo, cache, cfg.JwtKey)
	otpService := service.NewOtpService(authRepo, cache, mg, cfg.MailgunDomain)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	otpHandler := handlers.NewOtpHandler(otpService)

	// Initialize Gin router
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Add a ping endpoint to test connectivity
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Map routes
	router.POST("/auth/signup", authHandler.SignUp)
	router.POST("/auth/login", authHandler.Login)
	router.POST("/auth/refresh", authHandler.Refresh)
	router.POST("/auth/send-email-otp", otpHandler.SendEmailOtp)
	router.POST("/auth/verify-email-otp", otpHandler.VerifyEmailOtp)
	router.POST("/auth/forgot-password", otpHandler.SendResetOtp)
	router.POST("/auth/verify-reset-otp", otpHandler.VerifyResetOtp)
	router.POST("/auth/reset-password", authHandler.ResetPassword)
	router.POST("/auth/logout", authHandler.VerifyToken(), authHandler.Logout)
	router.DELETE("/auth/delete", authHandler.VerifyToken(), authHandler.Delete)

	// Define and start server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start server
	go func() {
		log.Println("Starting the server on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // block here until a signal arrives

	log.Println("Shutting down server...")

	// Give in-flight requests up to 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close redis and db connections
	redisCache.Close()
	pool.Close()

	log.Println("Server stopped.")
}
