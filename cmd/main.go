package main

import (
	"log"
	"template/cmd/api"
	"template/configs"
	"template/db"
)

func main() {
	storage, err := db.NewPostgresStorage(configs.Envs)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":"+configs.Envs.Port, storage)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
