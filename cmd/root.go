/*
Copyright © 2025 K V Naresh Karthigeyan nareshkarthigeyan.2005@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "revly",
	Short: "AI-powered code review assistant",
	Long:  "Revly is a CLI tool that uses LLMs to analyze git diffs and suggest code improvements.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Revly CLI - AI Code Review Assistant\n\nTry `revly review` to get started.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.revly.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


