/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"
	"github.com/joho/godotenv"
	"github.com/nareshkarthigeyan/revly/cmd"
)

func main() {

	 err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	
	cmd.Execute()
}
