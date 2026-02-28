// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var _ OptionsSet = &OutputFile{}

// OutputFile provides output file configuration for Carabiner applications.
//
// Deprecated: Use github.com/carabiner-dev/command/output.Options instead.
type OutputFile struct {
	config     *OptionsSetConfig
	OutputPath string
}

// Config returns the flag configuration for output options.
//
// Deprecated: Use github.com/carabiner-dev/command/output.Options instead.
func (oo *OutputFile) Config() *OptionsSetConfig {
	if oo.config == nil {
		oo.config = &OptionsSetConfig{
			Flags: map[string]FlagConfig{
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

// AddFlags adds the options flags to a command.
//
// Deprecated: Use github.com/carabiner-dev/command/output.Options instead.
func (oo *OutputFile) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(
		&oo.OutputPath, oo.config.LongFlag("output"),
		oo.config.ShortFlag("output"), "",
		oo.config.HelpText("output"),
	)
}

// Validate checks the output options.
//
// Deprecated: Use github.com/carabiner-dev/command/output.Options instead.
func (oo *OutputFile) Validate() error {
	return nil
}

// GetWriter returns a writer for the output.
//
// Deprecated: Use github.com/carabiner-dev/command/output.Options instead.
func (oo *OutputFile) GetWriter() (io.Writer, error) {
	if oo.OutputPath == "" {
		return os.Stdout, nil
	}

	f, err := os.Create(oo.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("opening output file: %w", err)
	}
	return f, nil
}
