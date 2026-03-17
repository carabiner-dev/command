// SPDX-FileCopyrightText: Copyright 2026 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package keys

import (
	"errors"
	"fmt"
	"os"

	"github.com/carabiner-dev/signer/key"
	"github.com/spf13/cobra"

	"github.com/carabiner-dev/command"
)

var _ command.OptionsSet = &Options{}

// Options provides key file configuration for Carabiner applications.
type Options struct {
	config         *command.OptionsSetConfig
	PublicKeyPaths []string
	extraKeys      []key.PublicKeyProvider
}

// Config returns the flag configuration for key options.
func (ko *Options) Config() *command.OptionsSetConfig {
	if ko.config == nil {
		ko.config = &command.OptionsSetConfig{
			Flags: map[string]command.FlagConfig{
				"key": {
					Short: "k",
					Long:  "key",
					Help:  "path to public key files",
				},
			},
		}
	}
	return ko.config
}

// AddFlags adds the key flags to a command.
func (ko *Options) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringSliceVarP(
		&ko.PublicKeyPaths,
		ko.Config().LongFlag("key"),
		ko.Config().ShortFlag("key"),
		[]string{},
		ko.Config().HelpText("key"),
	)
}

// Validate checks that the key files exist.
func (ko *Options) Validate() error {
	errs := []error{}
	for _, p := range ko.PublicKeyPaths {
		if _, err := os.Stat(p); err != nil {
			if os.IsNotExist(err) {
				errs = append(errs, fmt.Errorf("checking key %q: %w", p, err))
			}
		}
	}
	return errors.Join(errs...)
}

// AddKeys appends pre-parsed key providers. These are included in the
// result of ParseKeys alongside any keys parsed from file paths.
func (ko *Options) AddKeys(providers ...key.PublicKeyProvider) {
	ko.extraKeys = append(ko.extraKeys, providers...)
}

// ParseKeys parses the key files and returns a slice of public key providers.
// Any keys previously added via AddKeys are included in the result.
func (ko *Options) ParseKeys() ([]key.PublicKeyProvider, error) {
	parser := key.NewParser()
	r := append([]key.PublicKeyProvider{}, ko.extraKeys...)
	for _, path := range ko.PublicKeyPaths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading key file: %w", err)
		}
		k, err := parser.ParsePublicKeyProvider(data)
		if err != nil {
			return nil, fmt.Errorf("parsing key %q: %w", path, err)
		}
		r = append(r, k)
	}
	return r, nil
}
