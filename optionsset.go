// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"github.com/spf13/cobra"
)

// OptionsSet is an interface that defines the functions options set need
// to implement to make them reusable.
type OptionsSet interface {
	AddFlags(*cobra.Command)
	Validate() error
	Config() *OptionsSetConfig
}

// OptionsSetConfig configures a flag
type OptionsSetConfig struct {
	FlagPrefix string
	Flags      map[string]FlagConfig
}

func (c *OptionsSetConfig) LongFlag(id string) string {
	if _, ok := c.Flags[id]; !ok {
		return ""
	}

	p := ""
	if c.FlagPrefix != "" {
		p = c.FlagPrefix + "-"
	}
	return p + c.Flags[id].Long
}

func (c *OptionsSetConfig) ShortFlag(id string) string {
	if _, ok := c.Flags[id]; !ok {
		return ""
	}

	return c.Flags[id].Short
}

func (c *OptionsSetConfig) HelpText(id string) string {
	if _, ok := c.Flags[id]; !ok {
		return ""
	}

	return c.Flags[id].Help
}

// FlagConfig holds the configuration of a flag
type FlagConfig struct {
	Short string
	Long  string
	Help  string
}
