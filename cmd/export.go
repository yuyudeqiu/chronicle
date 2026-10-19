package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/yuyudeqiu/chronicle/internal/config"
)

var exportCmd = &cobra.Command{
	Use:   "export [output-path]",
	Short: "Export the entire database to a file",
	Long: `Export the SQLite database file to the specified path.
If no path is given, exports to chronicle-export-{timestamp}.db in the current directory.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dbPath := config.GetDBPath()

		// Check source exists
		if _, err := os.Stat(dbPath); err != nil {
			return fmt.Errorf("stat database: %w", err)
		}

		// Determine output path
		outputPath := ""
		if len(args) > 0 {
			outputPath = args[0]
		} else {
			timestamp := time.Now().Format("20060102_150405")
			outputPath = fmt.Sprintf("chronicle-export-%s.db", timestamp)
		}

		// Ensure parent directory exists
		outputDir := filepath.Dir(outputPath)
		if outputDir != "." {
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("cannot create directory %s: %w", outputDir, err)
			}
		}

		// Copy file
		sourceFile, err := os.Open(dbPath)
		if err != nil {
			return fmt.Errorf("cannot open source database: %w", err)
		}
		defer func() { _ = sourceFile.Close() }()

		destFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("cannot create output file: %w", err)
		}
		defer func() { _ = destFile.Close() }()

		written, err := io.Copy(destFile, sourceFile)
		if err != nil {
			return fmt.Errorf("copy failed: %w", err)
		}

		if err := destFile.Close(); err != nil {
			return fmt.Errorf("close output file: %w", err)
		}

		fmt.Printf("✅ Database exported successfully\n")
		fmt.Printf("   Source: %s (%d bytes)\n", dbPath, written)
		fmt.Printf("   Dest:   %s\n", outputPath)
		return nil
	},
}

var importCmd = &cobra.Command{
	Use:   "import <source-path>",
	Short: "Import a database file, replacing the current data",
	Long: `Import a previously exported SQLite database file, replacing all current data.
WARNING: This will overwrite the existing database. All current data will be lost.
Use --force to skip confirmation.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sourcePath := args[0]
		dbPath := config.GetDBPath()

		// Check source exists
		if _, err := os.Stat(sourcePath); err != nil {
			return fmt.Errorf("stat import source: %w", err)
		}

		force, _ := cmd.Flags().GetBool("force")

		if !force {
			fmt.Printf("⚠️  WARNING: This will REPLACE the current database at:\n")
			fmt.Printf("   %s\n\n", dbPath)
			fmt.Printf("   All current data will be lost!\n")
			fmt.Printf("   Use --force to skip this warning.\n\n")
			fmt.Printf("   Are you sure? (yes/no): ")

			var response string
			if _, err := fmt.Fscanln(cmd.InOrStdin(), &response); err != nil {
				return fmt.Errorf("read import confirmation: %w", err)
			}
			if response != "yes" {
				fmt.Println("Import cancelled.")
				return nil
			}
		}

		// Create backup of current database before import
		backupPath := dbPath + ".pre-import-backup"
		srcFile, err := os.Open(dbPath)
		if err == nil {
			defer func() { _ = srcFile.Close() }()
			bakFile, err := os.Create(backupPath)
			if err == nil {
				if _, copyErr := io.Copy(bakFile, srcFile); copyErr == nil {
					fmt.Printf("   Auto-backup saved to: %s\n", backupPath)
				}
				_ = bakFile.Close()
			}
		}

		// Copy source to database path
		sourceFile, err := os.Open(sourcePath)
		if err != nil {
			return fmt.Errorf("cannot open source file: %w", err)
		}
		defer func() { _ = sourceFile.Close() }()

		destFile, err := os.Create(dbPath)
		if err != nil {
			return fmt.Errorf("cannot write to database path: %w", err)
		}
		defer func() { _ = destFile.Close() }()

		written, err := io.Copy(destFile, sourceFile)
		if err != nil {
			return fmt.Errorf("import failed: %w", err)
		}

		if err := destFile.Close(); err != nil {
			return fmt.Errorf("close imported database: %w", err)
		}

		fmt.Printf("\n✅ Database imported successfully\n")
		fmt.Printf("   Source: %s (%d bytes)\n", sourcePath, written)
		fmt.Printf("   Dest:   %s\n", dbPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd, importCmd)
	importCmd.Flags().Bool("force", false, "Skip confirmation prompt")
}
