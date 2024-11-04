package cmd

import (
	"github.com/leogps/lunar/pkg/utils"
	"github.com/leogps/lunar/ui"
	"github.com/spf13/cobra"
	"log/slog"
)

func init() {
	rootCmd.AddCommand(uiCmd)
}

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Starts Terminal UI",
	Long:  `Starts Terminal UI`,
	Run: func(cmd *cobra.Command, _ []string) {
		silent, _ := cmd.Flags().GetBool("silent")
		var level slog.Level
		if silent {
			level = slog.LevelInfo
		} else {
			level = slog.LevelDebug
		}
		utils.InitLogger(level)

		_ = ui.StartApp()
	},
}
