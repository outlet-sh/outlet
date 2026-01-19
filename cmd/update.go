package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/outlet-sh/outlet/internal/version"

	"github.com/spf13/cobra"
)

const (
	// GitHubReleasesAPI is the URL for checking GitHub releases
	GitHubReleasesAPI = "https://api.github.com/repos/localrivet/outlet/releases/latest"
	// UpdateTimeout is the timeout for update operations
	UpdateTimeout = 5 * time.Minute
)

// GitHubRelease represents a GitHub release response
type GitHubRelease struct {
	TagName     string         `json:"tag_name"`
	Name        string         `json:"name"`
	Body        string         `json:"body"`
	PublishedAt string         `json:"published_at"`
	Assets      []GitHubAsset  `json:"assets"`
	HTMLURL     string         `json:"html_url"`
}

// GitHubAsset represents a release asset
type GitHubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// UpdateResponse represents the processed update information
type UpdateResponse struct {
	Version     string `json:"version"`
	ReleaseDate string `json:"release_date"`
	DownloadURL string `json:"download_url"`
	Checksum    string `json:"checksum"` // SHA256 (from checksums.txt asset)
	Changelog   string `json:"changelog"`
}

var (
	checkOnly   bool
	forceUpdate bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for and install updates",
	Long: `Check for available updates and optionally install them.

Outlet is open-source software. Updates are fetched directly from GitHub releases.

Examples:
  outlet update          # Check and install updates
  outlet update --check  # Check only, don't install
  outlet update --force  # Force update even if on latest`,
	RunE: runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVarP(&checkOnly, "check", "c", false, "Check for updates without installing")
	updateCmd.Flags().BoolVarP(&forceUpdate, "force", "f", false, "Force update even if on latest version")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	fmt.Printf("Outlet Update Checker\n")
	fmt.Printf("Current version: %s\n\n", version.Version)

	// Check for updates from GitHub
	fmt.Print("Checking for updates... ")
	update, err := checkForUpdates()
	if err != nil {
		fmt.Println("FAILED")
		return fmt.Errorf("update check failed: %w", err)
	}

	if update == nil {
		fmt.Println("Already on latest version")
		return nil
	}

	fmt.Println("UPDATE AVAILABLE")
	fmt.Printf("\n  New version: %s\n", update.Version)
	fmt.Printf("  Released: %s\n", update.ReleaseDate)
	if update.Changelog != "" {
		fmt.Printf("\n  Changelog:\n%s\n", indentText(update.Changelog, "    "))
	}

	if checkOnly {
		fmt.Println("\nRun 'outlet update' without --check to install.")
		return nil
	}

	// Confirm update
	if !forceUpdate {
		fmt.Print("\nInstall this update? [y/N]: ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("Update cancelled.")
			return nil
		}
	}

	// Download and install update
	fmt.Print("\nDownloading update... ")
	binaryPath, err := downloadUpdate(update)
	if err != nil {
		fmt.Println("FAILED")
		return fmt.Errorf("download failed: %w", err)
	}
	fmt.Println("OK")

	// Verify checksum if available
	if update.Checksum != "" {
		fmt.Print("Verifying checksum... ")
		if err := verifyChecksum(binaryPath, update.Checksum); err != nil {
			fmt.Println("FAILED")
			os.Remove(binaryPath)
			return fmt.Errorf("checksum verification failed: %w", err)
		}
		fmt.Println("OK")
	}

	// Install update
	fmt.Print("Installing update... ")
	if err := installUpdate(binaryPath); err != nil {
		fmt.Println("FAILED")
		return fmt.Errorf("installation failed: %w", err)
	}
	fmt.Println("OK")

	fmt.Printf("\nSuccessfully updated to version %s\n", update.Version)
	fmt.Println("Please restart outlet to use the new version.")

	return nil
}

func checkForUpdates() (*UpdateResponse, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("GET", GitHubReleasesAPI, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", fmt.Sprintf("Outlet/%s (%s/%s)", version.Version, runtime.GOOS, runtime.GOARCH))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// No releases yet
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned %d: %s", resp.StatusCode, string(body))
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub response: %w", err)
	}

	// Clean up version tag (remove 'v' prefix if present)
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := strings.TrimPrefix(version.Version, "v")

	// Compare versions (simple string comparison - assumes semver)
	if !forceUpdate && latestVersion <= currentVersion {
		return nil, nil
	}

	// Find the appropriate asset for this platform
	assetName := fmt.Sprintf("outlet_%s_%s", runtime.GOOS, runtime.GOARCH)
	var downloadURL string
	var checksumURL string

	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, assetName) && !strings.HasSuffix(asset.Name, ".txt") {
			downloadURL = asset.BrowserDownloadURL
		}
		if asset.Name == "checksums.txt" {
			checksumURL = asset.BrowserDownloadURL
		}
	}

	if downloadURL == "" {
		return nil, fmt.Errorf("no release asset found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	// Try to get checksum
	checksum := ""
	if checksumURL != "" {
		checksum = fetchChecksum(checksumURL, assetName)
	}

	return &UpdateResponse{
		Version:     latestVersion,
		ReleaseDate: release.PublishedAt,
		DownloadURL: downloadURL,
		Checksum:    checksum,
		Changelog:   release.Body,
	}, nil
}

func fetchChecksum(checksumURL, assetName string) string {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(checksumURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	// Parse checksums.txt format: "checksum  filename"
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if strings.Contains(line, assetName) {
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				return parts[0]
			}
		}
	}

	return ""
}

func downloadUpdate(update *UpdateResponse) (string, error) {
	client := &http.Client{Timeout: UpdateTimeout}

	resp, err := client.Get(update.DownloadURL)
	if err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download server returned %d", resp.StatusCode)
	}

	// Create temp file
	tmpFile, err := os.CreateTemp("", "outlet-update-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	// Download with progress (simple for now)
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("download interrupted: %w", err)
	}

	return tmpFile.Name(), nil
}

func verifyChecksum(filePath, expectedChecksum string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return err
	}

	actualChecksum := hex.EncodeToString(hasher.Sum(nil))
	if !strings.EqualFold(actualChecksum, expectedChecksum) {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	return nil
}

func installUpdate(newBinaryPath string) error {
	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Resolve symlinks
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	// Make new binary executable
	if err := os.Chmod(newBinaryPath, 0755); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// Create backup of current binary
	backupPath := execPath + ".backup"
	if err := os.Rename(execPath, backupPath); err != nil {
		return fmt.Errorf("failed to backup current binary: %w", err)
	}

	// Move new binary to exec path
	if err := os.Rename(newBinaryPath, execPath); err != nil {
		// Restore backup on failure
		os.Rename(backupPath, execPath)
		return fmt.Errorf("failed to install new binary: %w", err)
	}

	// Remove backup
	os.Remove(backupPath)

	return nil
}

func indentText(text, indent string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = indent + line
		}
	}
	return strings.Join(lines, "\n")
}
