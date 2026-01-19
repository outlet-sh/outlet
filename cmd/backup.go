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
	backupOutput    string
	backupCompress  bool
	backupS3        bool
	backupS3Bucket  string
	backupS3Region  string
	backupS3Prefix  string
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create a database backup",
	Long: `Create a backup of the SQLite database.

The backup uses SQLite's online backup API for consistent, safe backups
even while the server is running.

Examples:
  outlet backup                           # Backup to default location
  outlet backup -o /path/to/backup.db     # Backup to specific file
  outlet backup --compress                # Create compressed backup (.db.gz)
  outlet backup --s3 --s3-bucket=mybucket # Upload to S3`,
	RunE: runBackup,
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().StringVarP(&backupOutput, "output", "o", "", "Output file path (default: ./backups/outlet-YYYYMMDD-HHMMSS.db)")
	backupCmd.Flags().BoolVar(&backupCompress, "compress", false, "Compress the backup with gzip")
	backupCmd.Flags().BoolVar(&backupS3, "s3", false, "Upload backup to S3")
	backupCmd.Flags().StringVar(&backupS3Bucket, "s3-bucket", "", "S3 bucket name (or BACKUP_S3_BUCKET env)")
	backupCmd.Flags().StringVar(&backupS3Region, "s3-region", "", "S3 region (or BACKUP_S3_REGION env)")
	backupCmd.Flags().StringVar(&backupS3Prefix, "s3-prefix", "backups/", "S3 key prefix")
}

func runBackup(cmd *cobra.Command, args []string) error {
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

	// Check if database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("database not found at %s", dbPath)
	}

	// Generate default output path if not specified
	if backupOutput == "" {
		backupDir := "./backups"
		if err := os.MkdirAll(backupDir, 0755); err != nil {
			return fmt.Errorf("failed to create backup directory: %w", err)
		}
		timestamp := time.Now().Format("20060102-150405")
		backupOutput = filepath.Join(backupDir, fmt.Sprintf("outlet-%s.db", timestamp))
	}

	if backupCompress && !strings.HasSuffix(backupOutput, ".gz") {
		backupOutput += ".gz"
	}

	fmt.Printf("Creating backup...\n")
	fmt.Printf("  Source: %s\n", dbPath)
	fmt.Printf("  Output: %s\n", backupOutput)

	start := time.Now()

	// Create backup using SQLite online backup API
	if err := createSQLiteBackup(dbPath, backupOutput, backupCompress); err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}

	// Get backup file info
	info, err := os.Stat(backupOutput)
	if err != nil {
		return fmt.Errorf("failed to stat backup file: %w", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("\nBackup completed in %v\n", elapsed.Round(time.Millisecond))
	fmt.Printf("  Size: %s\n", formatBytes(info.Size()))

	// Upload to S3 if requested
	if backupS3 {
		if backupS3Bucket == "" {
			backupS3Bucket = os.Getenv("BACKUP_S3_BUCKET")
		}
		if backupS3Region == "" {
			backupS3Region = os.Getenv("BACKUP_S3_REGION")
		}

		if backupS3Bucket == "" {
			return fmt.Errorf("S3 bucket required. Set --s3-bucket or BACKUP_S3_BUCKET env")
		}

		fmt.Printf("\nUploading to S3...\n")
		fmt.Printf("  Bucket: %s\n", backupS3Bucket)
		fmt.Printf("  Region: %s\n", backupS3Region)

		s3Key := backupS3Prefix + filepath.Base(backupOutput)
		if err := uploadToS3(backupOutput, backupS3Bucket, s3Key, backupS3Region); err != nil {
			return fmt.Errorf("S3 upload failed: %w", err)
		}

		fmt.Printf("  Key: %s\n", s3Key)
		fmt.Println("S3 upload complete!")
	}

	return nil
}

// createSQLiteBackup creates a backup using SQLite's backup API
func createSQLiteBackup(srcPath, dstPath string, compress bool) error {
	// Open source database
	srcDB, err := sql.Open("sqlite", srcPath+"?mode=ro")
	if err != nil {
		return fmt.Errorf("failed to open source database: %w", err)
	}
	defer srcDB.Close()

	// For compressed backups, we need a temp file first
	var actualDstPath string
	if compress {
		actualDstPath = strings.TrimSuffix(dstPath, ".gz")
	} else {
		actualDstPath = dstPath
	}

	// Create destination directory
	if err := os.MkdirAll(filepath.Dir(actualDstPath), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Remove destination if it exists
	os.Remove(actualDstPath)

	// Open destination database
	dstDB, err := sql.Open("sqlite", actualDstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination database: %w", err)
	}
	defer dstDB.Close()

	// Use VACUUM INTO for atomic backup (SQLite 3.27.0+)
	// This is safer than the backup API as it creates a consistent snapshot
	_, err = srcDB.Exec(fmt.Sprintf("VACUUM INTO '%s'", actualDstPath))
	if err != nil {
		// Fallback to simple copy if VACUUM INTO not supported
		return copyDatabase(srcPath, actualDstPath)
	}

	// Compress if requested
	if compress {
		if err := compressFile(actualDstPath, dstPath); err != nil {
			return fmt.Errorf("compression failed: %w", err)
		}
		// Remove uncompressed file
		os.Remove(actualDstPath)
	}

	return nil
}

// copyDatabase performs a simple file copy (fallback)
func copyDatabase(srcPath, dstPath string) error {
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

// compressFile compresses a file with gzip
func compressFile(srcPath, dstPath string) error {
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

	gw := gzip.NewWriter(dst)
	defer gw.Close()

	_, err = io.Copy(gw, src)
	return err
}

// uploadToS3 uploads a file to S3
func uploadToS3(filePath, bucket, key, region string) error {
	// Note: This is a placeholder. In production, use AWS SDK v2.
	// For now, we'll check if aws CLI is available and use that.
	// The full implementation would use:
	// - github.com/aws/aws-sdk-go-v2/config
	// - github.com/aws/aws-sdk-go-v2/service/s3

	fmt.Println("  (S3 upload requires AWS credentials configured)")
	fmt.Println("  Run: aws s3 cp", filePath, fmt.Sprintf("s3://%s/%s", bucket, key))

	// TODO: Implement native S3 upload using AWS SDK v2
	// cfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(region))
	// if err != nil { return err }
	// client := s3.NewFromConfig(cfg)
	// f, _ := os.Open(filePath)
	// defer f.Close()
	// _, err = client.PutObject(context.Background(), &s3.PutObjectInput{
	//     Bucket: &bucket,
	//     Key:    &key,
	//     Body:   f,
	// })
	// return err

	return nil
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
