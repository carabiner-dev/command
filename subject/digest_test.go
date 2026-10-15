// SPDX-FileCopyrightText: Copyright 2026 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package subject

import (
	"errors"
	"strings"
	"testing"

	intoto "github.com/in-toto/attestation/go/v1"
)

const (
	sha1Hex   = "16bc446f94bcef9093785f81897d2cfcbbcbe238"
	sha256Hex = "7355804febd9ee8357b737b61ea7ca0f1b3ff5c898bc8f003a00beabba7e91fd"
)

func TestParseDigest(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name      string
		in        string
		algo      intoto.HashAlgorithm
		value     string
		notDigest bool
		mustErr   bool
	}{
		{"sha256", "sha256:" + sha256Hex, intoto.AlgorithmSHA256, sha256Hex, false, false},
		{"sha1", "sha1:" + sha1Hex, intoto.AlgorithmSHA1, sha1Hex, false, false},
		{"gitCommit", "gitCommit:" + sha1Hex, intoto.AlgorithmGitCommit, sha1Hex, false, false},
		{"gitTree", "gitTree:" + sha1Hex, intoto.AlgorithmGitTree, sha1Hex, false, false},
		{"gitBlob", "gitBlob:" + sha1Hex, intoto.AlgorithmGitBlob, sha1Hex, false, false},
		{"gitTag", "gitTag:" + sha1Hex, intoto.AlgorithmGitTag, sha1Hex, false, false},
		{"dirHash", "dirHash:" + sha256Hex, intoto.AlgorithmDirHash, sha256Hex, false, false},
		{"sha3_256", "sha3_256:" + sha256Hex, intoto.AlgorithmSHA3_256, sha256Hex, false, false},
		{"sha512_256", "sha512_256:" + sha256Hex, intoto.AlgorithmSHA512_256, sha256Hex, false, false},
		{"md5", "md5:" + sha256Hex[:32], intoto.AlgorithmMD5, sha256Hex[:32], false, false},
		{"sha512", "sha512:" + sha256Hex + sha256Hex, intoto.AlgorithmSHA512, sha256Hex + sha256Hex, false, false},
		{"case insensitive algorithm", "GITCOMMIT:" + sha1Hex, intoto.AlgorithmGitCommit, sha1Hex, false, false},
		{"uppercase hex is lowercased", "sha1:" + strings.ToUpper(sha1Hex), intoto.AlgorithmSHA1, sha1Hex, false, false},
		{"no colon", "somefile.txt", "", "", true, true},
		{"unknown algorithm", "sha3:" + sha256Hex, "", "", true, true},
		{"path with colon", "C:\\file", "", "", true, true},
		{"empty value", "sha256:", "", "", false, true},
		{"not hex", "sha256:" + sha256Hex[:63] + "z", "", "", false, true},
		{"wrong length", "sha256:" + sha1Hex, "", "", false, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			algo, value, err := ParseDigest(tt.in)
			if tt.mustErr {
				if err == nil {
					t.Fatalf("expected error parsing %q", tt.in)
				}
				if errors.Is(err, ErrNotDigest) != tt.notDigest {
					t.Fatalf("ErrNotDigest mismatch for %q: %v", tt.in, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if algo != tt.algo || value != tt.value {
				t.Fatalf("got %s:%s, want %s:%s", algo, value, tt.algo, tt.value)
			}
		})
	}
}

func TestParseDigestCoversSpec(t *testing.T) {
	t.Parallel()
	// Every algorithm the in-toto library knows must round trip through the
	// parser with its canonical name.
	for name, algo := range intoto.HashAlgorithms {
		value := strings.Repeat("ab", algo.HexLength())
		if value == "" {
			value = "abcd"
		}
		got, _, err := ParseDigest(name + ":" + value)
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		if got != algo {
			t.Errorf("%s: got %s", name, got)
		}
	}
}

func TestLooksLikeDigest(t *testing.T) {
	t.Parallel()
	for in, want := range map[string]bool{
		"sha256:" + sha256Hex:  true,
		"sha256:short":         true, // shape matters, not validity
		"GitCommit:" + sha1Hex: true,
		"sha3:" + sha256Hex:    false,
		"file.txt":             false,
		"dir/file:name":        false,
		"":                     false,
	} {
		if got := LooksLikeDigest(in); got != want {
			t.Errorf("LooksLikeDigest(%q) = %v, want %v", in, got, want)
		}
	}
}
