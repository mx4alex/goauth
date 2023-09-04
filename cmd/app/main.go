package main

import (
	"context"
	"goauth/internal/config"
	"goauth/internal/server"
	"goauth/internal/storage"
	"goauth/internal/usecase"
	"goauth/pkg/manager"
	"goauth/pkg/hash"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
)

// @title Authentication API
// @version 1.0
// @description API Server for Authentication

// @host localhost:8080
// @BasePath /

func main() {
	appConfig, err := config.New()
	if err != nil {
        log.Fatal("Error loading config file")
    }

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	authConfig := appConfig.Auth
	tokenManager := manager.NewManager(os.Getenv("SIGNING_KEY"), authConfig.AccessTokenTTL, authConfig.RefreshTokenTTL)
	passwordHasher := hash.NewSHA1Hasher(os.Getenv("HASH_SALT"))

	db :=  storage.ConnectDB(appConfig.MongoDB)
	appStorage := storage.NewUserStorage(db, appConfig.MongoDB)
	appInteractor := usecase.NewAuthInteractor(appStorage, tokenManager, passwordHasher)
	handlers := server.NewHandler(appInteractor)

	srv := new(server.Server)

	go func() {
		if err := srv.Run(appConfig.HostAddr, handlers.InitRoutes()); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()
	log.Print("App Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Print("App Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("error occured on server shutting down: %s", err.Error())
	}
}