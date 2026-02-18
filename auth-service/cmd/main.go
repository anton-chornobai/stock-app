package main

import (
	"auth-service/internal/application"
	"auth-service/internal/config"
	"auth-service/internal/infra"
	"auth-service/internal/infra/postgres"
	grpchandler "auth-service/internal/transport/grpc"
	"context"
	"log"
	"net"
	"os"
	"time"

	authpb "github.com/anton-chornobai/stock-protos/auth/gen"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Couldnt .env couldnt be loaded %v", err)
	}	
	secret := os.Getenv("SECRET") 
	if secret == "" {
		log.Fatalf("Secret couldnt is empty")
	}

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("couldnt load config, err: %v", err)
	}

	db, err := postgres.OpenDB(cfg.DB.ConnString())
	if err != nil {
		log.Fatalf("Failed to open DB %v", err)
	}

	userRepo := &infra.UserRepo{DB: db}
	authService := application.NewAuthService(userRepo, []byte(secret))

	//GRPC SET UP
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//handler with the implemented methods
	authHandler := grpchandler.NewAuthHandler(authService)

	gprcServer := grpc.NewServer()

	authpb.RegisterAuthServer(gprcServer, authHandler)

	go func() {
		if err := gprcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	//HTTP GIN SERVER START
	r := gin.Default()

	r.POST("/signup", func(c *gin.Context) {
		var req application.SignupRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5 * time.Second)
		defer cancel()
		token, err := authService.Signup(ctx, req)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"token": token,
			"exp":   time.Now().Add(15 * time.Minute).Unix(),
		})
	})
	r.POST("/login", func(c *gin.Context) {

	})
	r.GET("/profile", func(c *gin.Context) {

	})

	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal("couldnt set proxies")
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("couldnt start GIN server ")
	}
}
