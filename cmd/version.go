package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yuyudeqiu/chronicle/internal/buildinfo"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	// Skip database initialization for version information.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	RunE: func(cmd *cobra.Command, args []string) error {
		info := buildinfo.Current()
		if jsonOutput {
			return printJSON(cmd, info)
		}
		fmt.Println(info.String())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
