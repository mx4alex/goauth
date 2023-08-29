package main

import (
	"goauth/internal/server"
	"goauth/internal/usecase"
	"goauth/internal/storage"
	"log"
)

func main() {
	db :=  storage.ConnectDB()
	appStorage := storage.NewUserStorage(db, "users")
	appInteractor := usecase.NewAuthInteractor(appStorage)
	handlers := server.NewHandler(appInteractor)

	srv := new(server.Server)

	if err := srv.Run(":8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("%s", err.Error())
	}
}