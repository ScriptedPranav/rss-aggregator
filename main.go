package main

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
)

func main() {
	portString := os.Getenv("PORT")

	godotenv.Load()
	
	if portString == "" {
		log.Fatal("Port not found")
	}

	fmt.Println(portString)
}