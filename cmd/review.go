/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"
	"strings"

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

		out, err := exec.Command("git", "diff", "--cached").Output()
		if err != nil {
			fmt.Println("Error running git diff:", err)
			return
		}
		diff := string(out)
		if strings.TrimSpace(diff) == "" {
			fmt.Println("❌ No staged changes to review.")
			return
		}

		fmt.Println("🤖 Sending to LLM...")

		resp, err := llm.ReviewDiffWithLLM(diff)
		if err != nil {
			fmt.Println("Error from LLM:", err)
			return
		}

		fmt.Println("\n🧠 LLM Review Output:\n")
		fmt.Println(resp)
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
