package main

import (
	"log"
	"os"

	"github.com/muskelo/translator_bot/internal/sql"
)

func main() {
	engine, err := sql.New(os.Getenv("DB_CONNSTR"))
	if err != nil {
		log.Fatalf("Can't init db engine: %v", err.Error())
	}
	err = engine.Sync(new(sql.User))
	if err != nil {
		log.Fatalf("Can't migrate: %v", err.Error())
	}
}
