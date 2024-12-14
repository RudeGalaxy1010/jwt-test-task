package main

import (
	"log"

	"github.com/RudeGalaxy1010/jwt-test-task/internal/app/apiserver"
)

func main() {
	if err := apiserver.Start(); err != nil {
		log.Fatal(err)
	}
}
