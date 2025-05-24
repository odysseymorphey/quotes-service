package main

import (
	"github.com/odysseymorphey/quotes-service/internal/server"
)

type Mock struct{}

func main() {
	r := &Mock{}
	s := server.New(r)

	s.Run()
}
