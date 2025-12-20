package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lazycurl",
	Short: "A friendly TUI for curl",
	Long:  `LazyCurl is a terminal UI for API exploration and testing, inspired by lazygit.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default to running TUI if no command is passed
		tuiCmd.Run(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
