package main

import (
	"fmt"
	"log"
	"stock-app/internal/services/auth/application"
	"stock-app/internal/services/auth/config"
	"stock-app/internal/services/auth/infra"
	"stock-app/internal/services/auth/infra/postgres"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Couldnt .env couldnt be loaded %v", err)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("couldnt load config, err: %v", err)
	}
	fmt.Println(cfg.DB.ConnString())
	db, err := postgres.OpenDB(cfg.DB.ConnString())

	if err != nil {
		log.Fatalf("Failed to open DB %v", err)
	}

	userRepo := &infra.UserRepo{DB: db}
	authService := application.NewAuthService(userRepo)

	_ = authService
	
}
