// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"errors"
	"fmt"
	"os"

	"github.com/carabiner-dev/signer/key"
	"github.com/spf13/cobra"
)

var _ OptionsSet = &KeyOptions{}

type KeyOptions struct {
	PublicKeyPaths []string
}

// AddFlags adds the options flags to a command
func (ko *KeyOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringSliceVarP(
		&ko.PublicKeyPaths, "key", "k", []string{}, // this should be a public constant
		"path to public key files",
	)
}

// Verify checks the options. Key files are verified to check if they exist
func (ko *KeyOptions) Validate() error {
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

// ParseKeys parses the key files and returns a slice of public key
// providers.
func (ko *KeyOptions) ParseKeys() ([]key.PublicKeyProvider, error) {
	parser := key.NewParser()
	r := []key.PublicKeyProvider{}
	for _, path := range ko.PublicKeyPaths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading key file: %w", err)
		}
		k, err := parser.ParsePublicKey(data)
		if err != nil {
			return nil, fmt.Errorf("parsing key %q: %w", path, err)
		}
		r = append(r, k)
	}
	return r, nil
}
