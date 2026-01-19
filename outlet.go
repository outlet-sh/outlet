package main

import (
	_ "embed"

	"github.com/outlet-sh/outlet/cmd"

	_ "modernc.org/sqlite"
)

//go:embed etc/outlet.yaml
var embeddedConfig []byte

func main() {
	// Pass embedded config to cmd package
	cmd.EmbeddedConfig = embeddedConfig

	// Execute cobra root command
	cmd.Execute()
}
