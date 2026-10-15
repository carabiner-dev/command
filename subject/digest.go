// SPDX-FileCopyrightText: Copyright 2026 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package subject

import (
	"errors"
	"fmt"
	"strings"

	intoto "github.com/in-toto/attestation/go/v1"
)

// ErrNotDigest is returned by ParseDigest when the string does not have the
// shape of a digest (`algorithm:hex`) with a known algorithm. Callers use it
// to tell "this is not a digest at all" apart from "this is a broken digest".
var ErrNotDigest = errors.New("string is not an algorithm:value digest")

// canonicalAlgorithms maps the lowercase form of every algorithm name in the
// in-toto DigestSet spec to its canonical spelling (eg gitcommit -> gitCommit).
var canonicalAlgorithms = func() map[string]intoto.HashAlgorithm {
	m := make(map[string]intoto.HashAlgorithm, len(intoto.HashAlgorithms))
	for name, algo := range intoto.HashAlgorithms {
		m[strings.ToLower(name)] = algo
	}
	return m
}()

// Algorithm looks up an algorithm name case-insensitively and returns its
// canonical in-toto spelling. The boolean is false when the name is not part
// of the DigestSet spec.
func Algorithm(name string) (intoto.HashAlgorithm, bool) {
	algo, ok := canonicalAlgorithms[strings.ToLower(name)]
	return algo, ok
}

// LooksLikeDigest reports whether s has the shape `algorithm:value` with an
// algorithm from the in-toto DigestSet spec. It says nothing about whether the
// value is a valid digest: strings that look like digests are never treated as
// file paths, so a malformed digest is reported as such instead of silently
// falling through to a filesystem lookup.
func LooksLikeDigest(s string) bool {
	name, _, ok := strings.Cut(s, ":")
	if !ok {
		return false
	}
	_, ok = Algorithm(name)
	return ok
}

// ParseDigest parses a digest string in `algorithm:value` form. The algorithm
// must be one of the names defined in the in-toto DigestSet spec (matched
// case-insensitively) and the value must be hex encoded. For algorithms with a
// fixed digest size the value length is checked too.
//
// The returned algorithm is the canonical spec name and the value is
// lowercased, so digests parsed here compare equal regardless of how the user
// typed them.
func ParseDigest(s string) (intoto.HashAlgorithm, string, error) {
	name, value, ok := strings.Cut(s, ":")
	if !ok {
		return "", "", ErrNotDigest
	}
	algo, ok := Algorithm(name)
	if !ok {
		return "", "", fmt.Errorf("%w: unknown algorithm %q", ErrNotDigest, name)
	}

	value = strings.ToLower(value)
	if value == "" {
		return "", "", fmt.Errorf("%s digest has no value", algo)
	}
	for _, c := range value {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return "", "", fmt.Errorf("%s digest value is not hex encoded", algo)
		}
	}
	if l := algo.HexLength() * 2; l > 0 && len(value) != l {
		return "", "", fmt.Errorf("%s digest must be %d hex characters, got %d", algo, l, len(value))
	}
	return algo, value, nil
}
