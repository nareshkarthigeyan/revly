/*
Copyright © 2025 K V Naresh Karthigeyan nareshkarthigeyan.2005@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var Version = "0.1.0" // or inject via ldflags
// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "revly",
	Short: "Revly is an AI-powered code review CLI tool",
	Long:  "Revly is a CLI tool that uses LLMs to analyze \ngit diffs and suggest code improvements,\nreview both staged and unstaged changes,\nand provide actionable feedback on code quality.\nDo not wait for PR reviews, get instant feedback on your code changes.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Revly CLI - AI Code Review Assistant\n\nTry `revly review` to get started. \nFor more help, use `revly --help`.")
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

}


