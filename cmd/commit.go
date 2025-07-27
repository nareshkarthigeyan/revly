package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/nareshkarthigeyan/revly/internals/config"
	"github.com/nareshkarthigeyan/revly/internals/llm"
	"github.com/spf13/cobra"
)

var (
	push    bool
	dryRun  bool
	all     bool
	message string
)

var commitCmd = &cobra.Command{
	Use:   "commit [file/folder]",
	Short: "Stage changes, generate commit message via AI or custom input, commit and optionally push.",
	Long: `
'commit' is a convenient command to streamline your Git workflow by leveraging AI to auto-generate meaningful commit messages.

By default, it stages and commits all changes in the current directory. You can also specify a particular file or folder to commit.

This command:
- Stages the given path (or all changes with --all).
- Captures the Git diff of the staged changes.
- Sends the diff to an AI model to generate a descriptive commit message (unless -m is specified).
- Asks for your confirmation before committing (unless -m is specified).
- Optionally pushes to GitHub if --push is passed.

Flags:
--all, -a            Stage all changes recursively from the project root.
--push               Push the commit to the remote repository after committing.
--dry-run            Show the AI-generated commit message without actually committing or pushing.
--message, -m        Use a custom commit message instead of AI-generated one (disables confirmation).
`,
	Example: `
revly commit
	- Stages all changes (same as '.'), generates an AI commit message, and prompts for confirmation.

revly commit --all
	- Explicitly stages all changes and proceeds with the same workflow.

revly commit src/utils/
	- Stages only the 'src/utils' folder, generates an AI commit message for it, and commits.

revly commit --dry-run
	- Shows what the AI would generate as a commit message, but doesn’t actually commit anything.

revly commit --all --push
	- Stages everything, commits using an AI-generated message, and pushes it to the remote repo.

revly commit src/index.ts --push
	- Targets a specific file, commits it with an AI message, and pushes the result upstream.

revly commit -m "Fix typo in README"
	- Stages default ('.'), commits immediately with the given message without prompting or AI generation.
`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.GetConfig()
		if err != nil {
			// // DEBUG: Error getting config
			// fmt.Println("DEBUG: Error getting config:", err)
			return
		}

		var target string
		if all || len(args) == 0 {
			target = "."
		} else {
			target = args[0]
		}

		if _, err := os.Stat(target); os.IsNotExist(err) {
			log.Fatalf("Target %s does not exist: %v", target, err)
		}

		if target == "." || all {
			fmt.Println("\tStaging all changes...")
		} else {
			fmt.Printf("\tStaging changes for: %s...\n", target)
		}
		fmt.Printf("\tgit add %s\n", target)
		run("git", "add", target)

		diff := capture("git", "diff", "--cached", "--", target)
		if diff == "" {
			fmt.Println("No staged changes to commit for:", target)
			return
		}

		var msg string
		if message != "" {
			msg = message
		} else {
			aiMsg, err := llm.GetLLMResponse(diff)
			if err != nil {
				log.Fatalf("Error generating commit message: %v", err)
			}
			msg = strings.TrimSpace(aiMsg)
		}

		fmt.Printf("\tgit commit -m \"%s\"\n", msg)

		if dryRun {
			fmt.Println("\t\033[31mDry run mode: skipping commit and push.\033[0m")
			return
		}

		if message == "" {
			fmt.Print("\t\nDo you want to commit with this message? (Y/N): ")
			var confirm string
			fmt.Scanln(&confirm)
			if strings.ToLower(confirm) != "y" {
				fmt.Println("Commit aborted.")
				return
			}
		}

		run("git", "commit", "-m", msg)

		if cfg.Git.PushOnCommit || push {
		cmd := exec.Command("git", "push")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return
		}
		fmt.Println("\tPushed changes to remote repository.")
		} else {
			return
		}
	},
}

func init() {
	commitCmd.Flags().BoolVar(&push, "push", false, "Push after committing")
	commitCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Only show commit message")
	commitCmd.Flags().BoolVar(&all, "all", false, "Stage all changes (default false)")
	commitCmd.Flags().StringVarP(&message, "message", "m", "", "Use custom commit message (bypasses AI and confirmation)")
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