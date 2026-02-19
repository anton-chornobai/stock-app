package main

import (
	"auth-service/internal/application"
	"auth-service/internal/config"
	"auth-service/internal/infra"
	"auth-service/internal/infra/postgres"
	grpchandler "auth-service/internal/transport/grpc"
	httphandler "auth-service/internal/transport/http"
	"auth-service/internal/transport/http/handlers"
	"log"
	"log/slog"
	"net"
	"os"

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

	// SLOGGER

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	userRepo := &infra.UserRepo{DB: db}
	authService := application.NewAuthService(userRepo, []byte(secret), logger)

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
	httpHandler := handlers.NewAuthHandler(*authService, logger)
	r := gin.Default()

	r.Use(httphandler.ValidateJWT([]byte(secret)))
	r.POST("/signup", httpHandler.Signup)
	r.POST("/login", httpHandler.Login)
	
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal("couldnt set proxies")
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("couldnt start GIN server ")
	}
}
