package main

import (
	"context"
	"goauth/internal/config"
	"goauth/internal/server"
	"goauth/internal/storage"
	"goauth/internal/usecase"
	"goauth/pkg/manager"
	"goauth/pkg/hash"
	"log"
	"os"
	"os/signal"
)

func main() {
	appConfig, err := config.New()
	if err!= nil {
        log.Fatal(err)
    }
	authConfig := appConfig.Auth
	tokenManager := manager.NewManager(authConfig.SigningKey, authConfig.AccessTokenTTL, authConfig.RefreshTokenTTL)
	passwordHasher := hash.NewSHA1Hasher(authConfig.HashSalt)

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