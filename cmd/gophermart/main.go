package main

import (
	"log"
	"os"

	"github.com/sergeizaitcev/gophermart/internal/gophermart"
)

func main() {
	if err := gophermart.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
