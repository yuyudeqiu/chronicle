package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yuyudeqiu/chronicle/internal/service"
)

var archiveCmd = &cobra.Command{
	Use:   "archive <id>",
	Short: "Archive a task without deleting it",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := service.GetTask(args[0]); err != nil {
			return err
		}
		if err := service.ArchiveTask(args[0]); err != nil {
			return err
		}
		if jsonOutput {
			return printJSON(cmd, map[string]string{"id": args[0], "status": "archived"})
		} else {
			fmt.Printf("Task archived: %s\n", args[0])
		}
		return nil
	},
}

var unarchiveCmd = &cobra.Command{
	Use:   "unarchive <id>",
	Short: "Restore an archived task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := service.GetTask(args[0]); err != nil {
			return err
		}
		if err := service.UnarchiveTask(args[0]); err != nil {
			return err
		}
		if jsonOutput {
			return printJSON(cmd, map[string]string{"id": args[0], "status": "unarchived"})
		} else {
			fmt.Printf("Task unarchived: %s\n", args[0])
		}
		return nil
	},
}

var deleteWorklogCmd = &cobra.Command{
	Use:   "delete-worklog <id>",
	Short: "Delete a worklog by its ID (shown by get)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := service.DeleteWorklog(args[0]); err != nil {
			return err
		}
		if jsonOutput {
			return printJSON(cmd, map[string]string{"id": args[0], "status": "deleted"})
		} else {
			fmt.Printf("Worklog deleted: %s\n", args[0])
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(archiveCmd, unarchiveCmd, deleteWorklogCmd)
}
