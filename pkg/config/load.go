package config

import (
	"log"

	"github.com/joho/godotenv"
)

/*
Load .env file details into os environment
*/
func LoadEnvironment() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("env file loaded")
}
