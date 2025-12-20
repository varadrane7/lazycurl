package cmd

import (
	"fmt"
	"lazycurl/internal/curl"
	"lazycurl/internal/model"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [url]",
	Short: "Run a single request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		req := model.NewRequest()
		req.URL = url

		executor := curl.NewExecutor()
		fmt.Printf("Running GET %s...\n", url)

		resp := executor.Execute(req)

		if resp.Error != nil {
			fmt.Printf("Error: %v\n", resp.Error)
			return
		}

		fmt.Printf("Status: %d\n", resp.StatusCode)
		fmt.Printf("Time: %s\n", resp.TimeTaken)
		fmt.Printf("Body:\n%s\n", resp.Body)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
