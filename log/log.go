// SPDX-FileCopyrightText: Copyright 2026 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/chainguard-dev/clog"
	slogzap "github.com/samber/slog-zap/v2"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/carabiner-dev/command"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

var Levels = []string{LevelDebug, LevelInfo, LevelWarn, LevelError}

var _ command.OptionsSet = &Options{}

// Options provides logging configuration for Carabiner applications.
type Options struct {
	config   *command.OptionsSetConfig
	LogLevel string
}

// Config returns the flag configuration for logging options.
func (lo *Options) Config() *command.OptionsSetConfig {
	if lo.config == nil {
		lo.config = &command.OptionsSetConfig{
			Flags: map[string]command.FlagConfig{
				"log-level": {
					Short: "",
					Long:  "log-level",
					Help:  fmt.Sprintf("Log level %+v", Levels),
				},
			},
		}
	}
	return lo.config
}

// AddFlags adds the logging flags to a command.
func (lo *Options) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(
		&lo.LogLevel, lo.Config().LongFlag("log-level"), LevelWarn, lo.Config().HelpText("log-level"),
	)
}

// Validate checks that the log level is valid.
func (lo *Options) Validate() error {
	if lo.LogLevel == "" {
		lo.LogLevel = LevelWarn
		return nil
	}

	if !slices.Contains(Levels, lo.LogLevel) {
		return fmt.Errorf("invalid log level %q (must be one of: debug, info, warn, error)", lo.LogLevel)
	}
	return nil
}

// InitLogger creates and configures a structured logger with the specified level.
// The logger uses GCP-compatible field names for cloud deployment.
func (lo *Options) InitLogger() (*clog.Logger, error) {
	// Default to warn if not set
	level := lo.LogLevel
	if level == "" {
		level = LevelWarn
	}

	// Map log level to both slog and zap levels
	var slogLevel slog.Level
	var zapLevel zap.AtomicLevel
	switch level {
	case LevelDebug:
		slogLevel = slog.LevelDebug
		zapLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	case LevelInfo:
		slogLevel = slog.LevelInfo
		zapLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	case LevelWarn:
		slogLevel = slog.LevelWarn
		zapLevel = zap.NewAtomicLevelAt(zap.WarnLevel)
	case LevelError:
		slogLevel = slog.LevelError
		zapLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return nil, fmt.Errorf("invalid log level: %s (must be debug, info, warn, or error)", level)
	}

	// Configure production logger with GCP-compatible field names
	// This ensures structured JSON logging that GCP Cloud Logging can parse correctly
	config := zap.NewProductionConfig()
	config.Level = zapLevel

	// Use GCP-compatible field names
	config.EncoderConfig.LevelKey = "severity"                     // GCP expects "severity" not "level"
	config.EncoderConfig.TimeKey = "timestamp"                     // GCP prefers "timestamp"
	config.EncoderConfig.MessageKey = "message"                    // GCP expects "message" not "msg"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // "DEBUG", "INFO", not "debug", "info"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // ISO8601 format for timestamps

	zapLogger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer zapLogger.Sync() //nolint:errcheck

	return clog.New(
		slogzap.Option{Level: slogLevel, Logger: zapLogger}.NewZapHandler(),
	), nil
}
