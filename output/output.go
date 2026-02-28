// SPDX-FileCopyrightText: Copyright 2026 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package output

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/carabiner-dev/command"
)

var _ command.OptionsSet = &Options{}

// Options provides output file configuration for Carabiner applications.
type Options struct {
	config     *command.OptionsSetConfig
	OutputPath string
}

// Config returns the flag configuration for output options.
func (oo *Options) Config() *command.OptionsSetConfig {
	if oo.config == nil {
		oo.config = &command.OptionsSetConfig{
			Flags: map[string]command.FlagConfig{
				"output": {
					Short: "o",
					Long:  "output",
					Help:  "file path to write the output (defaults to STDOUT)",
				},
			},
		}
	}
	return oo.config
}

// AddFlags adds the output flags to a command.
func (oo *Options) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(
		&oo.OutputPath,
		oo.Config().LongFlag("output"),
		oo.Config().ShortFlag("output"),
		"",
		oo.Config().HelpText("output"),
	)
}

// Validate checks the output options.
func (oo *Options) Validate() error {
	return nil
}

// GetWriter returns a writer for the output. If no output path is set,
// returns os.Stdout. Otherwise creates and returns a file writer.
func (oo *Options) GetWriter() (io.Writer, error) {
	if oo.OutputPath == "" {
		return os.Stdout, nil
	}

	f, err := os.Create(oo.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("opening output file: %w", err)
	}
	return f, nil
}
