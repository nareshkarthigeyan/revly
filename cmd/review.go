/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os/exec"
	"strings"

	"github.com/fatih/color"
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
				color.Cyan("Fetching diff for commit <%s>...", target)
				diff, err = exec.Command("git", "show", target).Output()

			case staged:
				color.Cyan("Fetching staged diff...")
				diff, err = exec.Command("git", "diff", "--cached").Output()

			default:
				color.Cyan("Fetching working directory diff...")
				diff, err = exec.Command("git", "diff").Output()
			}

			if err != nil {
				color.Red("Error fetching diff: %v", err)
				return
			}

			if strings.TrimSpace(string(diff)) == "" {
				color.Yellow("No changes to review.")
				return
			}

			color.Green("🤖 Sending to AI...")

			resp, err := llm.ReviewDiffWithLLM(string(diff))
			if err != nil {
				color.Red("Error from AI: %v", err)
				return
			}

			color.Green("\n=== AI Review ===")
			color.White(resp)
			color.Green("\n=== END OF REVIEW ===")
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
