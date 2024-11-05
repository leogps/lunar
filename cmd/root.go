package cmd

import (
	"github.com/leogps/lunar/pkg/utils"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	// Use:   "yeet-cli [command] [flags]",
	Short: "lunar is a CLI tool to perform calculations for stocks.",
	Long:  `lunar is a CLI tool to perform calculations for stocks.`,
	Run: func(cmd *cobra.Command, _ []string) {
		// This will be executed when no subcommand is specified
		utils.LogInfo("lunar: please pass sub-command. See help for more details:")
		_ = cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		utils.LogError("error occurred", err)
		os.Exit(1)
	}
}
