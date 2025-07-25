/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/fatih/color"
	"github.com/nareshkarthigeyan/revly/internals/llm"
	"github.com/spf13/cobra"
)

var severityPatterns = map[*regexp.Regexp]func(string) string{
	regexp.MustCompile(`(?m)$begin:math:display$CRITICAL$end:math:display$`): func(_ string) string {
		return color.New(color.FgRed, color.Bold).Sprint("[CRITICAL]")
	},
	regexp.MustCompile(`(?m)$begin:math:display$WARNING$end:math:display$`): func(_ string) string {
		return color.New(color.FgYellow, color.Bold).Sprint("[WARNING]")
	},
	regexp.MustCompile(`(?m)$begin:math:display$INFO$end:math:display$`): func(_ string) string {
		return color.New(color.FgBlue).Sprint("[INFO]")
	},
}

func highlightSeverities(text string) string {
	for pattern, apply := range severityPatterns {
		text = pattern.ReplaceAllStringFunc(text, apply)
	}
	return text
}

// reviewCmd represents the review command
var reviewCmd = &cobra.Command{
	Use:   "review",
Short: "Run an AI-powered code review on your Git changes in the working directory",
Long: `
	By default, 'review' inspects your working directory diff. You can target specific sources using flags:

	--staged, -s        Review only staged changes (git diff --cached)
	--commit, -c <hash> Review a specific commit by hash
	--head              Review the latest commit (HEAD)

	If no flags are provided, it reviews unstaged changes in your working directory.`,

	Example: `
	
	revly review
    		- Reviews all current unstaged changes in your working directory.

	revly review -s
	revly review --staged
		- Reviews only the files staged for commit.

	revly review -c <commit-hash>
	revly review --commit <commit-hash>
		- Reviews a specific commit by its SHA.

	revly review -c
	revly review --commit
	revly review --head
		- Reviews the most recent commit (HEAD). If -c / --commit is provided without a value, HEAD is assumed.
`,
Run: func(cmd *cobra.Command, args []string) {
			commit, _ := cmd.Flags().GetString("commit")
			staged, _ := cmd.Flags().GetBool("staged")
			head, _ := cmd.Flags().GetBool("head")

			var diff []byte
			var err error

			switch {
			case head:
				color.Cyan("Fetching diff for latest commit (HEAD)...")
				diff, err = exec.Command("git", "show", "HEAD").Output()

			case commit != "":
				color.Cyan("Fetching diff for commit <%s>...", commit)
				diff, err = exec.Command("git", "show", commit).Output()

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


			color.Green("Sending to AI...")

			resp, err := llm.ReviewDiffWithLLM(string(diff))
			if err != nil {
				color.Red("Error from AI: %v", err)
				return
			}

			renderer, err := glamour.NewTermRenderer(
				glamour.WithAutoStyle(), // Auto adapts to terminal theme
			)
			if err != nil {
				log.Fatal(err)
			}

			highlighted := highlightSeverities(resp) // preprocess severity tags before glamour
			rendered, err := renderer.Render(highlighted)
			if err != nil {
				log.Fatal(err)
			}

			coloredOutput := highlightSeverities(rendered)

			color.Green("\n=== AI Review ===")
			fmt.Println(coloredOutput)
			color.Green("=== END OF REVIEW ===")
		},
}

func init() {
	rootCmd.AddCommand(reviewCmd)

	reviewCmd.Flags().StringP("commit", "c", "", "Review a specific commit (HEAD if no value given)")
	reviewCmd.Flags().Lookup("commit").NoOptDefVal = "HEAD"
	reviewCmd.Flags().BoolP("staged", "s", false, "Review only staged changes")
	reviewCmd.Flags().Bool("head", false, "Review the latest commit (HEAD)")
	
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reviewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reviewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
