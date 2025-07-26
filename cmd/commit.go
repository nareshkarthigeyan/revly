// cmd/commit.go
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/nareshkarthigeyan/revly/internals/llm"
	"github.com/spf13/cobra"
)

var (
	push    bool
	dryRun  bool
	all     bool
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Stage changes, review diffs, commit with relevant message and push to remote automatically.",
	Run: func(cmd *cobra.Command, args []string) {
		if all {
			run("git", "add", ".")
		}

		diff := capture("git", "diff", "--cached")
		if diff == "" {
			fmt.Println("No staged changes to commit.")
			return
		}

		prompt := fmt.Sprintf(diff)
		msg, err := llm.GetLLMResponse(prompt)
		if err != nil {
			log.Fatalf("Error generating commit message: %v", err)
		}
		msg = strings.TrimSpace(msg)

		fmt.Printf("Suggested commit message:\n \"%s\"\n\n", msg)

		if dryRun {
			fmt.Println("Dry run mode: skipping commit and push.")
			return
		}

		run("git", "commit", "-m", msg)

		if push {
			run("git", "push")
		}
	},
}

func init() {
	commitCmd.Flags().BoolVar(&push, "push", false, "Push after committing")
	commitCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Only show the commit message, don't commit")
	commitCmd.Flags().BoolVar(&all, "all", true, "Stage all changes (default: true)")
	rootCmd.AddCommand(commitCmd)
}


func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Command failed: %s %v\n%v", name, args, err)
	}
}

func capture(name string, args ...string) string {
	var out bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
	return out.String()
}