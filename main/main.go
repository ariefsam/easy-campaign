package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	log.Default().SetFlags(log.LstdFlags | log.Llongfile)
	// campaign.RunApplication()
}
