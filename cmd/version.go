package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Revly",
	Long:  "All software has versions. This is Revly's.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Revly v%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}