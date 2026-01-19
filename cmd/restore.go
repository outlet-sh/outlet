package cmd

import (
	"compress/gzip"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/outlet-sh/outlet/internal/config"

	_ "modernc.org/sqlite"
)

var (
	restoreForce    bool
	restoreBackupOriginal bool
)

var restoreCmd = &cobra.Command{
	Use:   "restore [backup-file]",
	Short: "Restore database from a backup",
	Long: `Restore the SQLite database from a backup file.

IMPORTANT: This will replace the current database with the backup.
The server should be stopped before running this command.

Examples:
  outlet restore ./backups/outlet-20240115-120000.db
  outlet restore ./backups/outlet.db.gz              # Restore from compressed backup
  outlet restore --force backup.db                   # Skip confirmation prompt`,
	Args: cobra.ExactArgs(1),
	RunE: runRestore,
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().BoolVarP(&restoreForce, "force", "f", false, "Skip confirmation prompt")
	restoreCmd.Flags().BoolVar(&restoreBackupOriginal, "backup-original", true, "Backup current database before restore")
}

func runRestore(cmd *cobra.Command, args []string) error {
	backupPath := args[0]

	// Load .env file
	godotenv.Load()

	var c config.Config
	if err := conf.LoadFromYamlBytes([]byte(os.ExpandEnv(string(EmbeddedConfig))), &c); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	dbPath := c.Database.Path
	if dbPath == "" {
		dbPath = "./data/outlet.db"
	}

	// Check if backup file exists
	backupInfo, err := os.Stat(backupPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupPath)
	}
	if err != nil {
		return fmt.Errorf("failed to stat backup file: %w", err)
	}

	fmt.Printf("Database Restore\n")
	fmt.Printf("================\n")
	fmt.Printf("  Backup file: %s (%s)\n", backupPath, formatBytes(backupInfo.Size()))
	fmt.Printf("  Target database: %s\n", dbPath)

	// Check if target database exists
	if _, err := os.Stat(dbPath); err == nil {
		fmt.Printf("  Current database exists: YES\n")
	} else {
		fmt.Printf("  Current database exists: NO (will create new)\n")
	}

	// Confirm restore
	if !restoreForce {
		fmt.Print("\nWARNING: This will REPLACE your current database!\n")
		fmt.Print("Continue with restore? [y/N]: ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("Restore cancelled.")
			return nil
		}
	}

	fmt.Println()

	// Verify backup file is a valid SQLite database
	fmt.Print("Verifying backup file... ")
	isCompressed := strings.HasSuffix(backupPath, ".gz")

	var verifyPath string
	if isCompressed {
		// Decompress to temp file for verification
		tmpFile, err := os.CreateTemp("", "outlet-restore-*.db")
		if err != nil {
			fmt.Println("FAILED")
			return fmt.Errorf("failed to create temp file: %w", err)
		}
		tmpFile.Close()
		verifyPath = tmpFile.Name()
		defer os.Remove(verifyPath)

		if err := decompressFile(backupPath, verifyPath); err != nil {
			fmt.Println("FAILED")
			return fmt.Errorf("failed to decompress backup: %w", err)
		}
	} else {
		verifyPath = backupPath
	}

	if err := verifySQLiteDatabase(verifyPath); err != nil {
		fmt.Println("FAILED")
		return fmt.Errorf("backup verification failed: %w", err)
	}
	fmt.Println("OK")

	// Backup current database if it exists
	if restoreBackupOriginal {
		if _, err := os.Stat(dbPath); err == nil {
			fmt.Print("Backing up current database... ")
			timestamp := time.Now().Format("20060102-150405")
			backupCurrentPath := dbPath + ".pre-restore-" + timestamp
			if err := copyFile(dbPath, backupCurrentPath); err != nil {
				fmt.Println("FAILED")
				return fmt.Errorf("failed to backup current database: %w", err)
			}
			fmt.Printf("OK (%s)\n", filepath.Base(backupCurrentPath))
		}
	}

	// Create target directory if needed
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Perform restore
	fmt.Print("Restoring database... ")
	start := time.Now()

	if isCompressed {
		if err := decompressFile(backupPath, dbPath); err != nil {
			fmt.Println("FAILED")
			return fmt.Errorf("restore failed: %w", err)
		}
	} else {
		if err := copyFile(backupPath, dbPath); err != nil {
			fmt.Println("FAILED")
			return fmt.Errorf("restore failed: %w", err)
		}
	}
	fmt.Println("OK")

	// Verify restored database
	fmt.Print("Verifying restored database... ")
	if err := verifySQLiteDatabase(dbPath); err != nil {
		fmt.Println("FAILED")
		return fmt.Errorf("restored database verification failed: %w", err)
	}

	// Run integrity check
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Println("FAILED")
		return fmt.Errorf("failed to open restored database: %w", err)
	}
	defer db.Close()

	var result string
	if err := db.QueryRow("PRAGMA integrity_check").Scan(&result); err != nil {
		fmt.Println("FAILED")
		return fmt.Errorf("integrity check failed: %w", err)
	}
	if result != "ok" {
		fmt.Printf("WARNING (%s)\n", result)
	} else {
		fmt.Println("OK")
	}

	elapsed := time.Since(start)
	restoredInfo, _ := os.Stat(dbPath)

	fmt.Printf("\nRestore completed in %v\n", elapsed.Round(time.Millisecond))
	fmt.Printf("  Database size: %s\n", formatBytes(restoredInfo.Size()))
	fmt.Println("\nYou can now start the outlet server.")

	return nil
}

// verifySQLiteDatabase checks if a file is a valid SQLite database
func verifySQLiteDatabase(path string) error {
	db, err := sql.Open("sqlite", path+"?mode=ro")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Try to read SQLite version to verify it's a valid database
	var version string
	if err := db.QueryRow("SELECT sqlite_version()").Scan(&version); err != nil {
		return fmt.Errorf("not a valid SQLite database: %w", err)
	}

	return nil
}

// decompressFile decompresses a gzip file
func decompressFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	gr, err := gzip.NewReader(src)
	if err != nil {
		return err
	}
	defer gr.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, gr)
	return err
}

// copyFile copies a file
func copyFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
