package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yuyudeqiu/chronicle/internal/model"
	"github.com/yuyudeqiu/chronicle/internal/service"
)

// Task Flags
var (
	category string
	desc     string
	links    string
	targets  string
	deadline string
	status   string
	title    string
)

var createCmd = &cobra.Command{
	Use:   "create <title>",
	Short: "Create a new task",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := args[0]
		deadlineTime := parseDeadline(deadline)

		req := model.CreateTaskReq{
			Title:       title,
			Category:    category,
			Description: desc,
			Targets:     targets,
			Links:       links,
			Deadline:    deadlineTime,
		}

		task, err := service.CreateTask(req)
		if err != nil {
			return err
		}

		if jsonOutput {
			return printJSON(cmd, task)
		} else {
			fmt.Printf("Task created: %s\n", task.ID)
			printTask(task)
		}
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list [status|archived]",
	Args:  cobra.MaximumNArgs(1),
	Short: "List tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		var tasks []model.ActiveTaskResp
		var err error

		if len(args) > 0 {
			queryStatus := args[0]
			switch queryStatus {
			case "archived":
				tasks, err = service.GetArchivedTasks()
			case "done":
				tasks, err = service.GetHistoryTasks()
			default:
				tasks, err = service.GetActiveTasks()
			}
		} else {
			tasks, err = service.GetActiveTasks()
		}

		if err != nil {
			return err
		}

		if len(tasks) == 0 {
			if jsonOutput {
				return printJSON(cmd, []model.ActiveTaskResp{})
			} else {
				fmt.Println("No tasks found")
			}
			return nil
		}

		if jsonOutput {
			return printJSON(cmd, tasks)
		} else {
			fmt.Printf("Found %d tasks:\n\n", len(tasks))
			for _, t := range tasks {
				fmt.Printf("  [%s] %s - %s\n", t.Status, t.Title, t.Category)
				fmt.Printf("    ID: %s\n\n", t.ID)
			}
		}
		return nil
	},
}

var getCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get task details by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		task, err := service.GetTask(taskID)
		if err != nil {
			return err
		}

		if jsonOutput {
			return printJSON(cmd, task)
		} else {
			printTask(task)

			if len(task.Logs) > 0 {
				fmt.Println("\nWorklogs:")
				for _, log := range task.Logs {
					fmt.Printf("  [%s] %s\n    ID: %s\n", log.CreatedAt.Format("2006-01-02 15:04"), log.LogText, log.ID)
				}
			}
		}
		return nil
	},
}

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update task by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		deadlineTime := parseDeadline(deadline)

		if status != "" {
			progressReq := model.UpdateProgressReq{
				NewStatus: status,
			}
			err := service.UpdateProgress(taskID, progressReq)
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(cmd, map[string]string{"id": taskID, "status": "updated"})
			} else {
				fmt.Printf("Task status updated: %s\n", taskID)
			}
			return nil
		}

		req := model.UpdateTaskReq{
			Title:       title,
			Category:    category,
			Description: desc,
			Targets:     targets,
			Links:       links,
			Deadline:    deadlineTime,
		}

		task, err := service.UpdateTask(taskID, req)
		if err != nil {
			return err
		}

		if jsonOutput {
			return printJSON(cmd, task)
		} else {
			fmt.Printf("Task updated: %s\n", taskID)
			printTask(task)
		}
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete task by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		err := service.DeleteTask(taskID)
		if err != nil {
			return err
		}

		if jsonOutput {
			return printJSON(cmd, map[string]string{"id": taskID, "status": "deleted"})
		} else {
			fmt.Printf("Task deleted: %s\n", taskID)
		}
		return nil
	},
}

var logCmd = &cobra.Command{
	Use:   "log <id> <message>",
	Short: "Add worklog to task",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		logText := strings.Join(args[1:], " ")

		req := model.UpdateProgressReq{
			LogText: logText,
		}

		err := service.UpdateProgress(taskID, req)
		if err != nil {
			return err
		}

		if jsonOutput {
			return printJSON(cmd, map[string]string{"id": taskID, "status": "worklog added"})
		} else {
			fmt.Printf("Worklog added to task: %s\n", taskID)
		}
		return nil
	},
}

var summaryCmd = &cobra.Command{
	Use:   "summary [date]",
	Short: "Get daily summary (default: today)",
	RunE: func(cmd *cobra.Command, args []string) error {
		dateStr := ""
		if len(args) > 0 {
			dateStr = args[0]
		}

		summary, err := service.GetDailySummary(dateStr)
		if err != nil {
			return err
		}

		if jsonOutput {
			return printJSON(cmd, summary)
		} else {
			fmt.Printf("Daily Summary for %s\n\n", summary.Date)

			if len(summary.Activities) == 0 {
				fmt.Println("No activities today")
				return nil
			}

			// Group by category
			categoryMap := make(map[string][]model.DailySummaryActivity)
			for _, a := range summary.Activities {
				categoryMap[a.Category] = append(categoryMap[a.Category], a)
			}

			for cat, activities := range categoryMap {
				fmt.Printf("### %s\n", cat)
				for _, a := range activities {
					fmt.Printf("%s %s\n", taskStatusIcon(a.Status), a.TaskTitle)
					for _, log := range a.TodayLogs {
						fmt.Printf("   - %s\n", log)
					}
				}
				fmt.Println()
			}
		}
		return nil
	},
}

var weeklySummaryCmd = &cobra.Command{
	Use:   "weekly-summary [date]",
	Short: "Get weekly summary (default: current week, date picks the containing week)",
	RunE: func(cmd *cobra.Command, args []string) error {
		dateStr := ""
		if len(args) > 0 {
			dateStr = args[0]
		}

		summary, err := service.GetWeeklySummary(dateStr)
		if err != nil {
			return err
		}

		if jsonOutput {
			return printJSON(cmd, summary)
		} else {
			fmt.Printf("Weekly Summary: %s ~ %s\n\n", summary.WeekStart, summary.WeekEnd)

			if len(summary.Activities) == 0 {
				fmt.Println("No activities this week")
				return nil
			}

			// Group by category
			categoryMap := make(map[string][]model.WeeklySummaryActivity)
			for _, a := range summary.Activities {
				categoryMap[a.Category] = append(categoryMap[a.Category], a)
			}

			for cat, activities := range categoryMap {
				fmt.Printf("### %s\n", cat)
				for _, a := range activities {
					fmt.Printf("%s %s\n", taskStatusIcon(a.Status), a.TaskTitle)
					for _, log := range a.WeekLogs {
						fmt.Printf("   - %s\n", log)
					}
				}
				fmt.Println()
			}

			fmt.Printf("=== Stats ===\n")
			fmt.Printf("  Total: %d | ✅ Done: %d | 🔄 In Progress: %d | 📝 Todo: %d | 🚫 Blocked: %d\n",
				summary.Stats.TotalTasks, summary.Stats.CompletedTasks,
				summary.Stats.InProgressTasks, summary.Stats.TodoTasks, summary.Stats.BlockedTasks)
		}
		return nil
	},
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get task statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		stats, err := service.GetStatsSummary()
		if err != nil {
			return err
		}

		if jsonOutput {
			return printJSON(cmd, stats)
		} else {
			fmt.Println("=== Task Statistics ===")
			fmt.Printf("Total Tasks: %d\n", stats.TotalTasks)
			fmt.Printf("Completed: %d\n", stats.CompletedTasks)
			fmt.Printf("In Progress: %d\n", stats.InProgressTasks)
			fmt.Printf("Blocked: %d\n", stats.BlockedTasks)
			fmt.Printf("Todo: %d\n", stats.TodoTasks)
			fmt.Printf("Completion Rate: %.1f%%\n", stats.CompletionRate*100)

			fmt.Println("\n=== By Category ===")
			for cat, count := range stats.ByCategory {
				fmt.Printf("  %s: %d\n", cat, count)
			}

			fmt.Println("\n=== Weekly Stats ===")
			for _, s := range stats.WeeklyStats {
				fmt.Printf("  %s: created=%d, completed=%d\n", s.Date, s.Created, s.Completed)
			}
		}
		return nil
	},
}

// Helpers
func taskStatusIcon(status string) string {
	switch status {
	case model.TaskStatusTodo:
		return "📝"
	case model.TaskStatusInProgress:
		return "🔄"
	case model.TaskStatusBlocked:
		return "🚫"
	case model.TaskStatusDone:
		return "✅"
	default:
		return ""
	}
}

func parseDeadline(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		fmt.Printf("Warning: invalid deadline format, ignoring: %s\n", s)
		return nil
	}
	return &t
}

func printTask(task *model.Task) {
	fmt.Println("\nTask Details:")
	fmt.Printf("  ID: %s\n", task.ID)
	fmt.Printf("  Title: %s\n", task.Title)
	fmt.Printf("  Category: %s\n", task.Category)
	fmt.Printf("  Status: %s\n", task.Status)
	if task.Deadline != nil {
		fmt.Printf("  Deadline: %s\n", task.Deadline.Format("2006-01-02 15:04"))
	}
	if task.Description != "" {
		fmt.Printf("  Description: %s\n", task.Description)
	}
	if task.Links != "" {
		fmt.Printf("  Links:\n")
		for _, link := range strings.Split(task.Links, "\n") {
			fmt.Printf("    - %s\n", link)
		}
	}
}

func init() {
	// Add subcommands to rootCmd
	rootCmd.AddCommand(createCmd, listCmd, getCmd, updateCmd, deleteCmd, logCmd, summaryCmd, weeklySummaryCmd, statsCmd)

	// Local flags for create and update
	createCmd.Flags().StringVarP(&category, "category", "c", "", "Task category")
	createCmd.Flags().StringVarP(&desc, "desc", "d", "", "Task description")
	createCmd.Flags().StringVarP(&links, "links", "l", "", "Task links (one per line)")
	createCmd.Flags().StringVarP(&targets, "target", "t", "", "Task targets")
	createCmd.Flags().StringVar(&deadline, "deadline", "", "Deadline (ISO8601 format)")

	updateCmd.Flags().StringVarP(&category, "category", "c", "", "Task category")
	updateCmd.Flags().StringVar(&title, "title", "", "Task title")
	updateCmd.Flags().StringVarP(&desc, "desc", "d", "", "Task description")
	updateCmd.Flags().StringVarP(&links, "links", "l", "", "Task links (one per line)")
	updateCmd.Flags().StringVarP(&targets, "target", "t", "", "Task targets")
	updateCmd.Flags().StringVar(&deadline, "deadline", "", "Deadline (ISO8601 format)")
	updateCmd.Flags().StringVar(&status, "new-status", "", "New status")
}
