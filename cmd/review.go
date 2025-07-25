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
		commit, _ := cmd.Flags().GetString("commit")
		staged, _ := cmd.Flags().GetBool("staged")

		var diff []byte
		var err error

		switch {
		case commit != "":
			target := commit
			if commit == "true" {
				target = "HEAD"
			}
			fmt.Printf("Fetching diff for commit <%s>...\n", target)
			diff, err = exec.Command("git", "show", target).Output()

		case staged:
			fmt.Println("Fetching staged diff...")
			diff, err = exec.Command("git", "diff", "--cached").Output()

		default:
			fmt.Println("Fetching working directory diff...")
			diff, err = exec.Command("git", "diff").Output()
		}

		if err != nil {
			fmt.Println("Error fetching diff:", err)
			return
		}

		if strings.TrimSpace(string(diff)) == "" {
			fmt.Println("No changes to review.")
			return
		}

		fmt.Println("Sending to AI...")

		resp, err := llm.ReviewDiffWithLLM(string(diff))
		if err != nil {
			fmt.Println("Error from AI:", err)
			return
		}

		fmt.Println("\nAI Review:")
		fmt.Println(resp)
	},
}

func init() {
	rootCmd.AddCommand(reviewCmd)

	reviewCmd.Flags().StringP("commit", "c", "", "Review a specific commit (HEAD if no value given)")
	reviewCmd.Flags().BoolP("staged", "s", false, "Review only staged changes")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reviewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reviewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
