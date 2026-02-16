package main

import (
	"log"
	"net"
	"auth-service/internal/config"
	"auth-service/internal/infra"
	"auth-service/internal/infra/postgres"
	"auth-service/internal/application"
	grpchandler "auth-service/internal/transport/grpc"

	authpb "github.com/anton-chornobai/stock-protos/auth/gen"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Couldnt .env couldnt be loaded %v", err)
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
	authService := application.NewAuthService(userRepo)

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

	type SignupRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	r.POST("/signup", func(c *gin.Context) {
		var req SignupRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		// authService.Signup()
	})
	r.POST("/login", func(c *gin.Context) {

	})
	r.GET("/profile", func(c *gin.Context) {

	})

	 r.SetTrustedProxies(nil) 

	r.Run(":8080")
}
