package cmd

import (
	"os"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
)

var (
	cfg = struct {
		LogLevel  string
		LogFormat string
		LogCaller bool
	}{}
	validLogLevels  = []string{"trace", "debug", "info", "warn", "error"}
	validLogFormats = []string{"plain", "color", "json"}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfg.LogLevel, "log-level", "info", "log level - allowed: "+strings.Join(validLogLevels, ","))
	rootCmd.PersistentFlags().StringVar(&cfg.LogFormat, "log-format", "color", "log format - allowed: "+strings.Join(validLogFormats, ","))
	rootCmd.PersistentFlags().BoolVar(&cfg.LogCaller, "log-caller", false, "include caller in log functions")
}

var rootCmd = &cobra.Command{
	Use:   `rssdownloader`,
	Short: `rss feed downloader`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// log format
		if !funk.ContainsString(validLogFormats, cfg.LogFormat) {
			log.Error().Str("current", cfg.LogFormat).Strs("valid", validLogFormats).Msg("invalid log format specified")
			os.Exit(1)
		}
		var logContext zerolog.Context
		if cfg.LogFormat == "plain" {
			logContext = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true}).With().Timestamp()
		} else if cfg.LogFormat == "color" {
			colorableOutput := colorable.NewColorableStdout()
			logContext = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: colorableOutput, NoColor: false}).With().Timestamp()
		} else if cfg.LogFormat == "json" {
			logContext = zerolog.New(os.Stderr).Output(os.Stderr).With().Timestamp()
		}
		if cfg.LogCaller {
			logContext = logContext.Caller()
		}
		log.Logger = logContext.Logger()

		// log time format
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

		// log level
		if !funk.ContainsString(validLogLevels, cfg.LogLevel) {
			log.Error().Str("current", cfg.LogLevel).Strs("valid", validLogLevels).Msg("invalid log level specified")
			os.Exit(1)
		}
		if cfg.LogLevel == "trace" {
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
		} else if cfg.LogLevel == "debug" {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else if cfg.LogLevel == "info" {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		} else if cfg.LogLevel == "warn" {
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		} else if cfg.LogLevel == "error" {
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		}

		// logging config
		log.Debug().Str("log-level", cfg.LogLevel).Str("log-format", cfg.LogFormat).Bool("log-caller", cfg.LogCaller).Msg("configured logging")
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
		os.Exit(0)
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
