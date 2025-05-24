package main

import (
	"github.com/aifedorov/gophermart/internal/config"
	"log"
)

func main() {
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
}
