/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/charmbracelet/glamour"
	"github.com/joho/godotenv"
	"github.com/nareshkarthigeyan/revly/cmd"
)

func checkFirstRun() {
	// Get config directory path: ~/.config/revly/
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	revlyDir := filepath.Join(configDir, "revly")
	firstRunFile := filepath.Join(revlyDir, "first_run")

	// Check if the file already exists
	if _, err := os.Stat(firstRunFile); os.IsNotExist(err) {
		// If not, show the message and create the file
		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
		)
		if err != nil {
			log.Fatal(err)
		}

		// Load .env if it exists
		err = godotenv.Load()
		var msg string
		if err != nil {
			msg = "No `.env` file found, proceeding without it.\nSet `OPENROUTER_KEY` in your environment variables.\n\nOr, run: `export OPENROUTER_KEY=your-api-key`"
		} else {
			msg = "`.env` file loaded successfully. Make sure `OPENROUTER_KEY` is set."
		}

		rendered, err := renderer.Render(msg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(rendered)

		// Create config directory if needed
		os.MkdirAll(revlyDir, 0755)

		// Write a simple flag file
		os.WriteFile(firstRunFile, []byte("shown"), 0644)
	}
}

func main() {
	checkFirstRun()
	cmd.Execute()
}
