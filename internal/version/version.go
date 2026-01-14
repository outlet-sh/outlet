package version

// Version information - set at build time via ldflags
// Use: -ldflags "-X outlet/internal/version.Version=v1.0.0"
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)
