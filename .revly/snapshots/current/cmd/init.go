package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/nareshkarthigeyan/revly/internals/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Revly config file",
	Long: `Scaffolds a default Revly configuration file.

This includes:
- LLM base URL set to OpenRouter.ai
- Adds a list of free default models for different tasks (review, commit names, etc.)
- Sets basic developer preferences like showDiff, pushOnCommit, etc.

You can customize everything in the generated config file.`,

	Run: func(cmd *cobra.Command, args []string) {
		createDefaultConfig()
	},
}

func createDefaultConfig() {
	configPaths := []string{"revly.config.toml"}

	fmt.Println("Revly: Initializing configuration...")

	for _, path := range configPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("Creating config at '%s'...\n", path)

			err := os.WriteFile(path, []byte(config.DefaultConfig), 0644)
			if err != nil {
				fmt.Printf("Failed to create %s: %v\n", path, err)
				return
			}

			fmt.Println("Done.")
			fmt.Println()
			fmt.Println("	Base URL set to: https://openrouter.ai/api/v1")
			fmt.Println("	Added default models for:")
			fmt.Println("   → Added default free OpenRouter Models for config:")
			fmt.Println("      - qwen/qwen3-coder:free")
			fmt.Println("      - qwen/qwen3-235b-a22b-2507:free")
			fmt.Println("      - moonshotai/kimi-k2:free")
			fmt.Println("      - cognitivecomputations/dolphin-mistral-24b-venice-edition:free")
			fmt.Println("      - tngtech/deepseek-r1t2-chimera:free")
			fmt.Println("      - moonshotai/kimi-dev-72b:free")
			fmt.Println("      - deepseek/deepseek-r1-0528-qwen3-8b:free")
			fmt.Println("      - tencent/hunyuan-a13b-instruct:free")
			fmt.Println("      - mistralai/mistral-small-3.2-24b-instruct:free")
			fmt.Println("      - deepseek/deepseek-r1-0528:free")
			fmt.Println("      - tngtech/deepseek-r1t-chimera:free")
			fmt.Println("      - microsoft/mai-ds-r1:free")
			fmt.Println("      - moonshotai/kimi-vl-a3b-thinking:free")
			fmt.Println("      - nvidia/llama-3.1-nemotron-ultra-253b-v1:free")
			fmt.Println("   Edit them to add your preferred models.")

			fmt.Println()
			fmt.Println("[Tip]: You can now edit this config to switch to custom endpoints, override model preferences, or tweak behavior flags like:")
			fmt.Println("   - provider")
			fmt.Println("   - base_url")
			fmt.Println("   - models")
			fmt.Println("   - showDiff")
			fmt.Println("   - pushOnCommitByDefault")
			fmt.Println("   - logLevel")
			fmt.Printf("\nTo customize, open '%s' in your editor.\n", path)
			return
		}
	}

	fmt.Println("A config file already exists. Skipping creation.")
	fmt.Println("To regenerate, delete the file and rerun `revly init`.")
}

func init() {
	rootCmd.AddCommand(initCmd)
}