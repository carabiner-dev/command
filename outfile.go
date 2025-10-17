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

type OutputFile struct {
	config     *OptionsSetConfig
	OutputPath string
}

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

// AddFlags adds the options flags to a command
func (oo *OutputFile) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(
		&oo.OutputPath, oo.config.LongFlag("output"),
		oo.config.ShortFlag("output"), "",
		oo.config.HelpText("output"),
	)
}

func (oo *OutputFile) Validate() error {
	return nil
}

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
