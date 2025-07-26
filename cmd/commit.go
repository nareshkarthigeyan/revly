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
	push   bool
	dryRun bool
	all    bool
)

var commitCmd = &cobra.Command{
	Use:   "commit [file/folder]",
	Short: "Stage changes, generate commit message via AI, commit and optionally push.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var target string
		if all || len(args) == 0 {
			// default to "." if no specific file/folder provided
			target = "."
		} else {
			target = args[0]
		}

		// Stage the specific target
		if target == "." || all {
			fmt.Println("\tStaging all changes...")
		} else {
			fmt.Printf("\tStaging changes for: %s...\n", target)
		}
		fmt.Printf("\tgit add %s\n", target)
		if _, err := os.Stat(target); os.IsNotExist(err) {
			log.Fatalf("Target %s does not exist: %v", target, err)
		}
		run("git", "add", target)

		// Capture diff of only that staged target
		diff := capture("git", "diff", "--cached", "--", target)
		if diff == "" {
			fmt.Println("No staged changes to commit for:", target)
			return
		}

		msg, err := llm.GetLLMResponse(diff)
		if err != nil {
			log.Fatalf("Error generating commit message: %v", err)
		}
		msg = strings.TrimSpace(msg)

		fmt.Printf("\tgit commit -m:\"%s\"\n", msg)

		if dryRun {
			fmt.Println("\t\033[31mDry run mode: skipping commit and push.\033[0m")
			return
		}
		fmt.Print("\t\nDo you want to commit with this message? (Y/N): ")
		var confirm string
		fmt.Scanln(&confirm)
		if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "Y" {
			fmt.Println("Commit aborted.")
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
	commitCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Only show commit message")
	commitCmd.Flags().BoolVar(&all, "all", false, "Stage all changes (default false)")
	rootCmd.AddCommand(commitCmd)
}

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
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