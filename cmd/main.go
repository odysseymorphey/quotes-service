package main

import (
	"fmt"
	"github.com/odysseymorphey/quotes-service/internal/server"
	"github.com/odysseymorphey/quotes-service/pkg/storage/postgres"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Mock struct{}

func main() {
	db, err := postgres.New(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	),
	)
	if err != nil {
		log.Fatalf("Can't open database: %v", err)
	}
	s := server.New(db)

	go s.Run()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	s.Stop()
	os.Exit(1)
}
