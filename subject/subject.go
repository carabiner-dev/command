// SPDX-FileCopyrightText: Copyright 2026 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

// Package subject provides a reusable options set to read attestation
// subjects from the command line. A subject can be given as a digest
// (`algorithm:value`, using any algorithm name from the in-toto DigestSet
// spec) or as the path to a file, which is hashed.
package subject

import (
	"errors"
	"fmt"
	"os"

	"github.com/carabiner-dev/hasher"
	intoto "github.com/in-toto/attestation/go/v1"
	"github.com/spf13/cobra"

	"github.com/carabiner-dev/command"
)

var _ command.OptionsSet = &Options{}

// Options reads subjects from the command line.
type Options struct {
	config *command.OptionsSetConfig

	// Subjects holds the raw values read from the flag: digests in
	// `algorithm:value` form or paths to files.
	Subjects []string

	// Algorithms lists the algorithms used to hash subjects given as files.
	// When empty, the hasher's defaults are used.
	Algorithms []intoto.HashAlgorithm
}

// Config returns the flag configuration for the subject options.
func (so *Options) Config() *command.OptionsSetConfig {
	if so.config == nil {
		so.config = &command.OptionsSetConfig{
			Flags: map[string]command.FlagConfig{
				"subject": {
					Short: "s",
					Long:  "subject",
					Help:  "subject digest (algorithm:value) or path to a file to hash",
				},
			},
		}
	}
	return so.config
}

// AddFlags adds the subject flags to a command.
func (so *Options) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringSliceVarP(
		&so.Subjects,
		so.Config().LongFlag("subject"),
		so.Config().ShortFlag("subject"),
		so.Subjects,
		so.Config().HelpText("subject"),
	)
}

// Validate checks that every subject is either a well formed digest or an
// existing file. Strings shaped like a digest are never looked up as files.
func (so *Options) Validate() error {
	errs := []error{}
	for _, s := range so.Subjects {
		if LooksLikeDigest(s) {
			if _, _, err := ParseDigest(s); err != nil {
				errs = append(errs, fmt.Errorf("invalid subject %q: %w", s, err))
			}
			continue
		}
		if _, err := os.Stat(s); err != nil {
			errs = append(errs, fmt.Errorf("subject %q is neither a digest nor a readable file: %w", s, err))
		}
	}
	return errors.Join(errs...)
}

// Digests returns the subjects given as digests, keyed by their raw string.
func (so *Options) Digests() (map[string]hasher.HashSet, error) {
	ret := map[string]hasher.HashSet{}
	for _, s := range so.Subjects {
		if !LooksLikeDigest(s) {
			continue
		}
		algo, value, err := ParseDigest(s)
		if err != nil {
			return nil, fmt.Errorf("invalid subject %q: %w", s, err)
		}
		ret[s] = hasher.HashSet{algo: value}
	}
	return ret, nil
}

// Files returns the subjects that are not digests, ie the paths to hash.
func (so *Options) Files() []string {
	ret := []string{}
	for _, s := range so.Subjects {
		if !LooksLikeDigest(s) {
			ret = append(ret, s)
		}
	}
	return ret
}

// HashSets resolves all subjects into a hasher.FileHashSet keyed by the raw
// subject string: digests are parsed and files are hashed.
func (so *Options) HashSets() (hasher.FileHashSet, error) {
	ret, err := so.Digests()
	if err != nil {
		return nil, err
	}

	files := so.Files()
	if len(files) == 0 {
		return ret, nil
	}

	h := hasher.New()
	if len(so.Algorithms) > 0 {
		h.Options.Algorithms = so.Algorithms
	}
	hashed, err := h.HashFiles(files)
	if err != nil {
		return nil, fmt.Errorf("hashing subject files: %w", err)
	}
	for path, hs := range *hashed {
		ret[path] = hs
	}
	return ret, nil
}

// ResourceDescriptors resolves all subjects into in-toto resource descriptors,
// in the order they were given on the command line.
func (so *Options) ResourceDescriptors() ([]*intoto.ResourceDescriptor, error) {
	sets, err := so.HashSets()
	if err != nil {
		return nil, err
	}
	ret := make([]*intoto.ResourceDescriptor, 0, len(so.Subjects))
	for _, s := range so.Subjects {
		hs, ok := sets[s]
		if !ok {
			return nil, fmt.Errorf("subject %q was not resolved", s)
		}
		ret = append(ret, hs.ToResourceDescriptor())
	}
	return ret, nil
}
