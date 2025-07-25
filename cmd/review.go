/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/nareshkarthigeyan/revly/internals/gitutils"
	"github.com/nareshkarthigeyan/revly/internals/llm"
	"github.com/spf13/cobra"
)

// reviewCmd represents the review command
var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Fetching git diff...")
		diff, err := gitutils.GetGitDiff()
		if err != nil {
			fmt.Println("❌ Failed to get git diff:", err)
			return
		}

		fmt.Println("🤖 Sending to LLM...")
		review, err := llm.ReviewDiffWithLLM(diff)
		if err != nil {
			fmt.Println("❌ LLM error:", err)
			return
		}

		fmt.Println("✅ Review complete:")
		fmt.Println("──────────────────────────────")
		fmt.Println(review)
		fmt.Println("──────────────────────────────")
	},
}

func init() {
	rootCmd.AddCommand(reviewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reviewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reviewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
