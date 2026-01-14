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

	"outlet/internal/version"

	"github.com/spf13/cobra"
)

const (
	// LicenseServerURL is the URL of the license/update server
	LicenseServerURL = "https://license.outlet.sh"
	// UpdateTimeout is the timeout for update operations
	UpdateTimeout = 5 * time.Minute
)

// UpdateResponse represents the response from the update check endpoint
type UpdateResponse struct {
	Version     string `json:"version"`
	ReleaseDate string `json:"release_date"`
	DownloadURL string `json:"download_url"`
	Checksum    string `json:"checksum"` // SHA256
	Changelog   string `json:"changelog"`
	Mandatory   bool   `json:"mandatory"`
}

// LicenseValidation represents a license key validation response
type LicenseValidation struct {
	Valid        bool      `json:"valid"`
	LicenseType  string    `json:"license_type"` // personal, business, enterprise
	ExpiresAt    time.Time `json:"expires_at"`
	MaxDomains   int       `json:"max_domains"`
	UpdatesUntil time.Time `json:"updates_until"`
	ErrorMessage string    `json:"error,omitempty"`
}

var (
	licenseKey  string
	checkOnly   bool
	forceUpdate bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for and install updates",
	Long: `Check for available updates and optionally install them.

Requires a valid license key for update access. Set the license key via:
  - OUTLET_LICENSE_KEY environment variable
  - --license-key flag

Examples:
  outlet update                    # Check and install updates
  outlet update --check            # Check only, don't install
  outlet update --force            # Force update even if on latest
  outlet update --license-key=xxx  # Use specific license key`,
	RunE: runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(&licenseKey, "license-key", "l", "", "License key for update access")
	updateCmd.Flags().BoolVarP(&checkOnly, "check", "c", false, "Check for updates without installing")
	updateCmd.Flags().BoolVarP(&forceUpdate, "force", "f", false, "Force update even if on latest version")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	// Get license key from flag or environment
	if licenseKey == "" {
		licenseKey = os.Getenv("OUTLET_LICENSE_KEY")
	}

	if licenseKey == "" {
		return fmt.Errorf("license key required. Set OUTLET_LICENSE_KEY or use --license-key flag")
	}

	fmt.Printf("Outlet.sh Update Checker\n")
	fmt.Printf("Current version: %s\n\n", version.Version)

	// Validate license
	fmt.Print("Validating license... ")
	license, err := validateLicense(licenseKey)
	if err != nil {
		fmt.Println("FAILED")
		return fmt.Errorf("license validation failed: %w", err)
	}

	if !license.Valid {
		fmt.Println("INVALID")
		return fmt.Errorf("invalid license: %s", license.ErrorMessage)
	}

	fmt.Println("OK")
	fmt.Printf("  License type: %s\n", license.LicenseType)
	fmt.Printf("  Updates until: %s\n\n", license.UpdatesUntil.Format("2006-01-02"))

	// Check if updates are still available for this license
	if time.Now().After(license.UpdatesUntil) {
		return fmt.Errorf("update access expired on %s. Please renew your license", license.UpdatesUntil.Format("2006-01-02"))
	}

	// Check for updates
	fmt.Print("Checking for updates... ")
	update, err := checkForUpdates(licenseKey)
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

	// Verify checksum
	fmt.Print("Verifying checksum... ")
	if err := verifyChecksum(binaryPath, update.Checksum); err != nil {
		fmt.Println("FAILED")
		os.Remove(binaryPath)
		return fmt.Errorf("checksum verification failed: %w", err)
	}
	fmt.Println("OK")

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

func validateLicense(key string) (*LicenseValidation, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("GET", LicenseServerURL+"/api/v1/license/validate", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("User-Agent", fmt.Sprintf("Outlet/%s (%s/%s)", version.Version, runtime.GOOS, runtime.GOARCH))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to license server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("license server returned %d: %s", resp.StatusCode, string(body))
	}

	var license LicenseValidation
	if err := json.NewDecoder(resp.Body).Decode(&license); err != nil {
		return nil, fmt.Errorf("failed to parse license response: %w", err)
	}

	return &license, nil
}

func checkForUpdates(key string) (*UpdateResponse, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	url := fmt.Sprintf("%s/api/v1/updates/check?version=%s&os=%s&arch=%s",
		LicenseServerURL, version.Version, runtime.GOOS, runtime.GOARCH)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("User-Agent", fmt.Sprintf("Outlet/%s (%s/%s)", version.Version, runtime.GOOS, runtime.GOARCH))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to update server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		// No update available
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("update server returned %d: %s", resp.StatusCode, string(body))
	}

	var update UpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&update); err != nil {
		return nil, fmt.Errorf("failed to parse update response: %w", err)
	}

	return &update, nil
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
