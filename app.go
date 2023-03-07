package main

import (
	"github.com/PhilippHeuer/rssdownloader/pkg/cmd"
	"github.com/rs/zerolog/log"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	status  = "clean"
)

// Init Hook
func init() {
	// pass version info the version cmd
	cmd.Version = version
	cmd.CommitHash = commit
	cmd.BuildAt = date
	cmd.RepositoryStatus = status
}

// CLI Main Entrypoint
func main() {
	cmdErr := cmd.Execute()
	if cmdErr != nil {
		log.Fatal().Err(cmdErr).Msg("cli error")
	}
}
