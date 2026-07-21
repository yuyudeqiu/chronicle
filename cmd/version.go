package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yuyudeqiu/chronicle/internal/buildinfo"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	// Skip root's PersistentPreRun (no DB needed for version)
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		info := buildinfo.Current()
		if jsonOutput {
			printJSON(info)
			return
		}
		fmt.Println(info.String())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
