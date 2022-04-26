package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"

	"github.com/rusinikita/gogoClub/airtable"
	"github.com/rusinikita/gogoClub/request"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	db := airtable.New()

	log.Println(db.Create(context.Background(), request.Request{
		UserID: 24,
		Name:   "Bla bla",
	}))
}
