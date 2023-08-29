package main

import (
	"goauth/internal/server"
	"goauth/internal/usecase"
	"goauth/internal/storage"
	"goauth/internal/config"
	"os"
	"os/signal"
	"context"
	"log"
)

func main() {
	appConfig, err := config.New()
	if err!= nil {
        log.Fatal(err)
    }

	db :=  storage.ConnectDB(appConfig.MongoDB)
	appStorage := storage.NewUserStorage(db, appConfig.MongoDB)
	appInteractor := usecase.NewAuthInteractor(appStorage, appConfig.Auth)
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