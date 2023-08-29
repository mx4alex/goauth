package main

import (
	"goauth/internal/server"
	"goauth/internal/usecase"
	"goauth/internal/storage"
	"goauth/internal/config"
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

	if err := srv.Run(appConfig.HostAddr , handlers.InitRoutes()); err != nil {
		log.Fatalf("%s", err.Error())
	}
}