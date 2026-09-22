package main

import (
	"context"
	"log"

	"github.com/levisantosp/atm-participa/api/db"
	"github.com/levisantosp/atm-participa/api/utils"
)

func main() {
	utils.LoadEnv("../.env")
	db.Connect()

	err := db.Client.Schema.Create(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
