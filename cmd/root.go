package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yuyudeqiu/chronicle/internal/service"
)

var rootCmd = &cobra.Command{
	Use:   "chronicle",
	Short: "Chronicle is a task management tool",
	Long:  `Chronicle is a task management tool with a CLI and a web interface.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize database
		service.InitDB("data/app.db")
	},
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}
