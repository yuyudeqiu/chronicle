package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yuyudeqiu/chronicle/internal/config"
	"github.com/yuyudeqiu/chronicle/internal/service"
)

var jsonOutput bool
var dataDir string

var rootCmd = &cobra.Command{
	Use:           "chronicle",
	Short:         "Chronicle is a task management tool",
	Long:          `Chronicle is a task management tool with a CLI and an HTTP API.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 初始化配置（命令行参数优先级最高）
		if dataDir != "" {
			config.DataDir = dataDir
		}
		if err := config.EnsureDataDir(); err != nil {
			return fmt.Errorf("create data directory: %w", err)
		}
		// 初始化数据库
		return service.InitDB(config.GetDBPath())
	},
}

// JSON output helper
func printJSON(cmd *cobra.Command, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}
	if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(data)); err != nil {
		return fmt.Errorf("write JSON: %w", err)
	}
	return nil
}

// Execute executes the root command and returns any error to the entry point.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "o", false, "Output in JSON format")
	rootCmd.PersistentFlags().StringVar(&dataDir, "data-dir", "", "Data directory (default: data/, or use CHRONICLE_DATA_DIR env var)")
}
